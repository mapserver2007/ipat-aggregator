package converter

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
)

type RaceForecastEntityConverter interface {
	RawToDataCache(input *raw_entity.RaceForecast) *data_cache_entity.RaceForecast
}

type raceForecastEntityConverter struct{}

func NewRaceForecastEntityConverter() RaceForecastEntityConverter {
	return &raceForecastEntityConverter{}
}

func (r *raceForecastEntityConverter) RawToDataCache(input *raw_entity.RaceForecast) *data_cache_entity.RaceForecast {
	forecasts := make([]*data_cache_entity.Forecast, 0, len(input.Forecasts))
	for _, forecast := range input.Forecasts {
		forecasts = append(forecasts, data_cache_entity.NewForecast(
			forecast.HorseNumber,
			forecast.TrainingComment,
			forecast.PreviousTrainingComment,
			forecast.HighlyRecommended,
			forecast.FavoriteNum,
			forecast.RivalNum,
			forecast.MarkerNum,
		))
	}

	return data_cache_entity.NewRaceForecast(
		input.RaceId,
		input.RaceDate,
		forecasts,
	)
}
