package prediction_usecase

import (
	"context"
	"sort"
	"sync"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
)

const oddsParallel = 5

func (p *prediction) Odds(ctx context.Context, input *PredictionInput) error {
	taskCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	predictionMarkers := input.PredictionMarkers
	predictionRaces := make([]*prediction_entity.Race, 0, len(predictionMarkers))

	var wg sync.WaitGroup
	errorCh := make(chan error, 1)
	resultCh := make(chan []*prediction_entity.Race, oddsParallel)
	chunkSize := (len(predictionMarkers) + oddsParallel - 1) / oddsParallel

	for i := 0; i < len(predictionMarkers); i += chunkSize {
		end := i + chunkSize
		if end > len(predictionMarkers) {
			end = len(predictionMarkers)
		}

		wg.Add(1)
		go func(markers []*marker_csv_entity.PredictionMarker) {
			defer wg.Done()
			localPredictionRaces := make([]*prediction_entity.Race, 0, len(markers))
			p.logger.Infof("prediction odds processing: %v/%v", end, len(predictionMarkers))

			var localError error
			for _, marker := range markers {
				select {
				case <-taskCtx.Done():
					return
				default:
					if localError != nil {
						continue
					}
					predictionRace, err := p.predictionOddsService.Get(taskCtx, marker.RaceId())
					if err != nil {
						localError = err
						select {
						case errorCh <- err:
							cancel()
						default:
						}
						continue
					}
					localPredictionRaces = append(localPredictionRaces, predictionRace)
				}
			}

			resultCh <- localPredictionRaces
		}(predictionMarkers[i:end])
	}

	wg.Wait()
	close(errorCh)
	close(resultCh)

	if err := <-errorCh; err != nil {
		return err
	}

	for races := range resultCh {
		predictionRaces = append(predictionRaces, races...)
	}

	placeCalculables, err := p.placeService.Create(ctx, input.AnalysisMarkers, input.Races)
	if err != nil {
		return err
	}

	raceTimeCalculables, err := p.raceTimeService.Create(ctx, input.Races, input.RaceTimes)
	if err != nil {
		return err
	}
	analysisRaceTimeMap, _, _ := p.raceTimeService.Convert(ctx, raceTimeCalculables)

	sort.Slice(predictionRaces, func(i, j int) bool {
		return predictionRaces[i].RaceId() < predictionRaces[j].RaceId()
	})

	firstPlaceMap, secondPlaceMap, thirdPlaceMap, raceCourseMap := p.predictionOddsService.ConvertAll(ctx, predictionRaces, predictionMarkers, placeCalculables, analysisRaceTimeMap)
	err = p.predictionOddsService.Write(ctx, firstPlaceMap, secondPlaceMap, thirdPlaceMap, raceCourseMap)
	if err != nil {
		return err
	}

	return nil
}
