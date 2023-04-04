package service

import (
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

	//raceMap := a.raceConverter.ConvertToRaceMapByRaceId(races)
	//recordsMap := a.getRecordsMapByRaceId(records, racingNumbers, betting_ticket_vo.Win)

	// TODO とりあえず券種別にメソッドを割るが、後々リファクタリングするかも

	return nil
}

func (a *Analyser) getPopularAnalyseForWin(recordsMap map[race_vo.RaceId][]*betting_ticket_entity.CsvEntity, raceMap map[race_vo.RaceId]*race_entity.Race) error {
	//popularAnalyseMap := map[int]*analyse_entity.PopularAnalyse{}
	//
	//for raceId, records := range recordsMap {
	//	race, ok := raceMap[raceId]
	//	if !ok {
	//		return fmt.Errorf("unknown raceId in raceMap: %s", raceId)
	//	}
	//	for _, raceResult := range race.RaceResults() {
	//		raceResult.
	//	}
	//}

	return nil

}

func (a *Analyser) getRecordsMapByRaceId(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, ticketType betting_ticket_vo.BettingTicket) map[race_vo.RaceId][]*betting_ticket_entity.CsvEntity {
	recordMap := map[race_vo.RaceId][]*betting_ticket_entity.CsvEntity{}
	racingNumberMap := a.raceConverter.ConvertToRacingNumberMap(racingNumbers)
	for _, record := range records {
		if record.BettingTicket() != ticketType {
			continue
		}
		key := race_vo.NewRacingNumberId(record.RaceDate(), record.RaceCourse())
		racingNumber, _ := racingNumberMap[key]
		raceId := a.raceConverter.GetRaceId(record, racingNumber)
		recordMap[*raceId] = append(recordMap[*raceId], record)
	}

	return recordMap
}
