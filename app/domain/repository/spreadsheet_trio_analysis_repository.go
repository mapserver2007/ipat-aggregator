package repository

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
)

type SpreadSheetTrioAnalysisRepository interface {
	Write(ctx context.Context, analysisData *spreadsheet_entity.AnalysisData, filters []filter.Id) error
	Style(ctx context.Context, analysisData *spreadsheet_entity.AnalysisData, filters []filter.Id) error
	Clear(ctx context.Context) error
}
