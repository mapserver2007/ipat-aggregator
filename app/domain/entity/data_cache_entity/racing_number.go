package data_cache_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type RacingNumber struct {
	raceDate     types.RaceDate
	raceCourseId types.RaceCourse
	round        int
	day          int
}

func NewRacingNumber(
	date int,
	round int,
	day int,
	raceCourseId string,
) *RacingNumber {
	return &RacingNumber{
		raceDate:     types.RaceDate(date),
		raceCourseId: types.RaceCourse(raceCourseId),
		round:        round,
		day:          day,
	}
}

func (r *RacingNumber) RaceDate() types.RaceDate {
	return r.raceDate
}

func (r *RacingNumber) Round() int {
	return r.round
}

func (r *RacingNumber) Day() int {
	return r.day
}

func (r *RacingNumber) RaceCourse() types.RaceCourse {
	return r.raceCourseId
}
