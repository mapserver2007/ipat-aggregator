package entity

import (
	betting_ticket_vo "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/value_object"
	race_vo "github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type CsvEntity struct {
	raceDate      race_vo.RaceDate
	entryNo       int
	raceCourse    race_vo.RaceCourse
	raceNo        int
	bettingTicket betting_ticket_vo.BettingTicket
	bettingResult betting_ticket_vo.BettingResult
	betNumber     betting_ticket_vo.BetNumber
	payment       int
	payout        int
}

func NewCsvEntity(
	raceDate race_vo.RaceDate,
	entryNo int,
	raceCourse race_vo.RaceCourse,
	raceNo int,
	bettingTicket betting_ticket_vo.BettingTicket,
	bettingResult betting_ticket_vo.BettingResult,
	betNumber betting_ticket_vo.BetNumber,
	payment int,
	payout int,
) *CsvEntity {
	return &CsvEntity{
		raceDate:      raceDate,
		entryNo:       entryNo,
		raceCourse:    raceCourse,
		raceNo:        raceNo,
		bettingTicket: bettingTicket,
		bettingResult: bettingResult,
		betNumber:     betNumber,
		payment:       payment,
		payout:        payout,
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

func (c *CsvEntity) BetNumber() betting_ticket_vo.BetNumber {
	return c.betNumber
}

func (c *CsvEntity) Payment() types.Payment {
	return types.Payment(c.payment)
}

func (c *CsvEntity) Payout() types.Payout {
	return types.Payout(c.payout)
}
