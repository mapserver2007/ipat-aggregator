package repository

import (
	"context"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
)

type JockeyResultRepository interface {
	Write(ctx context.Context, path string, jockeyResultInfo *raw_entity.JockeyResultInfo) error
	FetchUrls(ctx context.Context, url string) ([]string, error)
	FetchJockeyResults(ctx context.Context, url string) ([]*netkeiba_entity.JockeyResult, error)
}
