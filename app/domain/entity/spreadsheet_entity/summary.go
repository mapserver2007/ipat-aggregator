package spreadsheet_entity

import (
	"time"

	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type Summary struct {
	allTermResult             *TicketResult
	yearTermResult            *TicketResult
	monthTermResult           *TicketResult
	ticketResultMap           map[types.TicketType]*TicketResult
	gradeClassResultMap       map[types.GradeClass]*TicketResult
	courseCategoryResultMap   map[types.CourseCategory]*TicketResult
	distanceCategoryResultMap map[types.DistanceCategory]*TicketResult
	raceCourseResultMap       map[types.RaceCourse]*TicketResult
	monthlyResults            map[time.Time]*TicketResult
	dailyResults              map[time.Time]*TicketResult
}

func NewSummary(
	allTermResult *TicketResult,
	yearTermResult *TicketResult,
	monthTermResult *TicketResult,
	ticketResultMap map[types.TicketType]*TicketResult,
	gradeClassResultMap map[types.GradeClass]*TicketResult,
	courseCategoryResultMap map[types.CourseCategory]*TicketResult,
	distanceCategoryResultMap map[types.DistanceCategory]*TicketResult,
	raceCourseResultMap map[types.RaceCourse]*TicketResult,
	monthlyResults map[time.Time]*TicketResult,
	dailyResults map[time.Time]*TicketResult,
) *Summary {
	return &Summary{
		allTermResult:             allTermResult,
		yearTermResult:            yearTermResult,
		monthTermResult:           monthTermResult,
		ticketResultMap:           ticketResultMap,
		gradeClassResultMap:       gradeClassResultMap,
		courseCategoryResultMap:   courseCategoryResultMap,
		distanceCategoryResultMap: distanceCategoryResultMap,
		raceCourseResultMap:       raceCourseResultMap,
		monthlyResults:            monthlyResults,
		dailyResults:              dailyResults,
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

func (s *Summary) TicketResultMap() map[types.TicketType]*TicketResult {
	return s.ticketResultMap
}

func (s *Summary) GradeClassResultMap() map[types.GradeClass]*TicketResult {
	return s.gradeClassResultMap
}

func (s *Summary) CourseCategoryResultMap() map[types.CourseCategory]*TicketResult {
	return s.courseCategoryResultMap
}

func (s *Summary) DistanceCategoryResultMap() map[types.DistanceCategory]*TicketResult {
	return s.distanceCategoryResultMap
}

func (s *Summary) RaceCourseResultMap() map[types.RaceCourse]*TicketResult {
	return s.raceCourseResultMap
}

func (s *Summary) MonthlyResults() map[time.Time]*TicketResult {
	return s.monthlyResults
}

func (s *Summary) DailyResults() map[time.Time]*TicketResult {
	return s.dailyResults
}
