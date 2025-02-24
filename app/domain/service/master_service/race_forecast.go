package master_service

import (
	"context"
	"fmt"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/config"
)

const raceForecastFileName = "race_forecast.json"

type RaceForecast interface {
	Get(ctx context.Context) ([]*data_cache_entity.RaceForecast, error)
	CreateOrUpdate(ctx context.Context, raceForecasts []*data_cache_entity.RaceForecast) error
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
	rawRaceForecastInfo, err := r.raceForecastRepository.Read(ctx, fmt.Sprintf("%s/%s", config.CacheDir, raceForecastFileName))
	if err != nil {
		return nil, err
	}

	var raceForecasts []*data_cache_entity.RaceForecast
	if rawRaceForecastInfo != nil {
		for _, rawRaceForecast := range rawRaceForecastInfo.RaceForecasts {
			raceForecasts = append(raceForecasts, r.raceForecastEntityConverter.RawToDataCache(rawRaceForecast))
		}
	}

	return raceForecasts, nil
}

func (r *raceForecastService) CreateOrUpdate(
	ctx context.Context,
	raceForecasts []*data_cache_entity.RaceForecast,
) error {
	caches, err := r.Get(ctx)
	if err != nil {
		return err
	}

	caches = append(caches, raceForecasts...)

	raceIdMap := converter.ConvertToMap(caches, func(raceForecast *data_cache_entity.RaceForecast) types.RaceId {
		return raceForecast.RaceId()
	})

	newCaches := make([]*raw_entity.RaceForecast, 0, len(caches)+len(raceForecasts))

	for _, raceId := range service.SortedRaceIdKeys(raceIdMap) {
		newCaches = append(newCaches, r.raceForecastEntityConverter.DataCacheToRaw(raceIdMap[raceId]))
	}

	err = r.raceForecastRepository.Write(ctx, fmt.Sprintf("%s/%s", config.CacheDir, raceForecastFileName), &raw_entity.RaceForecastInfo{
		RaceForecasts: newCaches,
	})
	if err != nil {
		return err
	}

	cacheRaceDateMap := map[types.RaceDate]struct{}{}
	for _, raceForecast := range raceForecasts {
		if _, ok := cacheRaceDateMap[raceForecast.RaceDate()]; !ok {
			cacheRaceDateMap[raceForecast.RaceDate()] = struct{}{}
		}
	}

	return nil
}
