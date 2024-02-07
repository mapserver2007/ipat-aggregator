package data_cache_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type PayoutResult struct {
	ticketType types.TicketType
	numbers    []types.BetNumber
	odds       []string
	populars   []int
}

func NewPayoutResult(
	rawTicketType int,
	rawNumbers []string,
	odds []string,
	populars []int,
) *PayoutResult {
	numbers := make([]types.BetNumber, 0, len(rawNumbers))
	for _, rawNumber := range rawNumbers {
		numbers = append(numbers, types.NewBetNumber(rawNumber))
	}
	return &PayoutResult{
		ticketType: types.TicketType(rawTicketType),
		numbers:    numbers,
		odds:       odds,
		populars:   populars,
	}
}

func (p *PayoutResult) TicketType() types.TicketType {
	return p.ticketType
}

func (p *PayoutResult) Numbers() []types.BetNumber {
	return p.numbers
}

func (p *PayoutResult) Odds() []string {
	return p.odds
}

func (p *PayoutResult) Populars() []int {
	return p.populars
}
