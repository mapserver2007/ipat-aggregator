package entity

import (
	"github.com/mapserver2007/tools/baken/app/domain/race/value_object"
)

func NewRaceCourseSummary(raceCourseRates map[value_object.RaceCourse]ResultRate) RaceCourseSummary {
	return RaceCourseSummary{
		RaceCourseRates: raceCourseRates,
	}
}
