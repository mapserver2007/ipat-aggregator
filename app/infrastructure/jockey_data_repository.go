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

type jockeyDataRepository struct {
}

func NewJockeyDataRepository() repository.JockeyDataRepository {
	return &jockeyDataRepository{}
}

func (j *jockeyDataRepository) Read(ctx context.Context, fileName string) ([]*raw_entity.Jockey, []int, error) {
	rootPath, err := os.Getwd()
	if err != nil {
		return nil, nil, err
	}

	path, err := filepath.Abs(fmt.Sprintf("%s/cache/%s", rootPath, fileName))
	if err != nil {
		return nil, nil, err
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, err
	}

	var jockeyInfo *raw_entity.JockeyInfo
	if err := json.Unmarshal(bytes, &jockeyInfo); err != nil {
		return nil, nil, err
	}

	return jockeyInfo.Jockeys, jockeyInfo.ExcludeJockeyIds, nil
}

func (j *jockeyDataRepository) Write(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (j *jockeyDataRepository) Fetch(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
