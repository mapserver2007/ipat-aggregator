package infrastructure

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"google.golang.org/api/sheets/v4"
)

const (
	spreadSheetTicketSummaryFileName = "spreadsheet_ticket_summary.json"
)

type spreadSheetTicketSummaryRepository struct {
	client            *sheets.Service
	spreadSheetConfig *spreadsheet_entity.SpreadSheetConfig
}

func NewSpreadSheetTicketSummaryRepository() (repository.SpreadSheetTicketSummaryRepository, error) {
	ctx := context.Background()
	client, spreadSheetConfig, err := getSpreadSheetConfig(ctx, spreadSheetTicketSummaryFileName)
	if err != nil {
		return nil, err
	}

	return &spreadSheetTicketSummaryRepository{
		client:            client,
		spreadSheetConfig: spreadSheetConfig,
	}, nil
}

func (s spreadSheetTicketSummaryRepository) Write(
	ctx context.Context,
	summary *spreadsheet_entity.Summary,
) error {
	//TODO implement me
	panic("implement me")
}

func (s spreadSheetTicketSummaryRepository) Style(ctx context.Context, summary *spreadsheet_entity.Summary) error {
	//TODO implement me
	panic("implement me")
}

func (s spreadSheetTicketSummaryRepository) Clear(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
