package ticket_csv_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type RaceTicket struct {
	raceId types.RaceId
	ticket *Ticket
}

func NewRaceTicket(
	raceId types.RaceId,
	ticket *Ticket,
) *RaceTicket {
	return &RaceTicket{
		raceId: raceId,
		ticket: ticket,
	}
}

func (r *RaceTicket) RaceId() types.RaceId {
	return r.raceId
}

func (r *RaceTicket) Ticket() *Ticket {
	return r.ticket
}
