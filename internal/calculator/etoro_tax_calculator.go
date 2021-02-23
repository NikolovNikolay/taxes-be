package calculator

import (
	"github.com/sirupsen/logrus"
	"taxes-be/internal/conversion"
	"taxes-be/internal/core"
	"time"
)

type etoroTaxCalculator struct {
	es *conversion.ExchangeRateService
}

func NewEТoroTaxCalculator(es *conversion.ExchangeRateService) Calculator {
	return &etoroTaxCalculator{
		es: es,
	}
}

func (c etoroTaxCalculator) CalculateYear(report *core.Report, year int) error {
	var gtb float64
	var totalBuyAmount float64
	var totalSellAmount float64
	dividends := make(map[int64]map[string]*core.Dividend)

	for _, a := range report.Activities {

		if a.Type == core.RolloverFee {
			r := c.es.GetRateForDate(a.Date, core.Bgn)
			addToDividends(dividends, a, r)
			continue
		}

		or := c.es.GetRateForDate(a.OpenDate, core.Bgn)
		cr := c.es.GetRateForDate(a.ClosedDate, core.Bgn)

		if a.OpenDate.Year() == year && a.ClosedDate.Year() == year+1 {
			op := &core.OpenPosition{
				Date:        a.OpenDate,
				Units:       a.Units,
				PriceHome:   a.Amount * or,
				PriceOrigin: a.Amount,
				Token:       a.Token,
				Name:        a.Name,
			}
			report.OpenPositions = append(report.OpenPositions, op)
		}

		if a.OpenDate.Year() != year || a.ClosedDate.Year() != year {
			logrus.Debug("activity profit/loss not in current year: ", a)
			continue
		}

		opBuy := (a.Units * a.OpenRate) * or
		opSell := (a.Units * a.ClosedRate) * cr
		totalBuyAmount += opBuy
		totalSellAmount += opSell

		gtb += opSell - opBuy
	}

	report.Amounts = core.Amounts{
		TotalBuy:  totalBuyAmount,
		TotalSell: totalSellAmount,
	}
	report.Tax = gtb * 0.1

	postProcessOpenPositions(report, year)

	dm := summarizeDividends(dividends)
	report.Dividends = dm

	return nil
}

func postProcessOpenPositions(report *core.Report, year int) {
	tm := make(map[string][]*core.OpenPosition)
	for i := range report.OpenPositions {
		e := report.OpenPositions[i]
		_, found := tm[e.Token]
		if !found {
			tm[e.Token] = make([]*core.OpenPosition, 0)
			tm[e.Token] = append([]*core.OpenPosition{}, e)
			continue
		}
		tm[e.Token] = append(tm[e.Token], e)
	}

	report.OpenPositions = make([]*core.OpenPosition, 0)
	for k := range tm {
		positionsArray := tm[k]
		units := 0.0
		homePriceSum := 0.0
		priceSum := 0.0
		var date = time.Date(year, time.December, 31, 0, 0, 0, 0, time.Local)

		for i := range positionsArray {
			e := positionsArray[i]
			if e.Date.Unix() < date.Unix() {
				date = e.Date
			}
			units += e.Units
			homePriceSum += e.PriceHome
			priceSum += e.PriceOrigin
		}

		report.OpenPositions = append(report.OpenPositions, &core.OpenPosition{
			Date:        date,
			Units:       units,
			PriceHome:   homePriceSum / float64(len(positionsArray)),
			PriceOrigin: priceSum / float64(len(positionsArray)),
			Name:        positionsArray[0].Name,
			Token:       k,
		})
	}
}
