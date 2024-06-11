package prediction_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/shopspring/decimal"
)

type Odds struct {
	odds          decimal.Decimal
	popularNumber int
	horseNumber   types.HorseNumber
}

func NewOdds(
	odds string,
	popularNumber int,
	horseNumber types.HorseNumber,
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

func (o *Odds) HorseNumber() types.HorseNumber {
	return o.horseNumber
}
