package service

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"strconv"
	"time"
)

type TicketConverter interface {
	ConvertToTermTickets(ctx context.Context, tickets []*ticket_csv_entity.Ticket, from, to time.Time) []*ticket_csv_entity.Ticket
	ConvertToTicketTypeMap(ctx context.Context, tickets []*ticket_csv_entity.Ticket) map[types.TicketType][]*ticket_csv_entity.Ticket
	ConvertToGradeClassMap(ctx context.Context, tickets []*ticket_csv_entity.Ticket, racingNumbers []*data_cache_entity.RacingNumber, races []*data_cache_entity.Race) map[types.GradeClass][]*ticket_csv_entity.Ticket
	ConvertToCourseCategoryMap(ctx context.Context, tickets []*ticket_csv_entity.Ticket, racingNumbers []*data_cache_entity.RacingNumber, races []*data_cache_entity.Race) map[types.CourseCategory][]*ticket_csv_entity.Ticket
	ConvertToDistanceCategoryMap(ctx context.Context, tickets []*ticket_csv_entity.Ticket, racingNumbers []*data_cache_entity.RacingNumber, races []*data_cache_entity.Race) map[types.DistanceCategory][]*ticket_csv_entity.Ticket
	ConvertToRaceCourseMap(ctx context.Context, tickets []*ticket_csv_entity.Ticket, racingNumbers []*data_cache_entity.RacingNumber, races []*data_cache_entity.Race) map[types.RaceCourse][]*ticket_csv_entity.Ticket
	ConvertToMonthlyMap(ctx context.Context, tickets []*ticket_csv_entity.Ticket) map[int][]*ticket_csv_entity.Ticket
	ConvertToRaceIdMap(ctx context.Context, tickets []*ticket_csv_entity.Ticket, racingNumbers []*data_cache_entity.RacingNumber) map[types.RaceId][]*ticket_csv_entity.Ticket
}

type ticketConverter struct {
	raceConverter RaceConverter
}

func NewTicketConverter(
	raceConverter RaceConverter,
) TicketConverter {
	return &ticketConverter{
		raceConverter: raceConverter,
	}
}

func (t *ticketConverter) ConvertToTermTickets(
	ctx context.Context,
	tickets []*ticket_csv_entity.Ticket,
	from, to time.Time,
) []*ticket_csv_entity.Ticket {
	var timeRangeTickets []*ticket_csv_entity.Ticket
	for _, ticket := range tickets {
		if ticket.RaceDate().Date().Unix() >= from.Unix() && ticket.RaceDate().Date().Unix() < to.Unix() {
			timeRangeTickets = append(timeRangeTickets, ticket)
		}
	}
	return timeRangeTickets
}

func (t *ticketConverter) ConvertToTicketTypeMap(
	ctx context.Context,
	tickets []*ticket_csv_entity.Ticket,
) map[types.TicketType][]*ticket_csv_entity.Ticket {
	return ConvertToSliceMap(tickets, func(ticket *ticket_csv_entity.Ticket) types.TicketType {
		return ticket.TicketType()
	})
}

func (t *ticketConverter) ConvertToGradeClassMap(
	ctx context.Context,
	tickets []*ticket_csv_entity.Ticket,
	racingNumbers []*data_cache_entity.RacingNumber,
	races []*data_cache_entity.Race,
) map[types.GradeClass][]*ticket_csv_entity.Ticket {
	raceMap := ConvertToMap(races, func(race *data_cache_entity.Race) types.RaceId {
		return race.RaceId()
	})
	racingNumberMap := ConvertToMap(racingNumbers, func(racingNumber *data_cache_entity.RacingNumber) types.RacingNumberId {
		return types.NewRacingNumberId(
			racingNumber.RaceDate(),
			racingNumber.RaceCourse(),
		)
	})

	return ConvertToSliceMap(tickets, func(ticket *ticket_csv_entity.Ticket) types.GradeClass {
		racingNumberId := types.NewRacingNumberId(ticket.RaceDate(), ticket.RaceCourse())
		racingNumber, ok := racingNumberMap[racingNumberId]
		if !ok && ticket.RaceCourse().JRA() {
			panic(fmt.Errorf("unknown racingNumberId: %s", racingNumberId))
		}
		raceId := t.raceConverter.GetRaceId(ctx, ticket, racingNumber)
		if race, ok := raceMap[raceId]; ok {
			return race.Class()
		}
		return types.NonGrade
	})
}

func (t *ticketConverter) ConvertToCourseCategoryMap(
	ctx context.Context,
	tickets []*ticket_csv_entity.Ticket,
	racingNumbers []*data_cache_entity.RacingNumber,
	races []*data_cache_entity.Race,
) map[types.CourseCategory][]*ticket_csv_entity.Ticket {
	raceMap := ConvertToMap(races, func(race *data_cache_entity.Race) types.RaceId {
		return race.RaceId()
	})
	racingNumberMap := ConvertToMap(racingNumbers, func(racingNumber *data_cache_entity.RacingNumber) types.RacingNumberId {
		return types.NewRacingNumberId(
			racingNumber.RaceDate(),
			racingNumber.RaceCourse(),
		)
	})

	return ConvertToSliceMap(tickets, func(ticket *ticket_csv_entity.Ticket) types.CourseCategory {
		racingNumberId := types.NewRacingNumberId(ticket.RaceDate(), ticket.RaceCourse())
		racingNumber, ok := racingNumberMap[racingNumberId]
		if !ok && ticket.RaceCourse().JRA() {
			panic(fmt.Errorf("unknown racingNumberId: %s", racingNumberId))
		}
		raceId := t.raceConverter.GetRaceId(ctx, ticket, racingNumber)
		if race, ok := raceMap[raceId]; ok {
			return race.CourseCategory()
		}
		return types.NonCourseCategory
	})
}

func (t *ticketConverter) ConvertToDistanceCategoryMap(
	ctx context.Context,
	tickets []*ticket_csv_entity.Ticket,
	racingNumbers []*data_cache_entity.RacingNumber,
	races []*data_cache_entity.Race,
) map[types.DistanceCategory][]*ticket_csv_entity.Ticket {
	raceMap := ConvertToMap(races, func(race *data_cache_entity.Race) types.RaceId {
		return race.RaceId()
	})
	racingNumberMap := ConvertToMap(racingNumbers, func(racingNumber *data_cache_entity.RacingNumber) types.RacingNumberId {
		return types.NewRacingNumberId(
			racingNumber.RaceDate(),
			racingNumber.RaceCourse(),
		)
	})

	return ConvertToSliceMap(tickets, func(ticket *ticket_csv_entity.Ticket) types.DistanceCategory {
		racingNumberId := types.NewRacingNumberId(ticket.RaceDate(), ticket.RaceCourse())
		racingNumber, ok := racingNumberMap[racingNumberId]
		if !ok && ticket.RaceCourse().JRA() {
			panic(fmt.Errorf("unknown racingNumberId: %s", racingNumberId))
		}
		raceId := t.raceConverter.GetRaceId(ctx, ticket, racingNumber)
		if race, ok := raceMap[raceId]; ok {
			return types.NewDistanceCategory(race.Distance(), race.CourseCategory())
		}
		return types.UndefinedDistanceCategory
	})
}

func (t *ticketConverter) ConvertToRaceCourseMap(
	ctx context.Context,
	tickets []*ticket_csv_entity.Ticket,
	racingNumbers []*data_cache_entity.RacingNumber,
	races []*data_cache_entity.Race,
) map[types.RaceCourse][]*ticket_csv_entity.Ticket {
	raceMap := ConvertToMap(races, func(race *data_cache_entity.Race) types.RaceId {
		return race.RaceId()
	})
	racingNumberMap := ConvertToMap(racingNumbers, func(racingNumber *data_cache_entity.RacingNumber) types.RacingNumberId {
		return types.NewRacingNumberId(
			racingNumber.RaceDate(),
			racingNumber.RaceCourse(),
		)
	})

	return ConvertToSliceMap(tickets, func(ticket *ticket_csv_entity.Ticket) types.RaceCourse {
		racingNumberId := types.NewRacingNumberId(ticket.RaceDate(), ticket.RaceCourse())
		racingNumber, ok := racingNumberMap[racingNumberId]
		if !ok && ticket.RaceCourse().JRA() {
			panic(fmt.Errorf("unknown racingNumberId: %s", racingNumberId))
		}
		raceId := t.raceConverter.GetRaceId(ctx, ticket, racingNumber)
		if race, ok := raceMap[raceId]; ok {
			return race.RaceCourseId()
		}
		return types.UnknownPlace
	})
}

func (t *ticketConverter) ConvertToMonthlyMap(
	ctx context.Context,
	tickets []*ticket_csv_entity.Ticket,
) map[int][]*ticket_csv_entity.Ticket {
	return ConvertToSliceMap(tickets, func(ticket *ticket_csv_entity.Ticket) int {
		key, _ := strconv.Atoi(fmt.Sprintf("%d%02d", ticket.RaceDate().Year(), ticket.RaceDate().Month()))
		return key
	})
}

func (t *ticketConverter) ConvertToRaceIdMap(
	ctx context.Context,
	tickets []*ticket_csv_entity.Ticket,
	racingNumbers []*data_cache_entity.RacingNumber,
) map[types.RaceId][]*ticket_csv_entity.Ticket {
	racingNumberMap := ConvertToMap(racingNumbers, func(racingNumber *data_cache_entity.RacingNumber) types.RacingNumberId {
		return types.NewRacingNumberId(
			racingNumber.RaceDate(),
			racingNumber.RaceCourse(),
		)
	})
	return ConvertToSliceMap(tickets, func(ticket *ticket_csv_entity.Ticket) types.RaceId {
		racingNumberId := types.NewRacingNumberId(ticket.RaceDate(), ticket.RaceCourse())
		racingNumber, ok := racingNumberMap[racingNumberId]
		if !ok && ticket.RaceCourse().JRA() {
			panic(fmt.Errorf("unknown racingNumberId: %s", racingNumberId))
		}
		return t.raceConverter.GetRaceId(ctx, ticket, racingNumber)
	})
}
