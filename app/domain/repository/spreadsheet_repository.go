package repository

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
)

type SpreadSheetRepository interface {
	WriteSummary(ctx context.Context, summary *spreadsheet_entity.Summary) error
	WriteTicketSummary(ctx context.Context, ticketSummaryMap map[int]*spreadsheet_entity.TicketSummary) error
}
