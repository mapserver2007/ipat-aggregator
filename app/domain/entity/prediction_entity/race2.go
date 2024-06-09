package prediction_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
)

type Race2 struct {
	raceId                 types.RaceId
	raceName               string
	raceNumber             int
	entries                int
	distance               int
	class                  types.GradeClass
	courseCategory         types.CourseCategory
	trackCondition         types.TrackCondition
	raceSexCondition       types.RaceSexCondition
	raceWeightCondition    types.RaceWeightCondition
	raceCourseId           types.RaceCourse
	url                    string
	raceResultHorseNumbers []int
	odds                   []*Odds
	predictionFilter       filter.Id
}

func NewRace2(
	raceId string,
	raceName string,
	raceNumber int,
	entries int,
	distance int,
	class int,
	courseCategory int,
	trackCondition int,
	raceSexCondition int,
	raceWeightCondition int,
	raceCourseId string,
	url string,
	raceResultHorseNumbers []int,
	odds []*Odds,
	predictionFilter filter.Id,
) *Race2 {
	return &Race2{
		raceId:                 types.RaceId(raceId),
		raceName:               raceName,
		raceNumber:             raceNumber,
		entries:                entries,
		distance:               distance,
		class:                  types.GradeClass(class),
		courseCategory:         types.CourseCategory(courseCategory),
		trackCondition:         types.TrackCondition(trackCondition),
		raceSexCondition:       types.RaceSexCondition(raceSexCondition),
		raceWeightCondition:    types.RaceWeightCondition(raceWeightCondition),
		raceCourseId:           types.RaceCourse(raceCourseId),
		url:                    url,
		raceResultHorseNumbers: raceResultHorseNumbers,
		odds:                   odds,
		predictionFilter:       predictionFilter,
	}
}

func (r *Race2) RaceId() types.RaceId {
	return r.raceId
}

func (r *Race2) RaceName() string {
	return r.raceName
}

func (r *Race2) RaceNumber() int {
	return r.raceNumber
}

func (r *Race2) Entries() int {
	return r.entries
}

func (r *Race2) Distance() int {
	return r.distance
}

func (r *Race2) Class() types.GradeClass {
	return r.class
}

func (r *Race2) CourseCategory() types.CourseCategory {
	return r.courseCategory
}

func (r *Race2) TrackCondition() types.TrackCondition {
	return r.trackCondition
}

func (r *Race2) RaceSexCondition() types.RaceSexCondition {
	return r.raceSexCondition
}

func (r *Race2) RaceWeightCondition() types.RaceWeightCondition {
	return r.raceWeightCondition
}

func (r *Race2) RaceCourseId() types.RaceCourse {
	return r.raceCourseId
}

func (r *Race2) Url() string {
	return r.url
}

func (r *Race2) RaceResultHorseNumbers() []int {
	return r.raceResultHorseNumbers
}

func (r *Race2) Odds() []*Odds {
	return r.odds
}

func (r *Race2) PredictionFilter() filter.Id {
	return r.predictionFilter
}
