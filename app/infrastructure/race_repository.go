package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure/gateway"
	"os"
	"path/filepath"
)

type raceRepository struct {
	netKeibaGateway gateway.NetKeibaGateway
}

func NewRaceRepository(
	netKeibaGateway gateway.NetKeibaGateway,
) repository.RaceRepository {
	return &raceRepository{
		netKeibaGateway: netKeibaGateway,
	}
}

func (r *raceRepository) List(
	ctx context.Context,
	path string,
) ([]string, error) {
	rootPath, err := os.Getwd()
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
func (r *raceRepository) Read(
	ctx context.Context,
	path string,
) ([]*raw_entity.Race, error) {
	races := make([]*raw_entity.Race, 0)
	rootPath, err := os.Getwd()
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
		return races, nil
	}

	var raceInfo *raw_entity.RaceInfo
	if err := json.Unmarshal(bytes, &raceInfo); err != nil {
		return nil, err
	}
	races = raceInfo.Races

	return races, nil
}

func (r *raceRepository) Write(
	ctx context.Context,
	path string,
	raceInfo *raw_entity.RaceInfo,
) error {
	var buffer bytes.Buffer
	enc := json.NewEncoder(&buffer)
	enc.SetEscapeHTML(false)
	err := enc.Encode(raceInfo)
	if err != nil {
		return err
	}

	rootPath, err := os.Getwd()
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

func (r *raceRepository) Fetch(
	ctx context.Context,
	url string,
) (*netkeiba_entity.Race, error) {
	race, err := r.netKeibaGateway.FetchRace(ctx, url)
	if err != nil {
		return nil, err
	}
	return race, nil
}
