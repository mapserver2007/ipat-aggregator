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
	for _, race := range races {
		marker, ok := predictionMarkerMap[race.RaceId()]
		if !ok {
			return fmt.Errorf("raceId %v not found in predictionMarkerMap", race.RaceId())
		}

		strictFilterId, simpleFilterId := p.filterService.CreatePredictionFilters(ctx, race)
		strictPredictionMarkerCombinationData := p.spreadSheetService.CreateMarkerCombinationAnalysisData(ctx, analysisData, strictFilterId)
		strictPredictionOddsRangeMap := p.spreadSheetService.CreateOddsRangeRaceCountMap(ctx, analysisData, strictFilterId)
		simplePredictionMarkerCombinationData := p.spreadSheetService.CreateMarkerCombinationAnalysisData(ctx, analysisData, simpleFilterId)
		simplePredictionOddsRangeMap := p.spreadSheetService.CreateOddsRangeRaceCountMap(ctx, analysisData, simpleFilterId)

		markerOddsRangeMap := p.spreadSheetService.CreatePredictionOdds(ctx, marker, race)

		strictPredictionData := spreadsheet_entity.NewPredictionData(strictPredictionMarkerCombinationData, strictPredictionOddsRangeMap, race)
		simplePredictionData := spreadsheet_entity.NewPredictionData(simplePredictionMarkerCombinationData, simplePredictionOddsRangeMap, race)

		err := p.spreadSheetRepository.Write(ctx, strictPredictionData, simplePredictionData, markerOddsRangeMap, race)
		if err != nil {
			return err
		}
	}

	return nil
}
