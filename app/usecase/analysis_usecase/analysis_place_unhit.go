package analysis_usecase

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/tospo_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

func (a *analysis) PlaceUnHit(ctx context.Context, input *AnalysisInput) error {
	unHitRaces := a.placeUnHitService.GetUnHitRaces(ctx, input.Markers, input.Races)

	horses, err := a.horseMasterService.Get(ctx)
	if err != nil {
		return err
	}

	raceForecasts, err := a.raceForecastService.Get(ctx)
	if err != nil {
		return err
	}

	calculables, err := a.placeService.Create(ctx, input.Markers, input.Races)
	if err != nil {
		return err
	}

	//unHitRaceMap := converter.ConvertToMap(unHitRaces, func(race *analysis_entity.Race) types.RaceId {
	//	return race.RaceId()
	//})

	cacheHorseMap := converter.ConvertToMap(horses, func(horse *data_cache_entity.Horse) types.HorseId {
		return horse.HorseId()
	})

	cacheRaceForecastMap := converter.ConvertToMap(raceForecasts, func(forecast *data_cache_entity.RaceForecast) types.RaceId {
		return forecast.RaceId()
	})

	fetchHorseMap := map[types.RaceDate][]*netkeiba_entity.Horse{}
	fetchRaceForecastMap := map[types.RaceId][]*tospo_entity.Forecast{}
	fetchTrainingCommentMap := map[types.RaceId][]*tospo_entity.TrainingComment{}
	unHitRaceRateMap := map[types.RaceId]map[types.HorseId][]float64{}

	for _, race := range unHitRaces {
		unHitRaceRateMap[race.RaceId()] = a.placeUnHitService.GetUnHitRaceRate(ctx, race, calculables)
		// 実用上は1レースで分析対象のオッズは1つ想定だが、仕様上は複数オッズも計算可能なのでループを回す
		for _, raceResult := range race.RaceResults() {
			cachedHorse, ok := cacheHorseMap[raceResult.HorseId()]
			if !ok || race.RaceDate() != cachedHorse.LatestRaceDate() {
				fetchHorse, err := a.placeUnHitService.FetchHorse(ctx, raceResult.HorseId())
				if err != nil {
					return err
				}
				if _, ok = fetchHorseMap[race.RaceDate()]; !ok {
					fetchHorseMap[race.RaceDate()] = make([]*netkeiba_entity.Horse, 0)
				}
				fetchHorseMap[race.RaceDate()] = append(fetchHorseMap[race.RaceDate()], fetchHorse)
			}
		}

		if _, ok := cacheRaceForecastMap[race.RaceId()]; !ok {
			fetchRaceForecasts, err := a.placeUnHitService.FetchRaceForecasts(ctx, race.RaceId())
			if err != nil {
				return err
			}
			fetchRaceForecastMap[race.RaceId()] = fetchRaceForecasts

			fetchTrainingComments, err := a.placeUnHitService.FetchTrainingComments(ctx, race.RaceId())
			if err != nil {
				return err
			}
			fetchTrainingCommentMap[race.RaceId()] = fetchTrainingComments
		}
	}

	if len(fetchHorseMap) > 0 {
		cacheHorses := make([]*data_cache_entity.Horse, 0, len(fetchHorseMap))
		for raceDate, newHorses := range fetchHorseMap {
			for _, fetchHorse := range newHorses {
				cacheHorse, err := a.horseEntityConverter.NetKeibaToDataCache(fetchHorse, raceDate)
				if err != nil {
					return err
				}
				cacheHorses = append(cacheHorses, cacheHorse)
			}
		}
		if err = a.horseMasterService.CreateOrUpdate(ctx, cacheHorses); err != nil {
			return err
		}
	}

	if len(fetchRaceForecastMap) > 0 && len(fetchTrainingCommentMap) > 0 {
		cacheRaceForecasts := make([]*data_cache_entity.RaceForecast, 0, len(fetchRaceForecastMap))
		for _, race := range unHitRaces {
			fetchRaceForecasts, ok := fetchRaceForecastMap[race.RaceId()]
			if !ok {
				continue
			}
			fetchTrainingComments, ok := fetchTrainingCommentMap[race.RaceId()]
			if !ok {
				continue
			}

			horseNumberMap := converter.ConvertToMap(race.RaceResults(), func(raceResult *analysis_entity.RaceResult) types.HorseNumber {
				return raceResult.HorseNumber()
			})

			forecasts := make([]*data_cache_entity.Forecast, 0, len(fetchRaceForecasts))
			for i := 0; i < len(fetchRaceForecasts); i++ {
				if _, ok = horseNumberMap[fetchRaceForecasts[i].HorseNumber()]; ok {
					forecast := a.raceForecastEntityConverter.TospoToDataCache(fetchRaceForecasts[i], fetchTrainingComments[i])
					forecasts = append(forecasts, forecast)
				}
			}
			cacheRaceForecasts = append(cacheRaceForecasts, data_cache_entity.NewRaceForecast(
				race.RaceId().String(),
				race.RaceDate().Value(),
				forecasts,
			))
		}

		if err = a.raceForecastService.CreateOrUpdate(ctx, cacheRaceForecasts); err != nil {
			return err
		}
	}

	cacheHorses, err := a.horseMasterService.Get(ctx)
	if err != nil {
		return err
	}

	cacheRaceForecasts, err := a.raceForecastService.Get(ctx)
	if err != nil {
		return err
	}

	_ = unHitRaces
	_ = cacheHorses
	_ = unHitRaceRateMap
	_ = cacheRaceForecasts

	return nil
}
