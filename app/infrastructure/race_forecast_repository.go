package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/tospo_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure/gateway"
	"os"
	"path/filepath"
)

type raceForecastRepository struct {
	tospoGateway gateway.TospoGateway
}

func NewRaceForecastRepository(
	tospoGateway gateway.TospoGateway,
) repository.RaceForecastRepository {
	return &raceForecastRepository{
		tospoGateway: tospoGateway,
	}
}

func (r *raceForecastRepository) List(
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

func (r *raceForecastRepository) Read(
	ctx context.Context,
	path string,
) ([]*raw_entity.RaceForecast, error) {
	raceForecasts := make([]*raw_entity.RaceForecast, 0)
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
		return raceForecasts, nil
	}

	var raceForecastInfo *raw_entity.RaceForecastInfo
	if err := json.Unmarshal(bytes, &raceForecastInfo); err != nil {
		return nil, err
	}
	raceForecasts = raceForecastInfo.RaceForecasts

	return raceForecasts, nil
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

func (r *raceForecastRepository) FetchRaceForecast(
	ctx context.Context,
	url string,
) ([]*tospo_entity.Forecast, error) {
	forecasts, err := r.tospoGateway.FetchForecast(ctx, url)
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
