package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/tospo_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure/file_gateway"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure/gateway"
)

type raceForecastRepository struct {
	tospoGateway  gateway.TospoGateway
	pathOptimizer file_gateway.PathOptimizer
}

func NewRaceForecastRepository(
	tospoGateway gateway.TospoGateway,
	pathOptimizer file_gateway.PathOptimizer,
) repository.RaceForecastRepository {
	return &raceForecastRepository{
		tospoGateway:  tospoGateway,
		pathOptimizer: pathOptimizer,
	}
}

func (r *raceForecastRepository) Read(
	ctx context.Context,
	path string,
) (*raw_entity.RaceForecastInfo, error) {
	rootPath, err := r.pathOptimizer.GetProjectRoot()
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

	var raceForecastInfo *raw_entity.RaceForecastInfo
	if err := json.Unmarshal(bytes, &raceForecastInfo); err != nil {
		return nil, err
	}

	return raceForecastInfo, nil
}

func (r *raceForecastRepository) Write(
	ctx context.Context,
	path string,
	forecastInfo *raw_entity.RaceForecastInfo,
) error {
	var buffer bytes.Buffer
	enc := json.NewEncoder(&buffer)
	enc.SetEscapeHTML(false)
	err := enc.Encode(forecastInfo)
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

func (r *raceForecastRepository) FetchRaceForecast(
	ctx context.Context,
	url string,
) ([]*tospo_entity.Forecast, error) {
	rootPath, err := r.pathOptimizer.GetProjectRoot()
	if err != nil {
		return nil, err
	}

	absPath, err := filepath.Abs(fmt.Sprintf("%s/secret/tospo_cookie.json", rootPath))
	if err != nil {
		return nil, err
	}

	cookieFile, err := os.ReadFile(absPath)
	if err != nil {
		return nil, err
	}

	var cookie raw_entity.TospoCookie
	if err := json.Unmarshal(cookieFile, &cookie); err != nil {
		return nil, err
	}

	forecasts, err := r.tospoGateway.FetchForecast(ctx, url, &cookie)
	if err != nil {
		return nil, err
	}

	return forecasts, nil
}

func (r *raceForecastRepository) FetchTrainingComment(
	ctx context.Context,
	url string,
) ([]*tospo_entity.TrainingComment, error) {
	trainingComments, err := r.tospoGateway.FetchTrainingComment(ctx, url)
	if err != nil {
		return nil, err
	}

	return trainingComments, nil
}

func (r *raceForecastRepository) FetchReporterMemo(
	ctx context.Context,
	url string,
) ([]*tospo_entity.ReporterMemo, error) {
	reporterMemos, err := r.tospoGateway.FetchReporterMemo(ctx, url)
	if err != nil {
		return nil, err
	}

	return reporterMemos, nil
}

func (r *raceForecastRepository) FetchPaddockComment(
	ctx context.Context,
	url string,
) ([]*tospo_entity.PaddockComment, error) {
	paddockComment, err := r.tospoGateway.FetchPaddockComment(ctx, url)
	if err != nil {
		return nil, err
	}

	return paddockComment, nil
}
