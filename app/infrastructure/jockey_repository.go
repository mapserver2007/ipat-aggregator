package infrastructure

import (
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

type jockeyRepository struct {
	netKeibaGateway gateway.NetKeibaGateway
}

func NewJockeyRepository(
	netKeibaGateway gateway.NetKeibaGateway,
) repository.JockeyRepository {
	return &jockeyRepository{
		netKeibaGateway: netKeibaGateway,
	}
}

func (j *jockeyRepository) Read(
	ctx context.Context,
	path string,
) (*raw_entity.JockeyInfo, error) {
	rootPath, err := os.Getwd()
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

	var jockeyInfo *raw_entity.JockeyInfo
	if err := json.Unmarshal(bytes, &jockeyInfo); err != nil {
		return nil, err
	}

	return jockeyInfo, nil
}

func (j *jockeyRepository) Write(
	ctx context.Context,
	path string,
	data *raw_entity.JockeyInfo,
) error {
	bytes, err := json.Marshal(data)
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

	err = os.WriteFile(filePath, bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (j *jockeyRepository) Fetch(
	ctx context.Context,
	url string,
) (*netkeiba_entity.Jockey, error) {
	jockey, err := j.netKeibaGateway.FetchJockey(ctx, url)
	if err != nil {
		return nil, err
	}
	return jockey, nil
}
