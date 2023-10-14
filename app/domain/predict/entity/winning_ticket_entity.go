package entity

import betting_ticket_vo "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/value_object"

type WinningTicketEntity struct {
	bettingTicket betting_ticket_vo.BettingTicket
	betNumber     betting_ticket_vo.BetNumber
	odds          string
	popular       int
	repayment     int
}

func NewWinningTicketEntity(
	bettingTicket betting_ticket_vo.BettingTicket,
	betNumber betting_ticket_vo.BetNumber,
	odds string,
	popular int,
	repayment int,
) *WinningTicketEntity {
	return &WinningTicketEntity{
		bettingTicket: bettingTicket,
		betNumber:     betNumber,
		odds:          odds,
		popular:       popular,
		repayment:     repayment,
	}
}

func (w *WinningTicketEntity) BettingTicket() betting_ticket_vo.BettingTicket {
	return w.bettingTicket
}

func (w *WinningTicketEntity) BetNumber() betting_ticket_vo.BetNumber {
	return w.betNumber
}

func (w *WinningTicketEntity) Odds() string {
	return w.odds
}

func (w *WinningTicketEntity) Popular() int {
	return w.popular
}

func (w *WinningTicketEntity) Repayment() int {
	return w.repayment
}
