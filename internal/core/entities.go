package core

import (
	"time"
)

type ActivityType string

const (
	BUY          ActivityType = "BUY"
	SELL         ActivityType = "SELL"
	CSD          ActivityType = "CSD"
	DIV          ActivityType = "DIV"
	DIVNRA       ActivityType = "DIVNRA"
	SSP          ActivityType = "SSP"
	DIVFT        ActivityType = "DIVFT"
	MAS          ActivityType = "MAS"
	CDEP         ActivityType = "CDEP"
	LINKED       ActivityType = "LINKED"
	ROLLOVER_FEE ActivityType = "Rollover Fee"
)

type Currency string

const (
	USD Currency = "USD"
	BGN Currency = "BGN"
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
	Date         time.Time
	PositionID   string
	Units        float64
	AmountHome   float64
	AmountOrigin float64
	Token        string
	Name         string
	Type         string
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

type Month string

const (
	Jan = "January"
	Feb = "February"
	Mar = "March"
	Apr = "April"
	May = "May"
	Jun = "Jun"
	Jul = "Jul"
	Aug = "August"
	Sep = "September"
	Oct = "October"
	Nov = "November"
	Dec = "December"
)
