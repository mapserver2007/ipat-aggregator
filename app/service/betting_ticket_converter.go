package service

import (
	"fmt"
	betting_ticket_entity "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/entity"
	betting_ticket_vo "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/value_object"
	race_entity "github.com/mapserver2007/ipat-aggregator/app/domain/race/entity"
	race_vo "github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
	"strconv"
)

type BettingTicketConverter struct{}

func NewBettingTicketConverter() BettingTicketConverter {
	return BettingTicketConverter{}
}

func (b *BettingTicketConverter) ConvertToBettingTicketRecordsMap(records []*betting_ticket_entity.CsvEntity) map[betting_ticket_vo.BettingTicket][]*betting_ticket_entity.CsvEntity {
	return ConvertToSliceMap(records, func(record *betting_ticket_entity.CsvEntity) betting_ticket_vo.BettingTicket {
		return record.BettingTicket()
	})
}

func (b *BettingTicketConverter) ConvertToMonthRecordsMap(records []*betting_ticket_entity.CsvEntity) map[int][]*betting_ticket_entity.CsvEntity {
	return ConvertToSliceMap(records, func(record *betting_ticket_entity.CsvEntity) int {
		key, _ := strconv.Atoi(fmt.Sprintf("%d%02d", record.RaceDate().Year(), record.RaceDate().Month()))
		return key
	})
}

func (b *BettingTicketConverter) ConvertToYearRecordsMap(records []*betting_ticket_entity.CsvEntity) map[int][]*betting_ticket_entity.CsvEntity {
	return ConvertToSliceMap(records, func(record *betting_ticket_entity.CsvEntity) int {
		return record.RaceDate().Year()
	})
}

func (b *BettingTicketConverter) ConvertToRaceClassRecordsMap(records []*betting_ticket_entity.CsvEntity, raceMap map[race_vo.RacingNumberId]*race_entity.Race) map[race_vo.GradeClass][]*betting_ticket_entity.CsvEntity {
	return ConvertToSliceMap(records, func(record *betting_ticket_entity.CsvEntity) race_vo.GradeClass {
		key := race_vo.NewRacingNumberId(record.RaceDate(), record.RaceCourse())
		if race, ok := raceMap[key]; ok {
			return race_vo.GradeClass(race.Class())
		}
		return race_vo.NonGrade
	})
}

func (b *BettingTicketConverter) ConvertToBettingTicketMap(bettingTicketDetails []*betting_ticket_entity.BettingTicketDetail) map[betting_ticket_vo.BettingTicket][]*betting_ticket_entity.BettingTicketDetail {
	return ConvertToSliceMap(bettingTicketDetails, func(bettingTicketDetail *betting_ticket_entity.BettingTicketDetail) betting_ticket_vo.BettingTicket {
		return bettingTicketDetail.BettingTicket()
	})
}

func (b *BettingTicketConverter) ConvertToPayoutResultMap(payoutResults []*race_entity.PayoutResult) map[betting_ticket_vo.BettingTicket]*race_entity.PayoutResult {
	return ConvertToMap(payoutResults, func(payoutResult *race_entity.PayoutResult) betting_ticket_vo.BettingTicket {
		return betting_ticket_vo.BettingTicket(payoutResult.TicketType())
	})
}
