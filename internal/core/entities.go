package core

import (
	"time"
)

type ActivityType string

const (
	Buy         ActivityType = "BUY"
	Sell        ActivityType = "SELL"
	Csd         ActivityType = "CSD"
	Div         ActivityType = "DIV"
	DivNra      ActivityType = "DIVNRA"
	Ssp         ActivityType = "SSP"
	FivFt       ActivityType = "DIVFT"
	Mas         ActivityType = "MAS"
	Cdep        ActivityType = "CDEP"
	Linked      ActivityType = "LINKED"
	RolloverFee ActivityType = "Rollover Fee"
)

type Currency string

const (
	Usd Currency = "USD"
	Bgn Currency = "BGN"
)

const (
	Unknown StatementType = iota + 1
	Revolut
	EToro
)

type StatementType int

type LinkedActivity struct {
	Activity
	OpenDate   time.Time
	ClosedDate time.Time
}

type OpenPosition struct {
	Date        time.Time
	PositionID  string
	Units       float64
	PriceHome   float64
	PriceOrigin float64
	Token       string
	Name        string
	Type        string
}

type DividendMeta struct {
	GrossAmount float64
	NetAmount   float64
	Tax         float64
}

type Amounts struct {
	TotalBuy  float64
	TotalSell float64
}

type Report struct {
	RequestID     string
	Activities    []LinkedActivity
	OpenPositions []*OpenPosition
	Dividends     DividendMeta
	Amounts       Amounts
	Tax           float64
	StartDate     *time.Time
	EndDate       *time.Time
	Deposits      float64
}

type Activity struct {
	Token      string
	Name       string
	Currency   Currency
	Type       ActivityType
	Date       time.Time
	Amount     float64
	Units      float64
	OpenRate   float64
	ClosedRate float64
}

type Dividend struct {
	Gross      float64
	Tax        float64
	TaxPercent float64
}
