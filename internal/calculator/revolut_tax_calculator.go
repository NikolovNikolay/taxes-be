package calculator

import (
	"github.com/sirupsen/logrus"
	"taxes-be/internal/conversion"
	"taxes-be/internal/core"
)

const (
	dividendTaxBound = 0.05
)

type Calculator interface {
	CalculateYear(report *core.Report, year int) error
}

type revolutTaxCalculator struct {
	es *conversion.ExchangeRateService
}

func NewRevolutTaxCalculator(es *conversion.ExchangeRateService) Calculator {
	return &revolutTaxCalculator{
		es: es,
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

		r := c.es.GetRateForDate(a.Date, core.BGN)
		if a.Type == core.BUY {
			totalBuyAmount += a.Amount * r
		} else if a.Type == core.SELL {
			totalSellAmount += a.Amount * r
		} else if a.Type == core.DIV || a.Type == core.DIVNRA {
			addToDividends(dividends, a, r)
			continue
		}
	}

	dm := summarizeDividends(dividends)
	report.Dividends = dm
	report.Amounts = core.Amounts{
		TotalSell: totalSellAmount,
		TotalBuy:  totalBuyAmount,
	}
	report.Tax = (totalSellAmount - totalBuyAmount) * 0.1

	return nil
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
