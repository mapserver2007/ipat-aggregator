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
}

func NewRace(
	raceId types.RaceId,
	raceName string,
	entries int,
	distance int,
	class types.GradeClass,
	courseCategory types.CourseCategory,
	trackCondition types.TrackCondition,
	raceSexCondition types.RaceSexCondition,
	raceWeightCondition types.RaceWeightCondition,
) *Race {
	return &Race{
		raceId:              raceId,
		raceName:            raceName,
		entries:             entries,
		distance:            distance,
		class:               class,
		courseCategory:      courseCategory,
		trackCondition:      trackCondition,
		raceSexCondition:    raceSexCondition,
		raceWeightCondition: raceWeightCondition,
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
