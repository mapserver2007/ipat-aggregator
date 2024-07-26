package repository

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/tospo_entity"
)

type RaceForecastRepository interface {
	List(ctx context.Context, path string) ([]string, error)
	Read(ctx context.Context, path string) ([]*raw_entity.RaceForecast, error)
	Write(ctx context.Context, path string, forecastInfo *raw_entity.RaceForecastInfo) error
	FetchRaceForecast(ctx context.Context, url string) ([]*tospo_entity.Forecast, error)
	FetchTrainingComment(ctx context.Context, url string) ([]*tospo_entity.TrainingComment, error)
}
