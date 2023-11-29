package data_cache_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type PayoutResult struct {
	ticketType types.TicketType
	number     types.BetNumber
	odds       string
	popular    int
}

func NewPayoutResult(
	rawTicketType string,
	rawNumber string,
	odds string,
	popular int,
) *PayoutResult {
	return &PayoutResult{
		ticketType: types.NewTicketType(rawTicketType),
		number:     types.NewBetNumber(rawNumber),
		odds:       odds,
		popular:    popular,
	}
}

func (p *PayoutResult) TicketType() types.TicketType {
	return p.ticketType
}

func (p *PayoutResult) Number() types.BetNumber {
	return p.number
}

func (p *PayoutResult) Odds() string {
	return p.odds
}

func (p *PayoutResult) Popular() int {
	return p.popular
}
