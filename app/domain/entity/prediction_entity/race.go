package prediction_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type Race struct {
	raceId              types.RaceId
	raceName            string
	entries             int
	distance            int
	class               types.GradeClass
	courseCategory      types.CourseCategory
	trackCondition      types.TrackCondition
	raceSexCondition    types.RaceSexCondition
	raceWeightCondition types.RaceWeightCondition
	raceCourseId        types.RaceCourse
	odds                []*Odds
}

func NewRace(
	raceId string,
	raceName string,
	entries int,
	distance int,
	class int,
	courseCategory int,
	trackCondition int,
	raceSexCondition int,
	raceWeightCondition int,
	raceCourseId string,
	odds []*Odds,
) *Race {

	return &Race{
		raceId:              types.RaceId(raceId),
		raceName:            raceName,
		entries:             entries,
		distance:            distance,
		class:               types.GradeClass(class),
		courseCategory:      types.CourseCategory(courseCategory),
		trackCondition:      types.TrackCondition(trackCondition),
		raceSexCondition:    types.RaceSexCondition(raceSexCondition),
		raceWeightCondition: types.RaceWeightCondition(raceWeightCondition),
		raceCourseId:        types.NewRaceCourse(raceCourseId),
		odds:                odds,
	}
}

func (r *Race) RaceId() types.RaceId {
	return r.raceId
}

func (r *Race) RaceName() string {
	return r.raceName
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

func (r *Race) Odds() []*Odds {
	return r.odds
}
