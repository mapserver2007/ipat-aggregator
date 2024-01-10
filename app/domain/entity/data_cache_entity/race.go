package data_cache_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type Race struct {
	raceId         types.RaceId
	raceDate       types.RaceDate
	raceNumber     int
	raceCourseId   types.RaceCourse
	raceName       string
	url            string
	time           string
	startTime      string
	entries        int
	distance       int
	class          types.GradeClass
	courseCategory types.CourseCategory
	trackCondition string
	raceResults    []*RaceResult
	payoutResults  []*PayoutResult
}

func NewRace(
	raceId string,
	raceDate int,
	raceNumber int,
	raceCourseId string,
	raceName string,
	url string,
	time string,
	startTime string,
	entries int,
	distance int,
	class int,
	courseCategory int,
	trackCondition string,
	raceResults []*RaceResult,
	payoutResults []*PayoutResult,
) *Race {
	return &Race{
		raceId:         types.RaceId(raceId),
		raceDate:       types.RaceDate(raceDate),
		raceNumber:     raceNumber,
		raceCourseId:   types.RaceCourse(raceCourseId),
		raceName:       raceName,
		url:            url,
		time:           time,
		startTime:      startTime,
		entries:        entries,
		distance:       distance,
		class:          types.GradeClass(class),
		courseCategory: types.CourseCategory(courseCategory),
		trackCondition: trackCondition,
		raceResults:    raceResults,
		payoutResults:  payoutResults,
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

func (r *Race) RaceCourseId() types.RaceCourse {
	return r.raceCourseId
}

func (r *Race) RaceName() string {
	return r.raceName
}

func (r *Race) Url() string {
	return r.url
}

func (r *Race) Time() string {
	return r.time
}

func (r *Race) StartTime() string {
	return r.startTime
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

func (r *Race) TrackCondition() string {
	return r.trackCondition
}

func (r *Race) RaceResults() []*RaceResult {
	return r.raceResults
}

func (r *Race) PayoutResults() []*PayoutResult {
	return r.payoutResults
}
