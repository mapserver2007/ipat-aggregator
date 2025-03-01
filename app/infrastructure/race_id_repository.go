package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure/file_gateway"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure/gateway"
)

type raceIdRepository struct {
	netKeibaGateway gateway.NetKeibaGateway
	pathOptimizer   file_gateway.PathOptimizer
}

func NewRaceIdRepository(
	netKeibaGateway gateway.NetKeibaGateway,
	pathOptimizer file_gateway.PathOptimizer,
) repository.RaceIdRepository {
	return &raceIdRepository{
		netKeibaGateway: netKeibaGateway,
		pathOptimizer:   pathOptimizer,
	}
}

func (r *raceIdRepository) Read(
	ctx context.Context,
	path string,
) (*raw_entity.RaceIdInfo, error) {
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
		return nil, nil
	}

	var raceIdInfo *raw_entity.RaceIdInfo
	if err := json.Unmarshal(bytes, &raceIdInfo); err != nil {
		return nil, err
	}

	return raceIdInfo, nil
}

func (r *raceIdRepository) Write(
	ctx context.Context,
	path string,
	data *raw_entity.RaceIdInfo,
) error {
	bytes, err := json.Marshal(data)
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

	err = os.WriteFile(filePath, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (r *raceIdRepository) Fetch(
	ctx context.Context,
	url string,
) ([]string, error) {
	rawRaceIds, err := r.netKeibaGateway.FetchRaceId(ctx, url)
	if err != nil {
		return nil, err
	}

	return rawRaceIds, nil
}
