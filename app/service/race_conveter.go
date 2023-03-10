package service

import (
	"fmt"
	betting_ticket_entity "github.com/mapserver2007/tools/baken/app/domain/betting_ticket/entity"
	race_entity "github.com/mapserver2007/tools/baken/app/domain/race/entity"
	race_vo "github.com/mapserver2007/tools/baken/app/domain/race/value_object"
)

type RaceConverter struct {
	raceNumberMap map[string]*race_entity.RacingNumber
}

func NewRaceConverter(
	racingNumbers []*race_entity.RacingNumber,
) RaceConverter {
	raceNumberMap := map[string]*race_entity.RacingNumber{}
	for _, racingNumber := range racingNumbers {
		key := fmt.Sprintf("%d_%d", racingNumber.Date, racingNumber.RaceCourseId)
		raceNumberMap[key] = racingNumber
	}
	return RaceConverter{
		raceNumberMap: raceNumberMap,
	}
}

func (r *RaceConverter) GetRaceId(record *betting_ticket_entity.CsvEntity) (*race_vo.RaceId, error) {
	var raceId race_vo.RaceId
	organizer := record.RaceCourse.Organizer()

	switch organizer {
	case race_vo.JRA:
		key := fmt.Sprintf("%d_%d", record.RaceDate, record.RaceCourse.Value())
		racingNumber, ok := r.raceNumberMap[key]
		if !ok {
			return nil, fmt.Errorf("undefined key: %s", key)
		}
		rawRaceId := fmt.Sprintf("%d%02d%02d%02d%02d", record.RaceDate.Year(), racingNumber.RaceCourseId, racingNumber.Round, racingNumber.Day, record.RaceNo)
		raceId = race_vo.RaceId(rawRaceId)
	case race_vo.NAR:
		rawRaceId := fmt.Sprintf("%d%02d%02d%02d%02d", record.RaceDate.Year(), record.RaceCourse.Value(), record.RaceDate.Month(), record.RaceDate.Day(), record.RaceNo)
		raceId = race_vo.RaceId(rawRaceId)
	case race_vo.OverseaOrganizer:
		raceCourseIdForOversea := race_vo.ConvertToOverseaRaceCourseId(record.RaceCourse)
		rawRaceId := fmt.Sprintf("%d%s%02d%02d%02d", record.RaceDate.Year(), raceCourseIdForOversea, record.RaceDate.Month(), record.RaceDate.Day(), record.RaceNo)
		raceId = race_vo.RaceId(rawRaceId)
	}

	return &raceId, nil
}