package calculator

import (
	"github.com/sirupsen/logrus"
	"math"
	"taxes-be/internal/conversion"
	"taxes-be/internal/core"
	"time"
)

const (
	dividendTaxBound = 0.05
)

type Calculator interface {
	CalculateYear(report *core.Report, year int) error
}

type positionMetadata struct {
	boughtUnits  float64
	soldUnits    float64
	lastDate     time.Time
	name         string
	homePrice    float64
	foreignPrice float64
	recordCount  int
}

func NewPositionMetadata(name string) *positionMetadata {
	pm := &positionMetadata{
		name: name,
	}
	return pm
}

type revolutTaxCalculator struct {
	es              *conversion.ExchangeRateService
	openPositionMap map[string]*positionMetadata
}

func NewRevolutTaxCalculator(es *conversion.ExchangeRateService) Calculator {
	return &revolutTaxCalculator{
		es:              es,
		openPositionMap: make(map[string]*positionMetadata),
	}
}

func (c revolutTaxCalculator) CalculateYear(report *core.Report, year int) error {
	var totalBuyAmount float64
	var totalSellAmount float64
	dividends := make(map[int64]map[string]*core.Dividend)

	for _, a := range report.Activities {

		if a.Date.Year() != year {
			logrus.Debug("activity profit/loss not in current year", a)
			continue
		}

		r := c.es.GetRateForDate(a.Date, core.Bgn)

		var opmd *positionMetadata
		if a.Type == core.Buy || a.Type == core.Sell {
			if _, ok := c.openPositionMap[a.Token]; !ok {
				c.openPositionMap[a.Token] = NewPositionMetadata(a.Name)
			}

			opmd = c.openPositionMap[a.Token]
			opmd.recordCount = opmd.recordCount + 1
			opmd.lastDate = a.Date
		}

		if a.Type == core.Buy {
			totalBuyAmount += a.Amount * r
			opmd.boughtUnits = opmd.boughtUnits + a.Units
			opmd.foreignPrice = opmd.foreignPrice + a.OpenRate
			opmd.homePrice = opmd.homePrice + (a.OpenRate * r)
		} else if a.Type == core.Sell {
			totalSellAmount += a.Amount * r
			opmd.soldUnits = opmd.soldUnits + (a.Units * -1)
			opmd.foreignPrice = opmd.foreignPrice + a.ClosedRate
			opmd.homePrice = opmd.homePrice + (a.ClosedRate * r)
		} else if a.Type == core.Div || a.Type == core.DivNra {
			addToDividends(dividends, a, r)
			continue
		}
	}

	postProcessPositionMetadata(report, c.openPositionMap, year)
	dm := summarizeDividends(dividends)
	report.Dividends = dm
	report.Amounts = core.Amounts{
		TotalSell: totalSellAmount,
		TotalBuy:  totalBuyAmount,
	}
	report.Tax = (totalSellAmount - totalBuyAmount) * 0.1

	return nil
}

func postProcessPositionMetadata(report *core.Report, pm map[string]*positionMetadata, year int) {
	op := make([]*core.OpenPosition, 0)
	for token, p := range pm {
		if p.soldUnits > p.boughtUnits || math.Abs(p.soldUnits-p.boughtUnits) <= 0.1 {
			continue
		}

		op = append(op, &core.OpenPosition{
			Date:        p.lastDate,
			Units:       p.boughtUnits - p.soldUnits,
			PriceHome:   p.homePrice / float64(p.recordCount),
			PriceOrigin: p.foreignPrice / float64(p.recordCount),
			Token:       token,
			Name:        p.name,
		})
	}
	report.OpenPositions = op
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

	if a.Type == core.Div || a.Type == core.RolloverFee {
		d[t][a.Token].Gross = rate * a.Amount
	} else if a.Type == core.DivNra {
		d[t][a.Token].Tax = rate * a.Amount
	}

	if d[t][a.Token].Tax > 0 && d[t][a.Token].Gross > 0 {
		d[t][a.Token].TaxPercent = d[t][a.Token].Tax / d[t][a.Token].Gross
	}
}

func summarizeDividends(d map[int64]map[string]*core.Dividend) core.DividendMeta {
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

			if div.TaxPercent < dividendTaxBound {
				boundTax := div.Gross * dividendTaxBound
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
