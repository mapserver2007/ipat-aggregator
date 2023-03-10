package repository

import (
	"context"
	betting_ticket_entity "github.com/mapserver2007/tools/baken/app/domain/betting_ticket/entity"
	race_entity "github.com/mapserver2007/tools/baken/app/domain/race/entity"
)

type RaceDB interface {
	ReadRaceResult(ctx context.Context, fileName string) (*race_entity.RaceInfo, error)
	ReadRacingNumber(ctx context.Context, fileName string) (*race_entity.RacingNumberInfo, error)
	UpdateRaceResult(ctx context.Context, fileName string, racingNumbers []*race_entity.RacingNumber, entities []*betting_ticket_entity.CsvEntity) error
	UpdateRacingNumber(ctx context.Context, fileName string, entities []*betting_ticket_entity.CsvEntity) error
}
