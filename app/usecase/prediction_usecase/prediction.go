package prediction_usecase

import (
	"context"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/analysis_service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/prediction_service"
	"github.com/sirupsen/logrus"
)

type Prediction interface {
	Odds(ctx context.Context, input *PredictionInput) error
	CheckList(ctx context.Context, input *PredictionInput) error
	Sync(ctx context.Context) error
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
	logger                          *logrus.Logger
}

func NewPrediction(
	predictionOddsService prediction_service.Odds,
	predictionPlaceCandidateService prediction_service.PlaceCandidate,
	predictionMarkerSyncService prediction_service.MarkerSync,
	placeService analysis_service.Place,
	logger *logrus.Logger,
) Prediction {
	return &prediction{
		predictionOddsService:           predictionOddsService,
		predictionPlaceCandidateService: predictionPlaceCandidateService,
		predictionMarkerSyncService:     predictionMarkerSyncService,
		placeService:                    placeService,
		logger:                          logger,
	}
}
