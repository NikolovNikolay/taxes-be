package statements

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"taxes-be/internal/calculator"
	"taxes-be/internal/conversion"
	"taxes-be/internal/core"
	"taxes-be/internal/parser"
	"taxes-be/internal/reader"
	"taxes-be/utils/files"

	"sort"
	"strings"
	"time"
)

const (
	pdf  = "pdf"
	xlsx = "xlsx"
	xls  = "xls"
)

var (
	excelReader      = reader.NewExcelReader()
	pdfReader        = reader.NewPDFReader()
	supportedFormats = map[string]reader.Reader{
		pdf:  pdfReader,
		xlsx: excelReader,
		xls:  excelReader,
	}
)

type ReportProcessor struct {
	year     int
	report   *core.Report
	deposits float64
	rType    core.StatementType
}

func NewReportProcessor(year int, id string) *ReportProcessor {
	return &ReportProcessor{
		year: year,
		report: &core.Report{
			RequestID: id,
		},
	}
}

func (p *ReportProcessor) ParseLines(fileName string, sType int) error {
	var lines []string
	var err error
	ext := files.GetExtensionFromName(fileName)
	if r, ok := supportedFormats[ext]; ok {
		logrus.Info(fmt.Sprintf(`Reading file "%s"`, fileName))
		lines, err = r.Read(fileName)
		if err != nil {
			return err
		}
	} else {
		logrus.Warn(fmt.Sprintf(`File extension not supported: "%s", "%s"`, fileName, ext))
		return fmt.Errorf("file format not supported: %s", fileName)
	}

	p.rType = getStatementType(lines)
	if p.rType == core.Unknown || int(p.rType) != sType {
		return fmt.Errorf("can't process statement - ambiguous type")
	}

	pf := parser.NewParserFactory()
	parsFactory, err := pf.Build(p.rType)
	if err != nil {
		return err
	}

	if p.rType == core.Revolut && p.year <= 0 {
		return fmt.Errorf("invalid year")
	}
	r, err := parsFactory.Parse(lines)
	if err != nil {
		return err
	}
	p.deposits += r.Deposits

	p.report.Activities = append(p.report.Activities, r.Activities...)
	if p.rType == core.EToro {
		p.report.OpenPositions = append(p.report.OpenPositions, r.OpenPositions...)
	}

	return nil
}

func (p *ReportProcessor) CalculateTaxes() error {
	sort.Slice(p.report.Activities, func(i, j int) bool {
		return p.report.Activities[i].Date.UnixNano() < p.report.Activities[j].Date.UnixNano()
	})

	var s, e time.Time

	s, e = getRange(p.report)
	rs := conversion.NewExchangeRateService(
		s.AddDate(-1, 0, 0).Format("2006-01-02"),
		e.Format("2006-01-02"),
	)

	var tc calculator.Calculator
	if p.rType == core.Revolut {
		tc = calculator.NewRevolutTaxCalculator(rs)
	} else if p.rType == core.EToro {
		tc = calculator.NewEТoroTaxCalculator(rs)
	} else {
		return fmt.Errorf("can't process statement - type is unknown")
	}

	err := tc.CalculateYear(p.report, p.year)
	if err != nil {
		return err
	}

	return nil
}

func getStatementType(lines []string) core.StatementType {
	for _, l := range lines {
		if strings.Contains(l, "Revolut Trading Ltd") {
			return core.Revolut
		}
		if strings.Contains(l, "eToro (Europe) Ltd") || strings.Contains(l, "eToro (UK) Ltd") {
			return core.EToro
		}
	}
	return core.Unknown
}

func getRange(report *core.Report) (start, end time.Time) {
	a := report.Activities
	if len(a) == 0 {
		return time.Now(), time.Now()
	}

	start = a[0].Date
	end = a[len(a)-1].Date

	return
}
