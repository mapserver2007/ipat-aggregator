package service

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
)

type TicketSummaryService interface {
	CreateSummary(ctx context.Context, tickets []*ticket_csv_entity.Ticket) error
}

type ticketSummaryService struct {
}

func NewTicketSummaryService() TicketSummaryService {
	return &ticketSummaryService{}
}

func (t ticketSummaryService) CreateSummary(
	ctx context.Context,
	tickets []*ticket_csv_entity.Ticket,
) error {
	return nil
}
