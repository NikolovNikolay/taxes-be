package parser

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"taxes-be/internal/core"
	"time"
)

const (
	activityStart = "ACTIVITY"
	dateLayout    = "01/02/2006"
	activityCols  = 8
)

type revolutStatementParser struct {
}

func newRevolutStatementParser() Parser {
	return &revolutStatementParser{}
}

func (p *revolutStatementParser) Parse(lines []string) (*core.Report, error) {
	report := &core.Report{}

	expectActivities := false
	inActivity := false
	wait := false
	isSkipCol := false
	expectDeposit := false
	var deposits float64

	var a core.LinkedActivity
	currentCol := 1
	for _, l := range lines {
		if strings.Contains(l, "Deposits") {
			expectDeposit = true
			continue
		}
		if expectDeposit {
			dep, err := parseFloat(l)
			if err != nil {
				continue
			}
			deposits += dep
			expectDeposit = false
		}
		if wait {
			if currentCol == activityCols {
				wait = false
				currentCol = 1
				continue
			}
			currentCol++
			continue
		}

		if currentCol == getEndCol() {
			currentCol = 1
			//if !isSkipCol {
			report.Activities = append(report.Activities, a)
			//}
			isSkipCol = false
			inActivity = false
		}

		if strings.Contains(l, activityStart) {
			expectActivities = true
			wait = true
			continue
		}

		if expectActivities {
			if !inActivity {
				od, err := parseDate(l)
				if err == nil {
					inActivity = true
					a = core.LinkedActivity{}
					a.Activity.OpenDate = od
					currentCol++
				}
				continue
			} else {
				switch currentCol {
				case 2:
					t, err := parseDate(l)
					if err != nil {
						return nil, err
					}
					a.Date = t
					currentCol++
					continue
				case 3:
					c := parseCurrency(l)
					a.Currency = c
					currentCol++
					continue
				case 4:
					at := parseActivityType(l)
					a.Type = at
					currentCol++

					if at == core.Cdep || at == core.Csd {
						isSkipCol = true
					}
					continue
				case 5:
					if a.Type != core.Cdep && a.Type != core.Csd {
						s := strings.Split(l, " ")
						a.Token = s[0]

						s = strings.Split(l, " - ")
						a.Name = s[1]
					}
					currentCol++
					continue
				case 6:
					units, err := parseFloat(l)
					if err != nil {
						logrus.Debug(fmt.Sprintf("could not parse number from string: %s", l))
						continue
					}
					a.Units = units
					currentCol++
					if isSkipCol {
						i := fmt.Sprintf("%.2f", a.Units)
						v, _ := strconv.ParseFloat(i, 2)
						a.Amount = v
						a.Units = 0
						currentCol = getEndCol()
					}
					continue
				case 7:
					price, err := parseFloat(l)
					if err != nil {
						return nil, err
					}
					if a.Type == core.Buy {
						a.OpenRate = price
					} else if a.Type == core.Sell {
						a.ClosedRate = price
					}
					currentCol++
					continue
				case 8:
					if isSkipCol {
						i := fmt.Sprintf("%.2f", a.Units)
						v, _ := strconv.ParseFloat(i, 2)
						a.Amount = v
					} else {
						amount, err := parseFloat(l)
						if err != nil {
							return nil, err
						}
						i := fmt.Sprintf("%.2f", amount)
						v, _ := strconv.ParseFloat(i, 2)
						a.Amount = v
					}
					currentCol++
					continue
				default:
					currentCol++
					continue
				}
			}
		}
	}
	return report, nil
}

func getEndCol() int {
	return activityCols + 1
}

func parseDate(l string) (time.Time, error) {
	return time.Parse(dateLayout, l)
}

func parseCurrency(l string) core.Currency {
	return core.Currency(l)
}

func parseActivityType(l string) core.ActivityType {
	return core.ActivityType(l)
}

func parseFloat(l string) (float64, error) {
	r := strings.NewReplacer(",", "", "(", "", ")", "", "$", "")
	l = r.Replace(l)
	f, err := strconv.ParseFloat(l, 32)
	return f, err
}
