package usecase

import (
	"context"
	"fmt"
	betting_ticket_entity "github.com/mapserver2007/tools/baken/app/domain/betting_ticket/entity"
	race_entity "github.com/mapserver2007/tools/baken/app/domain/race/entity"
	raw_race_entity "github.com/mapserver2007/tools/baken/app/domain/race/raw_entity"
	race_vo "github.com/mapserver2007/tools/baken/app/domain/race/value_object"
	"github.com/mapserver2007/tools/baken/app/repository"
	"github.com/mapserver2007/tools/baken/app/service"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	raceListUrlForJRA       = "https://race.netkeiba.com/top/race_list_sub.html?kaisai_date=%d"
	raceResultUrlForJRA     = "https://race.netkeiba.com/race/result.html?race_id=%s&organizer=%d"
	raceResultUrlForNAR     = "https://nar.netkeiba.com/race/result.html?race_id=%s&organizer=%d"
	raceResultUrlForOversea = "https://race.netkeiba.com/race/result.html?race_id=%s&organizer=%d"
)

type DataCache struct {
	csvReader     service.CsvReader
	raceDB        repository.RaceDB
	raceFetcher   service.RaceFetcher
	raceConverter service.RaceConverter
}

func NewDataCache(
	csvReader service.CsvReader,
	raceDB repository.RaceDB,
	raceFetcher service.RaceFetcher,
	raceConverter service.RaceConverter,
) DataCache {
	return DataCache{
		csvReader:     csvReader,
		raceDB:        raceDB,
		raceFetcher:   raceFetcher,
		raceConverter: raceConverter,
	}
}

func (d *DataCache) ReadAndUpdate(ctx context.Context) (
	[]*betting_ticket_entity.CsvEntity,
	*race_entity.RacingNumberInfo,
	*race_entity.RaceInfo,
	error,
) {
	records, err := d.readCsv(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	rawRacingNumbers, rawRaces, err := d.readCache(ctx)
	// エラーだった場合はファイルが空だった場合なので無視
	if err != nil {
		log.Println(ctx, "racing_number.json is empty")
	}

	racingNumberParams, err := d.getRacingNumberRequestParams(rawRacingNumbers, records)
	if err != nil {
		return nil, nil, nil, err
	}

	var newRawRacingNumbers []*raw_race_entity.RacingNumber
	log.Println(ctx, "update racing_number.json ...")
	for _, param := range racingNumberParams {
		time.Sleep(time.Second * 1)
		rawRacingNumbersNetkeiba, err := d.raceFetcher.FetchRacingNumbers(ctx, param.Url())
		if err != nil {
			return nil, nil, nil, err
		}
		for _, rawRacingNumber := range rawRacingNumbersNetkeiba {
			newRawRacingNumbers = append(newRawRacingNumbers, d.raceConverter.ConvertFromRawRacingNumberNetkeibaToRawRacingNumberCsv(rawRacingNumber))
		}
	}
	newRawRacingNumberInfo := &raw_race_entity.RacingNumberInfo{RacingNumbers: newRawRacingNumbers}

	raceParams, err := d.getRaceRequestParams(rawRaces, newRawRacingNumbers, records)
	if err != nil {
		return nil, nil, nil, err
	}

	var newRawRaces []*raw_race_entity.Race
	log.Println(ctx, "update race_result.json ...")
	for _, param := range raceParams {
		time.Sleep(time.Second * 1)
		rawRace, err := d.raceFetcher.FetchRace(ctx, param.Url())
		if err != nil {
			return nil, nil, nil, err
		}
		newRawRaces = append(newRawRaces, d.raceConverter.ConvertFromRawRaceNetkeibaToRawRaceCsv(rawRace, param.RaceId(), param.Record()))
	}

	rawRaces = append(rawRaces, newRawRaces...)
	newRawRaceInfo := &raw_race_entity.RaceInfo{Races: rawRaces}

	err = d.writeCache(ctx, newRawRaceInfo, newRawRacingNumberInfo)
	if err != nil {
		return nil, nil, nil, err
	}

	rawRacingNumbers, rawRaces, err = d.readCache(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	racingNumbers := d.raceConverter.ConvertFromRawRacingNumbersCsvToRacingNumbers(rawRacingNumbers)
	races := d.raceConverter.ConvertFromRawRacesCsvToRaces(rawRaces)

	return records, race_entity.NewRacingNumberInfo(racingNumbers), race_entity.NewRaceInfo(races), nil
}

func (d *DataCache) getRacingNumberRequestParams(
	racingNumbers []*raw_race_entity.RacingNumber,
	records []*betting_ticket_entity.CsvEntity,
) ([]*racingNumberRequestParam, error) {
	racingNumberCache := map[race_vo.RacingNumberId]bool{}
	if racingNumbers != nil {
		for _, racingNumber := range racingNumbers {
			racingNumberId := race_vo.NewRacingNumberId(
				race_vo.RaceDate(racingNumber.Date),
				race_vo.RaceCourse(racingNumber.RaceCourseId),
			)
			racingNumberCache[racingNumberId] = true
		}
	}
	racingNumberRequestParams := make([]*racingNumberRequestParam, 0)

	for _, record := range records {
		// JRA以外は日付からレース番号の特定が可能のため処理しない
		if record.RaceCourse().Organizer() != race_vo.JRA {
			continue
		}

		racingNumberId := race_vo.NewRacingNumberId(
			record.RaceDate(),
			record.RaceCourse(),
		)

		if _, ok := racingNumberCache[racingNumberId]; ok {
			continue
		}
		racingNumberCache[racingNumberId] = true

		racingNumberRequestParams = append(racingNumberRequestParams, createRacingNumberRequestParam(
			fmt.Sprintf(raceListUrlForJRA, int(racingNumberId.Date())),
			record,
		))
	}

	return racingNumberRequestParams, nil
}

func (d *DataCache) getRaceRequestParams(
	races []*raw_race_entity.Race,
	racingNumbers []*raw_race_entity.RacingNumber,
	records []*betting_ticket_entity.CsvEntity,
) ([]*raceRequestParam, error) {
	raceMap := d.raceConverter.ConvertToRawRaceMap(races)
	racingNumberMap := d.raceConverter.ConvertToRawRacingNumberMap(racingNumbers)
	raceRequestParams := make([]*raceRequestParam, 0)

	for _, record := range records {
		var (
			url          string
			racingNumber *race_entity.RacingNumber
		)
		organizer := record.RaceCourse().Organizer()
		if organizer == race_vo.JRA {
			racingNumberId := race_vo.NewRacingNumberId(
				record.RaceDate(),
				record.RaceCourse(),
			)

			rawRacingNumber, ok := racingNumberMap[racingNumberId]
			if !ok {
				return nil, fmt.Errorf("unknown racingNumberId: %s", string(racingNumberId))
			}

			racingNumber = race_entity.NewRacingNumber(
				rawRacingNumber.Date,
				rawRacingNumber.Round,
				rawRacingNumber.Day,
				rawRacingNumber.RaceCourseId,
			)
		}

		raceId := d.raceConverter.GetRaceId(record, racingNumber)
		if _, ok := raceMap[string(*raceId)]; ok {
			continue
		}

		switch organizer {
		case race_vo.JRA:
			url = fmt.Sprintf(raceResultUrlForJRA, *raceId, organizer)
		case race_vo.NAR:
			url = fmt.Sprintf(raceResultUrlForNAR, *raceId, organizer)
		case race_vo.OverseaOrganizer:
			url = fmt.Sprintf(raceResultUrlForOversea, *raceId, organizer)
		}

		raceRequestParams = append(raceRequestParams, createRaceRequestParam(
			url,
			raceId,
			record,
		))
	}

	return raceRequestParams, nil
}

func (d *DataCache) readCsv(ctx context.Context) ([]*betting_ticket_entity.CsvEntity, error) {
	rootPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	dirPath, err := filepath.Abs(rootPath + "/csv")
	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	csvReader := service.NewCsvReader()
	var results []*betting_ticket_entity.CsvEntity
	for _, file := range files {
		filePath := fmt.Sprintf("%s/%s", dirPath, file.Name())
		if filepath.Ext(filePath) != ".csv" {
			continue
		}
		csvEntities, err := csvReader.Read(ctx, filePath)
		if err != nil {
			return nil, err
		}
		results = append(results, csvEntities...)
	}

	return results, nil
}

func (d *DataCache) readCache(ctx context.Context) ([]*raw_race_entity.RacingNumber, []*raw_race_entity.Race, error) {
	rawRacingNumberInfo, err := d.raceDB.ReadRacingNumberInfo(ctx)
	if err != nil {
		return nil, nil, err
	}

	rawRaceInfo, err := d.raceDB.ReadRaceInfo(ctx)
	if err != nil {
		return nil, nil, err
	}

	return rawRacingNumberInfo.RacingNumbers, rawRaceInfo.Races, nil
}

func (d *DataCache) writeCache(ctx context.Context, raceInfo *raw_race_entity.RaceInfo, racingNumberInfo *raw_race_entity.RacingNumberInfo) error {
	err := d.raceDB.WriteRaceInfo(ctx, raceInfo)
	if err != nil {
		return err
	}

	err = d.raceDB.WriteRacingNumberInfo(ctx, racingNumberInfo)
	if err != nil {
		return err
	}

	return nil
}

type racingNumberRequestParam struct {
	url    string
	record *betting_ticket_entity.CsvEntity
}

func createRacingNumberRequestParam(
	url string,
	record *betting_ticket_entity.CsvEntity,
) *racingNumberRequestParam {
	return &racingNumberRequestParam{
		url:    url,
		record: record,
	}
}

func (r *racingNumberRequestParam) Url() string {
	return r.url
}

func (r *racingNumberRequestParam) Record() *betting_ticket_entity.CsvEntity {
	return r.record
}

type raceRequestParam struct {
	url    string
	raceId *race_vo.RaceId
	record *betting_ticket_entity.CsvEntity
}

func createRaceRequestParam(
	url string,
	raceId *race_vo.RaceId,
	record *betting_ticket_entity.CsvEntity,
) *raceRequestParam {
	return &raceRequestParam{
		url:    url,
		raceId: raceId,
		record: record,
	}
}

func (r *raceRequestParam) Url() string {
	return r.url
}

func (r *raceRequestParam) RaceId() *race_vo.RaceId {
	return r.raceId
}

func (r *raceRequestParam) Record() *betting_ticket_entity.CsvEntity {
	return r.record
}
