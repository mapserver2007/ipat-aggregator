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

type jockeyResultRepository struct {
	netKeibaGateway gateway.NetKeibaGateway
	pathOptimizer   file_gateway.PathOptimizer
}

func NewJockeyResultRepository(
	netKeibaGateway gateway.NetKeibaGateway,
	pathOptimizer file_gateway.PathOptimizer,
) repository.JockeyResultRepository {
	return &jockeyResultRepository{
		netKeibaGateway: netKeibaGateway,
		pathOptimizer:   pathOptimizer,
	}
}

func (j *jockeyResultRepository) Write(
	ctx context.Context,
	path string,
	jockeyResultInfo *raw_entity.JockeyResultInfo,
) error {
	var buffer bytes.Buffer
	enc := json.NewEncoder(&buffer)
	enc.SetEscapeHTML(false)
	err := enc.Encode(jockeyResultInfo)
	if err != nil {
		return err
	}

	rootPath, err := j.pathOptimizer.GetProjectRoot()
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

func (j *jockeyResultRepository) FetchUrls(
	ctx context.Context,
	url string,
) ([]string, error) {
	return j.netKeibaGateway.FetchJockeyResultUrls(ctx, url)
}

func (j *jockeyResultRepository) FetchJockeyResults(
	ctx context.Context,
	url string,
) ([]*netkeiba_entity.JockeyResult, error) {
	return j.netKeibaGateway.FetchJockeyResults(ctx, url)
}
