package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"os"
	"path/filepath"
)

type raceDataRepository struct {
}

func NewRaceDataRepository() repository.RaceDataRepository {
	return &raceDataRepository{}
}

func (r *raceDataRepository) Read(ctx context.Context, fileName string) ([]*raw_entity.Race, error) {
	rootPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	path, err := filepath.Abs(fmt.Sprintf("%s/cache/%s", rootPath, fileName))
	if err != nil {
		return nil, err
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raceInfo *raw_entity.RaceInfo
	if err := json.Unmarshal(bytes, &raceInfo); err != nil {
		return nil, err
	}

	return raceInfo.Races, nil
}

func (r *raceDataRepository) Write(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (r *raceDataRepository) Fetch(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
