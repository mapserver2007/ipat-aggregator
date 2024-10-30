package prediction_usecase

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/config"
)

func (p *prediction) CheckList(ctx context.Context, input *PredictionInput) error {
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

			predictionHorse, err := p.predictionPlaceCandidateService.GetHorse(ctx, horse.HorseId())
			if err != nil {
				return err
			}

			predictionJockey, err := p.predictionPlaceCandidateService.GetJockey(ctx, horse.JockeyId())
			if err != nil {
				return err
			}

			predictionTrainer, err := p.predictionPlaceCandidateService.GetTrainer(ctx, predictionHorse.TrainerId())
			if err != nil {
				return err
			}

			predictionCheckList := p.predictionPlaceCandidateService.Convert(
				ctx,
				predictionRace,
				predictionHorse,
				predictionJockey,
				predictionTrainer,
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

	return nil
}
