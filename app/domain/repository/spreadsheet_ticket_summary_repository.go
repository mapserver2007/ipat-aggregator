package repository

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
)

type SpreadSheetTicketSummaryRepository interface {
	Write(ctx context.Context, summary *spreadsheet_entity.Summary) error
	Style(ctx context.Context, summary *spreadsheet_entity.Summary) error
	Clear(ctx context.Context) error
}
