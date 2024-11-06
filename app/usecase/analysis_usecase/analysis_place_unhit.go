package analysis_usecase

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

func (a *analysis) PlaceUnHit(ctx context.Context, input *AnalysisInput) error {
	unHitRaces := a.placeUnHitService.GetUnHitRaces(ctx, input.Markers, input.Races)

	horses, err := a.horseMasterService.Get(ctx)
	if err != nil {
		return err
	}

	calculables, err := a.placeService.Create(ctx, input.Markers, input.Races)
	if err != nil {
		return err
	}
	_ = calculables

	cacheHorseMap := converter.ConvertToMap(horses, func(horse *data_cache_entity.Horse) types.HorseId {
		return horse.HorseId()
	})

	fetchHorseMap := map[types.RaceDate][]*netkeiba_entity.Horse{}
	for _, race := range unHitRaces {
		for _, raceResult := range race.RaceResults() {
			cachedHorse, ok := cacheHorseMap[raceResult.HorseId()]
			if !ok || race.RaceDate() != cachedHorse.LatestRaceDate() {
				fetchHorse, err := a.placeUnHitService.FetchHorse(ctx, raceResult.HorseId())
				if err != nil {
					return err
				}
				if _, ok := fetchHorseMap[race.RaceDate()]; !ok {
					fetchHorseMap[race.RaceDate()] = make([]*netkeiba_entity.Horse, 0)
				}
				fetchHorseMap[race.RaceDate()] = append(fetchHorseMap[race.RaceDate()], fetchHorse)
			}
		}
	}

	if len(fetchHorseMap) > 0 {
		for raceDate, fetchHorses := range fetchHorseMap {
			for _, fetchHorse := range fetchHorses {
				cacheHorse, err := a.horseEntityConverter.NetKeibaToDataCache(fetchHorse, raceDate)
				if err != nil {
					return err
				}
				cacheHorseMap[cacheHorse.HorseId()] = cacheHorse
			}
		}
		cacheHorses := make([]*data_cache_entity.Horse, 0, len(cacheHorseMap))
		for _, cacheHorse := range cacheHorseMap {
			cacheHorses = append(cacheHorses, cacheHorse)
		}
		if err = a.horseMasterService.CreateOrUpdate(ctx, cacheHorses); err != nil {
			return err
		}
	}

	cacheHorses, err := a.horseMasterService.Get(ctx)
	if err != nil {
		return err
	}

	_ = unHitRaces
	_ = cacheHorses

	return nil
}
