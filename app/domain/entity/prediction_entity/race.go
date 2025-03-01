package prediction_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
)

type Race struct {
	raceId                 types.RaceId
	raceName               string
	raceDate               types.RaceDate
	raceNumber             int
	entries                int
	distance               int
	class                  types.GradeClass
	courseCategory         types.CourseCategory
	trackCondition         types.TrackCondition
	raceSexCondition       types.RaceSexCondition
	raceWeightCondition    types.RaceWeightCondition
	raceCourse             types.RaceCourse
	url                    string
	raceEntryHorses        []*RaceEntryHorse
	raceResultHorseNumbers []types.HorseNumber
	odds                   []*Odds
	predictionFilters      []filter.AttributeId
}

func NewRace(
	raceId string,
	raceName string,
	raceDate int,
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
	raceEntryHorses []*RaceEntryHorse,
	rawRaceResultHorseNumbers []int,
	odds []*Odds,
	predictionFilters []filter.AttributeId,
) *Race {
	raceResultHorseNumbers := make([]types.HorseNumber, len(rawRaceResultHorseNumbers))
	for idx := range rawRaceResultHorseNumbers {
		raceResultHorseNumbers[idx] = types.HorseNumber(rawRaceResultHorseNumbers[idx])
	}

	return &Race{
		raceId:                 types.RaceId(raceId),
		raceName:               raceName,
		raceDate:               types.RaceDate(raceDate),
		raceNumber:             raceNumber,
		entries:                entries,
		distance:               distance,
		class:                  types.GradeClass(class),
		courseCategory:         types.CourseCategory(courseCategory),
		trackCondition:         types.TrackCondition(trackCondition),
		raceSexCondition:       types.RaceSexCondition(raceSexCondition),
		raceWeightCondition:    types.RaceWeightCondition(raceWeightCondition),
		raceCourse:             types.RaceCourse(raceCourseId),
		url:                    url,
		raceEntryHorses:        raceEntryHorses,
		raceResultHorseNumbers: raceResultHorseNumbers,
		odds:                   odds,
		predictionFilters:      predictionFilters,
	}
}

func (r *Race) RaceId() types.RaceId {
	return r.raceId
}

func (r *Race) RaceName() string {
	return r.raceName
}

func (r *Race) RaceDate() types.RaceDate {
	return r.raceDate
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

func (r *Race) RaceCourse() types.RaceCourse {
	return r.raceCourse
}

func (r *Race) Url() string {
	return r.url
}

func (r *Race) RaceEntryHorses() []*RaceEntryHorse {
	return r.raceEntryHorses
}

func (r *Race) RaceResultHorseNumbers() []types.HorseNumber {
	return r.raceResultHorseNumbers
}

func (r *Race) Odds() []*Odds {
	return r.odds
}

func (r *Race) PredictionFilters() []filter.AttributeId {
	return r.predictionFilters
}
