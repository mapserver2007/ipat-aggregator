package data_cache_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type RaceResult struct {
	orderNo       int
	horseName     string
	bracketNumber int
	horseNumber   types.HorseNumber
	jockeyId      types.JockeyId
	odds          string
	popularNumber int
}

func NewRaceResult(
	orderNo int,
	horseName string,
	bracketNumber int,
	horseNumber int,
	jockeyId string,
	odds string,
	popularNumber int,
) *RaceResult {
	return &RaceResult{
		orderNo:       orderNo,
		horseName:     horseName,
		bracketNumber: bracketNumber,
		horseNumber:   types.HorseNumber(horseNumber),
		jockeyId:      types.JockeyId(jockeyId),
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

func (r *RaceResult) Odds() string {
	return r.odds
}

func (r *RaceResult) PopularNumber() int {
	return r.popularNumber
}
