package list_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type Ticket struct {
	ticket  *ticket_csv_entity.Ticket
	number  types.BetNumber
	odds    string
	popular int
}

func NewTicket(
	ticket *ticket_csv_entity.Ticket,
	number types.BetNumber,
	odds string,
	popular int,
) *Ticket {
	return &Ticket{
		ticket:  ticket,
		number:  number,
		odds:    odds,
		popular: popular,
	}
}

func (t *Ticket) RaceDate() types.RaceDate {
	return t.ticket.RaceDate()
}

func (t *Ticket) EntryNo() int {
	return t.ticket.EntryNo()
}

func (t *Ticket) RaceCourse() types.RaceCourse {
	return t.ticket.RaceCourse()
}

func (t *Ticket) RaceNo() int {
	return t.ticket.RaceNo()
}

func (t *Ticket) BetNumber() types.BetNumber {
	return t.ticket.BetNumber()
}

func (t *Ticket) TicketType() types.TicketType {
	return t.ticket.TicketType()
}

func (t *Ticket) TicketResult() types.TicketResult {
	return t.ticket.TicketResult()
}

func (t *Ticket) Payment() types.Payment {
	return t.ticket.Payment()
}

func (t *Ticket) Payout() types.Payout {
	return t.ticket.Payout()
}

func (t *Ticket) Number() types.BetNumber {
	return t.number
}

func (t *Ticket) Odds() string {
	return t.odds
}

func (t *Ticket) Popular() int {
	return t.popular
}
