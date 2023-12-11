package service

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type RaceConverter interface {
	GetRaceId(ctx context.Context, ticket *ticket_csv_entity.Ticket, racingNumber *netkeiba_entity.RacingNumber) types.RaceId
	ConvertToTicketMap(ctx context.Context, tickets []*ticket_csv_entity.Ticket, rawRacingNumbers []*raw_entity.RacingNumber) map[types.RaceId]*ticket_csv_entity.Ticket
	ConvertToRawRaceMap(ctx context.Context, races []*raw_entity.Race) map[types.RaceId]*raw_entity.Race
	ConvertToRawRacingNumberMap(ctx context.Context, racingNumbers []*raw_entity.RacingNumber) map[types.RacingNumberId]*raw_entity.RacingNumber
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
	rawRacingNumbers []*raw_entity.RacingNumber,
) map[types.RaceId]*ticket_csv_entity.Ticket {
	racingNumberMap := r.ConvertToRawRacingNumberMap(ctx, rawRacingNumbers)
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
				rawRacingNumber.Date,
				rawRacingNumber.Round,
				rawRacingNumber.Day,
				rawRacingNumber.RaceCourseId,
			)
		}

		return r.GetRaceId(ctx, ticket, racingNumber)
	})
}

func (r *raceConverter) ConvertToRawRaceMap(ctx context.Context, races []*raw_entity.Race) map[types.RaceId]*raw_entity.Race {
	return ConvertToMap(races, func(race *raw_entity.Race) types.RaceId {
		return types.RaceId(race.RaceId)
	})
}

func (r *raceConverter) ConvertToRawRacingNumberMap(ctx context.Context, racingNumbers []*raw_entity.RacingNumber) map[types.RacingNumberId]*raw_entity.RacingNumber {
	return ConvertToMap(racingNumbers, func(racingNumber *raw_entity.RacingNumber) types.RacingNumberId {
		return types.NewRacingNumberId(
			types.RaceDate(racingNumber.Date),
			types.RaceCourse(racingNumber.RaceCourseId),
		)
	})
}
