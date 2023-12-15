package service

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"strconv"
	"time"
)

type TicketAggregator interface {
	AllPayment(ctx context.Context, tickets []*ticket_csv_entity.Ticket) types.Payment
	AllPayout(ctx context.Context, tickets []*ticket_csv_entity.Ticket) types.Payout
	MonthPayment(ctx context.Context, tickets []*ticket_csv_entity.Ticket) types.Payment
	MonthPayout(ctx context.Context, tickets []*ticket_csv_entity.Ticket) types.Payout
	YearPayment(ctx context.Context, tickets []*ticket_csv_entity.Ticket) types.Payment
	YearPayout(ctx context.Context, tickets []*ticket_csv_entity.Ticket) types.Payout
}

type ticketAggregator struct {
	ticketConverter TicketConverter
}

func NewTicketAggregator(
	ticketConverter TicketConverter,
) TicketAggregator {
	return &ticketAggregator{
		ticketConverter: ticketConverter,
	}
}

func (t *ticketAggregator) AllPayment(ctx context.Context, tickets []*ticket_csv_entity.Ticket) types.Payment {
	payment, _ := t.getSumAmount(tickets)
	return payment
}

func (t *ticketAggregator) AllPayout(ctx context.Context, tickets []*ticket_csv_entity.Ticket) types.Payout {
	_, payout := t.getSumAmount(tickets)
	return payout
}

func (t *ticketAggregator) MonthPayment(ctx context.Context, tickets []*ticket_csv_entity.Ticket) types.Payment {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	key, _ := strconv.Atoi(fmt.Sprintf("%d%02d", year, month))
	recordsGroup := t.ticketConverter.ConvertToMonthTicketsMap(ctx, tickets)

	if recordsForMonth, ok := recordsGroup[key]; ok {
		payment, _ := t.getSumAmount(recordsForMonth)
		return payment
	}

	return types.Payment(0)
}

func (t *ticketAggregator) MonthPayout(ctx context.Context, tickets []*ticket_csv_entity.Ticket) types.Payout {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	key, _ := strconv.Atoi(fmt.Sprintf("%d%02d", year, month))
	ticketsGroup := t.ticketConverter.ConvertToMonthTicketsMap(ctx, tickets)

	if ticketsForMonth, ok := ticketsGroup[key]; ok {
		_, payout := t.getSumAmount(ticketsForMonth)
		return payout
	}

	return types.Payout(0)
}

func (t *ticketAggregator) YearPayment(ctx context.Context, tickets []*ticket_csv_entity.Ticket) types.Payment {
	now := time.Now()
	key := now.Year()
	ticketsGroup := t.ticketConverter.ConvertToYearTicketsMap(ctx, tickets)

	if ticketsForYear, ok := ticketsGroup[key]; ok {
		payment, _ := t.getSumAmount(ticketsForYear)
		return payment
	}

	return types.Payment(0)
}

func (t *ticketAggregator) YearPayout(ctx context.Context, tickets []*ticket_csv_entity.Ticket) types.Payout {
	now := time.Now()
	key := now.Year()
	ticketsGroup := t.ticketConverter.ConvertToYearTicketsMap(ctx, tickets)

	if ticketsForYear, ok := ticketsGroup[key]; ok {
		_, payout := t.getSumAmount(ticketsForYear)
		return payout
	}

	return types.Payout(0)
}

func (t *ticketAggregator) getSumAmount(tickets []*ticket_csv_entity.Ticket) (types.Payment, types.Payout) {
	var (
		sumPayment int
		sumPayout  int
	)
	for _, ticket := range tickets {
		sumPayment += ticket.Payment().Value()
		sumPayout += ticket.Payout().Value()
	}

	return types.Payment(sumPayment), types.Payout(sumPayout)
}
