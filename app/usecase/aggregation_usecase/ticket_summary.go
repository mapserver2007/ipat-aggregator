package aggregation_usecase

import (
	"context"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/aggregation_service"
)

type TicketSummary interface {
	Execute(ctx context.Context, input *TicketSummaryInput) error
}

type TicketSummaryInput struct {
	Tickets []*ticket_csv_entity.RaceTicket
}

type ticketSummary struct {
	ticketSummaryService aggregation_service.TicketSummary
}

func NewTicketSummary(
	ticketSummaryService aggregation_service.TicketSummary,
) TicketSummary {
	return &ticketSummary{
		ticketSummaryService: ticketSummaryService,
	}
}

func (m *ticketSummary) Execute(ctx context.Context, input *TicketSummaryInput) error {
	ticketSummaryMap := m.ticketSummaryService.Create(ctx, input.Tickets)
	err := m.ticketSummaryService.Write(ctx, ticketSummaryMap)
	if err != nil {
		return err
	}

	return nil
}
