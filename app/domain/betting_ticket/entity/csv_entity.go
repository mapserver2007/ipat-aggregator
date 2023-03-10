package entity

import (
	betting_ticket_vo "github.com/mapserver2007/tools/baken/app/domain/betting_ticket/value_object"
	race_vo "github.com/mapserver2007/tools/baken/app/domain/race/value_object"
)

type CsvEntity struct {
	RaceDate      race_vo.RaceDate
	EntryNo       int
	RaceCourse    race_vo.RaceCourse
	RaceNo        int
	BettingTicket betting_ticket_vo.BettingTicket
	BettingResult betting_ticket_vo.BettingResult
	Winning       bool
	BetNumber     betting_ticket_vo.BetNumber
	Payment       int
	Repayment     int
}

func NewCsvEntity(
	raceDate race_vo.RaceDate,
	entryNo int,
	raceCourse race_vo.RaceCourse,
	raceNo int,
	bettingTicket betting_ticket_vo.BettingTicket,
	bettingResult betting_ticket_vo.BettingResult,
	winning bool,
	betNumber betting_ticket_vo.BetNumber,
	payment int,
	repayment int,
) *CsvEntity {
	return &CsvEntity{
		RaceDate:      raceDate,
		EntryNo:       entryNo,
		RaceCourse:    raceCourse,
		RaceNo:        raceNo,
		BettingTicket: bettingTicket,
		BettingResult: bettingResult,
		Winning:       winning,
		BetNumber:     betNumber,
		Payment:       payment,
		Repayment:     repayment,
	}
}
