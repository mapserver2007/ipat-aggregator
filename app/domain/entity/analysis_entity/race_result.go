package analysis_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/shopspring/decimal"
)

type RaceResult struct {
	orderNo       int
	horseId       types.HorseId
	horseName     string
	horseNumber   types.HorseNumber
	jockeyId      types.JockeyId
	jockeyName    string
	odds          decimal.Decimal
	popularNumber int
}

func NewRaceResult(
	orderNo int,
	horseId types.HorseId,
	horseName string,
	horseNumber types.HorseNumber,
	jockeyId types.JockeyId,
	jockeyName string,
	odds decimal.Decimal,
	popularNumber int,
) *RaceResult {
	return &RaceResult{
		orderNo:       orderNo,
		horseId:       horseId,
		horseName:     horseName,
		horseNumber:   horseNumber,
		jockeyId:      jockeyId,
		jockeyName:    jockeyName,
		odds:          odds,
		popularNumber: popularNumber,
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

func (r *RaceResult) HorseNumber() types.HorseNumber {
	return r.horseNumber
}

func (r *RaceResult) JockeyId() types.JockeyId {
	return r.jockeyId
}

func (r *RaceResult) JockeyName() string {
	return r.jockeyName
}

func (r *RaceResult) Odds() decimal.Decimal {
	return r.odds
}

func (r *RaceResult) PopularNumber() int {
	return r.popularNumber
}
