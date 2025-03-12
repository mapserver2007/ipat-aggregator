package data_cache_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/shopspring/decimal"
)

type RaceResult struct {
	orderNo        int
	horseId        types.HorseId
	horseName      string
	bracketNumber  int
	horseNumber    types.HorseNumber
	jockeyId       types.JockeyId
	odds           decimal.Decimal
	popularNumber  int
	jockeyWeight   string
	horseWeight    int
	horseWeightAdd int
}

func NewRaceResult(
	orderNo int,
	horseId string,
	horseName string,
	bracketNumber int,
	horseNumber int,
	jockeyId string,
	odds string,
	popularNumber int,
	jockeyWeight string,
	horseWeight int,
	horseWeightAdd int,
) *RaceResult {
	decimalOdds, _ := decimal.NewFromString(odds)
	return &RaceResult{
		orderNo:        orderNo,
		horseId:        types.HorseId(horseId),
		horseName:      horseName,
		bracketNumber:  bracketNumber,
		horseNumber:    types.HorseNumber(horseNumber),
		jockeyId:       types.JockeyId(jockeyId),
		odds:           decimalOdds,
		popularNumber:  popularNumber,
		jockeyWeight:   jockeyWeight,
		horseWeight:    horseWeight,
		horseWeightAdd: horseWeightAdd,
	}
}

func (r *RaceResult) OrderNo() int {
	return r.orderNo
}

func (r *RaceResult) HorseId() types.HorseId {
	return r.horseId
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

func (r *RaceResult) JockeyWeight() string {
	return r.jockeyWeight
}

func (r *RaceResult) HorseWeight() int {
	return r.horseWeight
}

func (r *RaceResult) HorseWeightAdd() int {
	return r.horseWeightAdd
}
