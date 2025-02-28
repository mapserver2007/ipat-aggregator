package spreadsheet_entity

import (
	"time"

	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type Summary struct {
	allTermResult                    *TicketResult
	yearTermResult                   *TicketResult
	monthTermResult                  *TicketResult
	weekTermResult                   *TicketResult
	ticketResultMap                  map[types.TicketType]*TicketResult
	ticketYearlyResultMap            map[types.TicketType]*TicketResult
	ticketMonthlyResultMap           map[types.TicketType]*TicketResult
	gradeClassResultMap              map[types.GradeClass]*TicketResult
	gradeClassYearlyResultMap        map[types.GradeClass]*TicketResult
	gradeClassMonthlyResultMap       map[types.GradeClass]*TicketResult
	courseCategoryResultMap          map[types.CourseCategory]*TicketResult
	distanceCategoryResultMap        map[types.DistanceCategory]*TicketResult
	distanceCategoryYearlyResultMap  map[types.DistanceCategory]*TicketResult
	distanceCategoryMonthlyResultMap map[types.DistanceCategory]*TicketResult
	raceCourseResultMap              map[types.RaceCourse]*TicketResult
	raceCourseYearlyResultMap        map[types.RaceCourse]*TicketResult
	raceCourseMonthlyResultMap       map[types.RaceCourse]*TicketResult
	yearlyResults                    map[time.Time]*TicketResult
	monthlyResults                   map[time.Time]*TicketResult
	weeklyResults                    map[time.Time]*TicketResult
}

func NewSummary(
	allTermResult *TicketResult,
	yearTermResult *TicketResult,
	monthTermResult *TicketResult,
	weekTermResult *TicketResult,
	ticketResultMap map[types.TicketType]*TicketResult,
	ticketYearlyResultMap map[types.TicketType]*TicketResult,
	ticketMonthlyResultMap map[types.TicketType]*TicketResult,
	gradeClassResultMap map[types.GradeClass]*TicketResult,
	gradeClassYearlyResultMap map[types.GradeClass]*TicketResult,
	gradeClassMonthlyResultMap map[types.GradeClass]*TicketResult,
	courseCategoryResultMap map[types.CourseCategory]*TicketResult,
	distanceCategoryResultMap map[types.DistanceCategory]*TicketResult,
	distanceCategoryYearlyResultMap map[types.DistanceCategory]*TicketResult,
	distanceCategoryMonthlyResultMap map[types.DistanceCategory]*TicketResult,
	raceCourseResultMap map[types.RaceCourse]*TicketResult,
	raceCourseYearlyResultMap map[types.RaceCourse]*TicketResult,
	raceCourseMonthlyResultMap map[types.RaceCourse]*TicketResult,
	yearlyResults map[time.Time]*TicketResult,
	monthlyResults map[time.Time]*TicketResult,
	weeklyResults map[time.Time]*TicketResult,
) *Summary {
	return &Summary{
		allTermResult:                    allTermResult,
		yearTermResult:                   yearTermResult,
		monthTermResult:                  monthTermResult,
		weekTermResult:                   weekTermResult,
		ticketResultMap:                  ticketResultMap,
		ticketYearlyResultMap:            ticketYearlyResultMap,
		ticketMonthlyResultMap:           ticketMonthlyResultMap,
		gradeClassResultMap:              gradeClassResultMap,
		gradeClassYearlyResultMap:        gradeClassYearlyResultMap,
		gradeClassMonthlyResultMap:       gradeClassMonthlyResultMap,
		courseCategoryResultMap:          courseCategoryResultMap,
		distanceCategoryResultMap:        distanceCategoryResultMap,
		distanceCategoryYearlyResultMap:  distanceCategoryYearlyResultMap,
		distanceCategoryMonthlyResultMap: distanceCategoryMonthlyResultMap,
		raceCourseResultMap:              raceCourseResultMap,
		raceCourseYearlyResultMap:        raceCourseYearlyResultMap,
		raceCourseMonthlyResultMap:       raceCourseMonthlyResultMap,
		yearlyResults:                    yearlyResults,
		monthlyResults:                   monthlyResults,
		weeklyResults:                    weeklyResults,
	}
}

func (s *Summary) AllTermResult() *TicketResult {
	return s.allTermResult
}

func (s *Summary) YearTermResult() *TicketResult {
	return s.yearTermResult
}

func (s *Summary) MonthTermResult() *TicketResult {
	return s.monthTermResult
}

func (s *Summary) WeekTermResult() *TicketResult {
	return s.weekTermResult
}

func (s *Summary) TicketResultMap() map[types.TicketType]*TicketResult {
	return s.ticketResultMap
}

func (s *Summary) TicketYearlyResultMap() map[types.TicketType]*TicketResult {
	return s.ticketYearlyResultMap
}

func (s *Summary) TicketMonthlyResultMap() map[types.TicketType]*TicketResult {
	return s.ticketMonthlyResultMap
}

func (s *Summary) GradeClassResultMap() map[types.GradeClass]*TicketResult {
	return s.gradeClassResultMap
}

func (s *Summary) GradeClassYearlyResultMap() map[types.GradeClass]*TicketResult {
	return s.gradeClassYearlyResultMap
}

func (s *Summary) GradeClassMonthlyResultMap() map[types.GradeClass]*TicketResult {
	return s.gradeClassMonthlyResultMap
}

func (s *Summary) CourseCategoryResultMap() map[types.CourseCategory]*TicketResult {
	return s.courseCategoryResultMap
}

func (s *Summary) DistanceCategoryResultMap() map[types.DistanceCategory]*TicketResult {
	return s.distanceCategoryResultMap
}

func (s *Summary) DistanceCategoryYearlyResultMap() map[types.DistanceCategory]*TicketResult {
	return s.distanceCategoryYearlyResultMap
}

func (s *Summary) DistanceCategoryMonthlyResultMap() map[types.DistanceCategory]*TicketResult {
	return s.distanceCategoryMonthlyResultMap
}

func (s *Summary) RaceCourseResultMap() map[types.RaceCourse]*TicketResult {
	return s.raceCourseResultMap
}

func (s *Summary) RaceCourseYearlyResultMap() map[types.RaceCourse]*TicketResult {
	return s.raceCourseYearlyResultMap
}

func (s *Summary) RaceCourseMonthlyResultMap() map[types.RaceCourse]*TicketResult {
	return s.raceCourseMonthlyResultMap
}

func (s *Summary) YearlyResults() map[time.Time]*TicketResult {
	return s.yearlyResults
}

func (s *Summary) MonthlyResults() map[time.Time]*TicketResult {
	return s.monthlyResults
}

func (s *Summary) WeeklyResults() map[time.Time]*TicketResult {
	return s.weeklyResults
}
