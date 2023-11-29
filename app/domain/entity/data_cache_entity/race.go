package data_cache_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type Race struct {
	raceId         types.RaceId
	raceDate       int
	raceNumber     int
	raceCourseId   types.RaceCourse
	raceName       string
	url            string
	time           string
	startTime      string
	entries        int
	distance       int
	class          int
	courseCategory int
	trackCondition string
	raceResults    []*RaceResult
	payoutResults  []*PayoutResult
}

func NewRace(
	raceId string,

) Race {

	return Race{
		raceId:         "",
		raceDate:       0,
		raceNumber:     0,
		raceCourseId:   0,
		raceName:       "",
		url:            "",
		time:           "",
		startTime:      "",
		entries:        0,
		distance:       0,
		class:          0,
		courseCategory: 0,
		trackCondition: "",
		raceResults:    nil,
		payoutResults:  nil,
	}
}
