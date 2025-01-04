package analysis_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
)

type Race struct {
	raceId              types.RaceId
	raceDate            types.RaceDate
	raceNumber          int
	raceCourse          types.RaceCourse
	raceName            string
	url                 string
	entries             int
	distance            int
	class               types.GradeClass
	courseCategory      types.CourseCategory
	trackCondition      types.TrackCondition
	raceWeightCondition types.RaceWeightCondition
	raceResults         []*RaceResult
	markers             []*Marker
	analysisFilters     []filter.AttributeId
}

func NewRace(
	raceId types.RaceId,
	raceDate types.RaceDate,
	raceNumber int,
	raceCourse types.RaceCourse,
	raceName string,
	url string,
	entries int,
	distance int,
	class types.GradeClass,
	courseCategory types.CourseCategory,
	trackCondition types.TrackCondition,
	raceWeightCondition types.RaceWeightCondition,
	raceResults []*RaceResult,
	markers []*Marker,
	analysisFilters []filter.AttributeId,
) *Race {
	return &Race{
		raceId:              raceId,
		raceDate:            raceDate,
		raceNumber:          raceNumber,
		raceCourse:          raceCourse,
		raceName:            raceName,
		url:                 url,
		entries:             entries,
		distance:            distance,
		class:               class,
		courseCategory:      courseCategory,
		trackCondition:      trackCondition,
		raceWeightCondition: raceWeightCondition,
		raceResults:         raceResults,
		markers:             markers,
		analysisFilters:     analysisFilters,
	}
}

func (r *Race) RaceId() types.RaceId {
	return r.raceId
}

func (r *Race) RaceDate() types.RaceDate {
	return r.raceDate
}

func (r *Race) RaceNumber() int {
	return r.raceNumber
}

func (r *Race) RaceCourse() types.RaceCourse {
	return r.raceCourse
}

func (r *Race) RaceName() string {
	return r.raceName
}

func (r *Race) Url() string {
	return r.url
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

func (r *Race) RaceWeightCondition() types.RaceWeightCondition {
	return r.raceWeightCondition
}

func (r *Race) RaceResults() []*RaceResult {
	return r.raceResults
}

func (r *Race) Markers() []*Marker {
	return r.markers
}

func (r *Race) AnalysisFilters() []filter.AttributeId {
	return r.analysisFilters
}
