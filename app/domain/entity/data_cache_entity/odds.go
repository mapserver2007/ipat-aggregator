package data_cache_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type Odds struct {
	number        types.BetNumber
	popularNumber int
	odds          string
}

func NewOdds(
	rawNumber string,
	popularNumber int,
	odds string,
) *Odds {
	return &Odds{
		number:        types.NewBetNumber(rawNumber),
		popularNumber: popularNumber,
		odds:          odds,
	}
}

func (o *Odds) Number() types.BetNumber {
	return o.number
}

func (o *Odds) PopularNumber() int {
	return o.popularNumber

}

func (o *Odds) Odds() string {
	return o.odds
}
