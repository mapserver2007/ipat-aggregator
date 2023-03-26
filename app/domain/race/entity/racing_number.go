package entity

import (
	race_vo "github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
	"strconv"
)

type RacingNumber struct {
	date         int
	round        int
	day          int
	raceCourseId int
}

func NewRacingNumber(
	date int,
	round int,
	day int,
	raceCourseId int,
) *RacingNumber {
	return &RacingNumber{
		date:         date,
		round:        round,
		day:          day,
		raceCourseId: raceCourseId,
	}
}

func (r *RacingNumber) Date() race_vo.RaceDate {
	return race_vo.NewRaceDate(strconv.Itoa(r.date))
}

func (r *RacingNumber) Round() int {
	return r.round
}

func (r *RacingNumber) Day() int {
	return r.day
}

func (r *RacingNumber) RaceCourseId() race_vo.RaceCourse {
	return race_vo.RaceCourse(r.raceCourseId)
}
