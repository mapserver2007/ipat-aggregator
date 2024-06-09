package prediction_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
)

type Race struct {
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

func NewRace(
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
) *Race {
	return &Race{
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

func (r *Race) RaceId() types.RaceId {
	return r.raceId
}

func (r *Race) RaceName() string {
	return r.raceName
}

func (r *Race) RaceNumber() int {
	return r.raceNumber
}

func (r *Race) Entries() int {
	return r.entries
}

func (r *Race) Distance() int {
	return r.distance
}

func (r *Race) Class() types.GradeClass {
	return r.class
}

func (r *Race) CourseCategory() types.CourseCategory {
	return r.courseCategory
}

func (r *Race) TrackCondition() types.TrackCondition {
	return r.trackCondition
}

func (r *Race) RaceSexCondition() types.RaceSexCondition {
	return r.raceSexCondition
}

func (r *Race) RaceWeightCondition() types.RaceWeightCondition {
	return r.raceWeightCondition
}

func (r *Race) RaceCourseId() types.RaceCourse {
	return r.raceCourseId
}

func (r *Race) Url() string {
	return r.url
}

func (r *Race) RaceResultHorseNumbers() []int {
	return r.raceResultHorseNumbers
}

func (r *Race) Odds() []*Odds {
	return r.odds
}

func (r *Race) PredictionFilter() filter.Id {
	return r.predictionFilter
}
