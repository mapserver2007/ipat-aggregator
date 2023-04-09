package service

import (
	"context"
	raw_jockey_entity "github.com/mapserver2007/ipat-aggregator/app/domain/jockey/raw_entity"
	raw_race_entity "github.com/mapserver2007/ipat-aggregator/app/domain/race/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/repository"
)

type RaceFetcher struct {
	raceClient repository.RaceClient
}

func NewRaceFetcher(
	raceClient repository.RaceClient,
) RaceFetcher {
	return RaceFetcher{
		raceClient: raceClient,
	}
}

func (f *RaceFetcher) FetchRacingNumbers(ctx context.Context, url string) ([]*raw_race_entity.RawRacingNumberNetkeiba, error) {
	rawRacingNumbers, err := f.raceClient.GetRacingNumbers(ctx, url)
	if err != nil {
		return nil, err
	}
	return rawRacingNumbers, nil
}

func (f *RaceFetcher) FetchRace(ctx context.Context, url string) (*raw_race_entity.RawRaceNetkeiba, error) {
	rawRace, err := f.raceClient.GetRaceResult(ctx, url)
	if err != nil {
		return nil, err
	}
	return rawRace, nil
}

func (f *RaceFetcher) FetchJockey(ctx context.Context, url string) (*raw_jockey_entity.RawJockeyNetkeiba, error) {
	rawJockey, err := f.raceClient.GetJockey(ctx, url)
	if err != nil {
		return nil, err
	}
	return rawJockey, nil
}
