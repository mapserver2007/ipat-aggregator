package netkeiba_entity

import "fmt"

type JockeyResult struct {
	jockeyResultId   string
	jockeyId         string
	raceId           string
	raceDate         int
	raceCourseId     string
	odds             string
	orderNo          int
	horseId          string
	courseCategoryId int
	distance         int
}

func NewJockeyResult(
	jockeyId string,
	raceId string,
	raceDate int,
	raceCourseId string,
	odds string,
	orderNo int,
	horseId string,
	courseCategoryId int,
	distance int,
) *JockeyResult {
	return &JockeyResult{
		jockeyResultId:   fmt.Sprintf("%s%s", jockeyId, raceId),
		jockeyId:         jockeyId,
		raceId:           raceId,
		raceDate:         raceDate,
		raceCourseId:     raceCourseId,
		odds:             odds,
		orderNo:          orderNo,
		horseId:          horseId,
		courseCategoryId: courseCategoryId,
		distance:         distance,
	}
}

func (j *JockeyResult) JockeyResultId() string {
	return j.jockeyResultId
}

func (j *JockeyResult) JockeyId() string {
	return j.jockeyId
}

func (j *JockeyResult) RaceId() string {
	return j.raceId
}

func (j *JockeyResult) RaceDate() int {
	return j.raceDate
}

func (j *JockeyResult) RaceCourseId() string {
	return j.raceCourseId
}

func (j *JockeyResult) Odds() string {
	return j.odds
}

func (j *JockeyResult) OrderNo() int {
	return j.orderNo
}

func (j *JockeyResult) HorseId() string {
	return j.horseId
}

func (j *JockeyResult) CourseCategoryId() int {
	return j.courseCategoryId
}

func (j *JockeyResult) Distance() int {
	return j.distance
}
