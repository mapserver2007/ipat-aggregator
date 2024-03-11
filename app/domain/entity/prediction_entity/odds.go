package prediction_entity

import "github.com/shopspring/decimal"

type Odds struct {
	odds          decimal.Decimal
	popularNumber int
	horseNumber   int
}

func NewOdds(
	odds string,
	popularNumber int,
	horseNumber int,
) *Odds {
	decimalOdds, _ := decimal.NewFromString(odds)
	return &Odds{
		odds:          decimalOdds,
		popularNumber: popularNumber,
		horseNumber:   horseNumber,
	}
}

func (o *Odds) Odds() decimal.Decimal {
	return o.odds
}

func (o *Odds) PopularNumber() int {
	return o.popularNumber
}

func (o *Odds) HorseNumber() int {
	return o.horseNumber
}
