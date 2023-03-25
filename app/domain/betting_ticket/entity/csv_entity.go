package entity

import (
	betting_ticket_vo "github.com/mapserver2007/tools/baken/app/domain/betting_ticket/value_object"
	race_vo "github.com/mapserver2007/tools/baken/app/domain/race/value_object"
)

type CsvEntity struct {
	raceDate      race_vo.RaceDate
	entryNo       int
	raceCourse    race_vo.RaceCourse
	raceNo        int
	bettingTicket betting_ticket_vo.BettingTicket
	bettingResult betting_ticket_vo.BettingResult
	winning       bool
	betNumber     betting_ticket_vo.BetNumber
	payment       int
	repayment     int
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
		raceDate:      raceDate,
		entryNo:       entryNo,
		raceCourse:    raceCourse,
		raceNo:        raceNo,
		bettingTicket: bettingTicket,
		bettingResult: bettingResult,
		winning:       winning,
		betNumber:     betNumber,
		payment:       payment,
		repayment:     repayment,
	}
}

func (c *CsvEntity) RaceDate() race_vo.RaceDate {
	return c.raceDate
}

func (c *CsvEntity) EntryNo() int {
	return c.entryNo
}

func (c *CsvEntity) RaceCourse() race_vo.RaceCourse {
	return c.raceCourse
}

func (c *CsvEntity) RaceNo() int {
	return c.raceNo
}

func (c *CsvEntity) BettingTicket() betting_ticket_vo.BettingTicket {
	return c.bettingTicket
}

func (c *CsvEntity) BettingResult() betting_ticket_vo.BettingResult {
	return c.bettingResult
}

func (c *CsvEntity) Winning() bool {
	return c.winning
}

func (c *CsvEntity) BetNumber() betting_ticket_vo.BetNumber {
	return c.betNumber
}

func (c *CsvEntity) Payment() int {
	return c.payment
}

func (c *CsvEntity) Repayment() int {
	return c.repayment
}
