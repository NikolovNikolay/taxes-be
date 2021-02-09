package calculator

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"sort"
	"taxes-be/internal/conversion"
	"taxes-be/internal/core"
	"time"
)

const (
	dateLayout = "2006-01-02"
)

type Calculator interface {
	Calculate(report *core.Report) (float64, error)
	CalculateYear(report *core.Report, year int) (float64, error)
}

type balance struct {
	date           time.Time
	operations     int
	balance        float64
	taxableBalance float64
	rate           float64
	boughtUnits    float64
	soldUnits      float64
}

type revolutTaxCalculator struct {
	es          *conversion.ExchangeRateService
	dayMap      map[string]*balance
	tokenMap    map[string]map[string]*balance
	activityMap map[string]map[string][]core.LinkedActivity
}

func NewRevolutTaxCalculator(es *conversion.ExchangeRateService) Calculator {
	return &revolutTaxCalculator{
		es:          es,
		dayMap:      map[string]*balance{},
		tokenMap:    map[string]map[string]*balance{},
		activityMap: map[string]map[string][]core.LinkedActivity{},
	}
}

func (c revolutTaxCalculator) CalculateYear(report *core.Report, year int) (float64, error) {
	var totalBuyAmount float64
	var totalSellAmount float64
	dividents := make(map[int64]map[string]*core.Dividend)

	for _, a := range report.Activities {
		d := a.Date.Format(dateLayout)

		if a.Date.Year() != year {
			logrus.Info("activity profit/loss not in current year", a)
			continue
		}

		if _, ok := c.tokenMap[a.Token]; !ok {
			c.tokenMap[a.Token] = map[string]*balance{}
			c.activityMap[a.Token] = map[string][]core.LinkedActivity{}
		}

		r := c.es.GetRateForDate(a.Date, core.BGN)
		if a.Type == core.BUY {
			totalBuyAmount += a.Amount * r
		} else if a.Type == core.SELL {
			totalSellAmount += a.Amount * r
		} else if a.Type == core.DIV || a.Type == core.DIVNRA {
			addToDividends(dividents, a, r)
			continue
		}

		c.tokenMap[a.Token][d] = &balance{
			date: a.Date,
			rate: r,
		}

		if _, ok := c.activityMap[a.Token][d]; !ok {
			c.activityMap[a.Token][d] = []core.LinkedActivity{}
		}

		c.activityMap[a.Token][d] = append(c.activityMap[a.Token][d], a)
	}

	dm := summarizeDividends(dividents)

	logrus.Info("Total BUY amount: ", totalBuyAmount)
	logrus.Info("Total SELL amount: ", totalSellAmount)
	logrus.Info("Total TAX: ", (totalSellAmount-totalBuyAmount)*0.1)

	logrus.Infoln("")
	logrus.Info("Dividends GROSS amount: ", dm.NetAmount)
	logrus.Info("Dividends NET amount: ", dm.NetAmount)
	logrus.Info("Dividends TAX amount: ", dm.Tax)

	return (totalSellAmount - totalBuyAmount) * 0.1, nil
}

func (c revolutTaxCalculator) Calculate(report *core.Report) (float64, error) {
	var gtb float64
	for _, a := range report.Activities {
		d := a.Date.Format(dateLayout)

		if _, ok := c.tokenMap[a.Token]; !ok {
			c.tokenMap[a.Token] = map[string]*balance{}
			c.activityMap[a.Token] = map[string][]core.LinkedActivity{}
		}

		r := c.es.GetRateForDate(a.Date, core.BGN)
		c.tokenMap[a.Token][d] = &balance{
			date: a.Date,
			rate: r,
		}

		if _, ok := c.activityMap[a.Token][d]; !ok {
			c.activityMap[a.Token][d] = []core.LinkedActivity{}
		}

		c.activityMap[a.Token][d] = append(c.activityMap[a.Token][d], a)
	}

	for token, am := range c.activityMap {
		dKeys := make([]string, 0, len(am))
		for k := range am {
			dKeys = append(dKeys, k)
		}

		sort.Strings(dKeys)
		isFirstEntry := true
		for i, date := range dKeys {
			for _, a := range c.activityMap[token][date] {

				if isFirstEntry && a.Type == core.SELL {
					c.tokenMap[token][date].taxableBalance += a.Amount
					logrus.Info(fmt.Sprintf("Closed a last year position: %s [%s]", a.Token, a.Date.String()))
				}

				b := c.tokenMap[token][date]

				if a.Type == core.SELL {
					b.soldUnits -= a.Units
					b.balance += a.Amount * b.rate
				} else if a.Type == core.BUY {
					b.boughtUnits += a.Units
					b.balance -= a.Amount * b.rate
				}
				isFirstEntry = false
			}

			if i == len(dKeys)-1 {
				var bal float64
				var bUnits float64
				var sUnits float64
				for _, b := range c.tokenMap[token] {
					bal += b.balance
					bUnits += b.boughtUnits
					sUnits += b.soldUnits
					if b.taxableBalance > 0 {
						bal += b.taxableBalance
					}
				}

				if bUnits > 0 && sUnits == 0 {
					continue
				} else if bUnits > 0 && sUnits < bUnits && bUnits-sUnits > 0.05 {
					b := c.activityMap[token][date][len(c.activityMap[token][date])-1]
					var price float64
					if b.ClosedRate > 0 {
						price = b.ClosedRate
					} else {
						price = b.OpenRate
					}
					bal += price * (bUnits - sUnits)
				}

				gtb += bal
			}
		}
	}

	return gtb * 0.1, nil
}

func addToDividends(d map[int64]map[string]*core.Dividend, a core.LinkedActivity, rate float64) {
	t := a.Date.Unix()
	_, found := d[t]
	if !found {
		d[t] = make(map[string]*core.Dividend)
	}

	if _, ok := d[t][a.Token]; !ok {
		d[t][a.Token] = &core.Dividend{}
	}

	if a.Type == core.DIV || a.Type == core.ROLLOVER_FEE {
		d[t][a.Token].Gross = rate * a.Amount
	} else if a.Type == core.DIVNRA {
		d[t][a.Token].Tax = rate * a.Amount
	}

	if d[t][a.Token].Tax > 0 && d[t][a.Token].Gross > 0 {
		d[t][a.Token].TaxPercent = d[t][a.Token].Tax / d[t][a.Token].Gross
	}
}

func summarizeDividends(d map[int64]map[string]*core.Dividend) core.DividendMeta {
	const divTaxBound = 0.05
	grossDiv := 0.0
	netDiv := 0.0
	taxDiv := 0.0
	for date := range d {
		m := d[date]
		for token := range m {
			div := m[token]
			if div.Gross < 0 {
				continue
			}
			grossDiv += div.Gross

			if div.TaxPercent < divTaxBound {
				boundTax := div.Gross * divTaxBound
				div.Tax += boundTax - div.Tax
			}

			taxDiv += div.Tax
			netDiv += div.Gross - div.Tax
		}
	}

	return core.DividendMeta{
		GrossAmount: grossDiv,
		NetAmount:   netDiv,
		Tax:         taxDiv,
	}
}
