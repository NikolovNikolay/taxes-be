package sendgrid

import (
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"sort"
	"taxes-be/internal/core"
)

const (
	fromName    = "Nikolay from Digitools-it"
	fromAddress = "nikolov@digitools-it.com"
	subject     = "Annual tax calculations for %d"
)

type Mailer struct {
	apiKey string
}

func NewMailer(apiKey string) *Mailer {
	return &Mailer{
		apiKey: apiKey,
	}
}

func (m *Mailer) SendReportMail(year int, toFullName, toEmail string, report *core.Report) error {
	from := mail.NewEmail(fromName, fromAddress)
	sub := fmt.Sprintf(subject, year)
	to := mail.NewEmail(toFullName, toEmail)
	p := message.NewPrinter(language.Bulgarian)
	plainTextContent := fmt.Sprintf(
		`Hi %s,

Here are the results of the tax calculations you requested with inquiry ID %s (%d):

----------------------------------------
Operations (BGN):
----------------------------------------
Total sell amount: %s
Total buy amount:  %s
Total tax:         %s

----------------------------------------
Dividends (BGN):
----------------------------------------
Net amount:   %s
Gross amount: %s
Total tax:    %s

----------------------------------------
Transferred positions %d (%d):
----------------------------------------
%s

Regards,
Nikolay`,
		toFullName,
		report.RequestID,
		year,
		p.Sprintf("%.2f", report.Amounts.TotalSell),
		p.Sprintf("%.2f", report.Amounts.TotalBuy),
		p.Sprintf("%.2f", report.Tax),
		p.Sprintf("%.2f", report.Dividends.NetAmount),
		p.Sprintf("%.2f", report.Dividends.GrossAmount),
		p.Sprintf("%.2f", report.Dividends.Tax),
		year,
		len(report.OpenPositions),
		buildOpenPositions(report.OpenPositions, p),
	)
	msg := mail.NewSingleEmail(from, sub, to, plainTextContent, "")
	client := sendgrid.NewSendClient(m.apiKey)
	_, err := client.Send(msg)

	if err != nil {
		return err
	}

	return nil
}

func buildOpenPositions(op []*core.OpenPosition, p *message.Printer) string {
	sort.Slice(op, func(i, j int) bool {
		return op[i].Date.UnixNano() < op[j].Date.UnixNano()
	})
	positions := ""
	for i := range op {
		positions += fmt.Sprintf(
			`%d | Date: %v | Token: %s | Name: %s | Avg price (USD): %s | Avg price (BGN): %s | Units: %s |
---------------------
`,
			i+1,
			op[i].Date.Format("02-01-2006"),
			op[i].Token,
			op[i].Name,
			p.Sprintf("%.2f", op[i].AmountOrigin),
			p.Sprintf("%.2f", op[i].AmountHome),
			p.Sprintf("%.5f", op[i].Units),
		)
	}
	return positions
}
