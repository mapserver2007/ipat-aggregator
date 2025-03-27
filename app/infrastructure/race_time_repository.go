package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure/file_gateway"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure/gateway"
)

type raceTimeRepository struct {
	netKeibaGateway gateway.NetKeibaGateway
	pathOptimizer   file_gateway.PathOptimizer
}

func NewRaceTimeRepository(
	netKeibaGateway gateway.NetKeibaGateway,
	pathOptimizer file_gateway.PathOptimizer,
) repository.RaceTimeRepository {
	return &raceTimeRepository{
		netKeibaGateway: netKeibaGateway,
		pathOptimizer:   pathOptimizer,
	}
}

func (r *raceTimeRepository) List(
	ctx context.Context,
	path string,
) ([]string, error) {
	rootPath, err := r.pathOptimizer.GetProjectRoot()
	if err != nil {
		return nil, err
	}

	absPath, err := filepath.Abs(fmt.Sprintf("%s/%s", rootPath, path))
	if err != nil {
		return nil, err
	}

	pattern := filepath.Join(absPath, "*.json")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	fileNames := make([]string, 0, len(files))
	for _, file := range files {
		fileNames = append(fileNames, filepath.Base(file))
	}

	return fileNames, nil
}

func (r *raceTimeRepository) Read(
	ctx context.Context,
	path string,
) ([]*raw_entity.RaceTime, error) {
	raceTimes := make([]*raw_entity.RaceTime, 0)
	rootPath, err := r.pathOptimizer.GetProjectRoot()
	if err != nil {
		return nil, err
	}

	filePath, err := filepath.Abs(fmt.Sprintf("%s/%s", rootPath, path))
	if err != nil {
		return nil, err
	}

	// ファイルが存在しない場合はエラーは返さず処理を継続する
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return raceTimes, nil
	}

	var raceTimeInfo *raw_entity.RaceTimeInfo
	if err := json.Unmarshal(bytes, &raceTimeInfo); err != nil {
		return nil, err
	}
	raceTimes = raceTimeInfo.RaceTimes

	return raceTimes, nil
}

func (r *raceTimeRepository) Write(
	ctx context.Context,
	path string,
	raceTimeInfo *raw_entity.RaceTimeInfo,
) error {
	var buffer bytes.Buffer
	enc := json.NewEncoder(&buffer)
	enc.SetEscapeHTML(false)
	err := enc.Encode(raceTimeInfo)
	if err != nil {
		return err
	}

	rootPath, err := r.pathOptimizer.GetProjectRoot()
	if err != nil {
		return err
	}

	filePath, err := filepath.Abs(fmt.Sprintf("%s/%s", rootPath, path))
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, buffer.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}

func (r *raceTimeRepository) Fetch(
	ctx context.Context,
	url string,
) (*netkeiba_entity.RaceTime, error) {
	raceTime, err := r.netKeibaGateway.FetchRaceTime(ctx, url)
	if err != nil {
		return nil, err
	}
	return raceTime, nil
}
