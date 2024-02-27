package spreadsheet_usecase

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
)

type ticketSummaryUseCase struct {
	ticketSummaryService service.TicketSummaryService
}

func NewTicketSummaryUseCase(
	ticketSummaryService service.TicketSummaryService,
) *ticketSummaryUseCase {
	return &ticketSummaryUseCase{
		ticketSummaryService: ticketSummaryService,
	}
}

func (t *ticketSummaryUseCase) Write(
	ctx context.Context,
	tickets []*ticket_csv_entity.Ticket,
) error {

	t.ticketSummaryService.CreateSummary(ctx, tickets)

	return nil
}
