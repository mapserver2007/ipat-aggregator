package service

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"strconv"
)

type TicketConverter interface {
	ConvertToMonthTicketsMap(ctx context.Context, tickets []*ticket_csv_entity.Ticket) map[int][]*ticket_csv_entity.Ticket
	ConvertToYearTicketsMap(ctx context.Context, tickets []*ticket_csv_entity.Ticket) map[int][]*ticket_csv_entity.Ticket
}

type ticketConverter struct{}

func NewTicketConverter() TicketConverter {
	return &ticketConverter{}
}

func (t *ticketConverter) ConvertToMonthTicketsMap(
	ctx context.Context,
	tickets []*ticket_csv_entity.Ticket,
) map[int][]*ticket_csv_entity.Ticket {
	return ConvertToSliceMap(tickets, func(ticket *ticket_csv_entity.Ticket) int {
		key, _ := strconv.Atoi(fmt.Sprintf("%d%02d", ticket.RaceDate().Year(), ticket.RaceDate().Month()))
		return key
	})
}

func (t *ticketConverter) ConvertToYearTicketsMap(
	ctx context.Context,
	tickets []*ticket_csv_entity.Ticket,
) map[int][]*ticket_csv_entity.Ticket {
	return ConvertToSliceMap(tickets, func(ticket *ticket_csv_entity.Ticket) int {
		return ticket.RaceDate().Year()
	})
}
