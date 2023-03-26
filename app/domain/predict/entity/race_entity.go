package entity

import (
	race_entity "github.com/mapserver2007/ipat-aggregator/app/domain/race/entity"
	race_vo "github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
)

type RaceEntity struct {
	raceId         race_vo.RaceId
	raceNumber     int
	raceName       string
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

func NewRaceEntity(
	raceId race_vo.RaceId,
	raceNumber int,
	raceName string,
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
) *RaceEntity {
	return &RaceEntity{
		raceId:         raceId,
		raceNumber:     raceNumber,
		raceName:       raceName,
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

func (r *RaceEntity) RaceId() race_vo.RaceId {
	return r.raceId
}

func (r *RaceEntity) RaceNumber() int {
	return r.raceNumber
}

func (r *RaceEntity) RaceName() string {
	return r.raceName
}

func (r *RaceEntity) Class() race_vo.GradeClass {
	return r.class
}

func (r *RaceEntity) RaceCourse() race_vo.RaceCourse {
	return r.raceCourse
}

func (r *RaceEntity) RaceDate() race_vo.RaceDate {
	return r.raceDate
}

func (r *RaceEntity) CourseCategory() race_vo.CourseCategory {
	return r.courseCategory
}

func (r *RaceEntity) Distance() int {
	return r.distance
}

func (r *RaceEntity) TrackCondition() string {
	return r.trackCondition
}

func (r *RaceEntity) Payment() int {
	return r.payment
}

func (r *RaceEntity) Repayment() int {
	return r.repayment
}

func (r *RaceEntity) Url() string {
	return r.url
}

func (r *RaceEntity) RaceResults() []*race_entity.RaceResult {
	return r.raceResults
}
