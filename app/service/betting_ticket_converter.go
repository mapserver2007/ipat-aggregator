package service

import (
	"fmt"
	betting_ticket_entity "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/entity"
	betting_ticket_vo "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/value_object"
	race_entity "github.com/mapserver2007/ipat-aggregator/app/domain/race/entity"
	race_vo "github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
	"strconv"
)

type BettingTicketConverter struct {
	raceConverter RaceConverter
}

func NewBettingTicketConverter(
	raceConverter RaceConverter,
) BettingTicketConverter {
	return BettingTicketConverter{
		raceConverter: raceConverter,
	}
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

func (b *BettingTicketConverter) ConvertToRaceIdRecordsMap(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race) map[race_vo.RaceId][]*betting_ticket_entity.CsvEntity {
	raceMap := b.raceConverter.ConvertToRaceMapByRaceId(races)
	racingNumberMap := b.raceConverter.ConvertToRacingNumberMap(racingNumbers)
	return ConvertToSliceMap(records, func(record *betting_ticket_entity.CsvEntity) race_vo.RaceId {
		racingNumberId := race_vo.NewRacingNumberId(record.RaceDate(), record.RaceCourse())
		racingNumber, ok := racingNumberMap[racingNumberId]
		if !ok && record.RaceCourse().Organizer() == race_vo.JRA {
			panic(fmt.Errorf("unknown racingNumberId: %s", racingNumberId))
		}
		raceId := b.raceConverter.GetRaceId(record, racingNumber)
		if race, ok := raceMap[*raceId]; ok {
			return race.RaceId()
		}
		return ""
	})
}

func (b *BettingTicketConverter) ConvertToRaceClassRecordsMap(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race) map[race_vo.GradeClass][]*betting_ticket_entity.CsvEntity {
	raceMap := b.raceConverter.ConvertToRaceMapByRaceId(races)
	racingNumberMap := b.raceConverter.ConvertToRacingNumberMap(racingNumbers)
	return ConvertToSliceMap(records, func(record *betting_ticket_entity.CsvEntity) race_vo.GradeClass {
		racingNumberId := race_vo.NewRacingNumberId(record.RaceDate(), record.RaceCourse())
		racingNumber, ok := racingNumberMap[racingNumberId]
		if !ok && record.RaceCourse().Organizer() == race_vo.JRA {
			panic(fmt.Errorf("unknown racingNumberId: %s", racingNumberId))
		}
		raceId := b.raceConverter.GetRaceId(record, racingNumber)
		if race, ok := raceMap[*raceId]; ok {
			return race.Class()
		}
		return race_vo.NonGrade
	})
}

func (b *BettingTicketConverter) ConvertToDistanceCategoryRecordsMap(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race) map[race_vo.DistanceCategory][]*betting_ticket_entity.CsvEntity {
	raceMap := b.raceConverter.ConvertToRaceMapByRaceId(races)
	racingNumberMap := b.raceConverter.ConvertToRacingNumberMap(racingNumbers)
	return ConvertToSliceMap(records, func(record *betting_ticket_entity.CsvEntity) race_vo.DistanceCategory {
		racingNumberId := race_vo.NewRacingNumberId(record.RaceDate(), record.RaceCourse())
		racingNumber, ok := racingNumberMap[racingNumberId]
		if !ok && record.RaceCourse().Organizer() == race_vo.JRA {
			panic(fmt.Errorf("unknown racingNumberId: %s", racingNumberId))
		}
		raceId := b.raceConverter.GetRaceId(record, racingNumber)
		if race, ok := raceMap[*raceId]; ok {
			courseCategory := race.CourseCategory()
			return race_vo.NewDistanceCategory(race.Distance(), courseCategory)
		}
		return race_vo.UndefinedDistanceCategory
	})
}

func (b *BettingTicketConverter) ConvertToRaceCourseRecordsMap(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race) map[race_vo.RaceCourse][]*betting_ticket_entity.CsvEntity {
	raceMap := b.raceConverter.ConvertToRaceMapByRaceId(races)
	racingNumberMap := b.raceConverter.ConvertToRacingNumberMap(racingNumbers)
	return ConvertToSliceMap(records, func(record *betting_ticket_entity.CsvEntity) race_vo.RaceCourse {
		racingNumberId := race_vo.NewRacingNumberId(record.RaceDate(), record.RaceCourse())
		racingNumber, ok := racingNumberMap[racingNumberId]
		if !ok && record.RaceCourse().Organizer() == race_vo.JRA {
			panic(fmt.Errorf("unknown racingNumberId: %s", racingNumberId))
		}
		raceId := b.raceConverter.GetRaceId(record, racingNumber)
		if race, ok := raceMap[*raceId]; ok {
			return race.RaceCourseId()
		}
		return race_vo.UnknownPlace
	})
}

func (b *BettingTicketConverter) ConvertToBettingTicketMap(bettingTicketDetails []*betting_ticket_entity.BettingTicketDetail) map[betting_ticket_vo.BettingTicket][]*betting_ticket_entity.BettingTicketDetail {
	return ConvertToSliceMap(bettingTicketDetails, func(bettingTicketDetail *betting_ticket_entity.BettingTicketDetail) betting_ticket_vo.BettingTicket {
		return bettingTicketDetail.BettingTicket()
	})
}

func (b *BettingTicketConverter) ConvertToPayoutResultMap(payoutResults []*race_entity.PayoutResult) map[betting_ticket_vo.BettingTicket][]*race_entity.PayoutResult {
	return ConvertToSliceMap(payoutResults, func(payoutResult *race_entity.PayoutResult) betting_ticket_vo.BettingTicket {
		return payoutResult.TicketType()
	})
}
