package prediction_usecase

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/config"
)

func (p *prediction) Sync(ctx context.Context) error {
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

	return nil
}
