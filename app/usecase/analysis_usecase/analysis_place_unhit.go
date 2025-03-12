package analysis_usecase

import (
	"context"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/tospo_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/shopspring/decimal"
)

const (
	winOddsThreshold           = 10.0 // 赤オッズボーダー
	trioPopularNumberThreshold = 100  // 三連複100番人気
)

func (a *analysis) PlaceUnHit(
	ctx context.Context,
	input *AnalysisInput,
) error {
	unHitRaces := a.placeUnHitService.GetUnHitRaces(ctx, input.Markers, input.Races, input.Jockeys)

	winRedOdds, err := a.placeUnHitService.GetWinRedOdds(ctx, input.Odds.Win, decimal.NewFromFloat(winOddsThreshold))
	if err != nil {
		return err
	}

	winOddsFaults, err := a.placeUnHitService.GetWinOddsFaults(ctx, input.Odds.Win)
	if err != nil {
		return err
	}

	trioOdds, err := a.placeUnHitService.GetTrioOdds(ctx, input.Odds.Trio, trioPopularNumberThreshold)
	if err != nil {
		return err
	}

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

	unHitRaceMap := converter.ConvertToMap(unHitRaces, func(race *analysis_entity.Race) types.RaceId {
		return race.RaceId()
	})

	cacheHorseMap := converter.ConvertToMap(horses, func(horse *data_cache_entity.Horse) types.HorseId {
		return horse.HorseId()
	})

	cacheRaceForecastMap := converter.ConvertToMap(raceForecasts, func(forecast *data_cache_entity.RaceForecast) types.RaceId {
		return forecast.RaceId()
	})

	winMultiOddsMap := map[types.RaceId][]*analysis_entity.Odds{}
	for _, odds := range winRedOdds {
		if _, ok := winMultiOddsMap[odds.RaceId()]; !ok {
			winMultiOddsMap[odds.RaceId()] = make([]*analysis_entity.Odds, 0)
		}
		winMultiOddsMap[odds.RaceId()] = append(winMultiOddsMap[odds.RaceId()], odds)
	}

	winOddsFaultMap := map[types.RaceId][]*analysis_entity.OddsFault{}
	for _, oddsFault := range winOddsFaults {
		if oddsFault.OddsFaultNo() > 2 {
			continue
		}
		if _, ok := winOddsFaultMap[oddsFault.RaceId()]; !ok {
			winOddsFaultMap[oddsFault.RaceId()] = make([]*analysis_entity.OddsFault, 0)
		}
		winOddsFaultMap[oddsFault.RaceId()] = append(winOddsFaultMap[oddsFault.RaceId()], oddsFault)
	}

	trioOddsMap := converter.ConvertToMap(trioOdds, func(odds *analysis_entity.Odds) types.RaceId {
		return odds.RaceId()
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
			for i, forecast := range fetchRaceForecasts {
				if _, ok = horseNumberMap[forecast.HorseNumber()]; ok {
					forecasts = append(forecasts, a.raceForecastEntityConverter.TospoToDataCache(forecast, fetchTrainingComments[i]))
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

	cacheRaceForecasts, err := a.raceForecastService.Get(ctx)
	if err != nil {
		return err
	}

	unHitRaceForecastMap := map[types.RaceId]*data_cache_entity.RaceForecast{}
	for _, raceForecast := range cacheRaceForecasts {
		if _, ok := unHitRaceMap[raceForecast.RaceId()]; ok {
			unHitRaceForecastMap[raceForecast.RaceId()] = raceForecast
		}
	}

	cacheHorses, err := a.horseMasterService.Get(ctx)
	if err != nil {
		return err
	}

	horseMap := converter.ConvertToMap(cacheHorses, func(horse *data_cache_entity.Horse) types.HorseId {
		return horse.HorseId()
	})

	analysisPlaceUnhits, err := a.placeUnHitService.CreateUnhitRaces(
		ctx,
		unHitRaces,
		unHitRaceRateMap,
		unHitRaceForecastMap,
		horseMap,
		winMultiOddsMap,
		winOddsFaultMap,
		trioOddsMap,
	)
	if err != nil {
		return err
	}

	for _, analysisPlaceUnhit := range analysisPlaceUnhits {
		race := unHitRaceMap[analysisPlaceUnhit.RaceId()]
		a.placeUnHitService.CreateCheckPoints(ctx, race)
	}

	err = a.placeUnHitService.Write(ctx, analysisPlaceUnhits)
	if err != nil {
		return err
	}

	return nil
}
