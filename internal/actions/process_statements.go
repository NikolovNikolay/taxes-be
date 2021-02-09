package actions

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"taxes-be/internal/calculator"
	"taxes-be/internal/conversion"
	"taxes-be/internal/core"
	"taxes-be/internal/parser"
	"taxes-be/internal/reader"

	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	pdf  = ".pdf"
	xlsx = ".xlsx"
)

var (
	supportedFormats = map[string]reader.Reader{
		pdf:  reader.NewPDFReader(),
		xlsx: reader.NewExcelReader(),
	}
)

func ProcessStatements(path string, year int) error {
	var revolutReport = &core.Report{}
	var etoroReport = &core.Report{}

	var deposits float64
	if path != "" {
		err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			var lines []string
			ext := filepath.Ext(path)
			if r, ok := supportedFormats[ext]; ok {
				logrus.Info(fmt.Sprintf(`Reading file "%s"`, info.Name()))
				lines, err = r.Read(path)
				if err != nil {
					return err
				}
			} else {
				logrus.Warn(fmt.Sprintf(`File extension not supported: "%s", "%s"`, info.Name(), ext))
				return nil
			}

			st := getStatementType(lines)
			pf := parser.NewParserFactory()
			p, err := pf.Build(st)
			if err != nil {
				return err
			}

			if st == reader.Revolut && year <= 0 {
				return fmt.Errorf("invalid year flag")
			}
			r, err := p.Parse(lines)
			if err != nil {
				return err
			}
			deposits += r.Deposits

			if st == reader.Revolut {
				revolutReport.Activities = append(revolutReport.Activities, r.Activities...)
			} else if st == reader.EToro {
				etoroReport.Activities = append(etoroReport.Activities, r.Activities...)
				etoroReport.OpenPositions = append(etoroReport.OpenPositions, r.OpenPositions...)
			}
			return nil
		})

		if err != nil {
			return err
		}

		sort.Slice(revolutReport.Activities, func(i, j int) bool {
			return revolutReport.Activities[i].Date.UnixNano() < revolutReport.Activities[j].Date.UnixNano()
		})
		sort.Slice(etoroReport.Activities, func(i, j int) bool {
			return etoroReport.Activities[i].Date.UnixNano() < etoroReport.Activities[j].Date.UnixNano()
		})

		//var rr string
		//if len(revolutReport.Activities) > 0 {
		//	rStart := revolutReport.Activities[0].Date
		//	rEnd := revolutReport.Activities[len(revolutReport.Activities)-1].Date
		//	rr = fmt.Sprintf("%s-%s", rStart.String(), rEnd.String())
		//}

		//var er string
		//if len(etoroReport.Activities) > 0 {
		//	eStart := etoroReport.Activities[0].Date
		//	eEnd := etoroReport.Activities[len(etoroReport.Activities)-1].Date
		//	er = fmt.Sprintf("%s-%s", eStart.String(), eEnd.String())
		//}

		var s, e time.Time

		s, e = getRange(revolutReport, etoroReport)
		rs := conversion.NewExchangeRateService(
			s.AddDate(-1, 0, 0).Format("2006-01-02"),
			e.Format("2006-01-02"),
		)

		rtc := calculator.NewRevolutTaxCalculator(rs)
		etc := calculator.NewEТoroTaxCalculator(rs)

		_, err = rtc.CalculateYear(revolutReport, year)
		if err != nil {
			return err
		}

		_, err = etc.CalculateYear(etoroReport, year)
		if err != nil {
			return err
		}

		//logrus.Info(fmt.Sprintf(`Tax for period Revolut "%s": %f`, rr, rTax))
		//logrus.Info(fmt.Sprintf(`Tax for period eToro "%s": %f`, er, eTax))
		//logrus.Info(fmt.Sprintf(`Total Tax": %f`, rTax+eTax))
	}

	return nil
}

func getStatementType(lines []string) reader.StatementType {
	for _, l := range lines {
		if strings.Contains(l, "Revolut Trading Ltd") {
			return reader.Revolut
		}
		if strings.Contains(l, "eToro (Europe) Ltd") || strings.Contains(l, "eToro (UK) Ltd") {
			return reader.EToro
		}
	}
	return reader.Unknown
}

func getRange(revolutReport, etoroReport *core.Report) (start, end time.Time) {
	rActivities := revolutReport.Activities
	eActivities := etoroReport.Activities
	if len(rActivities) > 0 && len(eActivities) == 0 {
		return rActivities[0].Date, rActivities[len(rActivities)-1].Date
	} else if len(eActivities) > 0 && len(rActivities) == 0 {
		return eActivities[0].Date, eActivities[len(eActivities)-1].Date
	} else if len(rActivities) > 0 && len(eActivities) > 0 {
		es := eActivities[0].Date
		rs := rActivities[0].Date

		var s time.Time
		if es.UnixNano() < rs.UnixNano() {
			s = es
		} else {
			s = rs
		}

		ee := eActivities[len(eActivities)-1].Date
		re := rActivities[len(rActivities)-1].Date

		var e time.Time
		if ee.UnixNano() > re.UnixNano() {
			e = ee
		} else {
			e = re
		}

		return s, e
	} else {
		return time.Now(), time.Now()
	}
}
