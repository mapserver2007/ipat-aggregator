package usecase

import (
	"context"
	"fmt"
	betting_ticket_entity "github.com/mapserver2007/tools/baken/app/domain/betting_ticket/entity"
	race_entity "github.com/mapserver2007/tools/baken/app/domain/race/entity"
	race_vo "github.com/mapserver2007/tools/baken/app/domain/race/value_object"
	"github.com/mapserver2007/tools/baken/app/repository"
)

type FetchRaceInfo struct {
	raceClient repository.RaceClient
}

func NewFetchRaceInfo(
	raceClient repository.RaceClient,
) *FetchRaceInfo {
	return &FetchRaceInfo{
		raceClient: raceClient,
	}
}

func (f *FetchRaceInfo) FetchRaceInfo(ctx context.Context, racingNumbers []*race_entity.RacingNumber, records []*betting_ticket_entity.CsvEntity) error {
	raceNumberMap := map[string]*race_entity.RacingNumber{}
	for _, racingNumber := range racingNumbers {
		key := fmt.Sprintf("%d_%d", racingNumber.Date(), racingNumber.RaceCourseId())
		raceNumberMap[key] = racingNumber
	}

	for _, record := range records {
		rawRaceId, err := getRaceId(record, raceNumberMap)
		if err != nil {
			return err
		}
		
	}

}

func getRaceId(record *betting_ticket_entity.CsvEntity, raceNumberMap map[string]*race_entity.RacingNumber) (string, error) {
	rawRaceId := ""
	organizer := record.RaceCourse.Organizer()

	switch organizer {
	case race_vo.JRA:
		key := fmt.Sprintf("%d_%d", record.RaceDate, record.RaceCourse.Value())
		racingNumber, ok := raceNumberMap[key]
		if !ok {
			return rawRaceId, fmt.Errorf("undefined key: %s", key)
		}
		rawRaceId = fmt.Sprintf("%d%02d%02d%02d%02d", record.RaceDate.Year(), racingNumber.RaceCourseId(), racingNumber.Round(), racingNumber.Day(), record.RaceNo)
	case race_vo.NAR:
		rawRaceId = fmt.Sprintf("%d%02d%02d%02d%02d", record.RaceDate.Year(), record.RaceCourse.Value(), record.RaceDate.Month(), record.RaceDate.Day(), record.RaceNo)
	case race_vo.OverseaOrganizer:
		raceCourseIdForOversea := race_vo.ConvertToOverseaRaceCourseId(record.RaceCourse)
		rawRaceId = fmt.Sprintf("%d%s%02d%02d%02d", record.RaceDate.Year(), raceCourseIdForOversea, record.RaceDate.Month(), record.RaceDate.Day(), record.RaceNo)
	}

	return rawRaceId, nil
}
