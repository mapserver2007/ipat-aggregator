package infrastructure

import (
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

type horseRepository struct {
	netKeibaGateway gateway.NetKeibaGateway
	pathOptimizer   file_gateway.PathOptimizer
}

func NewHorseRepository(
	netKeibaGateway gateway.NetKeibaGateway,
	pathOptimizer file_gateway.PathOptimizer,
) repository.HorseRepository {
	return &horseRepository{
		netKeibaGateway: netKeibaGateway,
		pathOptimizer:   pathOptimizer,
	}
}

func (h *horseRepository) Read(
	ctx context.Context,
	path string,
) (*raw_entity.HorseInfo, error) {
	rootPath, err := h.pathOptimizer.GetProjectRoot()
	if err != nil {
		return nil, err
	}

	absPath, err := filepath.Abs(fmt.Sprintf("%s/%s", rootPath, path))
	if err != nil {
		return nil, err
	}

	// ファイルが存在しない場合はエラーは返さず処理を継続する
	bytes, err := os.ReadFile(absPath)
	if err != nil {
		return nil, nil
	}

	var horseInfo *raw_entity.HorseInfo
	if err := json.Unmarshal(bytes, &horseInfo); err != nil {
		return nil, err
	}

	return horseInfo, nil
}

func (h *horseRepository) Write(
	ctx context.Context,
	path string,
	data *raw_entity.HorseInfo,
) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	rootPath, err := h.pathOptimizer.GetProjectRoot()
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

func (h *horseRepository) Fetch(
	ctx context.Context,
	url string,
) (*netkeiba_entity.Horse, error) {
	horse, err := h.netKeibaGateway.FetchHorse(ctx, url)
	if err != nil {
		return nil, err
	}

	return horse, nil
}
