package repository

import (
	"context"
	betting_ticket_entity "github.com/mapserver2007/tools/baken/app/domain/betting_ticket/entity"
	"github.com/mapserver2007/tools/baken/app/domain/race/raw_entity"
)

type RaceClient interface {
	GetRacingNumbers(ctx context.Context, url string, entity *betting_ticket_entity.CsvEntity) ([]*raw_entity.RawRacingNumberNetkeiba, error)
	GetRaceResult(ctx context.Context, url string) (*raw_entity.RawRaceNetkeiba, error)
}
