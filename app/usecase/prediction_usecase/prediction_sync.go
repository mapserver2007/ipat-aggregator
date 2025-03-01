package prediction_usecase

import (
	"context"
	"sync"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/config"
)

const markerSyncParallel = 3

func (p *prediction) Sync(ctx context.Context) error {
	taskCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	raceDate, err := types.NewRaceDate(config.PredictionSyncRaceDate)
	if err != nil {
		return err
	}

	raceIds, err := p.predictionMarkerSyncService.GetRaceIds(ctx, raceDate)
	if err != nil {
		return err
	}

	var predictionMarkers []*prediction_entity.Marker

	var wg sync.WaitGroup
	errorCh := make(chan error, 1)
	resultCh := make(chan []*prediction_entity.Marker, markerSyncParallel)
	chunkSize := (len(raceIds) + markerSyncParallel - 1) / markerSyncParallel

	for i := 0; i < len(raceIds); i += chunkSize {
		end := i + chunkSize
		if end > len(raceIds) {
			end = len(raceIds)
		}

		wg.Add(1)
		go func(splitRaceIds []types.RaceId) {
			defer wg.Done()
			localPredictionMarkers := make([]*prediction_entity.Marker, 0, len(splitRaceIds))
			p.logger.Infof("prediction marker sync processing: %v/%v", end, len(raceIds))
			for _, raceId := range splitRaceIds {
				select {
				case <-taskCtx.Done():
					return
				default:
					markers, err := p.predictionMarkerSyncService.GetMarkers(ctx, raceId)
					if err != nil {
						select {
						case errorCh <- err:
							cancel()
						}
						return
					}

					localPredictionMarkers = append(localPredictionMarkers, markers...)
				}
			}

			resultCh <- localPredictionMarkers
		}(raceIds[i:end])
	}

	wg.Wait()
	close(errorCh)
	close(resultCh)

	if err := <-errorCh; err != nil {
		return err
	}

	for markers := range resultCh {
		predictionMarkers = append(predictionMarkers, markers...)
	}

	spreadSheetPredictionMarkers := p.predictionMarkerSyncService.Convert(ctx, predictionMarkers)
	err = p.predictionMarkerSyncService.Write(ctx, spreadSheetPredictionMarkers)
	if err != nil {
		return err
	}

	return nil
}
