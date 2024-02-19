package repository

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
)

type SpreadsheetListRepository interface {
	Write(ctx context.Context, rows []*spreadsheet_entity.Row) error
	Style(ctx context.Context, styles []*spreadsheet_entity.Style) error
	Clear(ctx context.Context) error
}
