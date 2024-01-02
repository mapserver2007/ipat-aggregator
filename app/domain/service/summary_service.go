package service

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"time"
)

type SummaryService interface {
	CreateSummary(ctx context.Context, tickets []*ticket_csv_entity.Ticket, racingNumbers []*data_cache_entity.RacingNumber, races []*data_cache_entity.Race) *spreadsheet_entity.Summary
}

type summaryService struct {
	ticketAggregator TicketAggregator
}

func NewSummaryService(
	ticketAggregator TicketAggregator,
) SummaryService {
	return &summaryService{
		ticketAggregator: ticketAggregator,
	}
}

func (s *summaryService) CreateSummary(
	ctx context.Context,
	tickets []*ticket_csv_entity.Ticket,
	racingNumbers []*data_cache_entity.RacingNumber,
	races []*data_cache_entity.Race,
) *spreadsheet_entity.Summary {
	now := time.Now()
	allFrom := time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local)
	allResult := s.ticketAggregator.TermResult(ctx, tickets, races, racingNumbers, allFrom, now)

	yearFrom := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	nextYear := now.AddDate(1, 0, 0)
	yearTo := time.Date(nextYear.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	yearResult := s.ticketAggregator.TermResult(ctx, tickets, races, racingNumbers, yearFrom, yearTo)

	monthFrom := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	nextMonth := now.AddDate(0, 1, 0)
	monthTo := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, time.Local)
	monthResult := s.ticketAggregator.TermResult(ctx, tickets, races, racingNumbers, monthFrom, monthTo)

	ticketResultMap := map[types.TicketType]*spreadsheet_entity.TicketResult{}
	ticketResultMap[types.Win] = s.ticketAggregator.TicketResult(ctx, tickets, races, racingNumbers, types.Win)
	ticketResultMap[types.Place] = s.ticketAggregator.TicketResult(ctx, tickets, races, racingNumbers, types.Place)
	ticketResultMap[types.Quinella] = s.ticketAggregator.TicketResult(ctx, tickets, races, racingNumbers, types.Quinella)
	ticketResultMap[types.Exacta] = s.ticketAggregator.TicketResult(ctx, tickets, races, racingNumbers, types.Exacta, types.ExactaWheelOfFirst)
	ticketResultMap[types.QuinellaPlace] = s.ticketAggregator.TicketResult(ctx, tickets, races, racingNumbers, types.QuinellaPlace, types.QuinellaPlaceWheel)
	ticketResultMap[types.Trio] = s.ticketAggregator.TicketResult(ctx, tickets, races, racingNumbers, types.Trio, types.TrioFormation, types.TrifectaWheelOfFirst)
	ticketResultMap[types.Trifecta] = s.ticketAggregator.TicketResult(ctx, tickets, races, racingNumbers, types.Trifecta, types.TrifectaFormation, types.TrifectaWheelOfFirst, types.TrifectaWheelOfSecondMulti)
	ticketResultMap[types.AllTicketType] = s.ticketAggregator.TicketResult(ctx, tickets, races, racingNumbers, types.Win, types.Place, types.Quinella,
		types.Exacta, types.ExactaWheelOfFirst, types.QuinellaPlace, types.QuinellaPlaceWheel,
		types.Trio, types.TrioFormation, types.TrifectaWheelOfFirst,
		types.Trifecta, types.TrifectaFormation, types.TrifectaWheelOfFirst, types.TrifectaWheelOfSecondMulti)

	gradeClassResultMap := map[types.GradeClass]*spreadsheet_entity.TicketResult{}
	gradeClassResultMap[types.Grade1] = s.ticketAggregator.GradeClassResult(ctx, tickets, races, racingNumbers, types.Grade1, types.JumpGrade1)
	gradeClassResultMap[types.Grade2] = s.ticketAggregator.GradeClassResult(ctx, tickets, races, racingNumbers, types.Grade2, types.JumpGrade2)
	gradeClassResultMap[types.Grade3] = s.ticketAggregator.GradeClassResult(ctx, tickets, races, racingNumbers, types.Grade3, types.JumpGrade3)
	gradeClassResultMap[types.Jpn1] = s.ticketAggregator.GradeClassResult(ctx, tickets, races, racingNumbers, types.Jpn1)
	gradeClassResultMap[types.Jpn2] = s.ticketAggregator.GradeClassResult(ctx, tickets, races, racingNumbers, types.Jpn2)
	gradeClassResultMap[types.Jpn3] = s.ticketAggregator.GradeClassResult(ctx, tickets, races, racingNumbers, types.Jpn3)
	gradeClassResultMap[types.OpenClass] = s.ticketAggregator.GradeClassResult(ctx, tickets, races, racingNumbers, types.OpenClass, types.ListedClass, types.LocalGrade)
	gradeClassResultMap[types.ThreeWinClass] = s.ticketAggregator.GradeClassResult(ctx, tickets, races, racingNumbers, types.ThreeWinClass)
	gradeClassResultMap[types.TwoWinClass] = s.ticketAggregator.GradeClassResult(ctx, tickets, races, racingNumbers, types.TwoWinClass)
	gradeClassResultMap[types.OneWinClass] = s.ticketAggregator.GradeClassResult(ctx, tickets, races, racingNumbers, types.OneWinClass)
	gradeClassResultMap[types.Maiden] = s.ticketAggregator.GradeClassResult(ctx, tickets, races, racingNumbers, types.Maiden, types.JumpMaiden)
	gradeClassResultMap[types.MakeDebut] = s.ticketAggregator.GradeClassResult(ctx, tickets, races, racingNumbers, types.MakeDebut)

	courseCategoryResultMap := map[types.CourseCategory]*spreadsheet_entity.TicketResult{}
	courseCategories := []types.CourseCategory{
		types.Turf, types.Dirt, types.Jump,
	}
	for _, courseCategory := range courseCategories {
		courseCategoryResultMap[courseCategory] = s.ticketAggregator.CourseCategoryResult(ctx, tickets, races, racingNumbers, courseCategory)
	}

	distanceCategoryResultMap := map[types.DistanceCategory]*spreadsheet_entity.TicketResult{}
	distanceCategories := []types.DistanceCategory{
		types.TurfSprint, types.TurfMile, types.TurfIntermediate, types.TurfLong, types.TurfExtended,
		types.DirtSprint, types.DirtMile, types.DirtIntermediate, types.DirtLong, types.DirtExtended,
		types.JumpAllDistance,
	}
	for _, distanceCategory := range distanceCategories {
		distanceCategoryResultMap[distanceCategory] = s.ticketAggregator.DistanceCategoryResult(ctx, tickets, races, racingNumbers, distanceCategory)
	}

	raceCourseResultMap := map[types.RaceCourse]*spreadsheet_entity.TicketResult{}
	raceCourses := []types.RaceCourse{
		types.Sapporo, types.Hakodate, types.Fukushima, types.Niigata, types.Tokyo, types.Nakayama, types.Chukyo, types.Kyoto, types.Hanshin, types.Kokura,
		types.Monbetsu, types.Morioka, types.Urawa, types.Hunabashi, types.Ooi, types.Kawasaki, types.Kanazawa, types.Nagoya, types.Sonoda, types.Kouchi, types.Saga,
	}
	for _, raceCourse := range raceCourses {
		raceCourseResultMap[raceCourse] = s.ticketAggregator.RaceCourseResult(ctx, tickets, races, racingNumbers, raceCourse)
	}
	raceCourseResultMap[types.Overseas] = s.ticketAggregator.RaceCourseResult(ctx, tickets, races, racingNumbers, types.Longchamp, types.Deauville, types.Shatin, types.Meydan, types.SantaAnitaPark)

	monthlyResult := s.ticketAggregator.MonthlyResult(ctx, tickets, races, racingNumbers)

	return spreadsheet_entity.NewSummary(
		allResult,
		monthResult,
		yearResult,
		ticketResultMap,
		gradeClassResultMap,
		courseCategoryResultMap,
		distanceCategoryResultMap,
		raceCourseResultMap,
		monthlyResult,
	)
}
