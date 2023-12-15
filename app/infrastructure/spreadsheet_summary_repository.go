package infrastructure

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
)

type spreadsheetSummaryRepository struct {
}

func NewSpreadSheetSummaryRepository() repository.SpreadSheetSummaryRepository {
	return &spreadsheetSummaryRepository{}
}

func (s *spreadsheetSummaryRepository) Write(
	ctx context.Context,
	summary *spreadsheet_entity.Summary,
) error {
	//TODO implement me
	fmt.Println(summary)
	return nil
}

func (s *spreadsheetSummaryRepository) Style(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
