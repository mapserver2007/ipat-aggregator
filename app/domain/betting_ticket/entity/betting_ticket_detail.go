package entity

import betting_ticket_vo "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/value_object"

type BettingTicketDetail struct {
	bettingTicket betting_ticket_vo.BettingTicket
	betNumber     betting_ticket_vo.BetNumber
	payment       int
	repayment     int
	winning       bool
}

type PredictionForHorse struct {
	first  string
	second string
}

func NewBettingTicketDetail(
	bettingTicket betting_ticket_vo.BettingTicket,
	betNumber betting_ticket_vo.BetNumber,
	payment int,
	repayment int,
	winning bool,
) *BettingTicketDetail {
	return &BettingTicketDetail{
		bettingTicket: bettingTicket,
		betNumber:     betNumber,
		payment:       payment,
		repayment:     repayment,
		winning:       winning,
	}
}

func (b *BettingTicketDetail) BettingTicket() betting_ticket_vo.BettingTicket {
	return b.bettingTicket
}

func (b *BettingTicketDetail) BetNumber() betting_ticket_vo.BetNumber {
	return b.betNumber
}

func (b *BettingTicketDetail) Payment() int {
	return b.payment
}

func (b *BettingTicketDetail) Repayment() int {
	return b.repayment
}

func (b *BettingTicketDetail) Winning() bool {
	return b.winning
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
