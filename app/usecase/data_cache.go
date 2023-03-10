package usecase

import (
	"context"
	"fmt"
	betting_ticket_entity "github.com/mapserver2007/tools/baken/app/domain/betting_ticket/entity"
	race_entity "github.com/mapserver2007/tools/baken/app/domain/race/entity"
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
	racingNumberInfo, err := d.raceDB.ReadRacingNumber(ctx, racingNumberFileName)
	if err != nil {
		return nil, nil, err
	}

	raceInfo, err := d.raceDB.ReadRaceResult(ctx, raceResultFileName)
	if err != nil {
		return nil, nil, err
	}

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
	racingNumberInfo, err := d.raceDB.ReadRacingNumber(ctx, racingNumberFileName)
	racingNumbers := racingNumberInfo.RacingNumbers

	err = d.raceDB.UpdateRaceResult(ctx, raceResultFileName, racingNumbers, entities)
	if err != nil {
		return fmt.Errorf("update race_result.json failed: %w", err)
	}
	log.Println(ctx, "update race_result.json done!")

	return nil
}
