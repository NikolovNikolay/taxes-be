package conversion

import (
	"encoding/csv"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"strconv"
	"taxes-be/internal/core"
	"time"
)

const (
	dateLayout = "2006-01-02"
)

type ExchangeRateService struct {
	rates map[string]map[string]float64
}

func NewExchangeRateService() *ExchangeRateService {
	s := &ExchangeRateService{
	}
	s.getExchangeRates()
	return s
}

func (s *ExchangeRateService) GetRateForDate(date time.Time, currency core.Currency) float64 {
	var ok bool
	var rate float64
	for !ok {
		d := date.Format(dateLayout)
		rate, ok = s.rates[d][string(currency)]
		if !ok {
			date = date.AddDate(0, 0, -1)
			continue
		}
		return rate
	}
	return 0
}

func (s *ExchangeRateService) getExchangeRates() {
	csvfile, err := os.Open("resources/usd_bgn_rates.csv")
	if err != nil {
		logrus.Fatalln("couldn't open the rates csv file", err)
	}

	// Parse the file
	rdr := csv.NewReader(csvfile)
	rates := map[string]map[string]float64{}
	for {
		// Read each record from csv
		record, err := rdr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logrus.Fatal(err)
		}

		r, _ := strconv.ParseFloat(record[1], 64)
		if rates[record[0]] == nil {
			rates[record[0]] = map[string]float64{}
		}
		rates[record[0]][string(core.Bgn)] = r
		s.rates = rates
	}
}
