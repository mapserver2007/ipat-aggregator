package data_cache_usecase

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
)

const (
	racingNumberFileName = "racing_number.json"
	raceResultFileName   = "race_result.json"
	jockeyFileName       = "jockey.json"
)

type dataCacheUseCase struct {
	racingNumberRepository repository.RacingNumberDataRepository
	raceDataRepository     repository.RaceDataRepository
	jockeyDataRepository   repository.JockeyDataRepository
}

func NewDataCacheUseCase(
	racingNumberRepository repository.RacingNumberDataRepository,
	raceDataRepository repository.RaceDataRepository,
	jockeyDataRepository repository.JockeyDataRepository,
) *dataCacheUseCase {
	return &dataCacheUseCase{
		racingNumberRepository: racingNumberRepository,
		raceDataRepository:     raceDataRepository,
		jockeyDataRepository:   jockeyDataRepository,
	}
}

func (d *dataCacheUseCase) Read(ctx context.Context) ([]*raw_entity.RacingNumber, []*raw_entity.Race, []*raw_entity.Jockey, []int, error) {
	rawRacingNumbers, err := d.racingNumberRepository.Read(ctx, racingNumberFileName)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	rawRaces, err := d.raceDataRepository.Read(ctx, raceResultFileName)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	rawJockeys, excludeJockeyIds, err := d.jockeyDataRepository.Read(ctx, jockeyFileName)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return rawRacingNumbers, rawRaces, rawJockeys, excludeJockeyIds, nil
}

func (d *dataCacheUseCase) Write(
	ctx context.Context,
	tickets []*ticket_csv_entity.Ticket,
	rawRacingNumbers []*raw_entity.RacingNumber,
	rawRaces []*raw_entity.Race,
	rawJockeys []*raw_entity.Jockey,
	excludeJockeyIds []int,
) error {

	d.racingNumberRepository.Fetch(ctx, rawRacingNumbers, tickets)

	//TODO implement me
	return nil
}
