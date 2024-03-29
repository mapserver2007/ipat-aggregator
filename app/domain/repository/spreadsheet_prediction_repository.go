package repository

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type SpreadSheetPredictionRepository interface {
	Write(ctx context.Context, strictPredictionDataList, simplePredictionDataList []*spreadsheet_entity.PredictionData, markerOddsRangeMapList []map[types.Marker]*prediction_entity.OddsRange) error
	Style(ctx context.Context, markerOddsRangeMapList []map[types.Marker]*prediction_entity.OddsRange) error
	Clear(ctx context.Context) error
}
