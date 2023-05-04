package entity

import race_vo "github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"

type Race struct {
	raceId         string
	raceDate       int
	raceNumber     int
	raceCourseId   int
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
	raceDate int,
	raceNumber int,
	raceCourseId int,
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
		raceId:         raceId,
		raceDate:       raceDate,
		raceNumber:     raceNumber,
		raceCourseId:   raceCourseId,
		raceName:       raceName,
		url:            url,
		time:           time,
		startTime:      startTime,
		entries:        entries,
		distance:       distance,
		class:          class,
		courseCategory: courseCategory,
		trackCondition: trackCondition,
		raceResults:    raceResults,
		payoutResults:  payoutResults,
	}
}

func (r *Race) RaceId() race_vo.RaceId {
	return race_vo.RaceId(r.raceId)
}

func (r *Race) RaceDate() race_vo.RaceDate {
	return race_vo.RaceDate(r.raceDate)
}

func (r *Race) RaceNumber() int {
	return r.raceNumber
}

func (r *Race) RaceCourseId() race_vo.RaceCourse {
	return race_vo.RaceCourse(r.raceCourseId)
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

func (r *Race) Class() race_vo.GradeClass {
	return race_vo.GradeClass(r.class)
}

func (r *Race) CourseCategory() race_vo.CourseCategory {
	return race_vo.CourseCategory(r.courseCategory)
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
