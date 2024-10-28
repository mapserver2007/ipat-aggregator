package prediction_usecase

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/analysis_service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/prediction_service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/config"
)

type Prediction interface {
	Execute(ctx context.Context, input *PredictionInput) error
}

type PredictionInput struct {
	AnalysisMarkers   []*marker_csv_entity.AnalysisMarker
	PredictionMarkers []*marker_csv_entity.PredictionMarker
	Races             []*data_cache_entity.Race
}

type prediction struct {
	predictionOddsService           prediction_service.Odds
	predictionPlaceCandidateService prediction_service.PlaceCandidate
	predictionMarkerSyncService     prediction_service.MarkerSync
	placeService                    analysis_service.Place
}

func NewPrediction(
	predictionOddsService prediction_service.Odds,
	predictionPlaceCandidateService prediction_service.PlaceCandidate,
	predictionMarkerSyncService prediction_service.MarkerSync,
	placeService analysis_service.Place,
) Prediction {
	return &prediction{
		predictionOddsService:           predictionOddsService,
		predictionPlaceCandidateService: predictionPlaceCandidateService,
		predictionMarkerSyncService:     predictionMarkerSyncService,
		placeService:                    placeService,
	}
}

func (p *prediction) Execute(ctx context.Context, input *PredictionInput) error {
	if config.EnablePredictionOdds {
		predictionMarkers := input.PredictionMarkers
		predictionRaces := make([]*prediction_entity.Race, 0, len(predictionMarkers))

		for _, marker := range predictionMarkers {
			predictionRace, err := p.predictionOddsService.Get(ctx, marker.RaceId())
			if err != nil {
				return err
			}
			predictionRaces = append(predictionRaces, predictionRace)
		}

		calculables, err := p.placeService.Create(ctx, input.AnalysisMarkers, input.Races)
		if err != nil {
			return err
		}

		firstPlaceMap, secondPlaceMap, thirdPlaceMap, raceCourseMap := p.predictionOddsService.ConvertAll(ctx, predictionRaces, predictionMarkers, calculables)
		err = p.predictionOddsService.Write(ctx, firstPlaceMap, secondPlaceMap, thirdPlaceMap, raceCourseMap)
		if err != nil {
			return err
		}
	}

	if config.EnablePredictionCheckList {
		predictionMarkers := input.PredictionMarkers
		predictionRaces := make([]*prediction_entity.Race, 0, len(predictionMarkers))
		predictionCheckLists := make([]*spreadsheet_entity.PredictionCheckList, 0, len(predictionMarkers)*6)

		calculables, err := p.placeService.Create(ctx, input.AnalysisMarkers, input.Races)
		if err != nil {
			return err
		}

		for _, predictionMarker := range predictionMarkers {
			predictionRace, err := p.predictionPlaceCandidateService.GetRaceCard(ctx, predictionMarker.RaceId())
			if err != nil {
				return err
			}
			predictionRaces = append(predictionRaces, predictionRace)

			raceForecasts, err := p.predictionPlaceCandidateService.GetRaceForecasts(ctx, predictionMarker.RaceId())
			if err != nil {
				return err
			}

			horseNumberMap := converter.ConvertToMap(predictionRace.RaceEntryHorses(), func(horse *prediction_entity.RaceEntryHorse) types.HorseNumber {
				return horse.HorseNumber()
			})

			horseOddsMap := converter.ConvertToMap(predictionRace.Odds(), func(o *prediction_entity.Odds) types.HorseNumber {
				return o.HorseNumber()
			})

			raceForecastMap := converter.ConvertToMap(raceForecasts, func(forecast *prediction_entity.RaceForecast) types.HorseNumber {
				return forecast.HorseNumber()
			})

			horseNumbers := []types.HorseNumber{
				predictionMarker.Favorite(), predictionMarker.Rival(), predictionMarker.BrackTriangle(), predictionMarker.WhiteTriangle(), predictionMarker.Star(), predictionMarker.Check()}
			for idx, horseNumber := range horseNumbers {
				odds, ok := horseOddsMap[horseNumber]
				if !ok {
					return fmt.Errorf("odds not found: %d", horseNumber)
				}
				if odds.Odds().InexactFloat64() > config.PredictionCheckListWinLowerOdds {
					continue
				}

				marker, err := types.NewMarker(idx + 1)
				if err != nil {
					return err
				}
				horse, ok := horseNumberMap[horseNumber]
				if !ok {
					return fmt.Errorf("horse not found: %d", horseNumber)
				}
				raceForecast, ok := raceForecastMap[horseNumber]
				if !ok {
					return fmt.Errorf("race forecast not found: %d", horseNumber)
				}
				predictionHorse, err := p.predictionPlaceCandidateService.GetHorse(ctx, horse)
				if err != nil {
					return err
				}

				predictionCheckList := p.predictionPlaceCandidateService.Convert(
					ctx,
					predictionRace,
					predictionHorse,
					raceForecast,
					calculables,
					horseNumber,
					marker,
					p.predictionPlaceCandidateService.CreateCheckList(ctx, predictionRace, predictionHorse, raceForecast),
				)

				predictionCheckLists = append(predictionCheckLists, predictionCheckList)
			}
		}

		err = p.predictionPlaceCandidateService.Write(ctx, predictionCheckLists)
		if err != nil {
			return err
		}
	}

	if config.EnablePredictionSync {
		raceDate, err := types.NewRaceDate(config.PredictionSyncRaceDate)
		if err != nil {
			return err
		}

		raceIds, err := p.predictionMarkerSyncService.GetRaceIds(ctx, raceDate)
		if err != nil {
			return err
		}

		var predictionMarkers []*prediction_entity.Marker
		for _, raceId := range raceIds {
			markers, err := p.predictionMarkerSyncService.GetMarkers(ctx, raceId)
			if err != nil {
				return err
			}
			predictionMarkers = append(predictionMarkers, markers...)
		}

		spreadSheetPredictionMarkers := p.predictionMarkerSyncService.Convert(ctx, predictionMarkers)
		err = p.predictionMarkerSyncService.Write(ctx, spreadSheetPredictionMarkers)
		if err != nil {
			return err
		}
	}

	return nil
}
