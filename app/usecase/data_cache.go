package usecase

import (
	"context"
	"fmt"
	betting_ticket_entity "github.com/mapserver2007/tools/baken/app/domain/betting_ticket/entity"
	race_entity "github.com/mapserver2007/tools/baken/app/domain/race/entity"
	race_raw_entity "github.com/mapserver2007/tools/baken/app/domain/race/raw_entity"
	"github.com/mapserver2007/tools/baken/app/repository"
	"github.com/mapserver2007/tools/baken/app/service"
	"log"
	"os"
	"path/filepath"
)

const (
	racingNumberFileName = "racing_number.json"
	raceResultFileName   = "race_result.json"
)

type DataCache struct {
	csvReader service.CsvReader
	raceDB    repository.RaceDB
}

func NewDataCache(
	csvReader service.CsvReader,
	raceDB repository.RaceDB,
) DataCache {
	return DataCache{
		csvReader: csvReader,
		raceDB:    raceDB,
	}
}

func (d *DataCache) ReadAndUpdate(ctx context.Context) (
	[]*betting_ticket_entity.CsvEntity,
	*race_entity.RacingNumberInfo,
	*race_entity.RaceInfo,
	error,
) {
	entities, err := d.readCsv(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	err = d.updateCache(ctx, entities)
	if err != nil {
		return nil, nil, nil, err
	}

	raceNumberInfo, raceInfo, err := d.readCache(ctx)
	if err != nil {
		return nil, nil, nil, err
	}

	return entities, raceNumberInfo, raceInfo, nil
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

func (d *DataCache) readCache(ctx context.Context) (*race_entity.RacingNumberInfo, *race_entity.RaceInfo, error) {
	rawRacingNumberInfo, err := d.raceDB.ReadRacingNumber(ctx, racingNumberFileName)
	if err != nil {
		return nil, nil, err
	}

	rawRaceInfo, err := d.raceDB.ReadRaceResult(ctx, raceResultFileName)
	if err != nil {
		return nil, nil, err
	}

	racingNumberInfo := convertToRaceNumberInfo(rawRacingNumberInfo)
	raceInfo := convertToRaceInfo(rawRaceInfo)

	return racingNumberInfo, raceInfo, nil
}

func (d *DataCache) updateCache(ctx context.Context, entities []*betting_ticket_entity.CsvEntity) error {
	log.Println(ctx, "update racing_number.json ...")
	err := d.raceDB.UpdateRacingNumber(ctx, racingNumberFileName, entities)
	if err != nil {
		return fmt.Errorf("update racing_number.json failed: %w", err)
	}
	log.Println(ctx, "update racing_number.json done!")

	log.Println(ctx, "update race_result.json ...")
	rawRaceInfo, _ := d.raceDB.ReadRaceResult(ctx, raceResultFileName)

	rawRacingNumberInfo, err := d.raceDB.ReadRacingNumber(ctx, racingNumberFileName)
	racingNumberInfo := convertToRaceNumberInfo(rawRacingNumberInfo)
	racingNumbers := racingNumberInfo.RacingNumbers()

	err = d.raceDB.UpdateRaceResult(ctx, rawRaceInfo, racingNumbers, entities)
	if err != nil {
		return fmt.Errorf("update race_result.json failed: %w", err)
	}
	log.Println(ctx, "update race_result.json done!")

	return nil
}

func convertToRaceInfo(rawRaceInfo *race_raw_entity.RaceInfo) *race_entity.RaceInfo {
	var races []*race_entity.Race
	for _, rawRace := range rawRaceInfo.Races {
		race := race_entity.NewRace(
			rawRace.RaceId,
			rawRace.RaceDate,
			rawRace.RaceNumber,
			rawRace.RaceCourseId,
			rawRace.RaceName,
			rawRace.Url,
			rawRace.Time,
			rawRace.Entries,
			rawRace.Distance,
			rawRace.Class,
			rawRace.CourseCategory,
			rawRace.TrackCondition,
			convertToRaceResults(rawRace.RaceResults),
			convertToPayoutResults(rawRace.PayoutResults),
		)
		races = append(races, race)
	}

	return race_entity.NewRaceInfo(races)
}

func convertToRaceResults(rawRaceResults []*race_raw_entity.RaceResult) []*race_entity.RaceResult {
	var raceResults []*race_entity.RaceResult
	for _, rawRaceResult := range rawRaceResults {
		raceResult := race_entity.NewRaceResult(
			rawRaceResult.OrderNo,
			rawRaceResult.HorseName,
			rawRaceResult.BracketNumber,
			rawRaceResult.HorseNumber,
			rawRaceResult.Odds,
			rawRaceResult.PopularNumber,
		)
		raceResults = append(raceResults, raceResult)
	}

	return raceResults
}

func convertToPayoutResults(rawPayoutResults []*race_raw_entity.PayoutResult) []*race_entity.PayoutResult {
	var payoutResults []*race_entity.PayoutResult
	for _, rawPayoutResult := range rawPayoutResults {
		payoutResult := race_entity.NewPayoutResult(
			rawPayoutResult.TicketType,
			rawPayoutResult.Numbers,
			rawPayoutResult.Odds,
		)
		payoutResults = append(payoutResults, payoutResult)
	}

	return payoutResults
}

func convertToRaceNumberInfo(rawRacingNumberInfo *race_raw_entity.RacingNumberInfo) *race_entity.RacingNumberInfo {
	var racingNumbers []*race_entity.RacingNumber
	for _, rawRacingNumber := range rawRacingNumberInfo.RacingNumbers {
		racingNumber := race_entity.NewRacingNumber(
			rawRacingNumber.Date,
			rawRacingNumber.Round,
			rawRacingNumber.Day,
			rawRacingNumber.RaceCourseId,
		)
		racingNumbers = append(racingNumbers, racingNumber)
	}

	return race_entity.NewRacingNumberInfo(racingNumbers)
}

func convertToRawRaceInfo() {

}
