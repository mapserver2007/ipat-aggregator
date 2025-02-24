package repository

import (
	"context"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/umaca_csv_entity"
)

type UmacaTicketRepository interface {
	GetMaster(ctx context.Context, path string) ([]*umaca_csv_entity.UmacaMaster, error)
	List(ctx context.Context, path string) ([]string, error)
	Write(ctx context.Context, path string, data [][]string) error
}
