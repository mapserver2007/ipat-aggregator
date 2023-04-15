package entity

import (
	race_entity "github.com/mapserver2007/ipat-aggregator/app/domain/race/entity"
	race_vo "github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
)

type Race struct {
	raceId         race_vo.RaceId
	raceNumber     int
	raceName       string
	startTime      string
	class          race_vo.GradeClass
	raceCourse     race_vo.RaceCourse
	courseCategory race_vo.CourseCategory
	raceDate       race_vo.RaceDate
	distance       int
	trackCondition string
	payment        int
	repayment      int
	url            string
	raceResults    []*race_entity.RaceResult
}

func NewRace(
	raceId race_vo.RaceId,
	raceNumber int,
	raceName string,
	startTime string,
	class race_vo.GradeClass,
	raceCourse race_vo.RaceCourse,
	courseCategory race_vo.CourseCategory,
	raceDate race_vo.RaceDate,
	distance int,
	trackCondition string,
	payment int,
	repayment int,
	url string,
	raceResults []*race_entity.RaceResult,
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
		payment:        payment,
		repayment:      repayment,
		url:            url,
		raceResults:    raceResults,
	}
}

func (r *Race) RaceId() race_vo.RaceId {
	return r.raceId
}

func (r *Race) RaceNumber() int {
	return r.raceNumber
}

func (r *Race) RaceName() string {
	return r.raceName
}

func (r *Race) StartTime() string {
	return r.startTime
}

func (r *Race) Class() race_vo.GradeClass {
	return r.class
}

func (r *Race) RaceCourse() race_vo.RaceCourse {
	return r.raceCourse
}

func (r *Race) RaceDate() race_vo.RaceDate {
	return r.raceDate
}

func (r *Race) CourseCategory() race_vo.CourseCategory {
	return r.courseCategory
}

func (r *Race) Distance() int {
	return r.distance
}

func (r *Race) TrackCondition() string {
	return r.trackCondition
}

func (r *Race) Payment() int {
	return r.payment
}

func (r *Race) Repayment() int {
	return r.repayment
}

func (r *Race) Url() string {
	return r.url
}

func (r *Race) RaceResults() []*race_entity.RaceResult {
	return r.raceResults
}
