package converter

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/tospo_entity"
)

type RaceForecastEntityConverter interface {
	RawToDataCache(input *raw_entity.RaceForecast) *data_cache_entity.RaceForecast
	DataCacheToRaw(input *data_cache_entity.RaceForecast) *raw_entity.RaceForecast
	TospoToDataCache(input1 *tospo_entity.Forecast, input2 *tospo_entity.TrainingComment) *data_cache_entity.Forecast
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

func (r *raceForecastEntityConverter) DataCacheToRaw(input *data_cache_entity.RaceForecast) *raw_entity.RaceForecast {
	rawForecasts := make([]*raw_entity.Forecast, 0, len(input.Forecasts()))
	for _, rawForecast := range input.Forecasts() {
		rawForecasts = append(rawForecasts, &raw_entity.Forecast{
			HorseNumber:             rawForecast.HorseNumber().Value(),
			TrainingComment:         rawForecast.TrainingComment(),
			PreviousTrainingComment: rawForecast.PreviousTrainingComment(),
			HighlyRecommended:       rawForecast.HighlyRecommended(),
			FavoriteNum:             rawForecast.FavoriteNum(),
			RivalNum:                rawForecast.RivalNum(),
			MarkerNum:               rawForecast.MarkerNum(),
		})
	}

	return &raw_entity.RaceForecast{
		RaceId:    input.RaceId().String(),
		RaceDate:  input.RaceDate().Value(),
		Forecasts: rawForecasts,
	}
}

func (r *raceForecastEntityConverter) TospoToDataCache(input1 *tospo_entity.Forecast, input2 *tospo_entity.TrainingComment) *data_cache_entity.Forecast {
	return data_cache_entity.NewForecast(
		input1.HorseNumber().Value(),
		input2.TrainingComment(),
		input2.PreviousTrainingComment(),
		input2.IsHighlyRecommended(),
		input1.FavoriteNum(),
		input1.RivalNum(),
		input1.MarkerNum(),
	)
}
