package entity

import betting_ticket_vo "github.com/mapserver2007/tools/baken/app/domain/betting_ticket/value_object"

type WinningTicketEntity struct {
	BettingTicket betting_ticket_vo.BettingTicket
	BetNumber     betting_ticket_vo.BetNumber
	Odds          string
	Repayment     int
}

func NewWinningTicketEntity(
	bettingTicket betting_ticket_vo.BettingTicket,
	betNumber betting_ticket_vo.BetNumber,
	odds string,
	repayment int,
) *WinningTicketEntity {
	return &WinningTicketEntity{
		BettingTicket: bettingTicket,
		BetNumber:     betNumber,
		Odds:          odds,
		Repayment:     repayment,
	}
}
