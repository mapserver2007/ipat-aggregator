package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	raw_race_entity "github.com/mapserver2007/tools/baken/app/domain/race/raw_entity"
	"github.com/mapserver2007/tools/baken/app/repository"
	"os"
	"path/filepath"
)

const (
	racingNumberFileName = "racing_number.json"
	raceResultFileName   = "race_result.json"
)

type RaceDB struct {
	raceClient repository.RaceClient
}

func NewRaceDB(
	raceClient repository.RaceClient,
) repository.RaceDB {
	return &RaceDB{
		raceClient: raceClient,
	}
}

func (r *RaceDB) ReadRaceInfo(ctx context.Context) (*raw_race_entity.RaceInfo, error) {
	bytes, err := r.readFile(raceResultFileName)
	if err != nil {
		return nil, err
	}

	var raceInfo *raw_race_entity.RaceInfo
	if err := json.Unmarshal(bytes, &raceInfo); err != nil {
		return nil, err
	}

	return raceInfo, nil
}

func (r *RaceDB) ReadRacingNumberInfo(ctx context.Context) (*raw_race_entity.RacingNumberInfo, error) {
	bytes, err := r.readFile(racingNumberFileName)
	if err != nil {
		return nil, err
	}

	var racingNumberInfo *raw_race_entity.RacingNumberInfo
	if err := json.Unmarshal(bytes, &racingNumberInfo); err != nil {
		return nil, err
	}

	return racingNumberInfo, nil
}

func (r *RaceDB) WriteRaceInfo(ctx context.Context, raceInfo *raw_race_entity.RaceInfo) error {
	bytes, err := json.Marshal(raceInfo)
	if err != nil {
		return err
	}

	err = r.writeFile(raceResultFileName, bytes)
	if err != nil {
		return fmt.Errorf("update %s failed: %w", raceResultFileName, err)
	}

	return nil
}

// WriteRacingNumberInfo JRAの場合はDateからIDが特定できないので開催場所、日のデータをキャッシュしておいて
// 変換処理をする必要がある。NAR、海外はDateから特定可能
func (r *RaceDB) WriteRacingNumberInfo(ctx context.Context, racingNumberInfo *raw_race_entity.RacingNumberInfo) error {
	bytes, err := json.Marshal(racingNumberInfo)
	if err != nil {
		return err
	}

	err = r.writeFile(racingNumberFileName, bytes)
	if err != nil {
		return fmt.Errorf("update %s failed: %w", racingNumberFileName, err)
	}

	return nil
}

func (r *RaceDB) readFile(fileName string) ([]byte, error) {
	rootPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	filePath, err := filepath.Abs(fmt.Sprintf("%s/cache/%s", rootPath, fileName))
	if err != nil {
		return nil, err
	}

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (r *RaceDB) writeFile(fileName string, data []byte) error {
	rootPath, err := os.Getwd()
	if err != nil {
		return err
	}

	filePath, err := filepath.Abs(fmt.Sprintf("%s/cache/%s", rootPath, fileName))
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
