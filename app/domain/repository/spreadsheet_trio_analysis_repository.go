package repository

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
)

type SpreadSheetTrioAnalysisRepository interface {
	Write(ctx context.Context, analysisData *spreadsheet_entity.AnalysisData, races []*data_cache_entity.Race, odds []*data_cache_entity.Odds, markers []*marker_csv_entity.AnalysisMarker) error
	Style(ctx context.Context, analysisData *spreadsheet_entity.AnalysisData) error
	Clear(ctx context.Context) error
}
