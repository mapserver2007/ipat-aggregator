package prediction_usecase

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
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
	placeService                    analysis_service.Place
}

func NewPrediction(
	predictionOddsService prediction_service.Odds,
	predictionPlaceCandidateService prediction_service.PlaceCandidate,
	placeService analysis_service.Place,
) Prediction {
	return &prediction{
		predictionOddsService:           predictionOddsService,
		predictionPlaceCandidateService: predictionPlaceCandidateService,
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

		firstPlaceMap, secondPlaceMap, thirdPlaceMap, raceCourseMap := p.predictionOddsService.Convert(ctx, predictionRaces, predictionMarkers, calculables)
		err = p.predictionOddsService.Write(ctx, firstPlaceMap, secondPlaceMap, thirdPlaceMap, raceCourseMap)
		if err != nil {
			return err
		}
	}

	if config.EnablePredictionCheckList {
		predictionMarkers := input.PredictionMarkers
		predictionHorses := make([]*prediction_entity.Horse, 0, len(predictionMarkers))
		raceId := predictionMarkers[0].RaceId()
		predictionRace, err := p.predictionPlaceCandidateService.GetRaceCard(ctx, raceId)
		if err != nil {
			return err
		}

		horseNumberMap := converter.ConvertToMap(predictionRace.RaceEntryHorses(), func(horse *prediction_entity.RaceEntryHorse) types.HorseNumber {
			return horse.HorseNumber()
		})

		for _, marker := range predictionMarkers {
			horseNumbers := []types.HorseNumber{
				marker.Favorite(), marker.Rival(), marker.BrackTriangle(), marker.WhiteTriangle(), marker.Star(), marker.Check()}
			for _, horseNumber := range horseNumbers {
				horse, ok := horseNumberMap[horseNumber]
				if !ok {
					return fmt.Errorf("horse not found: %d", horseNumber)
				}
				predictionHorse, err := p.predictionPlaceCandidateService.GetHorse(ctx, horse)
				if err != nil {
					return err
				}
				predictionHorses = append(predictionHorses, predictionHorse)
			}
		}
	}

	return nil
}
