package repository

import (
	"context"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
)

type JockeyResultRepository interface {
	List(ctx context.Context, path string) ([]string, error)
	Read(ctx context.Context, path string) ([]*raw_entity.JockeyResult, error)
	Write(ctx context.Context, path string, jockeyResultInfo *raw_entity.JockeyResultInfo) error
	FetchJockeyResultCounts(ctx context.Context, url string) ([]*raw_entity.JockeyResultCount, error)
	FetchJockeyResultUrls(ctx context.Context, url string) ([]string, error)
	FetchJockeyResults(ctx context.Context, url string) ([]*netkeiba_entity.JockeyResult, error)
}
