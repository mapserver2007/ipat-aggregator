package repository

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
)

type TicketCsvRepository interface {
	Read(ctx context.Context, path string) ([]*ticket_csv_entity.Ticket, error)
}
