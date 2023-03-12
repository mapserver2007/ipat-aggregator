package entity

import race_vo "github.com/mapserver2007/tools/baken/app/domain/race/value_object"

type Race struct {
	raceId         string
	raceDate       int
	raceNumber     int
	raceCourseId   int
	raceName       string
	url            string
	time           string
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

func (r *Race) RaceDate() int {
	return r.raceDate
}

func (r *Race) RaceNumber() int {
	return r.raceNumber
}

func (r *Race) RaceCourseId() int {
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

func (r *Race) Entries() int {
	return r.entries
}

func (r *Race) Distance() int {
	return r.distance
}

func (r *Race) Class() int {
	return r.class
}

func (r *Race) CourseCategory() int {
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
