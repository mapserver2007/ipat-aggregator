package entity

import (
	"github.com/mapserver2007/tools/baken/app/domain/betting_ticket/value_object"
)

func NewBettingTicketSummary(bettingTicketRates map[value_object.BettingTicket]ResultRate) BettingTicketSummary {
	return BettingTicketSummary{
		BettingTicketRates: bettingTicketRates,
	}
}
