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

func (j *jockeyResultRepository) List(
	ctx context.Context,
	path string,
) ([]string, error) {
	rootPath, err := j.pathOptimizer.GetProjectRoot()
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

func (j *jockeyResultRepository) Read(
	ctx context.Context,
	path string,
) ([]*raw_entity.JockeyResult, error) {
	jockeyResults := make([]*raw_entity.JockeyResult, 0)
	rootPath, err := j.pathOptimizer.GetProjectRoot()
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
		return jockeyResults, nil
	}

	var raceInfo *raw_entity.JockeyResultInfo
	if err := json.Unmarshal(bytes, &raceInfo); err != nil {
		return nil, err
	}
	jockeyResults = raceInfo.JockeyResults

	return jockeyResults, nil
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

func (j *jockeyResultRepository) FetchJockeyResultCounts(
	ctx context.Context,
	url string,
) ([]*raw_entity.JockeyResultCount, error) {
	return j.netKeibaGateway.FetchJockeyResultCounts(ctx, url)
}

func (j *jockeyResultRepository) FetchJockeyResultUrls(
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
