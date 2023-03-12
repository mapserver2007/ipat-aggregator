package repository

import (
	"context"
	betting_ticket_entity "github.com/mapserver2007/tools/baken/app/domain/betting_ticket/entity"
	race_entity "github.com/mapserver2007/tools/baken/app/domain/race/entity"
	race_vo "github.com/mapserver2007/tools/baken/app/domain/race/value_object"
)

type RaceClient interface {
	GetRacingNumber(ctx context.Context, entity *betting_ticket_entity.CsvEntity) ([]*race_entity.RacingNumber, error)
	GetRaceResult(ctx context.Context, raceId race_vo.RaceId, entity *betting_ticket_entity.CsvEntity) *race_entity.Race
}
