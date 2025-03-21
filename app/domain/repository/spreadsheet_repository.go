package repository

import (
	"context"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
)

type SpreadSheetRepository interface {
	WriteSummaryV2(ctx context.Context, summary *spreadsheet_entity.Summary) error
	WriteTicketSummary(ctx context.Context, ticketSummaryMap map[int]*spreadsheet_entity.TicketSummary) error
	WriteList(ctx context.Context, listRows []*spreadsheet_entity.ListRow) error
	WriteAnalysisPlace(ctx context.Context,
		firstPlaceMap, secondPlaceMap, thirdPlaceMap map[types.Marker]map[filter.AttributeId]*spreadsheet_entity.AnalysisPlace,
		filters []filter.AttributeId,
	) error
	WriteAnalysisPlaceAllIn(ctx context.Context,
		placeAllInMap1 map[filter.AttributeId]*spreadsheet_entity.AnalysisPlaceAllIn,
		placeAllInMap2 map[filter.MarkerCombinationId]*spreadsheet_entity.AnalysisPlaceAllIn,
		attributeFilters []filter.AttributeId,
		markerCombinationFilters []filter.MarkerCombinationId,
	) error
	WriteAnalysisPlaceUnhit(ctx context.Context, analysisPlaceUnhits []*spreadsheet_entity.AnalysisPlaceUnhit) error
	WritePredictionOdds(ctx context.Context,
		firstPlaceMap, secondPlaceMap, thirdPlaceMap map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace,
		raceCourseMap map[types.RaceCourse][]types.RaceId,
	) error
	WritePredictionCheckList(ctx context.Context, predictionCheckLists []*spreadsheet_entity.PredictionCheckList) error
	WritePredictionMarker(ctx context.Context, predictionMarkers []*spreadsheet_entity.PredictionMarker) error
}
