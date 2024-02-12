package repository

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/list_entity"
)

type SpreadsheetListRepository interface {
	Write(ctx context.Context, rows []*list_entity.ListRow) error
}
