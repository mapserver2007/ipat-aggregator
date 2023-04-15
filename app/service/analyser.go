package service

import (
	"fmt"
	betting_ticket_entity "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/entity"
	betting_ticket_vo "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/value_object"
	race_entity "github.com/mapserver2007/ipat-aggregator/app/domain/race/entity"
	race_vo "github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
)

type Analyser struct {
	raceConverter RaceConverter
}

func NewAnalyser(
	raceConverter RaceConverter,
) *Analyser {
	return &Analyser{
		raceConverter: raceConverter,
	}
}

func (a *Analyser) Analyse(
	records []*betting_ticket_entity.CsvEntity,
	racingNumbers []*race_entity.RacingNumber,
	races []*race_entity.Race,
) error {

	raceMap := a.raceConverter.ConvertToRaceMapByRaceId(races)
	racingNumberMap := a.raceConverter.ConvertToRacingNumberMap(racingNumbers)

	// TODO とりあえず券種別にメソッドを割るが、後々リファクタリングするかも
	recordsByWin := a.getRecordsByWin(records)
	for _, record := range recordsByWin {
		racingNumberId := race_vo.NewRacingNumberId(record.RaceDate(), record.RaceCourse())
		racingNumber, ok := racingNumberMap[racingNumberId]
		if !ok && record.RaceCourse().Organizer() == race_vo.JRA {
			panic(fmt.Errorf("unknown racingNumberId: %s", racingNumberId))
		}
		raceId := a.raceConverter.GetRaceId(record, racingNumber)
		if race, ok := raceMap[*raceId]; ok {
			fmt.Println(race)
		}
	}

	return nil
}

func (a *Analyser) getRecordsByWin(records []*betting_ticket_entity.CsvEntity) []*betting_ticket_entity.CsvEntity {
	var recordsByWin []*betting_ticket_entity.CsvEntity
	for _, record := range records {
		if record.BettingTicket() != betting_ticket_vo.Win {
			continue
		}
		recordsByWin = append(recordsByWin, record)
	}

	return recordsByWin
}

func (a *Analyser) getRecordsByPopular(records []*betting_ticket_entity.CsvEntity) []*betting_ticket_entity.CsvEntity {
	var recordsByPopular []*betting_ticket_entity.CsvEntity
	for _, record := range records {
		record.BettingTicket()
	}

	return recordsByPopular
}

//func (a *Analyser) getRecordsMapByRaceId(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, ticketType betting_ticket_vo.BettingTicket) map[race_vo.RaceId][]*betting_ticket_entity.CsvEntity {
//	recordMap := map[race_vo.RaceId][]*betting_ticket_entity.CsvEntity{}
//	//racingNumberMap := a.raceConverter.ConvertToRacingNumberMap(racingNumbers)
//	//for _, record := range records {
//	//	if record.BettingTicket() != ticketType {
//	//		continue
//	//	}
//	//	key := race_vo.NewRacingNumberId(record.RaceDate(), record.RaceCourse())
//	//	racingNumber, _ := racingNumberMap[key]
//	//	raceId := a.raceConverter.GetRaceId(record, racingNumber)
//	//	recordMap[*raceId] = append(recordMap[*raceId], record)
//	//}
//
//	return recordMap
//}
