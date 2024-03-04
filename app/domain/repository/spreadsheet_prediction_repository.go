package repository

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type SpreadSheetPredictionRepository interface {
	Write(ctx context.Context, strictPredictionData, simplePredictionData *spreadsheet_entity.PredictionData, markerOddsRangeMap map[types.Marker]types.OddsRangeType, race *prediction_entity.Race) error
	Style(ctx context.Context) error
	Clear(ctx context.Context) error
}
