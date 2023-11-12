package usecase

import (
	"context"
	"fmt"
	betting_ticket_entity "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/entity"
	jockey_entity "github.com/mapserver2007/ipat-aggregator/app/domain/jockey/entity"
	raw_jockey_entity "github.com/mapserver2007/ipat-aggregator/app/domain/jockey/raw_entity"
	jockey_vo "github.com/mapserver2007/ipat-aggregator/app/domain/jockey/value_object"
	race_entity "github.com/mapserver2007/ipat-aggregator/app/domain/race/entity"
	raw_race_entity "github.com/mapserver2007/ipat-aggregator/app/domain/race/raw_entity"
	race_vo "github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
	"github.com/mapserver2007/ipat-aggregator/app/repository"
	"github.com/mapserver2007/ipat-aggregator/app/service"
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
	jockeyUrl               = "https://db.netkeiba.com/jockey/%s/"
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
) *DataCache {
	return &DataCache{
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
	*jockey_entity.JockeyInfo,
	error,
) {
	records, err := d.readCsv(ctx)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	rawRacingNumbers, rawRaces, rawJockeys, excludeRawJockeyIds := d.readCache(ctx)
	racingNumberParams, err := d.getRacingNumberRequestParams(rawRacingNumbers, records)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	var newRawRacingNumbers []*raw_race_entity.RacingNumber
	log.Println(ctx, "updating racing_number.json")
	for _, param := range racingNumberParams {
		time.Sleep(time.Second * 1)
		log.Println(ctx, "fetch from "+param.Url())
		rawRacingNumbersNetkeiba, err := d.raceFetcher.FetchRacingNumbers(ctx, param.Url())
		if err != nil {
			return nil, nil, nil, nil, err
		}
		for _, rawRacingNumber := range rawRacingNumbersNetkeiba {
			newRawRacingNumbers = append(newRawRacingNumbers, d.raceConverter.ConvertFromRawRacingNumberNetkeibaToRawRacingNumberCsv(rawRacingNumber))
		}
	}
	rawRacingNumbers = append(rawRacingNumbers, newRawRacingNumbers...)
	rawRacingNumberInfo := &raw_race_entity.RacingNumberInfo{RacingNumbers: rawRacingNumbers}

	raceParams, err := d.getRaceRequestParams(rawRaces, rawRacingNumbers, records)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	var newRawRaces []*raw_race_entity.Race
	log.Println(ctx, "updating race_result.json")
	for _, param := range raceParams {
		time.Sleep(time.Second * 1)
		log.Println(ctx, "fetch from "+param.Url())
		rawRace, err := d.raceFetcher.FetchRace(ctx, param.Url())
		if err != nil {
			return nil, nil, nil, nil, err
		}
		newRawRaces = append(newRawRaces, d.raceConverter.ConvertFromRawRaceNetkeibaToRawRaceCsv(rawRace, param.RaceId(), param.Record()))
	}

	rawRaces = append(rawRaces, newRawRaces...)
	rawRaceInfo := &raw_race_entity.RaceInfo{Races: rawRaces}

	jockeyParams := d.getJockeyRequestParams(rawJockeys, excludeRawJockeyIds)
	var (
		newRawJockeys       []*raw_jockey_entity.Jockey
		newExcludeJockeyIds []int
	)
	for _, param := range jockeyParams {
		time.Sleep(time.Second * 1)
		log.Println(ctx, "fetch from "+param.Url())
		rawJockey, err := d.raceFetcher.FetchJockey(ctx, param.Url())
		if err != nil {
			log.Printf("GetJockey error: %v", err)
			continue
		}
		if rawJockey.Name() == "" {
			newExcludeJockeyIds = append(newExcludeJockeyIds, rawJockey.Id())
		} else {
			newRawJockey := raw_jockey_entity.Jockey{
				JockeyId:   rawJockey.Id(),
				JockeyName: rawJockey.Name(),
			}
			newRawJockeys = append(newRawJockeys, &newRawJockey)
		}
	}

	rawJockeys = append(rawJockeys, newRawJockeys...)
	excludeRawJockeyIds = append(excludeRawJockeyIds, newExcludeJockeyIds...)

	rawJockeyInfo := &raw_jockey_entity.JockeyInfo{Jockeys: rawJockeys, ExcludeJockeyIds: excludeRawJockeyIds}

	err = d.writeCache(ctx, rawRaceInfo, rawRacingNumberInfo, rawJockeyInfo)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	racingNumbers := d.raceConverter.ConvertFromRawRacingNumbersCsvToRacingNumbers(rawRacingNumbers)
	jockeys := d.raceConverter.ConvertFromRawJockeysToJockeys(rawJockeys)
	races := d.raceConverter.ConvertFromRawRacesCsvToRaces(rawRaces, jockeys)

	return records, race_entity.NewRacingNumberInfo(racingNumbers), race_entity.NewRaceInfo(races), jockey_entity.NewJockeyInfo(jockeys), nil
}

func (d *DataCache) getRacingNumberRequestParams(
	racingNumbers []*raw_race_entity.RacingNumber,
	records []*betting_ticket_entity.CsvEntity,
) ([]*racingNumberRequestParam, error) {
	racingNumberMap := d.raceConverter.ConvertToRawRacingNumberMap(racingNumbers)
	racingNumberDateCache := map[race_vo.RaceDate]*racingNumberRequestParam{}
	for _, record := range records {
		// JRA以外は日付からレース番号の特定が可能のため処理しない
		if record.RaceCourse().Organizer() != race_vo.JRA {
			continue
		}
		racingNumberId := race_vo.NewRacingNumberId(
			record.RaceDate(),
			record.RaceCourse(),
		)
		if _, ok := racingNumberMap[racingNumberId]; ok {
			continue
		}
		if _, ok := racingNumberDateCache[racingNumberId.Date()]; ok {
			continue
		}
		racingNumberDateCache[racingNumberId.Date()] = createRacingNumberRequestParam(
			fmt.Sprintf(raceListUrlForJRA, int(racingNumberId.Date())),
			record,
		)
	}

	racingNumberRequestParams := make([]*racingNumberRequestParam, 0)
	for _, param := range racingNumberDateCache {
		racingNumberRequestParams = append(racingNumberRequestParams, param)
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
	raceIdCache := map[race_vo.RaceId]*raceRequestParam{}

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
		if _, ok := raceIdCache[*raceId]; ok {
			continue
		}

		switch organizer {
		case race_vo.JRA:
			url = fmt.Sprintf(raceResultUrlForJRA, *raceId, organizer)
		case race_vo.NAR:
			url = fmt.Sprintf(raceResultUrlForNAR, *raceId, organizer)
		case race_vo.OverseaOrganizer:
			url = fmt.Sprintf(raceResultUrlForOversea, *raceId, organizer)
		default:
			return nil, fmt.Errorf("undefined organizer: race_date %d, race_no %d", record.RaceDate(), record.RaceNo())
		}

		raceIdCache[*raceId] = createRaceRequestParam(
			url,
			raceId,
			record,
		)
	}

	for _, param := range raceIdCache {
		raceRequestParams = append(raceRequestParams, param)
	}

	return raceRequestParams, nil
}

func (d *DataCache) getJockeyRequestParams(jockeys []*raw_jockey_entity.Jockey, excludeRawJockeyIds []int) []*jockeyRequestParam {
	beginIdForJRA := 422
	endIdForJRA := 2000
	beginIdForNARandOversea := 5000
	endIdForNARandOversea := 5999

	jockeysMap := map[int]bool{}
	for _, jockey := range jockeys {
		jockeysMap[jockey.JockeyId] = true
	}

	excludeJockeyIdsMap := map[int]jockey_vo.JockeyId{}
	for _, rawJockeyId := range excludeRawJockeyIds {
		excludeJockeyIdsMap[rawJockeyId] = jockey_vo.JockeyId(rawJockeyId)
	}

	params := make([]*jockeyRequestParam, 0)
	for i := beginIdForJRA; i <= endIdForJRA; i++ {
		// 除外リストに含まれてたら何もしない
		if _, ok := excludeJockeyIdsMap[i]; ok {
			continue
		}
		// 既に取得済みの場合は何もしない
		if _, ok := jockeysMap[i]; ok {
			continue
		}
		jockeyId := jockey_vo.JockeyId(i)
		params = append(params, createJockeyRequestParam(fmt.Sprintf(jockeyUrl, jockeyId.Format()), jockeyId))
	}
	for i := beginIdForNARandOversea; i <= endIdForNARandOversea; i++ {
		// 除外リストに含まれてたら何もしない
		if _, ok := excludeJockeyIdsMap[i]; ok {
			continue
		}
		// 既に取得済みの場合は何もしない
		if _, ok := jockeysMap[i]; ok {
			continue
		}
		jockeyId := jockey_vo.JockeyId(i)
		params = append(params, createJockeyRequestParam(fmt.Sprintf(jockeyUrl, jockeyId.Format()), jockeyId))
	}

	return params
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

func (d *DataCache) readCache(ctx context.Context) ([]*raw_race_entity.RacingNumber, []*raw_race_entity.Race, []*raw_jockey_entity.Jockey, []int) {
	var (
		racingNumbers    []*raw_race_entity.RacingNumber
		races            []*raw_race_entity.Race
		jockeys          []*raw_jockey_entity.Jockey
		excludeJockeyIds []int
	)
	rawRacingNumberInfo, err := d.raceDB.ReadRacingNumberInfo(ctx)
	if err == nil {
		racingNumbers = rawRacingNumberInfo.RacingNumbers
	}

	rawRaceInfo, err := d.raceDB.ReadRaceInfo(ctx)
	if err == nil {
		races = rawRaceInfo.Races
	}

	rawJockeyInfo, err := d.raceDB.ReadJockeyInfo(ctx)
	if err == nil {
		jockeys = rawJockeyInfo.Jockeys
		excludeJockeyIds = rawJockeyInfo.ExcludeJockeyIds
	}

	return racingNumbers, races, jockeys, excludeJockeyIds
}

func (d *DataCache) writeCache(
	ctx context.Context,
	raceInfo *raw_race_entity.RaceInfo,
	racingNumberInfo *raw_race_entity.RacingNumberInfo,
	jockeyInfo *raw_jockey_entity.JockeyInfo,
) error {
	err := d.raceDB.WriteRaceInfo(ctx, raceInfo)
	if err != nil {
		return err
	}

	err = d.raceDB.WriteRacingNumberInfo(ctx, racingNumberInfo)
	if err != nil {
		return err
	}

	err = d.raceDB.WriteJockeyInfo(ctx, jockeyInfo)

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

type jockeyRequestParam struct {
	url      string
	jockeyId jockey_vo.JockeyId
	valid    bool
}

func createJockeyRequestParam(
	url string,
	jockeyId jockey_vo.JockeyId,
) *jockeyRequestParam {
	return &jockeyRequestParam{
		url:      url,
		jockeyId: jockeyId,
	}
}

func (j *jockeyRequestParam) Url() string {
	return j.url
}

func (j *jockeyRequestParam) JockeyId() jockey_vo.JockeyId {
	return j.jockeyId
}
