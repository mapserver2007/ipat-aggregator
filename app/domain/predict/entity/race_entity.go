package entity

import (
	race_entity "github.com/mapserver2007/tools/baken/app/domain/race/entity"
	race_vo "github.com/mapserver2007/tools/baken/app/domain/race/value_object"
)

type RaceEntity struct {
	RaceId         race_vo.RaceId
	RaceNumber     int
	RaceName       string
	Class          race_vo.GradeClass
	RaceCourse     race_vo.RaceCourse
	CourseCategory race_vo.CourseCategory
	RaceDate       race_vo.RaceDate
	Distance       int
	TrackCondition string
	Payment        int
	Repayment      int
	Url            string
	RaceResults    []*race_entity.RaceResult
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
		RaceId:         raceId,
		RaceNumber:     raceNumber,
		RaceName:       raceName,
		Class:          class,
		RaceCourse:     raceCourse,
		CourseCategory: courseCategory,
		RaceDate:       raceDate,
		Distance:       distance,
		TrackCondition: trackCondition,
		Payment:        payment,
		Repayment:      repayment,
		Url:            url,
		RaceResults:    raceResults,
	}
}
