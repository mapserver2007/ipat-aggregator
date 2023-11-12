package entity

import (
	betting_ticket_vo "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/value_object"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type BettingTicketDetail struct {
	bettingTicket betting_ticket_vo.BettingTicket
	bettingResult betting_ticket_vo.BettingResult
	betNumber     betting_ticket_vo.BetNumber
	payment       types.Payment
	payout        types.Payout
}

type PredictionForHorse struct {
	first  string
	second string
}

func NewBettingTicketDetail(
	bettingTicket betting_ticket_vo.BettingTicket,
	bettingResult betting_ticket_vo.BettingResult,
	betNumber betting_ticket_vo.BetNumber,
	payment types.Payment,
	payout types.Payout,
) *BettingTicketDetail {
	return &BettingTicketDetail{
		bettingTicket: bettingTicket,
		bettingResult: bettingResult,
		betNumber:     betNumber,
		payment:       payment,
		payout:        payout,
	}
}

func (b *BettingTicketDetail) BettingTicket() betting_ticket_vo.BettingTicket {
	return b.bettingTicket
}

func (b *BettingTicketDetail) BetNumber() betting_ticket_vo.BetNumber {
	return b.betNumber
}

func (b *BettingTicketDetail) Payment() types.Payment {
	return b.payment
}

func (b *BettingTicketDetail) Payout() types.Payout {
	return b.payout
}

func (b *BettingTicketDetail) BettingResult() betting_ticket_vo.BettingResult {
	return b.bettingResult
}

func NewPredictionForHorse(first, second string) *PredictionForHorse {
	return &PredictionForHorse{
		first:  first,
		second: second,
	}
}

func (p *PredictionForHorse) First() string {
	return p.first
}

func (p *PredictionForHorse) Second() string {
	return p.second
}
