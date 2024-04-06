package service

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"time"
)

type TicketAggregator interface {
	TermResult(ctx context.Context, tickets []*ticket_csv_entity.Ticket, races []*data_cache_entity.Race, racingNumbers []*data_cache_entity.RacingNumber, from, to time.Time) *spreadsheet_entity.TicketResult
	TicketResult(ctx context.Context, tickets []*ticket_csv_entity.Ticket, races []*data_cache_entity.Race, racingNumbers []*data_cache_entity.RacingNumber, ticketTypes ...types.TicketType) *spreadsheet_entity.TicketResult
	GradeClassResult(ctx context.Context, tickets []*ticket_csv_entity.Ticket, races []*data_cache_entity.Race, racingNumbers []*data_cache_entity.RacingNumber, gradeClasses ...types.GradeClass) *spreadsheet_entity.TicketResult
	CourseCategoryResult(ctx context.Context, tickets []*ticket_csv_entity.Ticket, races []*data_cache_entity.Race, racingNumbers []*data_cache_entity.RacingNumber, courseCategories ...types.CourseCategory) *spreadsheet_entity.TicketResult
	DistanceCategoryResult(ctx context.Context, tickets []*ticket_csv_entity.Ticket, races []*data_cache_entity.Race, racingNumbers []*data_cache_entity.RacingNumber, distanceCategories ...types.DistanceCategory) *spreadsheet_entity.TicketResult
	RaceCourseResult(ctx context.Context, tickets []*ticket_csv_entity.Ticket, races []*data_cache_entity.Race, racingNumbers []*data_cache_entity.RacingNumber, raceCourses ...types.RaceCourse) *spreadsheet_entity.TicketResult
	MonthlyResult(ctx context.Context, tickets []*ticket_csv_entity.Ticket, races []*data_cache_entity.Race, racingNumbers []*data_cache_entity.RacingNumber) map[int]*spreadsheet_entity.TicketResult
}

type ticketAggregator struct {
	ticketConverter TicketConverter
}

func (t *ticketAggregator) MonthlyResult(ctx context.Context, tickets []*ticket_csv_entity.Ticket, races []*data_cache_entity.Race, racingNumbers []*data_cache_entity.RacingNumber) map[int]*spreadsheet_entity.TicketResult {
	ticketsGroup := t.ticketConverter.ConvertToMonthlyMap(ctx, tickets)
	ticketResultMap := map[int]*spreadsheet_entity.TicketResult{}
	for month, monthTickets := range ticketsGroup {
		ticketResultMap[month] = t.TicketResult(ctx, monthTickets, races, racingNumbers, types.Win, types.Place, types.Quinella, types.QuinellaWheel,
			types.Exacta, types.ExactaWheelOfFirst, types.QuinellaPlace, types.QuinellaPlaceWheel,
			types.Trio, types.TrioFormation, types.TrioWheelOfFirst, types.TrioWheelOfSecond, types.TrioBox,
			types.Trifecta, types.TrifectaFormation, types.TrifectaWheelOfFirst, types.TrifectaWheelOfSecond, types.TrifectaWheelOfFirstMulti, types.TrifectaWheelOfSecondMulti)
	}
	return ticketResultMap
}

func NewTicketAggregator(
	ticketConverter TicketConverter,
) TicketAggregator {
	return &ticketAggregator{
		ticketConverter: ticketConverter,
	}
}

func (t *ticketAggregator) TermResult(ctx context.Context, tickets []*ticket_csv_entity.Ticket, races []*data_cache_entity.Race, racingNumbers []*data_cache_entity.RacingNumber, from, to time.Time) *spreadsheet_entity.TicketResult {
	tickets = t.ticketConverter.ConvertToTermTickets(ctx, tickets, from, to)
	var hitTickets []*ticket_csv_entity.Ticket
	for _, ticket := range tickets {
		if ticket.TicketResult() == types.TicketHit {
			hitTickets = append(hitTickets, ticket)
		}
	}
	hitCount := len(hitTickets)
	maxPayout := 0
	minPayout := 0
	for _, ticket := range hitTickets {
		if maxPayout < ticket.Payout().Value() {
			maxPayout = ticket.Payout().Value()
		}
		if minPayout == 0 || minPayout > ticket.Payout().Value() {
			minPayout = ticket.Payout().Value()
		}
	}

	ticketsByRaceId := t.ticketConverter.ConvertToRaceIdMap(ctx, tickets, racingNumbers)
	payment, payout := t.getSumAmount(tickets)
	betCount := types.BetCount(len(tickets))
	raceCount := types.RaceCount(len(ticketsByRaceId))
	averagePayout := 0
	if len(hitTickets) > 0 {
		averagePayout = int(float64(payout) / float64(len(hitTickets)))
	}

	return spreadsheet_entity.NewTicketResult(
		raceCount,
		betCount,
		types.HitCount(hitCount),
		payment,
		payout,
		types.Payout(averagePayout),
		types.Payout(maxPayout),
		types.Payout(minPayout),
	)
}

func (t *ticketAggregator) TicketResult(ctx context.Context, tickets []*ticket_csv_entity.Ticket, races []*data_cache_entity.Race, racingNumbers []*data_cache_entity.RacingNumber, ticketTypes ...types.TicketType) *spreadsheet_entity.TicketResult {
	ticketsGroup := t.ticketConverter.ConvertToTicketTypeMap(ctx, tickets)
	var (
		mergedTickets []*ticket_csv_entity.Ticket
		hitTickets    []*ticket_csv_entity.Ticket
	)
	for _, ticketType := range ticketTypes {
		if ticketsByTicketType, ok := ticketsGroup[ticketType]; ok {
			for _, ticket := range ticketsByTicketType {
				if ticket.TicketResult() == types.TicketHit {
					hitTickets = append(hitTickets, ticket)
				}
			}
			mergedTickets = append(mergedTickets, ticketsByTicketType...)
		}
	}
	hitCount := len(hitTickets)
	maxPayout := 0
	minPayout := 0
	for _, ticket := range hitTickets {
		if maxPayout < ticket.Payout().Value() {
			maxPayout = ticket.Payout().Value()
		}
		if minPayout == 0 || minPayout > ticket.Payout().Value() {
			minPayout = ticket.Payout().Value()
		}
	}

	ticketsByRaceId := t.ticketConverter.ConvertToRaceIdMap(ctx, mergedTickets, racingNumbers)
	payment, payout := t.getSumAmount(mergedTickets)
	betCount := types.BetCount(len(mergedTickets))
	raceCount := types.RaceCount(len(ticketsByRaceId))
	averagePayout := 0
	if len(hitTickets) > 0 {
		averagePayout = int(float64(payout) / float64(len(hitTickets)))
	}

	return spreadsheet_entity.NewTicketResult(
		raceCount,
		betCount,
		types.HitCount(hitCount),
		payment,
		payout,
		types.Payout(averagePayout),
		types.Payout(maxPayout),
		types.Payout(minPayout),
	)
}

func (t *ticketAggregator) GradeClassResult(ctx context.Context, tickets []*ticket_csv_entity.Ticket, races []*data_cache_entity.Race, racingNumbers []*data_cache_entity.RacingNumber, gradeClasses ...types.GradeClass) *spreadsheet_entity.TicketResult {
	ticketsGroup := t.ticketConverter.ConvertToGradeClassMap(ctx, tickets, racingNumbers, races)
	var (
		mergedTickets []*ticket_csv_entity.Ticket
		hitTickets    []*ticket_csv_entity.Ticket
	)
	for _, gradeClass := range gradeClasses {
		if ticketsByGradeClass, ok := ticketsGroup[gradeClass]; ok {
			for _, ticket := range ticketsByGradeClass {
				if ticket.TicketResult() == types.TicketHit {
					hitTickets = append(hitTickets, ticket)
				}
			}
			mergedTickets = append(mergedTickets, ticketsByGradeClass...)
		}
	}
	hitCount := len(hitTickets)
	maxPayout := 0
	minPayout := 0
	for _, ticket := range hitTickets {
		if maxPayout < ticket.Payout().Value() {
			maxPayout = ticket.Payout().Value()
		}
		if minPayout == 0 || minPayout > ticket.Payout().Value() {
			minPayout = ticket.Payout().Value()
		}
	}

	ticketsByRaceId := t.ticketConverter.ConvertToRaceIdMap(ctx, mergedTickets, racingNumbers)
	payment, payout := t.getSumAmount(mergedTickets)
	betCount := types.BetCount(len(mergedTickets))
	raceCount := types.RaceCount(len(ticketsByRaceId))
	averagePayout := 0
	if len(hitTickets) > 0 {
		averagePayout = int(float64(payout) / float64(len(hitTickets)))
	}

	return spreadsheet_entity.NewTicketResult(
		raceCount,
		betCount,
		types.HitCount(hitCount),
		payment,
		payout,
		types.Payout(averagePayout),
		types.Payout(maxPayout),
		types.Payout(minPayout),
	)
}

func (t *ticketAggregator) CourseCategoryResult(ctx context.Context, tickets []*ticket_csv_entity.Ticket, races []*data_cache_entity.Race, racingNumbers []*data_cache_entity.RacingNumber, courseCategories ...types.CourseCategory) *spreadsheet_entity.TicketResult {
	ticketsGroup := t.ticketConverter.ConvertToCourseCategoryMap(ctx, tickets, racingNumbers, races)
	var (
		mergedTickets []*ticket_csv_entity.Ticket
		hitTickets    []*ticket_csv_entity.Ticket
	)
	for _, courseCategory := range courseCategories {
		if ticketsByCourseCategory, ok := ticketsGroup[courseCategory]; ok {
			for _, ticket := range ticketsByCourseCategory {
				if ticket.TicketResult() == types.TicketHit {
					hitTickets = append(hitTickets, ticket)
				}
			}
			mergedTickets = append(mergedTickets, ticketsByCourseCategory...)
		}
	}
	hitCount := len(hitTickets)
	maxPayout := 0
	minPayout := 0
	for _, ticket := range hitTickets {
		if maxPayout < ticket.Payout().Value() {
			maxPayout = ticket.Payout().Value()
		}
		if minPayout == 0 || minPayout > ticket.Payout().Value() {
			minPayout = ticket.Payout().Value()
		}
	}

	ticketsByRaceId := t.ticketConverter.ConvertToRaceIdMap(ctx, mergedTickets, racingNumbers)
	payment, payout := t.getSumAmount(mergedTickets)
	betCount := types.BetCount(len(mergedTickets))
	raceCount := types.RaceCount(len(ticketsByRaceId))
	averagePayout := 0
	if len(hitTickets) > 0 {
		averagePayout = int(float64(payout) / float64(len(hitTickets)))
	}

	return spreadsheet_entity.NewTicketResult(
		raceCount,
		betCount,
		types.HitCount(hitCount),
		payment,
		payout,
		types.Payout(averagePayout),
		types.Payout(maxPayout),
		types.Payout(minPayout),
	)
}

func (t *ticketAggregator) DistanceCategoryResult(ctx context.Context, tickets []*ticket_csv_entity.Ticket, races []*data_cache_entity.Race, racingNumbers []*data_cache_entity.RacingNumber, distanceCategories ...types.DistanceCategory) *spreadsheet_entity.TicketResult {
	ticketsGroup := t.ticketConverter.ConvertToDistanceCategoryMap(ctx, tickets, racingNumbers, races)
	var (
		mergedTickets []*ticket_csv_entity.Ticket
		hitTickets    []*ticket_csv_entity.Ticket
	)
	for _, distanceCategory := range distanceCategories {
		if ticketsByDistanceCategory, ok := ticketsGroup[distanceCategory]; ok {
			for _, ticket := range ticketsByDistanceCategory {
				if ticket.TicketResult() == types.TicketHit {
					hitTickets = append(hitTickets, ticket)
				}
			}
			mergedTickets = append(mergedTickets, ticketsByDistanceCategory...)
		}
	}
	hitCount := len(hitTickets)
	maxPayout := 0
	minPayout := 0
	for _, ticket := range hitTickets {
		if maxPayout < ticket.Payout().Value() {
			maxPayout = ticket.Payout().Value()
		}
		if minPayout == 0 || minPayout > ticket.Payout().Value() {
			minPayout = ticket.Payout().Value()
		}
	}

	ticketsByRaceId := t.ticketConverter.ConvertToRaceIdMap(ctx, mergedTickets, racingNumbers)
	payment, payout := t.getSumAmount(mergedTickets)
	betCount := types.BetCount(len(mergedTickets))
	raceCount := types.RaceCount(len(ticketsByRaceId))
	averagePayout := 0
	if len(hitTickets) > 0 {
		averagePayout = int(float64(payout) / float64(len(hitTickets)))
	}

	return spreadsheet_entity.NewTicketResult(
		raceCount,
		betCount,
		types.HitCount(hitCount),
		payment,
		payout,
		types.Payout(averagePayout),
		types.Payout(maxPayout),
		types.Payout(minPayout),
	)
}

func (t *ticketAggregator) RaceCourseResult(ctx context.Context, tickets []*ticket_csv_entity.Ticket, races []*data_cache_entity.Race, racingNumbers []*data_cache_entity.RacingNumber, raceCourses ...types.RaceCourse) *spreadsheet_entity.TicketResult {
	ticketsGroup := t.ticketConverter.ConvertToRaceCourseMap(ctx, tickets, racingNumbers, races)
	var (
		mergedTickets []*ticket_csv_entity.Ticket
		hitTickets    []*ticket_csv_entity.Ticket
	)
	for _, raceCourse := range raceCourses {
		if ticketsByRaceCourse, ok := ticketsGroup[raceCourse]; ok {
			for _, ticket := range ticketsByRaceCourse {
				if ticket.TicketResult() == types.TicketHit {
					hitTickets = append(hitTickets, ticket)
				}
			}
			mergedTickets = append(mergedTickets, ticketsByRaceCourse...)
		}
	}
	hitCount := len(hitTickets)
	maxPayout := 0
	minPayout := 0
	for _, ticket := range hitTickets {
		if maxPayout < ticket.Payout().Value() {
			maxPayout = ticket.Payout().Value()
		}
		if minPayout == 0 || minPayout > ticket.Payout().Value() {
			minPayout = ticket.Payout().Value()
		}
	}

	ticketsByRaceId := t.ticketConverter.ConvertToRaceIdMap(ctx, mergedTickets, racingNumbers)
	payment, payout := t.getSumAmount(mergedTickets)
	betCount := types.BetCount(len(mergedTickets))
	raceCount := types.RaceCount(len(ticketsByRaceId))
	averagePayout := 0
	if len(hitTickets) > 0 {
		averagePayout = int(float64(payout) / float64(len(hitTickets)))
	}

	return spreadsheet_entity.NewTicketResult(
		raceCount,
		betCount,
		types.HitCount(hitCount),
		payment,
		payout,
		types.Payout(averagePayout),
		types.Payout(maxPayout),
		types.Payout(minPayout),
	)
}

func (t *ticketAggregator) getSumAmount(tickets []*ticket_csv_entity.Ticket) (types.Payment, types.Payout) {
	var (
		sumPayment int
		sumPayout  int
	)
	for _, ticket := range tickets {
		sumPayment += ticket.Payment().Value()
		sumPayout += ticket.Payout().Value()
	}

	return types.Payment(sumPayment), types.Payout(sumPayout)
}
