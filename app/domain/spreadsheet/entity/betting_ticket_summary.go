package entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/value_object"
)

func NewBettingTicketSummary(bettingTicketRates map[value_object.BettingTicket]ResultRate) BettingTicketSummary {
	return BettingTicketSummary{
		BettingTicketRates: bettingTicketRates,
	}
}
