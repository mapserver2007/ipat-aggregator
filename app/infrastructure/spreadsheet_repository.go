package infrastructure

import (
	"context"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure/gateway"
)

type spreadSheetRepository struct {
	summaryGateway             gateway.SpreadSheetSummaryGateway
	ticketSummaryGateway       gateway.SpreadSheetTicketSummaryGateway
	listGateway                gateway.SpreadSheetListGateway
	analysisPlaceGateway       gateway.SpreadSheetAnalysisPlaceGateway
	analysisPlaceAllInGateway  gateway.SpreadSheetAnalysisPlaceAllInGateway
	analysisPlaceUnhitGateway  gateway.SpreadSheetAnalysisPlaceUnhitGateway
	predictionOddsGateway      gateway.SpreadSheetPredictionOddsGateway
	predictionCheckListGateway gateway.SpreadSheetPredictionCheckListGateway
	predictionMarkerGateway    gateway.SpreadSheetPredictionMarkerGateway
}

func NewSpreadSheetRepository(
	summaryGateway gateway.SpreadSheetSummaryGateway,
	ticketSummaryGateway gateway.SpreadSheetTicketSummaryGateway,
	listGateway gateway.SpreadSheetListGateway,
	analysisPlaceGateway gateway.SpreadSheetAnalysisPlaceGateway,
	analysisPlaceAllInGateway gateway.SpreadSheetAnalysisPlaceAllInGateway,
	analysisPlaceUnhitGateway gateway.SpreadSheetAnalysisPlaceUnhitGateway,
	predictionOddsGateway gateway.SpreadSheetPredictionOddsGateway,
	predictionCheckListGateway gateway.SpreadSheetPredictionCheckListGateway,
	predictionMarkerGateway gateway.SpreadSheetPredictionMarkerGateway,
) repository.SpreadSheetRepository {
	return &spreadSheetRepository{
		summaryGateway:             summaryGateway,
		ticketSummaryGateway:       ticketSummaryGateway,
		listGateway:                listGateway,
		analysisPlaceGateway:       analysisPlaceGateway,
		analysisPlaceAllInGateway:  analysisPlaceAllInGateway,
		analysisPlaceUnhitGateway:  analysisPlaceUnhitGateway,
		predictionOddsGateway:      predictionOddsGateway,
		predictionCheckListGateway: predictionCheckListGateway,
		predictionMarkerGateway:    predictionMarkerGateway,
	}
}

func (s *spreadSheetRepository) WriteSummaryV2(
	ctx context.Context,
	summary *spreadsheet_entity.Summary,
) error {
	if err := s.summaryGateway.ClearV2(ctx); err != nil {
		return err
	}
	if err := s.summaryGateway.WriteV2(ctx, summary); err != nil {
		return err
	}
	if err := s.summaryGateway.StyleV2(ctx, summary); err != nil {
		return err
	}

	return nil
}

func (s *spreadSheetRepository) WriteTicketSummary(
	ctx context.Context,
	ticketSummaryMap map[int]*spreadsheet_entity.TicketSummary,
) error {
	err := s.ticketSummaryGateway.Clear(ctx)
	if err != nil {
		return err
	}
	err = s.ticketSummaryGateway.Write(ctx, ticketSummaryMap)
	if err != nil {
		return err
	}
	err = s.ticketSummaryGateway.Style(ctx, ticketSummaryMap)
	if err != nil {
		return err
	}

	return nil
}

func (s *spreadSheetRepository) WriteList(
	ctx context.Context,
	listRows []*spreadsheet_entity.ListRow,
) error {
	err := s.listGateway.Clear(ctx)
	if err != nil {
		return err
	}

	err = s.listGateway.Write(ctx, listRows)
	if err != nil {
		return err
	}

	err = s.listGateway.Style(ctx, listRows)
	if err != nil {
		return err
	}

	return nil
}

func (s *spreadSheetRepository) WriteAnalysisPlace(
	ctx context.Context,
	firstPlaceMap,
	secondPlaceMap,
	thirdPlaceMap map[types.Marker]map[filter.AttributeId]*spreadsheet_entity.AnalysisPlace,
	filters []filter.AttributeId,
) error {
	err := s.analysisPlaceGateway.Clear(ctx)
	if err != nil {
		return err
	}

	err = s.analysisPlaceGateway.Write(ctx, firstPlaceMap, secondPlaceMap, thirdPlaceMap, filters)
	if err != nil {
		return err
	}

	err = s.analysisPlaceGateway.Style(ctx, firstPlaceMap, secondPlaceMap, thirdPlaceMap, filters)
	if err != nil {
		return err
	}

	return nil
}

func (s *spreadSheetRepository) WriteAnalysisPlaceAllIn(
	ctx context.Context,
	placeAllInMap1 map[filter.AttributeId]*spreadsheet_entity.AnalysisPlaceAllIn,
	placeAllInMap2 map[filter.MarkerCombinationId]*spreadsheet_entity.AnalysisPlaceAllIn,
	attributeFilters []filter.AttributeId,
	markerCombinationFilters []filter.MarkerCombinationId,
) error {
	err := s.analysisPlaceAllInGateway.Clear(ctx)
	if err != nil {
		return err
	}

	err = s.analysisPlaceAllInGateway.Write(ctx, placeAllInMap1, placeAllInMap2, attributeFilters, markerCombinationFilters)
	if err != nil {
		return err
	}

	err = s.analysisPlaceAllInGateway.Style(ctx, placeAllInMap1, placeAllInMap2, attributeFilters, markerCombinationFilters)
	if err != nil {
		return err
	}

	return nil
}

func (s *spreadSheetRepository) WriteAnalysisPlaceUnhit(
	ctx context.Context,
	analysisPlaceUnhits []*spreadsheet_entity.AnalysisPlaceUnhit,
) error {
	// err := s.analysisPlaceUnhitGateway.Clear(ctx)
	// if err != nil {
	// 	return err
	// }

	err := s.analysisPlaceUnhitGateway.Write(ctx, analysisPlaceUnhits)
	if err != nil {
		return err
	}

	return nil
}

func (s *spreadSheetRepository) WritePredictionOdds(
	ctx context.Context,
	firstPlaceMap,
	secondPlaceMap,
	thirdPlaceMap map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace,
	raceCourseMap map[types.RaceCourse][]types.RaceId,
) error {
	err := s.predictionOddsGateway.Clear(ctx)
	if err != nil {
		return err
	}

	err = s.predictionOddsGateway.Write(ctx, firstPlaceMap, secondPlaceMap, thirdPlaceMap, raceCourseMap)
	if err != nil {
		return err
	}

	err = s.predictionOddsGateway.Style(ctx, firstPlaceMap, secondPlaceMap, thirdPlaceMap, raceCourseMap)
	if err != nil {
		return err
	}

	return nil
}

func (s *spreadSheetRepository) WritePredictionCheckList(
	ctx context.Context,
	predictionCheckLists []*spreadsheet_entity.PredictionCheckList,
) error {
	err := s.predictionCheckListGateway.Clear(ctx)
	if err != nil {
		return err
	}

	err = s.predictionCheckListGateway.Write(ctx, predictionCheckLists)
	if err != nil {
		return err
	}

	err = s.predictionCheckListGateway.Style(ctx, predictionCheckLists)
	if err != nil {
		return err
	}

	return nil
}

func (s *spreadSheetRepository) WritePredictionMarker(
	ctx context.Context,
	predictionMarkers []*spreadsheet_entity.PredictionMarker,
) error {
	err := s.predictionMarkerGateway.Clear(ctx)
	if err != nil {
		return err
	}

	err = s.predictionMarkerGateway.Write(ctx, predictionMarkers)
	if err != nil {
		return err
	}

	return nil
}
