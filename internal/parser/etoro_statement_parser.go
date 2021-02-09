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
	accountDetailsSheet        = "Account Details"
	closedPositionsSheet       = "Closed Positions"
	transactionsReportSheet    = "Transactions Report"
	financialSummarySheet      = "Financial Summary"
	closedPositionsColCount    = 17
	transactionsReportColCount = 9
)

var eToroDateLayouts = []string{"02.01.2006", "02/01/2006", "2.1.2006", "2006-01-02"}

type eToroStatementParserLinked struct {
}

func newEToroStatementParserLinked() Parser {
	return &eToroStatementParserLinked{}
}

func (p *eToroStatementParserLinked) Parse(lines []string) (*core.Report, error) {
	var currentSheet string
	var currentPositionID string
	var currentToken string

	expectDeposits := false
	expectStartDate := false
	expectEndDate := false
	expectClosedPositionsHeaders := true
	expectTransactionsReportHeaders := true
	currentIndex := 0

	var dividentActivity *core.LinkedActivity

	var cpm = map[string][]string{}
	report := &core.Report{}
	for _, l := range lines {
		updateStatus(l, &currentSheet)

		switch currentSheet {
		case accountDetailsSheet:
			if l == "Deposits" {
				expectDeposits = true
			} else if expectDeposits {
				dep, err := parseEТoroFloat(l)
				if err != nil {
					continue
				}
				report.Deposits += dep
			}

			if l == "Start Date" {
				expectStartDate = true
			} else if expectStartDate {
				dc := strings.Split(l, " ")
				d, err := parseEToroDate(dc[0])
				if err != nil {
					continue
				}
				report.StartDate = &d
				expectStartDate = false
			}

			if l == "End Date" {
				expectEndDate = true
			} else if expectEndDate {
				dc := strings.Split(l, " ")
				d, err := parseEToroDate(dc[0])
				if err != nil {
					continue
				}
				report.EndDate = &d
				expectEndDate = false
			}
		case closedPositionsSheet:
			currentIndex++
			if expectClosedPositionsHeaders {
				if currentIndex == closedPositionsColCount+1 {
					expectClosedPositionsHeaders = false
					currentIndex = 0
					continue
				}
				continue
			}

			if currentIndex == 1 {
				currentPositionID = l
				cpm[currentPositionID] = append([]string{}, l)
			} else if currentIndex == 2 {
				nameComponents := strings.Split(l, " ")
				var fullName string
				for i := range nameComponents {
					if i == 0 {
						continue
					}
					fullName += nameComponents[i] + " "
				}
				cpm[currentPositionID] = append(cpm[currentPositionID], strings.Trim(fullName, " "))
			} else if currentIndex < closedPositionsColCount {
				cpm[currentPositionID] = append(cpm[currentPositionID], l)
			} else {
				cpm[currentPositionID] = append(cpm[currentPositionID], l)
				currentIndex = 0
			}
		case transactionsReportSheet:
			currentIndex++
			if expectTransactionsReportHeaders {
				if currentIndex == transactionsReportColCount+1 {
					expectTransactionsReportHeaders = false
					currentIndex = 0
					continue
				}
				continue
			}

			if currentIndex == 1 {
				d, err := parseEToroDate(strings.Split(l, " ")[0])
				if err != nil {
					return nil, err
				}
				dividentActivity = &core.LinkedActivity{}
				dividentActivity.Date = d

			} else if currentIndex == 3 {
				if strings.Trim(l, " ") == string(core.ROLLOVER_FEE) {
					dividentActivity.Type = core.ROLLOVER_FEE
				}
			} else if currentIndex == 4 {
				if dividentActivity != nil && dividentActivity.Type == core.ROLLOVER_FEE {
					if strings.Trim(l, " ") != "Payment caused by dividend" {
						logrus.Info("Expected dividend but got: ", l)
						dividentActivity = nil
					}
				} else {
					currentToken = strings.Split(l, "/")[0]
				}
			} else if currentIndex == 5 {
				currentPositionID = l
				if len(cpm[currentPositionID]) < closedPositionsColCount+1 {
					cpm[currentPositionID] = append(cpm[currentPositionID], currentToken)
				}
			} else if currentIndex == 6 && dividentActivity != nil {
				am, err := parseEТoroFloat(l)
				if err != nil {
					return nil, err
				}
				dividentActivity.Amount = am
				if dividentActivity.Type == core.ROLLOVER_FEE {
					dividentActivity.Token = fmt.Sprintf("%f_%d", dividentActivity.Amount, dividentActivity.Date.Unix())
					report.Activities = append(report.Activities, *dividentActivity)
				}
				dividentActivity = nil
			} else if currentIndex < transactionsReportColCount {
				continue
			} else {
				currentIndex = 0
			}
		}
	}

	for k, v := range cpm {
		if _, err := strconv.ParseFloat(k, 64); err == nil && len(v) == closedPositionsColCount+1 {
			a := core.LinkedActivity{}

			a.Name = v[1]

			amount, err := parseEТoroFloat(v[3])
			if err != nil {
				logrus.Debug(fmt.Sprintf("could not parse number from string: %s", v[3]))
				continue
			}
			a.Amount = amount

			units, err := parseEТoroFloat(v[4])
			if err != nil {
				logrus.Debug(fmt.Sprintf("could not parse number from string: %s", v[4]))
				continue
			}
			a.Units = units

			a.Type = core.LINKED
			a.Currency = core.USD

			or, err := parseEТoroFloat(v[5])
			if err != nil {
				logrus.Debug(fmt.Sprintf("could not parse number from string: %s", v[5]))
				continue
			}
			a.OpenRate = or

			cr, err := parseEТoroFloat(v[6])
			if err != nil {
				logrus.Debug(fmt.Sprintf("could not parse number from string: %s", v[6]))
				continue
			}
			a.ClosedRate = cr

			od, err := parseEToroDate(strings.Split(v[9], " ")[0])
			if err != nil {
				logrus.Debug(fmt.Sprintf("could not parse date from string: %s", v[9]))
				continue
			}
			a.OpenDate = od

			cd, err := parseEToroDate(strings.Split(v[10], " ")[0])
			if err != nil {
				logrus.Debug(fmt.Sprintf("could not parse date from string: %s", v[10]))
				continue
			}

			a.ClosedDate = cd
			a.Date = cd
			a.Token = v[len(v)-1]
			report.Activities = append(report.Activities, a)
		}
	}

	return report, nil
}

func updateStatus(l string, currentSheet *string) {
	if l == accountDetailsSheet {
		*currentSheet = accountDetailsSheet
	} else if l == closedPositionsSheet {
		*currentSheet = closedPositionsSheet
	} else if l == transactionsReportSheet {
		*currentSheet = transactionsReportSheet
	} else if l == financialSummarySheet {
		*currentSheet = financialSummarySheet
	}
}

func parseEТoroFloat(l string) (float64, error) {
	r := strings.NewReplacer(",", ".", "(", "", ")", "", "$", "")
	l = r.Replace(l)
	f, err := strconv.ParseFloat(l, 32)
	return f, err
}

func parseEToroDate(l string) (time.Time, error) {
	var err error
	var t time.Time
	for i := range eToroDateLayouts {
		t, err = time.Parse(eToroDateLayouts[i], l)
		if err == nil {
			return t, nil
		}
	}

	return t, err
}
