package data_cache_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/shopspring/decimal"
)

type JockeyResult struct {
	jockeyResultId types.JockeyResultId
	jockeyId       types.JockeyId
	raceId         types.RaceId
	raceDate       types.RaceDate
	raceCourseId   types.RaceCourse
	odds           decimal.Decimal
	orderNo        int
	horseId        types.HorseId
	courseCategory types.CourseCategory
	distance       int
}

func NewJockeyResult(
	jockeyResultId string,
	jockeyId string,
	raceId string,
	raceDate int,
	raceCourseId string,
	odds string,
	orderNo int,
	horseId string,
	courseCategory int,
	distance int,
) *JockeyResult {
	decimalOdds, _ := decimal.NewFromString(odds)
	return &JockeyResult{
		jockeyResultId: types.JockeyResultId(jockeyResultId),
		jockeyId:       types.JockeyId(jockeyId),
		raceId:         types.RaceId(raceId),
		raceDate:       types.RaceDate(raceDate),
		raceCourseId:   types.RaceCourse(raceCourseId),
		odds:           decimalOdds,
		orderNo:        orderNo,
		horseId:        types.HorseId(horseId),
		courseCategory: types.CourseCategory(courseCategory),
		distance:       distance,
	}
}

func (j *JockeyResult) JockeyResultId() types.JockeyResultId {
	return j.jockeyResultId
}

func (j *JockeyResult) JockeyId() types.JockeyId {
	return j.jockeyId
}

func (j *JockeyResult) RaceId() types.RaceId {
	return j.raceId
}

func (j *JockeyResult) RaceDate() types.RaceDate {
	return j.raceDate
}

func (j *JockeyResult) RaceCourseId() types.RaceCourse {
	return j.raceCourseId
}

func (j *JockeyResult) Odds() decimal.Decimal {
	return j.odds
}

func (j *JockeyResult) OrderNo() int {
	return j.orderNo
}

func (j *JockeyResult) HorseId() types.HorseId {
	return j.horseId
}

func (j *JockeyResult) CourseCategory() types.CourseCategory {
	return j.courseCategory
}

func (j *JockeyResult) Distance() int {
	return j.distance
}
