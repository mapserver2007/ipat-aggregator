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
	summaryGateway       gateway.SpreadSheetSummaryGateway
	ticketSummaryGateway gateway.SpreadSheetTicketSummaryGateway
	listGateway          gateway.SpreadSheetListGateway
	analysisPlaceGateway gateway.SpreadSheetAnalysisPlaceGateway
}

func NewSpreadSheetRepository(
	summaryGateway gateway.SpreadSheetSummaryGateway,
	ticketSummaryGateway gateway.SpreadSheetTicketSummaryGateway,
	listGateway gateway.SpreadSheetListGateway,
	analysisPlaceGateway gateway.SpreadSheetAnalysisPlaceGateway,
) repository.SpreadSheetRepository {
	return &spreadSheetRepository{
		summaryGateway:       summaryGateway,
		ticketSummaryGateway: ticketSummaryGateway,
		listGateway:          listGateway,
		analysisPlaceGateway: analysisPlaceGateway,
	}
}

func (s *spreadSheetRepository) WriteSummary(
	ctx context.Context,
	summary *spreadsheet_entity.Summary,
) error {
	err := s.summaryGateway.Clear(ctx)
	if err != nil {
		return err
	}
	err = s.summaryGateway.Write(ctx, summary)
	if err != nil {
		return err
	}
	err = s.summaryGateway.Style(ctx, summary)
	if err != nil {
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
	thirdPlaceMap map[types.Marker]map[filter.Id]*spreadsheet_entity.AnalysisPlace,
	filters []filter.Id,
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
