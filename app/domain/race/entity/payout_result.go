package entity

import betting_ticket_vo "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/value_object"

type PayoutResult struct {
	ticketType betting_ticket_vo.BettingTicket
	number     betting_ticket_vo.BetNumber
	odds       string
	popular    int
}

func NewPayoutResult(
	rawTicketType int,
	rawNumber string,
	odds string,
	popular int,
) *PayoutResult {
	var ticketType betting_ticket_vo.BettingTicket
	if rawTicketType > 0 && rawTicketType <= 15 {
		ticketType = betting_ticket_vo.BettingTicket(rawTicketType)
	} else {
		ticketType = betting_ticket_vo.UnknownTicket
	}

	return &PayoutResult{
		ticketType: ticketType,
		number:     betting_ticket_vo.BetNumber(rawNumber),
		odds:       odds,
		popular:    popular,
	}
}

func (p *PayoutResult) TicketType() betting_ticket_vo.BettingTicket {
	return p.ticketType
}

func (p *PayoutResult) Number() betting_ticket_vo.BetNumber {
	return p.number
}

func (p *PayoutResult) Odds() string {
	return p.odds
}

func (p *PayoutResult) Popular() int {
	return p.popular
}
