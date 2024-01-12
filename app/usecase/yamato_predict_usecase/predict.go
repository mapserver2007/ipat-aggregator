package yamato_predict_usecase

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
)

const (
	startDate = "20240101"
	endDate   = "20240110"
)

type predict struct {
	netKeibaService  service.NetKeibaService
	raceIdRepository repository.RaceIdDataRepository
}

func NewPredict(
	netKeibaService service.NetKeibaService,
	raceIdRepository repository.RaceIdDataRepository,
) *predict {
	return &predict{
		netKeibaService:  netKeibaService,
		raceIdRepository: raceIdRepository,
	}
}

func (p *predict) Predict(ctx context.Context) error {
	// TODO いろいろ集計データを作る処理
	return nil
}

func (p *predict) Fetch(ctx context.Context) error {
	return nil
}
