package list_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type Race struct {
	raceId         types.RaceId
	raceNumber     int
	raceName       string
	startTime      string
	class          types.GradeClass
	raceCourse     types.RaceCourse
	courseCategory types.CourseCategory
	raceDate       types.RaceDate
	distance       int
	trackCondition types.TrackCondition
	url            string
	raceResults    []*RaceResult
}

func NewRace(
	raceId types.RaceId,
	raceNumber int,
	raceName string,
	startTime string,
	class types.GradeClass,
	raceCourse types.RaceCourse,
	courseCategory types.CourseCategory,
	raceDate types.RaceDate,
	distance int,
	trackCondition types.TrackCondition,
	url string,
	raceResults []*RaceResult,
) *Race {
	return &Race{
		raceId:         raceId,
		raceNumber:     raceNumber,
		raceName:       raceName,
		startTime:      startTime,
		class:          class,
		raceCourse:     raceCourse,
		courseCategory: courseCategory,
		raceDate:       raceDate,
		distance:       distance,
		trackCondition: trackCondition,
		url:            url,
		raceResults:    raceResults,
	}
}

func (r *Race) RaceId() types.RaceId {
	return r.raceId
}

func (r *Race) RaceNumber() int {
	return r.raceNumber
}

func (r *Race) RaceName() string {
	return r.raceName
}

func (r *Race) Url() string {
	return r.url
}

func (r *Race) StartTime() string {
	return r.startTime
}

func (r *Race) RaceDate() types.RaceDate {
	return r.raceDate
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

func (r *Race) RaceResults() []*RaceResult {
	return r.raceResults
}
