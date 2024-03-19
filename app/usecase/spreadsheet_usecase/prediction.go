package spreadsheet_usecase

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type PredictionUseCase struct {
	spreadSheetRepository repository.SpreadSheetPredictionRepository
	filterService         service.FilterService
	spreadSheetService    service.SpreadSheetService
}

func NewPredictionUseCase(
	spreadSheetRepository repository.SpreadSheetPredictionRepository,
	filterService service.FilterService,
	spreadSheetService service.SpreadSheetService,
) *PredictionUseCase {
	return &PredictionUseCase{
		spreadSheetRepository: spreadSheetRepository,
		filterService:         filterService,
		spreadSheetService:    spreadSheetService,
	}
}

func (p *PredictionUseCase) Write(
	ctx context.Context,
	races []*prediction_entity.Race,
	predictionMarkerMap map[types.RaceId]*marker_csv_entity.PredictionMarker,
	analysisData *analysis_entity.Layer1,
) error {
	err := p.spreadSheetRepository.Clear(ctx)
	if err != nil {
		return err
	}

	var (
		strictPredictionDataList, simplePredictionDataList []*spreadsheet_entity.PredictionData
		markerOddsRangeMapList                             []map[types.Marker]*prediction_entity.OddsRange
	)
	for _, race := range races {
		marker, ok := predictionMarkerMap[race.RaceId()]
		if !ok {
			return fmt.Errorf("raceId %v not found in predictionMarkerMap", race.RaceId())
		}

		strictFilterId, simpleFilterId := p.filterService.CreatePredictionFilters(ctx, race)
		strictPredictionMarkerCombinationData := p.spreadSheetService.CreateMarkerCombinationAnalysisData(ctx, analysisData, strictFilterId)
		strictPredictionOddsRangeMap := p.spreadSheetService.CreateOddsRangeCountMap(ctx, analysisData, strictFilterId)
		simplePredictionMarkerCombinationData := p.spreadSheetService.CreateMarkerCombinationAnalysisData(ctx, analysisData, simpleFilterId)
		simplePredictionOddsRangeMap := p.spreadSheetService.CreateOddsRangeCountMap(ctx, analysisData, simpleFilterId)

		markerOddsRangeMap := p.spreadSheetService.CreatePredictionOdds(ctx, marker, race)
		markerOddsRangeMapList = append(markerOddsRangeMapList, markerOddsRangeMap)

		strictPredictionDataList = append(strictPredictionDataList,
			spreadsheet_entity.NewPredictionData(strictPredictionMarkerCombinationData, strictPredictionOddsRangeMap, race, strictFilterId))
		simplePredictionDataList = append(simplePredictionDataList,
			spreadsheet_entity.NewPredictionData(simplePredictionMarkerCombinationData, simplePredictionOddsRangeMap, race, simpleFilterId))
	}

	err = p.spreadSheetRepository.Write(ctx, strictPredictionDataList, simplePredictionDataList, markerOddsRangeMapList)
	if err != nil {
		return err
	}

	err = p.spreadSheetRepository.Style(ctx, markerOddsRangeMapList)
	if err != nil {
		return err
	}

	return nil
}
