package master_service

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/config"
	"sort"
)

const (
	analysisMarkerFileName = "analysis_marker.csv"
)

type AnalysisMarker interface {
	Get(ctx context.Context) ([]*marker_csv_entity.AnalysisMarker, error)
}

type analysisMarkerService struct {
	analysisMarkerRepository repository.AnalysisMarkerRepository
}

func NewAnalysisMarker(
	analysisMarkerRepository repository.AnalysisMarkerRepository,
) AnalysisMarker {
	return &analysisMarkerService{
		analysisMarkerRepository: analysisMarkerRepository,
	}
}

func (m *analysisMarkerService) Get(
	ctx context.Context,
) ([]*marker_csv_entity.AnalysisMarker, error) {
	markers, err := m.analysisMarkerRepository.Read(ctx, fmt.Sprintf("%s/%s", config.CsvDir, analysisMarkerFileName))
	if err != nil {
		return nil, err
	}

	sort.Slice(markers, func(i, j int) bool {
		return markers[i].RaceId() < markers[j].RaceId()
	})

	return markers, nil
}
