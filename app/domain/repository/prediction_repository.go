package repository

import "context"

type PredictionRepository interface {
	FetchOdds(ctx context.Context, url string) error
	FetchRace(ctx context.Context, url string) error
	FetchRaceResult(ctx context.Context, url string) error
}
