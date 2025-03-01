package repository

import (
	"context"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
)

type RaceRepository interface {
	List(ctx context.Context, path string) ([]string, error)
	Read(ctx context.Context, path string) ([]*raw_entity.Race, error)
	Write(ctx context.Context, path string, raceInfo *raw_entity.RaceInfo) error
	FetchRace(ctx context.Context, url string) (*netkeiba_entity.Race, error)
	FetchRaceCard(ctx context.Context, url string) (*netkeiba_entity.Race, error)
	FetchMarker(ctx context.Context, url string) ([]*netkeiba_entity.Marker, error)
}
