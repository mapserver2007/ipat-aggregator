package entity

import betting_ticket_vo "github.com/mapserver2007/tools/baken/app/domain/betting_ticket/value_object"

type RaceRecordEntity struct {
	Record []*RecordDetail
}

// RecordDetail レース単位のレース情報、購入馬券、本命対抗
type RecordDetail struct {
	Race                 Race
	BettingTicketDetails []*BettingTicketDetail
	PredictionForHorse   PredictionForHorse
}

type BettingTicketDetail struct {
	BettingTicket betting_ticket_vo.BettingTicket
	BetNumber     betting_ticket_vo.BetNumber
	Payment       int
	Repayment     int
	Winning       bool
}

type PredictionForHorse struct {
	First  string
	Second string
}

func NewBettingTicketDetail(
	bettingTicket betting_ticket_vo.BettingTicket,
	betNumber betting_ticket_vo.BetNumber,
	payment int,
	repayment int,
	winning bool,
) *BettingTicketDetail {
	return &BettingTicketDetail{
		BettingTicket: bettingTicket,
		BetNumber:     betNumber,
		Payment:       payment,
		Repayment:     repayment,
		Winning:       winning,
	}
}
