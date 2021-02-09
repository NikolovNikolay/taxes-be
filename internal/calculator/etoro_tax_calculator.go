package calculator

import (
	"github.com/sirupsen/logrus"
	"taxes-be/internal/conversion"
	"taxes-be/internal/core"
	"time"
)

type etoroTaxCalculator struct {
	es          *conversion.ExchangeRateService
	dayMap      map[string]*balance
	tokenMap    map[string]map[string]*balance
	activityMap map[string]map[string][]core.LinkedActivity
}

func NewEТoroTaxCalculator(es *conversion.ExchangeRateService) Calculator {
	return &etoroTaxCalculator{
		es:          es,
		dayMap:      map[string]*balance{},
		tokenMap:    map[string]map[string]*balance{},
		activityMap: map[string]map[string][]core.LinkedActivity{},
	}
}

func (c etoroTaxCalculator) Calculate(report *core.Report) (float64, error) {
	return 0, nil
}

func (c etoroTaxCalculator) CalculateYear(report *core.Report, year int) (float64, error) {
	var gtb float64
	var totalBuyAmount float64
	var totalSellAmount float64
	dividends := make(map[int64]map[string]*core.Dividend)

	for _, a := range report.Activities {

		if a.Type == core.ROLLOVER_FEE {
			r := c.es.GetRateForDate(a.Date, core.BGN)
			addToDividends(dividends, a, r)
			continue
		}

		or := c.es.GetRateForDate(a.OpenDate, core.BGN)
		cr := c.es.GetRateForDate(a.ClosedDate, core.BGN)

		if a.OpenDate.Year() == year && a.ClosedDate.Year() == year+1 {
			op := &core.OpenPosition{
				Date:         a.OpenDate,
				Units:        a.Units,
				AmountHome:   a.Amount * or,
				AmountOrigin: a.Amount,
				Token:        a.Token,
				Name:         a.Name,
			}
			report.OpenPositions = append(report.OpenPositions, op)
		}

		if a.OpenDate.Year() != year || a.ClosedDate.Year() != year {
			logrus.Info("activity profit/loss not in current year: ", a)
			continue
		}

		opBuy := (a.Units * a.OpenRate) * or
		opSell := (a.Units * a.ClosedRate) * cr
		totalBuyAmount += opBuy
		totalSellAmount += opSell

		gtb += opSell - opBuy
	}

	postProcessOpenPositions(report, year)

	dm := summarizeDividends(dividends)

	logrus.Info("Total BUY amount: ", totalBuyAmount)
	logrus.Info("Total SELL amount: ", totalSellAmount)
	logrus.Info("Total TAX: ", gtb*0.1)

	logrus.Infoln("")
	// TODO
	logrus.Info("Dividends GROSS amount: ", dm.GrossAmount)
	logrus.Info("Dividends NET amount: ", dm.NetAmount)
	logrus.Info("Dividends TAX amount: ", dm.Tax)

	logrus.Infoln("")
	printOpenPositions(report.OpenPositions)

	return gtb * 0.1, nil
}

func printOpenPositions(op []*core.OpenPosition) {
	logrus.Info("Open positions:")
	for i := range op {
		logrus.Printf(
			"Date: %v, Token: %s, Name: %s, Avg Foreign price: %f, Avg Home price: %f, Units: %f",
			op[i].Date, op[i].Token, op[i].Name, op[i].AmountOrigin, op[i].AmountHome, op[i].Units,
		)
	}
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
		homeAmountSum := 0.0
		amountSum := 0.0
		var date = time.Date(year, time.December, 31, 0, 0, 0, 0, time.Local)

		for i := range positionsArray {
			e := positionsArray[i]
			if e.Date.Unix() < date.Unix() {
				date = e.Date
			}
			units += e.Units
			homeAmountSum += e.AmountHome
			amountSum += e.AmountOrigin
		}

		report.OpenPositions = append(report.OpenPositions, &core.OpenPosition{
			Date:         date,
			Units:        units,
			AmountHome:   homeAmountSum / float64(len(positionsArray)),
			AmountOrigin: amountSum / float64(len(positionsArray)),
			Name:         positionsArray[0].Name,
			Token:        k,
		})
	}
}
