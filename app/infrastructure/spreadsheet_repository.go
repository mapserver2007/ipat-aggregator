package infrastructure

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure/gateway"
)

type spreadSummeryRepository struct {
	summaryGateway       gateway.SpreadSheetSummaryGateway
	ticketSummaryGateway gateway.SpreadSheetTicketSummaryGateway
	listGateway          gateway.SpreadSheetListGateway
}

func NewSpreadSummeryRepository(
	summaryGateway gateway.SpreadSheetSummaryGateway,
	ticketSummaryGateway gateway.SpreadSheetTicketSummaryGateway,
	listGateway gateway.SpreadSheetListGateway,
) repository.SpreadSheetRepository {
	return &spreadSummeryRepository{
		summaryGateway:       summaryGateway,
		ticketSummaryGateway: ticketSummaryGateway,
		listGateway:          listGateway,
	}
}

func (s *spreadSummeryRepository) WriteSummary(
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

func (s *spreadSummeryRepository) WriteTicketSummary(
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

func (s *spreadSummeryRepository) WriteList(
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
