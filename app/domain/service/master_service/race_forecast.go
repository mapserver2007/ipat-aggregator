package master_service

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/config"
	"time"
)

const (
	raceForecastUrl        = "https://tospo-keiba.jp/race/detail/%s/forecast"
	raceTrainingCommentUrl = "https://tospo-keiba.jp/race/detail/%s/comment"
	forecastFileName       = "forecast_%d.json"
)

type RaceForecast interface {
	Get(ctx context.Context) ([]*data_cache_entity.RaceForecast, error)
	CreateOrUpdate(ctx context.Context, races []*data_cache_entity.Race) error
}

type raceForecastService struct {
	raceForecastRepository      repository.RaceForecastRepository
	raceForecastEntityConverter converter.RaceForecastEntityConverter
}

func NewRaceForecast(
	raceForecastRepository repository.RaceForecastRepository,
	raceForecastEntityConverter converter.RaceForecastEntityConverter,
) RaceForecast {
	return &raceForecastService{
		raceForecastRepository:      raceForecastRepository,
		raceForecastEntityConverter: raceForecastEntityConverter,
	}
}

func (r *raceForecastService) Get(ctx context.Context) ([]*data_cache_entity.RaceForecast, error) {
	files, err := r.raceForecastRepository.List(ctx, fmt.Sprintf("%s/forecasts", config.CacheDir))
	if err != nil {
		return nil, err
	}

	var raceForecasts []*data_cache_entity.RaceForecast
	for _, file := range files {
		rawRaceForecasts, err := r.raceForecastRepository.Read(ctx, fmt.Sprintf("%s/forecasts/%s", config.CacheDir, file))
		if err != nil {
			return nil, err
		}
		for _, rawRaceForecast := range rawRaceForecasts {
			raceForecasts = append(raceForecasts, r.raceForecastEntityConverter.RawToDataCache(rawRaceForecast))
		}
	}

	return raceForecasts, nil
}

func (r *raceForecastService) CreateOrUpdate(
	ctx context.Context,
	races []*data_cache_entity.Race,
) error {
	raceForecasts, err := r.Get(ctx)
	if err != nil {
		return err
	}

	cacheRaceDateMap := map[types.RaceDate]struct{}{}
	for _, raceForecast := range raceForecasts {
		if _, ok := cacheRaceDateMap[raceForecast.RaceDate()]; !ok {
			cacheRaceDateMap[raceForecast.RaceDate()] = struct{}{}
		}
	}

	raceMap := map[types.RaceDate][]*data_cache_entity.Race{}
	for _, race := range races {
		_, ok := raceMap[race.RaceDate()]
		if !ok {
			raceMap[race.RaceDate()] = make([]*data_cache_entity.Race, 0)
		}
		raceMap[race.RaceDate()] = append(raceMap[race.RaceDate()], race)
	}

	for raceDate, races2 := range raceMap {
		if _, ok := cacheRaceDateMap[raceDate]; ok {
			continue
		}
		var rawRaceForecasts []*raw_entity.RaceForecast
		for _, race := range races2 {
			if race.Organizer() != types.JRA {
				continue
			}
			time.Sleep(time.Millisecond)
			forecasts, err := r.raceForecastRepository.FetchRaceForecast(ctx, fmt.Sprintf(raceForecastUrl, race.RaceId()))
			if err != nil {
				return err
			}

			time.Sleep(time.Millisecond)
			trainingComments, err := r.raceForecastRepository.FetchTrainingComment(ctx, fmt.Sprintf(raceTrainingCommentUrl, race.RaceId()))
			if err != nil {
				return err
			}

			rawForecasts := make([]*raw_entity.Forecast, 0, len(forecasts))
			for i := 0; i < len(forecasts); i++ {
				forecast := forecasts[i]
				trainingComment := "-"
				previousTrainingComment := "-"
				isHighlyRecommended := false

				if i >= 0 && i < len(trainingComments) {
					trainingComment = trainingComments[i].TrainingComment()
					previousTrainingComment = trainingComments[i].PreviousTrainingComment()
					isHighlyRecommended = trainingComments[i].IsHighlyRecommended()
				}

				rawForecasts = append(rawForecasts, &raw_entity.Forecast{
					HorseNumber:             forecast.HorseNumber().Value(),
					TrainingComment:         trainingComment,
					PreviousTrainingComment: previousTrainingComment,
					HighlyRecommended:       isHighlyRecommended,
					FavoriteNum:             forecast.FavoriteNum(),
					RivalNum:                forecast.RivalNum(),
					MarkerNum:               forecast.MarkerNum(),
				})
			}

			rawRaceForecasts = append(rawRaceForecasts, &raw_entity.RaceForecast{
				RaceId:    race.RaceId().String(),
				RaceDate:  race.RaceDate().Value(),
				Forecasts: rawForecasts,
			})
		}

		if len(rawRaceForecasts) == 0 {
			continue
		}

		err = r.raceForecastRepository.Write(
			ctx,
			fmt.Sprintf("%s/forecasts/%s", config.CacheDir, fmt.Sprintf(forecastFileName, raceDate)),
			&raw_entity.RaceForecastInfo{RaceForecasts: rawRaceForecasts},
		)
		if err != nil {
			return err
		}
	}

	return nil
}
