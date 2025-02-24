package prediction_service

import (
	"context"
	"fmt"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type MarkerSync interface {
	GetRaceIds(ctx context.Context, raceDate types.RaceDate) ([]types.RaceId, error)
	GetMarkers(ctx context.Context, raceId types.RaceId) ([]*prediction_entity.Marker, error)
	Convert(ctx context.Context, predictionMarkers []*prediction_entity.Marker) []*spreadsheet_entity.PredictionMarker
	Write(ctx context.Context, predictionMarkers []*spreadsheet_entity.PredictionMarker) error
}

type markerSync struct {
	raceIdRepository      repository.RaceIdRepository
	raceRepository        repository.RaceRepository
	spreadSheetRepository repository.SpreadSheetRepository
}

func NewMarkerSync(
	raceIdRepository repository.RaceIdRepository,
	raceRepository repository.RaceRepository,
	spreadSheetRepository repository.SpreadSheetRepository,
) MarkerSync {
	return &markerSync{
		raceIdRepository:      raceIdRepository,
		raceRepository:        raceRepository,
		spreadSheetRepository: spreadSheetRepository,
	}
}

func (m *markerSync) GetRaceIds(
	ctx context.Context,
	raceDate types.RaceDate,
) ([]types.RaceId, error) {
	rawRaceIds, err := m.raceIdRepository.Fetch(ctx, fmt.Sprintf(raceListUrlForJRA, raceDate))
	if err != nil {
		return nil, err
	}
	if len(rawRaceIds) == 0 {
		return nil, fmt.Errorf("race ids not found: %d", raceDate)
	}

	raceIds := make([]types.RaceId, len(rawRaceIds))
	for i, rawRaceId := range rawRaceIds {
		raceIds[i] = types.RaceId(rawRaceId)
	}

	return raceIds, nil
}

func (m *markerSync) GetMarkers(
	ctx context.Context,
	raceId types.RaceId,
) ([]*prediction_entity.Marker, error) {
	markers, err := m.raceRepository.FetchMarker(ctx, fmt.Sprintf(raceMarkerUrl, raceId))
	if err != nil {
		return nil, err
	}
	if len(markers) != 6 {
		return nil, nil
	}

	predictionMarkers := make([]*prediction_entity.Marker, 0, len(markers))
	for _, marker := range markers {
		predictionMarkers = append(predictionMarkers, prediction_entity.NewMarker(
			raceId,
			marker.HorseNumber(),
			marker.Marker(),
		))
	}

	return predictionMarkers, nil
}

func (m *markerSync) Convert(
	ctx context.Context,
	predictionMarkers []*prediction_entity.Marker,
) []*spreadsheet_entity.PredictionMarker {
	predictionMarkerMap := map[types.RaceId][]types.HorseNumber{}
	for _, marker := range predictionMarkers {
		if _, ok := predictionMarkerMap[marker.RaceId()]; !ok {
			predictionMarkerMap[marker.RaceId()] = make([]types.HorseNumber, 6)
		}
		switch marker.Marker() {
		case types.Favorite:
			predictionMarkerMap[marker.RaceId()][0] = marker.HorseNumber()
		case types.Rival:
			predictionMarkerMap[marker.RaceId()][1] = marker.HorseNumber()
		case types.BrackTriangle:
			predictionMarkerMap[marker.RaceId()][2] = marker.HorseNumber()
		case types.WhiteTriangle:
			predictionMarkerMap[marker.RaceId()][3] = marker.HorseNumber()
		case types.Star:
			predictionMarkerMap[marker.RaceId()][4] = marker.HorseNumber()
		case types.Check:
			predictionMarkerMap[marker.RaceId()][5] = marker.HorseNumber()
		}
	}

	spreadSheetMarkers := make([]*spreadsheet_entity.PredictionMarker, 0, len(predictionMarkerMap))
	for _, raceId := range service.SortedRaceIdKeys(predictionMarkerMap) {
		markers := predictionMarkerMap[raceId]
		spreadSheetMarkers = append(spreadSheetMarkers, spreadsheet_entity.NewPredictionMarker(
			raceId,
			markers[0],
			markers[1],
			markers[2],
			markers[3],
			markers[4],
			markers[5],
		))
	}

	return spreadSheetMarkers
}

func (m *markerSync) Write(
	ctx context.Context,
	predictionMarkers []*spreadsheet_entity.PredictionMarker,
) error {
	return m.spreadSheetRepository.WritePredictionMarker(ctx, predictionMarkers)
}
