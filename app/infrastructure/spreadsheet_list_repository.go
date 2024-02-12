package infrastructure

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/list_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"google.golang.org/api/sheets/v4"
)

const (
	spreadSheetListFileName2 = "spreadsheet_list.json"
)

type spreadSheetListRepository struct {
	client            *sheets.Service
	spreadSheetConfig *spreadsheet_entity.SpreadSheetConfig
}

func NewSpreadSheetListRepository() (repository.SpreadsheetListRepository, error) {
	ctx := context.Background()
	client, spreadSheetConfig, err := getSpreadSheetConfig(ctx, spreadSheetListFileName2)
	if err != nil {
		return nil, err
	}

	return &spreadSheetListRepository{
		client:            client,
		spreadSheetConfig: spreadSheetConfig,
	}, nil
}

func (s *spreadSheetListRepository) Write(ctx context.Context, rows []*list_entity.ListRow) error {
	return nil
}
