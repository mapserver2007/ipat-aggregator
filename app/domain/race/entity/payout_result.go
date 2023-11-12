package entity

import betting_ticket_vo "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/value_object"

type PayoutResult struct {
	ticketType int
	number     string
	odds       string
	popular    int
}

func NewPayoutResult(
	ticketType int,
	number string,
	odds string,
	popular int,
) *PayoutResult {
	return &PayoutResult{
		ticketType: ticketType,
		number:     number,
		odds:       odds,
		popular:    popular,
	}
}

func (p *PayoutResult) TicketType() betting_ticket_vo.BettingTicket {
	return betting_ticket_vo.BettingTicket(p.ticketType)
}

func (p *PayoutResult) Number() betting_ticket_vo.BetNumber {
	return betting_ticket_vo.BetNumber(p.number)
}

func (p *PayoutResult) Odds() string {
	return p.odds
}

func (p *PayoutResult) Popular() int {
	return p.popular
}
