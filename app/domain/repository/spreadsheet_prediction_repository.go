package repository

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
)

type SpreadSheetPredictionRepository interface {
	Write(ctx context.Context, strictPredictionDataList, simplePredictionDataList []*spreadsheet_entity.PredictionData) error
	Style(ctx context.Context, predictionDataSize int) error
	Clear(ctx context.Context) error
}
