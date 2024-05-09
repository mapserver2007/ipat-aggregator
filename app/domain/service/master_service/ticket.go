package master_service

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/config"
)

type Ticket interface {
	Get(ctx context.Context, races []*data_cache_entity.Race) ([]*ticket_csv_entity.RaceTicket, error)
}

type ticketService struct {
	ticketRepository repository.TicketRepository
}

func NewTicket(
	ticketRepository repository.TicketRepository,
) Ticket {
	return &ticketService{
		ticketRepository: ticketRepository,
	}
}

func (t *ticketService) Get(
	ctx context.Context,
	races []*data_cache_entity.Race,
) ([]*ticket_csv_entity.RaceTicket, error) {
	files, err := t.ticketRepository.List(ctx, config.CsvDir)
	if err != nil {
		return nil, err
	}

	raceDateMap := map[types.RaceDate][]*data_cache_entity.Race{}
	for _, race := range races {
		if _, ok := raceDateMap[race.RaceDate()]; !ok {
			raceDateMap[race.RaceDate()] = make([]*data_cache_entity.Race, 0)
		}
		raceDateMap[race.RaceDate()] = append(raceDateMap[race.RaceDate()], race)
	}

	raceTickets := make([]*ticket_csv_entity.RaceTicket, 0)
	for _, file := range files {
		tickets, err := t.ticketRepository.Read(ctx, fmt.Sprintf("%s/%s", config.CsvDir, file))
		if err != nil {
			return nil, err
		}
		for _, ticket := range tickets {
			if ticket.RaceCourse().NAR() {
				// NAR,海外はraceデータをキャッシュしてないのでraceIdを自力で構築する
				raceId := types.NewRaceIdForNAR(
					ticket.RaceDate().Year(),
					ticket.RaceDate().Month(),
					ticket.RaceDate().Day(),
					ticket.RaceCourse().Value(),
					ticket.RaceNo(),
				)
				raceTickets = append(raceTickets, ticket_csv_entity.NewRaceTicket(
					raceId,
					ticket,
				))
			} else if ticket.RaceCourse().Oversea() {
				// NAR,海外はraceデータをキャッシュしてないのでraceIdを自力で構築する
				raceId := types.NewRaceIdForOverseas(
					ticket.RaceDate().Year(),
					ticket.RaceDate().Month(),
					ticket.RaceDate().Day(),
					ticket.RaceCourse().Value(),
					ticket.RaceNo(),
				)
				raceTickets = append(raceTickets, ticket_csv_entity.NewRaceTicket(
					raceId,
					ticket,
				))
			} else if ticket.RaceCourse().JRA() {
				raceDateRaces, ok := raceDateMap[ticket.RaceDate()]
				if !ok {
					continue
				}
				for _, race := range raceDateRaces {
					// racingNumberのように完全に紐付けることはできないので、レースNo、開催場所から特定する
					if race.RaceNumber() == ticket.RaceNo() && race.RaceCourseId() == ticket.RaceCourse() {
						raceTickets = append(raceTickets, ticket_csv_entity.NewRaceTicket(
							race.RaceId(),
							ticket,
						))
					}
				}
			}
		}
	}

	return raceTickets, nil
}
