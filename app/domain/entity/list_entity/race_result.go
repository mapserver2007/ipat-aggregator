package list_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/shopspring/decimal"
)

type RaceResult struct {
	orderNo       int
	horseName     string
	bracketNumber int
	horseNumber   types.HorseNumber
	jockeyId      types.JockeyId
	odds          decimal.Decimal
	popularNumber int
}

func NewRaceResult(
	orderNo int,
	horseName string,
	bracketNumber int,
	horseNumber types.HorseNumber,
	jockeyId types.JockeyId,
	odds decimal.Decimal,
	popularNumber int,
) *RaceResult {
	return &RaceResult{
		orderNo:       orderNo,
		horseName:     horseName,
		bracketNumber: bracketNumber,
		horseNumber:   horseNumber,
		jockeyId:      jockeyId,
		odds:          odds,
		popularNumber: popularNumber,
	}
}

func (r *RaceResult) OrderNo() int {
	return r.orderNo
}

func (r *RaceResult) HorseName() string {
	return r.horseName
}

func (r *RaceResult) BracketNumber() int {
	return r.bracketNumber
}

func (r *RaceResult) HorseNumber() types.HorseNumber {
	return r.horseNumber
}

func (r *RaceResult) JockeyId() types.JockeyId {
	return r.jockeyId
}

func (r *RaceResult) Odds() decimal.Decimal {
	return r.odds
}

func (r *RaceResult) PopularNumber() int {
	return r.popularNumber
}
