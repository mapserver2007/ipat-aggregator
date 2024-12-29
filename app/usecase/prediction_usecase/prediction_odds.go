package prediction_usecase

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
	"sort"
	"sync"
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
			for _, marker := range markers {
				select {
				case <-taskCtx.Done():
					return
				default:
					predictionRace, err := p.predictionOddsService.Get(taskCtx, marker.RaceId())
					if err != nil {
						select {
						case errorCh <- err: // 最初に発生したエラーをチャネルに送信
							cancel() // すべてのスレッドを停止する
						}
						return
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

	calculables, err := p.placeService.Create(ctx, input.AnalysisMarkers, input.Races)
	if err != nil {
		return err
	}

	sort.Slice(predictionRaces, func(i, j int) bool {
		return predictionRaces[i].RaceId() < predictionRaces[j].RaceId()
	})

	firstPlaceMap, secondPlaceMap, thirdPlaceMap, raceCourseMap := p.predictionOddsService.ConvertAll(ctx, predictionRaces, predictionMarkers, calculables)
	err = p.predictionOddsService.Write(ctx, firstPlaceMap, secondPlaceMap, thirdPlaceMap, raceCourseMap)
	if err != nil {
		return err
	}

	return nil
}
