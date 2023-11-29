package repository

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
)

type RacingNumberDataRepository interface {
	Read(ctx context.Context, fileName string) ([]*raw_entity.RacingNumber, error)
	Write(ctx context.Context) error
	Fetch(ctx context.Context, racingNumbers []*raw_entity.RacingNumber, tickets []*ticket_csv_entity.Ticket) error
}
