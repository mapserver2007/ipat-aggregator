package entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
)

func NewCourseCategorySummary(courseCategoryRates map[value_object.CourseCategory]ResultRate) CourseCategorySummary {
	return CourseCategorySummary{
		CourseCategoryRates: courseCategoryRates,
	}
}
