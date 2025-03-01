package aggregation_service

import (
	"context"
	"fmt"
	"time"

	"github.com/jinzhu/now"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/summary_service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type Summary interface {
	Create(ctx context.Context,
		tickets []*ticket_csv_entity.RaceTicket,
		races []*data_cache_entity.Race,
	) *spreadsheet_entity.Summary
	WriteV2(ctx context.Context, data *spreadsheet_entity.Summary) error
}

type summaryService struct {
	termService             summary_service.Term
	ticketService           summary_service.Ticket
	classService            summary_service.Class
	courseCategoryService   summary_service.CourseCategory
	distanceCategoryService summary_service.DistanceCategory
	raceCourseService       summary_service.RaceCourse
	spreadSheetRepository   repository.SpreadSheetRepository
}

func NewSummary(
	termService summary_service.Term,
	ticketService summary_service.Ticket,
	classService summary_service.Class,
	courseCategoryService summary_service.CourseCategory,
	distanceCategoryService summary_service.DistanceCategory,
	raceCourseService summary_service.RaceCourse,
	spreadSheetRepository repository.SpreadSheetRepository,
) Summary {
	return &summaryService{
		termService:             termService,
		ticketService:           ticketService,
		classService:            classService,
		courseCategoryService:   courseCategoryService,
		distanceCategoryService: distanceCategoryService,
		raceCourseService:       raceCourseService,
		spreadSheetRepository:   spreadSheetRepository,
	}
}

func (s *summaryService) Create(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
	races []*data_cache_entity.Race,
) *spreadsheet_entity.Summary {
	allTermResult := s.getAllTermResult(ctx, tickets)
	yearTermResult := s.getYearTermResult(ctx, tickets)
	monthTermResult := s.getMonthTermResult(ctx, tickets)
	weekTermResult := s.getWeekTermResult(ctx, tickets)
	ticketResultMap := s.getTicketResultMap(ctx, tickets)
	ticketYearlyResultMap := s.getTicketYearlyResultMap(ctx, tickets)
	ticketMonthlyResultMap := s.getTicketMonthlyResultMap(ctx, tickets)
	classResultMap := s.getClassResultMap(ctx, tickets, races)
	classYearlyResultMap := s.getClassYearlyResultMap(ctx, tickets, races)
	classMonthlyResultMap := s.getClassMonthlyResultMap(ctx, tickets, races)
	courseCategoryResultMap := s.getCourseCategoryResultMap(ctx, tickets, races)
	distanceCategoryResultMap := s.getDistanceCategoryResultMap(ctx, tickets, races)
	distanceCategoryYearlyResultMap := s.getDistanceCategoryYearlyResultMap(ctx, tickets, races)
	distanceCategoryMonthlyResultMap := s.getDistanceCategoryMonthlyResultMap(ctx, tickets, races)
	raceCourseResultMap := s.getRaceCourseResultMap(ctx, tickets, races)
	raceCourseYearlyResultMap := s.getRaceCourseYearlyResultMap(ctx, tickets, races)
	raceCourseMonthlyResultMap := s.getRaceCourseMonthlyResultMap(ctx, tickets, races)
	yearlyResultMap := s.getYearlyResultMap(ctx, tickets)
	monthlyResultMap := s.getMonthlyResultMap(ctx, tickets)
	weeklyResultMap := s.getWeeklyResultMap(ctx, tickets)

	return spreadsheet_entity.NewSummary(
		allTermResult,
		yearTermResult,
		monthTermResult,
		weekTermResult,
		ticketResultMap,
		ticketYearlyResultMap,
		ticketMonthlyResultMap,
		classResultMap,
		classYearlyResultMap,
		classMonthlyResultMap,
		courseCategoryResultMap,
		distanceCategoryResultMap,
		distanceCategoryYearlyResultMap,
		distanceCategoryMonthlyResultMap,
		raceCourseResultMap,
		raceCourseYearlyResultMap,
		raceCourseMonthlyResultMap,
		yearlyResultMap,
		monthlyResultMap,
		weeklyResultMap,
	)
}

func (s *summaryService) WriteV2(
	ctx context.Context,
	data *spreadsheet_entity.Summary,
) error {
	return s.spreadSheetRepository.WriteSummaryV2(ctx, data)
}

func (s *summaryService) createTermResult(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
	from time.Time,
	to time.Time,
) *spreadsheet_entity.TicketResult {
	output := s.termService.Create(ctx, &summary_service.TermInput{
		Tickets:  tickets,
		DateFrom: from,
		DateTo:   to,
	})
	return spreadsheet_entity.NewTicketResult(
		output.RaceCount,
		output.BetCount,
		output.HitCount,
		output.Payment,
		output.Payout,
		output.AveragePayout,
		output.MaxPayout,
		output.MinPayout,
	)
}

func (s *summaryService) getAllTermResult(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
) *spreadsheet_entity.TicketResult {
	now := time.Now()
	allFrom := time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local)
	result := s.createTermResult(ctx, tickets, allFrom, now)

	return result
}

func (s *summaryService) getYearTermResult(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
) *spreadsheet_entity.TicketResult {
	now := time.Now()
	yearFrom := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	nextYear := now.AddDate(1, 0, 0)
	yearTo := time.Date(nextYear.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	result := s.createTermResult(ctx, tickets, yearFrom, yearTo)

	return result
}

func (s *summaryService) getMonthTermResult(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
) *spreadsheet_entity.TicketResult {
	now := time.Now()
	monthFrom := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	nextMonth := now.AddDate(0, 1, 0)
	monthTo := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, time.Local)
	result := s.createTermResult(ctx, tickets, monthFrom, monthTo)

	return result
}

func (s *summaryService) getWeekTermResult(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
) *spreadsheet_entity.TicketResult {
	currentTime := time.Now()

	lastWeekMonday := now.New(currentTime).Monday().AddDate(0, 0, -7)
	thisWeekMonday := now.New(currentTime).Monday()

	currentMonthStart := now.New(currentTime).BeginningOfMonth()
	if lastWeekMonday.Before(currentMonthStart) {
		lastWeekMonday = currentMonthStart
	}

	result := s.createTermResult(ctx, tickets, lastWeekMonday, thisWeekMonday)

	return result
}

func (s *summaryService) getTermRaceTicket(
	tickets []*ticket_csv_entity.RaceTicket,
	from time.Time,
	to time.Time,
) []*ticket_csv_entity.RaceTicket {
	termRaceTicketMap := []*ticket_csv_entity.RaceTicket{}
	for _, ticket := range tickets {
		if ticket.Ticket().RaceDate().Date().Unix() >= from.Unix() && ticket.Ticket().RaceDate().Date().Unix() < to.Unix() {
			termRaceTicketMap = append(termRaceTicketMap, ticket)
		}
	}

	return termRaceTicketMap
}

func (s *summaryService) getTicketResultMap(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
) map[types.TicketType]*spreadsheet_entity.TicketResult {
	ticketResultMap := map[types.TicketType]*spreadsheet_entity.TicketResult{}
	ticketResultMap[types.Win] = s.createTicketResult(ctx, tickets, []types.TicketType{types.Win})
	ticketResultMap[types.Place] = s.createTicketResult(ctx, tickets, []types.TicketType{types.Place})
	ticketResultMap[types.Quinella] = s.createTicketResult(ctx, tickets, []types.TicketType{types.Quinella, types.QuinellaWheel})
	ticketResultMap[types.Exacta] = s.createTicketResult(ctx, tickets, []types.TicketType{types.Exacta, types.ExactaWheelOfFirst})
	ticketResultMap[types.QuinellaPlace] = s.createTicketResult(ctx, tickets, []types.TicketType{types.QuinellaPlace, types.QuinellaPlaceWheel, types.QuinellaPlaceFormation})
	ticketResultMap[types.Trio] = s.createTicketResult(ctx, tickets, []types.TicketType{types.Trio, types.TrioFormation, types.TrioWheelOfFirst, types.TrioWheelOfSecond, types.TrioBox})
	ticketResultMap[types.Trifecta] = s.createTicketResult(ctx, tickets, []types.TicketType{types.Trifecta, types.TrifectaFormation, types.TrifectaWheelOfFirst, types.TrifectaWheelOfSecond, types.TrifectaWheelOfFirstMulti, types.TrifectaWheelOfSecondMulti})
	ticketResultMap[types.AllTicketType] = s.createTicketResult(ctx, tickets, []types.TicketType{types.Win, types.Place, types.Quinella, types.QuinellaWheel,
		types.Exacta, types.ExactaWheelOfFirst, types.QuinellaPlace, types.QuinellaPlaceWheel, types.QuinellaPlaceFormation,
		types.Trio, types.TrioFormation, types.TrioWheelOfFirst, types.TrioWheelOfSecond, types.TrioBox,
		types.Trifecta, types.TrifectaFormation, types.TrifectaWheelOfFirst, types.TrifectaWheelOfSecond, types.TrifectaWheelOfFirstMulti, types.TrifectaWheelOfSecondMulti})

	return ticketResultMap
}

func (s *summaryService) getTicketYearlyResultMap(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
) map[types.TicketType]*spreadsheet_entity.TicketResult {
	now := time.Now()
	yearFrom := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	nextYear := now.AddDate(1, 0, 0)
	yearTo := time.Date(nextYear.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	ticketYearlyResultMap := s.getTicketResultMap(ctx, s.getTermRaceTicket(tickets, yearFrom, yearTo))

	return ticketYearlyResultMap
}

func (s *summaryService) getTicketMonthlyResultMap(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
) map[types.TicketType]*spreadsheet_entity.TicketResult {
	now := time.Now()
	monthFrom := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	nextMonth := now.AddDate(0, 1, 0)
	monthTo := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, time.Local)
	ticketMonthlyResultMap := s.getTicketResultMap(ctx, s.getTermRaceTicket(tickets, monthFrom, monthTo))

	return ticketMonthlyResultMap
}

func (s *summaryService) getClassResultMap(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
	races []*data_cache_entity.Race,
) map[types.GradeClass]*spreadsheet_entity.TicketResult {
	classResultMap := map[types.GradeClass]*spreadsheet_entity.TicketResult{}
	classResultMap[types.Grade1] = s.createClassResult(ctx, tickets, races, []types.GradeClass{types.Grade1, types.JumpGrade1})
	classResultMap[types.Grade2] = s.createClassResult(ctx, tickets, races, []types.GradeClass{types.Grade2, types.JumpGrade2})
	classResultMap[types.Grade3] = s.createClassResult(ctx, tickets, races, []types.GradeClass{types.Grade3, types.JumpGrade3})
	classResultMap[types.Jpn1] = s.createClassResult(ctx, tickets, races, []types.GradeClass{types.Jpn1})
	classResultMap[types.Jpn2] = s.createClassResult(ctx, tickets, races, []types.GradeClass{types.Jpn2})
	classResultMap[types.Jpn3] = s.createClassResult(ctx, tickets, races, []types.GradeClass{types.Jpn3})
	classResultMap[types.OpenClass] = s.createClassResult(ctx, tickets, races, []types.GradeClass{types.OpenClass, types.ListedClass, types.LocalGrade})
	classResultMap[types.ThreeWinClass] = s.createClassResult(ctx, tickets, races, []types.GradeClass{types.ThreeWinClass})
	classResultMap[types.TwoWinClass] = s.createClassResult(ctx, tickets, races, []types.GradeClass{types.TwoWinClass})
	classResultMap[types.OneWinClass] = s.createClassResult(ctx, tickets, races, []types.GradeClass{types.OneWinClass})
	classResultMap[types.Maiden] = s.createClassResult(ctx, tickets, races, []types.GradeClass{types.Maiden, types.JumpMaiden})
	classResultMap[types.MakeDebut] = s.createClassResult(ctx, tickets, races, []types.GradeClass{types.MakeDebut})

	return classResultMap
}

func (s *summaryService) getClassYearlyResultMap(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
	races []*data_cache_entity.Race,
) map[types.GradeClass]*spreadsheet_entity.TicketResult {
	now := time.Now()
	yearFrom := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	nextYear := now.AddDate(1, 0, 0)
	yearTo := time.Date(nextYear.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	classYearlyResultMap := s.getClassResultMap(ctx, s.getTermRaceTicket(tickets, yearFrom, yearTo), races)

	return classYearlyResultMap
}

func (s *summaryService) getClassMonthlyResultMap(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
	races []*data_cache_entity.Race,
) map[types.GradeClass]*spreadsheet_entity.TicketResult {
	now := time.Now()
	monthFrom := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	nextMonth := now.AddDate(0, 1, 0)
	monthTo := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, time.Local)
	classMonthlyResultMap := s.getClassResultMap(ctx, s.getTermRaceTicket(tickets, monthFrom, monthTo), races)

	return classMonthlyResultMap
}

func (s *summaryService) getCourseCategoryResultMap(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
	races []*data_cache_entity.Race,
) map[types.CourseCategory]*spreadsheet_entity.TicketResult {
	courseCategoryResultMap := map[types.CourseCategory]*spreadsheet_entity.TicketResult{}
	courseCategories := []types.CourseCategory{
		types.Turf, types.Dirt, types.Jump,
	}
	for _, courseCategory := range courseCategories {
		courseCategoryResultMap[courseCategory] = s.createCourseCategoryResult(ctx, tickets, races, []types.CourseCategory{courseCategory})
	}

	return courseCategoryResultMap
}

func (s *summaryService) getDistanceCategoryResultMap(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
	races []*data_cache_entity.Race,
) map[types.DistanceCategory]*spreadsheet_entity.TicketResult {
	distanceCategoryResultMap := map[types.DistanceCategory]*spreadsheet_entity.TicketResult{}
	distanceCategories := []types.DistanceCategory{
		types.TurfSprint, types.TurfMile, types.TurfIntermediate, types.TurfLong, types.TurfExtended,
		types.DirtSprint, types.DirtMile, types.DirtIntermediate, types.DirtLong,
		types.JumpAllDistance,
	}
	for _, distanceCategory := range distanceCategories {
		distanceCategoryResultMap[distanceCategory] = s.createDistanceCategoryResult(ctx, tickets, races, []types.DistanceCategory{distanceCategory})
	}

	return distanceCategoryResultMap
}

func (s *summaryService) getDistanceCategoryYearlyResultMap(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
	races []*data_cache_entity.Race,
) map[types.DistanceCategory]*spreadsheet_entity.TicketResult {
	now := time.Now()
	yearFrom := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	nextYear := now.AddDate(1, 0, 0)
	yearTo := time.Date(nextYear.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	distanceCategoryYearlyResultMap := s.getDistanceCategoryResultMap(ctx, s.getTermRaceTicket(tickets, yearFrom, yearTo), races)

	return distanceCategoryYearlyResultMap
}

func (s *summaryService) getDistanceCategoryMonthlyResultMap(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
	races []*data_cache_entity.Race,
) map[types.DistanceCategory]*spreadsheet_entity.TicketResult {
	now := time.Now()
	monthFrom := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	nextMonth := now.AddDate(0, 1, 0)
	monthTo := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, time.Local)
	distanceCategoryMonthlyResultMap := s.getDistanceCategoryResultMap(ctx, s.getTermRaceTicket(tickets, monthFrom, monthTo), races)

	return distanceCategoryMonthlyResultMap
}

func (s *summaryService) getRaceCourseResultMap(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
	races []*data_cache_entity.Race,
) map[types.RaceCourse]*spreadsheet_entity.TicketResult {
	raceCourseResultMap := map[types.RaceCourse]*spreadsheet_entity.TicketResult{}
	raceCourses := []types.RaceCourse{
		types.Sapporo, types.Hakodate, types.Fukushima, types.Niigata, types.Tokyo, types.Nakayama, types.Chukyo, types.Kyoto, types.Hanshin, types.Kokura,
		types.Monbetsu, types.Morioka, types.Urawa, types.Hunabashi, types.Ooi, types.Kawasaki, types.Kanazawa, types.Nagoya, types.Sonoda, types.Kouchi, types.Saga,
	}
	for _, raceCourse := range raceCourses {
		raceCourseResultMap[raceCourse] = s.createRaceCourseResult(ctx, tickets, races, []types.RaceCourse{raceCourse})
	}
	raceCourseResultMap[types.Overseas] = s.createRaceCourseResult(ctx, tickets, races, []types.RaceCourse{types.Longchamp, types.Deauville, types.Shatin, types.Meydan, types.SantaAnitaPark, types.KingAbdulaziz, types.York, types.Delmar})

	return raceCourseResultMap
}

func (s *summaryService) getRaceCourseYearlyResultMap(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
	races []*data_cache_entity.Race,
) map[types.RaceCourse]*spreadsheet_entity.TicketResult {
	now := time.Now()
	yearFrom := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	nextYear := now.AddDate(1, 0, 0)
	yearTo := time.Date(nextYear.Year(), 1, 1, 0, 0, 0, 0, time.Local)
	raceCourseYearlyResultMap := s.getRaceCourseResultMap(ctx, s.getTermRaceTicket(tickets, yearFrom, yearTo), races)

	return raceCourseYearlyResultMap
}

func (s *summaryService) getRaceCourseMonthlyResultMap(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
	races []*data_cache_entity.Race,
) map[types.RaceCourse]*spreadsheet_entity.TicketResult {
	now := time.Now()
	monthFrom := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
	nextMonth := now.AddDate(0, 1, 0)
	monthTo := time.Date(nextMonth.Year(), nextMonth.Month(), 1, 0, 0, 0, 0, time.Local)
	raceCourseMonthlyResultMap := s.getRaceCourseResultMap(ctx, s.getTermRaceTicket(tickets, monthFrom, monthTo), races)

	return raceCourseMonthlyResultMap
}

func (s *summaryService) getYearlyResultMap(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
) map[time.Time]*spreadsheet_entity.TicketResult {
	dateTimeTicketMap := map[time.Time][]*ticket_csv_entity.RaceTicket{}
	for _, raceTicket := range tickets {
		dateStr := fmt.Sprintf("%d", raceTicket.Ticket().RaceDate().Value())
		dateTime, _ := time.Parse("20060102", dateStr)
		year := time.Date(dateTime.Year(), 1, 1, 0, 0, 0, 0, time.Local)
		if _, ok := dateTimeTicketMap[year]; !ok {
			dateTimeTicketMap[year] = make([]*ticket_csv_entity.RaceTicket, 0)
		}
		dateTimeTicketMap[year] = append(dateTimeTicketMap[year], raceTicket)
	}

	yearlyResultMap := map[time.Time]*spreadsheet_entity.TicketResult{}
	for currentYear, raceTickets := range dateTimeTicketMap {
		yearlyResultMap[currentYear] = s.createTermResult(ctx, raceTickets, currentYear, currentYear.AddDate(1, 0, 0))
	}

	return yearlyResultMap
}

func (s *summaryService) getMonthlyResultMap(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
) map[time.Time]*spreadsheet_entity.TicketResult {
	now := time.Now()
	dateTimeTicketMap := map[time.Time][]*ticket_csv_entity.RaceTicket{}
	for _, raceTicket := range tickets {
		dateStr := fmt.Sprintf("%d", raceTicket.Ticket().RaceDate().Value())
		dateTime, _ := time.Parse("20060102", dateStr)
		month := time.Date(dateTime.Year(), dateTime.Month(), 1, 0, 0, 0, 0, time.Local)
		if _, ok := dateTimeTicketMap[month]; !ok {
			dateTimeTicketMap[month] = make([]*ticket_csv_entity.RaceTicket, 0)
		}
		dateTimeTicketMap[month] = append(dateTimeTicketMap[month], raceTicket)
	}

	monthlyResultMap := map[time.Time]*spreadsheet_entity.TicketResult{}
	for currentMonth, raceTickets := range dateTimeTicketMap {
		monthlyResultMap[currentMonth] = s.createTermResult(ctx, raceTickets, currentMonth, now.AddDate(0, 1, 0))
	}

	return monthlyResultMap
}

func (s *summaryService) getWeeklyResultMap(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
) map[time.Time]*spreadsheet_entity.TicketResult {
	now := time.Now()
	dateTimeTicketMap := map[time.Time][]*ticket_csv_entity.RaceTicket{}
	for _, raceTicket := range tickets {
		dateStr := fmt.Sprintf("%d", raceTicket.Ticket().RaceDate().Value())
		dateTime, _ := time.Parse("20060102", dateStr)
		if dateTime.Month() == now.Month() && dateTime.Year() == now.Year() {
			if _, ok := dateTimeTicketMap[dateTime]; !ok {
				dateTimeTicketMap[dateTime] = make([]*ticket_csv_entity.RaceTicket, 0)
			}
			dateTimeTicketMap[dateTime] = append(dateTimeTicketMap[dateTime], raceTicket)
		}
	}

	weeklyResultMap := map[time.Time]*spreadsheet_entity.TicketResult{}
	for weekStart, raceTickets := range dateTimeTicketMap {
		weeklyResultMap[weekStart] = s.createTermResult(ctx, raceTickets, weekStart, weekStart.AddDate(0, 0, 1))
	}

	return weeklyResultMap
}

func (s *summaryService) createTicketResult(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
	ticketTypes []types.TicketType,
) *spreadsheet_entity.TicketResult {
	output := s.ticketService.Create(ctx, &summary_service.TicketInput{
		Tickets:     tickets,
		TicketTypes: ticketTypes,
	})
	return spreadsheet_entity.NewTicketResult(
		output.RaceCount,
		output.BetCount,
		output.HitCount,
		output.Payment,
		output.Payout,
		output.AveragePayout,
		output.MaxPayout,
		output.MinPayout,
	)
}

func (s *summaryService) createClassResult(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
	races []*data_cache_entity.Race,
	classes []types.GradeClass,
) *spreadsheet_entity.TicketResult {
	output := s.classService.Create(ctx, &summary_service.ClassInput{
		Tickets: tickets,
		Races:   races,
		Classes: classes,
	})
	return spreadsheet_entity.NewTicketResult(
		output.RaceCount,
		output.BetCount,
		output.HitCount,
		output.Payment,
		output.Payout,
		output.AveragePayout,
		output.MaxPayout,
		output.MinPayout,
	)
}

func (s *summaryService) createCourseCategoryResult(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
	races []*data_cache_entity.Race,
	courseCategories []types.CourseCategory,
) *spreadsheet_entity.TicketResult {
	output := s.courseCategoryService.Create(ctx, &summary_service.CourseCategoryInput{
		Tickets:          tickets,
		Races:            races,
		CourseCategories: courseCategories,
	})
	return spreadsheet_entity.NewTicketResult(
		output.RaceCount,
		output.BetCount,
		output.HitCount,
		output.Payment,
		output.Payout,
		output.AveragePayout,
		output.MaxPayout,
		output.MinPayout,
	)
}

func (s *summaryService) createDistanceCategoryResult(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
	races []*data_cache_entity.Race,
	distanceCategories []types.DistanceCategory,
) *spreadsheet_entity.TicketResult {
	output := s.distanceCategoryService.Create(ctx, &summary_service.DistanceCategoryInput{
		Tickets:            tickets,
		Races:              races,
		DistanceCategories: distanceCategories,
	})
	return spreadsheet_entity.NewTicketResult(
		output.RaceCount,
		output.BetCount,
		output.HitCount,
		output.Payment,
		output.Payout,
		output.AveragePayout,
		output.MaxPayout,
		output.MinPayout,
	)
}

func (s *summaryService) createRaceCourseResult(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
	races []*data_cache_entity.Race,
	raceCourses []types.RaceCourse,
) *spreadsheet_entity.TicketResult {
	output := s.raceCourseService.Create(ctx, &summary_service.RaceCourseInput{
		Tickets:     tickets,
		Races:       races,
		RaceCourses: raceCourses,
	})
	return spreadsheet_entity.NewTicketResult(
		output.RaceCount,
		output.BetCount,
		output.HitCount,
		output.Payment,
		output.Payout,
		output.AveragePayout,
		output.MaxPayout,
		output.MinPayout,
	)
}
