package prediction_usecase

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/config"
)

const checkListParallel = 5

func (p *prediction) CheckList(ctx context.Context, input *PredictionInput) error {
	taskCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	predictionMarkers := input.PredictionMarkers
	predictionCheckLists := make([]*spreadsheet_entity.PredictionCheckList, 0, len(predictionMarkers)*6)

	calculables, err := p.placeService.Create(ctx, input.AnalysisMarkers, input.Races)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errorCh := make(chan error, 1)
	resultCh := make(chan []*spreadsheet_entity.PredictionCheckList, checkListParallel)
	chunkSize := (len(predictionMarkers) + checkListParallel - 1) / checkListParallel

	for i := 0; i < len(predictionMarkers); i += chunkSize {
		end := i + chunkSize
		if end > len(predictionMarkers) {
			end = len(predictionMarkers)
		}

		wg.Add(1)
		go func(markers []*marker_csv_entity.PredictionMarker) {
			defer wg.Done()
			localCheckLists := make([]*spreadsheet_entity.PredictionCheckList, 0, len(markers)*6)
			p.logger.Infof("prediction checkLlst processing: %v/%v", end, len(predictionMarkers))
			for _, marker := range markers {
				select {
				case <-taskCtx.Done():
					return
				default:
					checkLists, err := p.createCheckList(taskCtx, calculables, marker)
					if err != nil {
						select {
						case errorCh <- err:
							cancel()
						}
						return
					}

					localCheckLists = append(localCheckLists, checkLists...)
				}
			}

			resultCh <- localCheckLists
		}(predictionMarkers[i:end])
	}

	wg.Wait()
	close(errorCh)
	close(resultCh)

	if err = <-errorCh; err != nil {
		return err
	}

	for checkLists := range resultCh {
		predictionCheckLists = append(predictionCheckLists, checkLists...)
	}

	sort.Slice(predictionCheckLists, func(i, j int) bool {
		return predictionCheckLists[i].RaceId() < predictionCheckLists[j].RaceId()
	})

	err = p.predictionPlaceCandidateService.Write(ctx, predictionCheckLists)
	if err != nil {
		return err
	}

	return nil
}

func (p *prediction) createCheckList(
	taskCtx context.Context,
	calculables []*analysis_entity.PlaceCalculable,
	marker *marker_csv_entity.PredictionMarker,
) ([]*spreadsheet_entity.PredictionCheckList, error) {
	predictionRace, err := p.predictionPlaceCandidateService.GetRaceCard(taskCtx, marker.RaceId())
	if err != nil {
		return nil, err
	}

	raceForecasts, err := p.predictionPlaceCandidateService.GetRaceForecasts(taskCtx, marker.RaceId())
	if err != nil {
		return nil, err
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
		marker.Favorite(), marker.Rival(), marker.BrackTriangle(), marker.WhiteTriangle(), marker.Star(), marker.Check()}

	predictionCheckLists := make([]*spreadsheet_entity.PredictionCheckList, 0, len(horseNumbers))

	for idx, horseNumber := range horseNumbers {
		odds, ok := horseOddsMap[horseNumber]
		if !ok {
			return nil, fmt.Errorf("invalid marker settings in %s", marker.RaceId())
		}

		if odds.Odds().InexactFloat64() > config.PredictionCheckListWinLowerOdds {
			continue
		}

		newMarker, err := types.NewMarker(idx + 1)
		if err != nil {
			return nil, err
		}

		horse, ok := horseNumberMap[horseNumber]
		if !ok {
			return nil, fmt.Errorf("horse not found: %d", horseNumber)
		}

		raceForecast, ok := raceForecastMap[horseNumber]
		if !ok {
			return nil, fmt.Errorf("race forecast not found: %d", horseNumber)
		}

		predictionHorse, err := p.predictionPlaceCandidateService.GetHorse(taskCtx, horse.HorseId())
		if err != nil {
			return nil, err
		}

		predictionJockey, err := p.predictionPlaceCandidateService.GetJockey(taskCtx, horse.JockeyId())
		if err != nil {
			return nil, err
		}

		predictionTrainer, err := p.predictionPlaceCandidateService.GetTrainer(taskCtx, predictionHorse.TrainerId())
		if err != nil {
			return nil, err
		}

		predictionCheckList := p.predictionPlaceCandidateService.Convert(
			taskCtx,
			predictionRace,
			predictionHorse,
			predictionJockey,
			predictionTrainer,
			raceForecast,
			calculables,
			horseNumber,
			newMarker,
			p.predictionPlaceCandidateService.CreateCheckList(taskCtx, predictionRace, predictionHorse, raceForecast),
		)

		predictionCheckLists = append(predictionCheckLists, predictionCheckList)
	}

	return predictionCheckLists, nil
}
