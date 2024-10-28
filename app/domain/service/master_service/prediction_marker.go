package master_service

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/config"
)

const (
	predictionMarkerFileName = "prediction_marker.csv"
)

type PredictionMarker interface {
	Get(ctx context.Context) ([]*marker_csv_entity.PredictionMarker, error)
	Push(ctx context.Context) error
}

type predictionMarkerService struct {
	predictionMarkerRepository repository.PredictionMarkerRepository
}

func NewPredictionMarker(
	predictionMarkerRepository repository.PredictionMarkerRepository,
) PredictionMarker {
	return &predictionMarkerService{
		predictionMarkerRepository: predictionMarkerRepository,
	}
}

func (p *predictionMarkerService) Get(
	ctx context.Context,
) ([]*marker_csv_entity.PredictionMarker, error) {
	markers, err := p.predictionMarkerRepository.Read(ctx, fmt.Sprintf("%s/%s", config.CsvDir, predictionMarkerFileName))
	if err != nil {
		return nil, err
	}

	return markers, nil
}

func (p *predictionMarkerService) Push(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
