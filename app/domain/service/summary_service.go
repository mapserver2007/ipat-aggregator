package service

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
)

type SummaryService interface {
	CreateShortSummary(ctx context.Context, tickets []*ticket_csv_entity.Ticket) (allShortSummary, monthShortSummary, yearShortSummary *spreadsheet_entity.ShortSummary)
	//CreateTicketSummary(ctx context.Context, tickets []*ticket_csv_entity.Ticket, racingNumbers []*)
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

func (s *summaryService) CreateShortSummary(
	ctx context.Context,
	tickets []*ticket_csv_entity.Ticket,
) (allShortSummary, monthShortSummary, yearShortSummary *spreadsheet_entity.ShortSummary) {
	allPayment := s.ticketAggregator.AllPayment(ctx, tickets)
	allPayout := s.ticketAggregator.AllPayout(ctx, tickets)
	allShortSummary = spreadsheet_entity.NewShortSummary(allPayment, allPayout)

	monthPayment := s.ticketAggregator.MonthPayment(ctx, tickets)
	monthPayout := s.ticketAggregator.MonthPayout(ctx, tickets)
	monthShortSummary = spreadsheet_entity.NewShortSummary(monthPayment, monthPayout)

	yearPayment := s.ticketAggregator.YearPayment(ctx, tickets)
	yearPayout := s.ticketAggregator.YearPayout(ctx, tickets)
	yearShortSummary = spreadsheet_entity.NewShortSummary(yearPayment, yearPayout)

	return
}
