package aggregation_service

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/summary_service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type TicketSummary interface {
	Create(ctx context.Context, tickets []*ticket_csv_entity.RaceTicket) map[int]*spreadsheet_entity.TicketSummary
	Write(ctx context.Context, ticketSummaryMap map[int]*spreadsheet_entity.TicketSummary) error
}

type ticketSummaryService struct {
	termService           summary_service.Term
	spreadSheetRepository repository.SpreadSheetRepository
	logger                *logrus.Logger
}

func NewTicketSummary(
	termService summary_service.Term,
	spreadSheetRepository repository.SpreadSheetRepository,
	logger *logrus.Logger,
) TicketSummary {
	return &ticketSummaryService{
		termService:           termService,
		spreadSheetRepository: spreadSheetRepository,
		logger:                logger,
	}
}

func (t *ticketSummaryService) Create(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
) map[int]*spreadsheet_entity.TicketSummary {
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

	monthlyResultMap := map[int]*spreadsheet_entity.TicketSummary{}
	for currentMonth, raceTickets := range dateTimeTicketMap {
		var (
			winRaceTickets           []*ticket_csv_entity.RaceTicket
			placeRaceTickets         []*ticket_csv_entity.RaceTicket
			quinellaRaceTickets      []*ticket_csv_entity.RaceTicket
			exactaRaceTickets        []*ticket_csv_entity.RaceTicket
			quinellaPlaceRaceTickets []*ticket_csv_entity.RaceTicket
			trioRaceTickets          []*ticket_csv_entity.RaceTicket
			trifectaRaceTickets      []*ticket_csv_entity.RaceTicket
		)
		for _, raceTicket := range raceTickets {
			switch raceTicket.Ticket().TicketType().OriginTicketType() {
			case types.Win:
				winRaceTickets = append(winRaceTickets, raceTicket)
			case types.Place:
				placeRaceTickets = append(placeRaceTickets, raceTicket)
			case types.Quinella:
				quinellaRaceTickets = append(quinellaRaceTickets, raceTicket)
			case types.Exacta:
				exactaRaceTickets = append(exactaRaceTickets, raceTicket)
			case types.QuinellaPlace:
				quinellaPlaceRaceTickets = append(quinellaPlaceRaceTickets, raceTicket)
			case types.Trio:
				trioRaceTickets = append(trioRaceTickets, raceTicket)
			case types.Trifecta:
				trifectaRaceTickets = append(trifectaRaceTickets, raceTicket)
			default:
				t.logger.Errorf("unknown ticket type in TicketSummary.Create")
			}
		}

		key, _ := strconv.Atoi(currentMonth.Format("200601"))
		nextMonth := currentMonth.AddDate(0, 1, 0)
		monthlyResultMap[key] = spreadsheet_entity.NewTicketSummary(
			t.createTermResult(ctx, winRaceTickets, currentMonth, nextMonth),
			t.createTermResult(ctx, placeRaceTickets, currentMonth, nextMonth),
			t.createTermResult(ctx, quinellaRaceTickets, currentMonth, nextMonth),
			t.createTermResult(ctx, exactaRaceTickets, currentMonth, nextMonth),
			t.createTermResult(ctx, quinellaPlaceRaceTickets, currentMonth, nextMonth),
			t.createTermResult(ctx, trioRaceTickets, currentMonth, nextMonth),
			t.createTermResult(ctx, trifectaRaceTickets, currentMonth, nextMonth),
		)
	}

	return monthlyResultMap
}

func (t *ticketSummaryService) Write(
	ctx context.Context,
	ticketSummaryMap map[int]*spreadsheet_entity.TicketSummary,
) error {
	return t.spreadSheetRepository.WriteTicketSummary(ctx, ticketSummaryMap)
}

func (t *ticketSummaryService) createTermResult(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
	from time.Time,
	to time.Time,
) *spreadsheet_entity.TicketResult {
	output := t.termService.Create(ctx, &summary_service.TermInput{
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
