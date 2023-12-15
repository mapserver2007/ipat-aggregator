package service

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type RaceConverter interface {
	GetRaceId(ctx context.Context, ticket *ticket_csv_entity.Ticket, racingNumber *netkeiba_entity.RacingNumber) types.RaceId
	ConvertToTicketMap(ctx context.Context, tickets []*ticket_csv_entity.Ticket, racingNumbers []*data_cache_entity.RacingNumber) map[types.RaceId]*ticket_csv_entity.Ticket
	ConvertToRawRaceMap(ctx context.Context, races []*data_cache_entity.Race) map[types.RaceId]*data_cache_entity.Race
	ConvertToRawRacingNumberMap(ctx context.Context, racingNumbers []*data_cache_entity.RacingNumber) map[types.RacingNumberId]*data_cache_entity.RacingNumber
}

type raceConverter struct{}

func NewRaceConverter() RaceConverter {
	return &raceConverter{}
}

func (r *raceConverter) GetRaceId(
	ctx context.Context,
	ticket *ticket_csv_entity.Ticket,
	racingNumber *netkeiba_entity.RacingNumber,
) types.RaceId {
	var raceId types.RaceId
	if ticket.RaceCourse().JRA() {
		raceId = types.NewRaceIdForJRA(
			ticket.RaceDate().Year(),
			racingNumber.Day(),
			racingNumber.RaceCourseId(),
			racingNumber.Round(),
			ticket.RaceNo(),
		)
	} else if ticket.RaceCourse().NAR() {
		raceId = types.NewRaceIdForNAR(
			ticket.RaceDate().Year(),
			ticket.RaceDate().Month(),
			ticket.RaceDate().Day(),
			ticket.RaceCourse().Value(),
			ticket.RaceNo(),
		)
	} else if ticket.RaceCourse().Oversea() {
		raceId = types.NewRaceIdForOverseas(
			ticket.RaceDate().Year(),
			ticket.RaceDate().Month(),
			ticket.RaceDate().Day(),
			ticket.RaceCourse().Value(),
			ticket.RaceNo(),
		)
	}
	return raceId
}

func (r *raceConverter) ConvertToTicketMap(
	ctx context.Context,
	tickets []*ticket_csv_entity.Ticket,
	racingNumbers []*data_cache_entity.RacingNumber,
) map[types.RaceId]*ticket_csv_entity.Ticket {
	racingNumberMap := r.ConvertToRawRacingNumberMap(ctx, racingNumbers)
	return ConvertToMap(tickets, func(ticket *ticket_csv_entity.Ticket) types.RaceId {
		var (
			racingNumber *netkeiba_entity.RacingNumber
		)
		if ticket.RaceCourse().JRA() {
			racingNumberId := types.NewRacingNumberId(
				ticket.RaceDate(),
				ticket.RaceCourse(),
			)
			rawRacingNumber, ok := racingNumberMap[racingNumberId]
			if !ok {
				panic(fmt.Sprintf("unknown racingNumberId: %s", string(racingNumberId)))
			}

			racingNumber = netkeiba_entity.NewRacingNumber(
				rawRacingNumber.RaceDate().Value(),
				rawRacingNumber.Round(),
				rawRacingNumber.Day(),
				rawRacingNumber.RaceCourse().Value(),
			)
		}

		return r.GetRaceId(ctx, ticket, racingNumber)
	})
}

func (r *raceConverter) ConvertToRawRaceMap(ctx context.Context, races []*data_cache_entity.Race) map[types.RaceId]*data_cache_entity.Race {
	return ConvertToMap(races, func(race *data_cache_entity.Race) types.RaceId {
		return race.RaceId()
	})
}

func (r *raceConverter) ConvertToRawRacingNumberMap(ctx context.Context, racingNumbers []*data_cache_entity.RacingNumber) map[types.RacingNumberId]*data_cache_entity.RacingNumber {
	return ConvertToMap(racingNumbers, func(racingNumber *data_cache_entity.RacingNumber) types.RacingNumberId {
		return types.NewRacingNumberId(
			racingNumber.RaceDate(),
			racingNumber.RaceCourse(),
		)
	})
}
