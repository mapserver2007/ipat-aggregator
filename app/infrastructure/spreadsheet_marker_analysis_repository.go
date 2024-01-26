package infrastructure

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"google.golang.org/api/sheets/v4"
	"log"
)

const (
	spreadSheetMarkerAnalysisFileName = "spreadsheet_marker_analysis.json"
)

type spreadSheetMarkerAnalysisRepository struct {
	client            *sheets.Service
	spreadSheetConfig *spreadsheet_entity.SpreadSheetConfig
}

func NewSpreadSheetMarkerAnalysisRepository() (repository.SpreadSheetMarkerAnalysisRepository, error) {
	ctx := context.Background()
	client, spreadSheetConfig, err := getSpreadSheetConfig(ctx, spreadSheetMarkerAnalysisFileName)
	if err != nil {
		return nil, err
	}

	return &spreadSheetMarkerAnalysisRepository{
		client:            client,
		spreadSheetConfig: spreadSheetConfig,
	}, nil
}

func (s *spreadSheetMarkerAnalysisRepository) Write(
	ctx context.Context,
) error {
	log.Println(ctx, "write marker analysis start")

	log.Println(ctx, "write marker analysis end")

	return nil
}
