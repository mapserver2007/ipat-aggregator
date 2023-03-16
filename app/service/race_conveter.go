package service

import (
	"fmt"
	betting_ticket_entity "github.com/mapserver2007/tools/baken/app/domain/betting_ticket/entity"
	raw_race_entity "github.com/mapserver2007/tools/baken/app/domain/race/raw_entity"
	race_vo "github.com/mapserver2007/tools/baken/app/domain/race/value_object"
)

const (
	raceListUrlForJRA       = "https://race.netkeiba.com/top/race_list_sub.html?kaisai_date=%d"
	raceResultUrlForJRA     = "https://race.netkeiba.com/race/result.html?race_id=%s&organizer=%d"
	raceResultUrlForNAR     = "https://nar.netkeiba.com/race/result.html?race_id=%s&organizer=%d"
	raceResultUrlForOversea = "https://race.netkeiba.com/race/result.html?race_id=%s&organizer=%d"
)

type RaceConverter struct{}

type RacingNumberRequestParam struct {
	url string
}

type RaceRequestParam struct {
	url    string
	raceId *race_vo.RaceId
	record *betting_ticket_entity.CsvEntity
}

func NewRaceConverter() RaceConverter {
	return RaceConverter{}
}

func (r *RaceConverter) GetRacingNumberRequestParams(
	racingNumbers []*raw_race_entity.RacingNumber,
	records []*betting_ticket_entity.CsvEntity,
) error {
	racingNumberMap := convertToRacingNumberMap(racingNumbers)
	racingNumberRequestParam := make([]*RacingNumberRequestParam, 0)

	for _, record := range records {
		// JRA以外は日付からレース番号の特定が可能のため処理しない
		if record.RaceCourse.Organizer() != race_vo.JRA {
			continue
		}

		racingNumberId := race_vo.NewRacingNumberId(
			record.RaceDate,
			record.RaceCourse,
		)

		if _, ok := racingNumberMap[racingNumberId]; ok {
			continue
		}

		racingNumberRequestParam = append(racingNumberRequestParam, createRacingNumberRequestParam(
			convertToRaceListUrl(racingNumberId.Date()),
		))

		//url := fmt.Sprintf(raceListUrlForJRA, int(entity.RaceDate))

	}

	return nil
}

func (r *RaceConverter) GetRaceRequestParams(
	races []*raw_race_entity.Race,
	racingNumbers []*raw_race_entity.RacingNumber,
	records []*betting_ticket_entity.CsvEntity,
) ([]*RaceRequestParam, error) {
	raceMap := convertToRaceMap(races)
	racingNumberMap := convertToRacingNumberMap(racingNumbers)
	raceRequestParams := make([]*RaceRequestParam, 0)

	for _, record := range records {
		racingNumberId := race_vo.NewRacingNumberId(
			record.RaceDate,
			record.RaceCourse,
		)

		racingNumber, ok := racingNumberMap[racingNumberId]
		if !ok {
			return nil, fmt.Errorf("unknown racingNumberId: %s", string(racingNumberId))
		}

		raceId, err := r.GetRaceId(record, racingNumber)
		if err != nil {
			return nil, err
		}
		if _, ok := raceMap[string(*raceId)]; ok {
			continue
		}
		raceRequestParams = append(raceRequestParams, createRaceRequestParam(
			convertToRaceResultUrl(*raceId, record.RaceCourse.Organizer()),
			raceId,
			record,
		))
	}

	return raceRequestParams, nil
}

func (r *RaceConverter) GetRaceId(
	record *betting_ticket_entity.CsvEntity,
	racingNumber *raw_race_entity.RacingNumber,
) (*race_vo.RaceId, error) {
	//racingNumberMap := convertToRacingNumberMap(racingNumbers)
	var raceId race_vo.RaceId
	organizer := record.RaceCourse.Organizer()

	switch organizer {
	case race_vo.JRA:
		//key := fmt.Sprintf("%d_%d", record.RaceDate, record.RaceCourse.Value())
		//racingNumber, ok := racingNumberMap[key]
		//if !ok {
		//	return nil, fmt.Errorf("undefined key: %s", key)
		//}
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

func createRacingNumberRequestParam(
	url string,
) *RacingNumberRequestParam {
	return &RacingNumberRequestParam{url: url}
}

func createRaceRequestParam(
	url string,
	raceId *race_vo.RaceId,
	record *betting_ticket_entity.CsvEntity,
) *RaceRequestParam {
	return &RaceRequestParam{
		url:    url,
		raceId: raceId,
		record: record,
	}
}

func (r *RaceRequestParam) Url() string {
	return r.url
}

func (r *RaceRequestParam) RaceId() *race_vo.RaceId {
	return r.raceId
}

func (r *RaceRequestParam) Record() *betting_ticket_entity.CsvEntity {
	return r.record
}

func convertToRaceListUrl(raceDate race_vo.RaceDate) string {
	return fmt.Sprintf(raceListUrlForJRA, int(raceDate))
}

func convertToRaceResultUrl(raceId race_vo.RaceId, organizer race_vo.Organizer) string {
	var url string
	switch organizer {
	case race_vo.JRA:
		url = fmt.Sprintf(raceResultUrlForJRA, raceId, organizer)
	case race_vo.NAR:
		url = fmt.Sprintf(raceResultUrlForNAR, raceId, organizer)
	case race_vo.OverseaOrganizer:
		url = fmt.Sprintf(raceResultUrlForOversea, raceId, organizer)
	}

	return url
}

func convertToRaceMap(races []*raw_race_entity.Race) map[string]*raw_race_entity.Race {
	raceMap := map[string]*raw_race_entity.Race{}
	for _, race := range races {
		raceMap[race.RaceId] = race
	}
	return raceMap
}

func convertToRacingNumberMap(racingNumbers []*raw_race_entity.RacingNumber) map[race_vo.RacingNumberId]*raw_race_entity.RacingNumber {
	raceNumberMap := map[race_vo.RacingNumberId]*raw_race_entity.RacingNumber{}
	for _, racingNumber := range racingNumbers {
		racingNumberId := race_vo.NewRacingNumberId(
			race_vo.RaceDate(racingNumber.Date),
			race_vo.RaceCourse(racingNumber.RaceCourseId),
		)
		raceNumberMap[racingNumberId] = racingNumber
	}
	return raceNumberMap
}
