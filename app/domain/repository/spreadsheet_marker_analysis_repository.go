package repository

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
)

type SpreadSheetMarkerAnalysisRepository interface {
	Write(ctx context.Context, analysisData *spreadsheet_entity.AnalysisData) error
	Style(ctx context.Context, analysisData *spreadsheet_entity.AnalysisData) error
	Clear(ctx context.Context) error
}
