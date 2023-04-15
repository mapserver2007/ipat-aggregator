package entity

import (
	jockey_vo "github.com/mapserver2007/ipat-aggregator/app/domain/jockey/value_object"
)

type RaceResult struct {
	orderNo       int
	horseName     string
	bracketNumber int
	horseNumber   int
	jockeyId      int
	jockeyName    string
	odds          string
	popularNumber int
}

func NewRaceResult(
	orderNo int,
	horseName string,
	bracketNumber int,
	horseNumber int,
	jockeyId int,
	jockeyName string,
	odds string,
	popularNumber int,
) *RaceResult {
	return &RaceResult{
		orderNo:       orderNo,
		horseName:     horseName,
		bracketNumber: bracketNumber,
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

func (r *RaceResult) HorseName() string {
	return r.horseName
}

func (r *RaceResult) BracketNumber() int {
	return r.bracketNumber
}

func (r *RaceResult) HorseNumber() int {
	return r.horseNumber
}

func (r *RaceResult) JockeyId() jockey_vo.JockeyId {
	return jockey_vo.JockeyId(r.jockeyId)
}

func (r *RaceResult) JockeyName() string {
	return r.jockeyName
}

func (r *RaceResult) Odds() string {
	return r.odds
}

func (r *RaceResult) PopularNumber() int {
	return r.popularNumber
}
