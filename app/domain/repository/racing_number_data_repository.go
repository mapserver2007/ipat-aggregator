package repository

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
)

type RacingNumberDataRepository interface {
	Read(ctx context.Context, fileName string) ([]*raw_entity.RacingNumber, error)
	Write(ctx context.Context, fileName string, racingNumberInfo *raw_entity.RacingNumberInfo) error
	Fetch(ctx context.Context, url string) ([]*netkeiba_entity.RacingNumber, error)
}
