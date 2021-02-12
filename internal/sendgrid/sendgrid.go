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
	fromName    = "Nikolay from Clerky"
	fromAddress = "nikolov89@gmail.com"
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
		`Hi there,

Here are the results for the tax calculations you requested 
for ID %s (%d):

----------------------------------------
Operations:
----------------------------------------
Total sell amount: %s
Total buy amount:  %s
Total tax:         %s

----------------------------------------
Dividents:
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
		report.RequestID,
		year,
		p.Sprint(report.Amounts.TotalSell),
		p.Sprint(report.Amounts.TotalBuy),
		p.Sprint(report.Tax),
		p.Sprint(report.Dividends.NetAmount),
		p.Sprint(report.Dividends.GrossAmount),
		p.Sprint(report.Dividends.Tax),
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
			`%d | Date: %v, Token: %s, Name: %s, Avg Foreign price: %s, Avg Home price: %s, Units: %s
---------------------
`,
			i+1,
			op[i].Date.Format("02-01-2006"),
			op[i].Token,
			op[i].Name,
			p.Sprint(op[i].AmountOrigin),
			p.Sprint(op[i].AmountHome),
			p.Sprint(op[i].Units),
		)
	}
	return positions
}
