package entity

import (
	race_vo "github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
	"strconv"
)

type RacingNumber struct {
	date         race_vo.RaceDate
	round        int
	day          int
	raceCourseId race_vo.RaceCourse
}

func NewRacingNumber(
	rawDate int,
	round int,
	day int,
	rawRaceCourseId int,
) *RacingNumber {
	date := race_vo.NewRaceDate(strconv.Itoa(rawDate))
	var raceCourseId race_vo.RaceCourse
	if rawRaceCourseId > 0 && rawRaceCourseId <= 99 {
		raceCourseId = race_vo.RaceCourse(rawRaceCourseId)
	} else {
		raceCourseId = race_vo.UnknownPlace
	}

	return &RacingNumber{
		date:         date,
		round:        round,
		day:          day,
		raceCourseId: raceCourseId,
	}
}

func (r *RacingNumber) Date() race_vo.RaceDate {
	return r.date
}

func (r *RacingNumber) Round() int {
	return r.round
}

func (r *RacingNumber) Day() int {
	return r.day
}

func (r *RacingNumber) RaceCourseId() race_vo.RaceCourse {
	return r.raceCourseId
}
