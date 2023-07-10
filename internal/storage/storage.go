package storage

import (
	"fmt"
)

type LatestRates struct {
	Base    string
	Date    string
	EUR     float64
	JPY     float64
	UAH     float64
}

type RawRates struct {
	Motd    AddMessage
	Success bool
	Base    string
	Date    string
	Rates   map[string]float64
}

type AddMessage struct {
	Msg string
	Url string
}

func (l *RawRates) Clean() *LatestRates {
	return &LatestRates{
		l.Base,
		l.Date,
		l.Rates["EUR"],
		l.Rates["JPY"],
		l.Rates["UAH"],
	}
}

func (l *LatestRates) Stringify() string {
	return fmt.Sprintf("Base=%s Date=%s EUR : %f JPY : %f UAH : %f", l.Base, l.Date, l.EUR, l.JPY, l.UAH) 
}