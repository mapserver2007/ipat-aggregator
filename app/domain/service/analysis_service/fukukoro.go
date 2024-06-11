package analysis_service

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/shopspring/decimal"
)

type Fukukoro interface {
	Create(ctx context.Context, races []*data_cache_entity.Race, odds decimal.Decimal) error
}

type fukukoroService struct {
}

func NewFukukoro() Fukukoro {
	return &fukukoroService{}
}

func (f *fukukoroService) Create(
	ctx context.Context,
	races []*data_cache_entity.Race,
	odds decimal.Decimal,
) error {
	// raceIdに対する、odds倍以下の複勝オッズに該当するデータを抽出
	// 抽出項目：複勝オッズ、単勝オッズ

	return nil
}
