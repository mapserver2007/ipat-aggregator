package infrastructure

import (
	"context"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure/gateway"
)

type trainerRepository struct {
	netKeibaGateway gateway.NetKeibaGateway
}

func NewTrainerRepository(
	netKeibaGateway gateway.NetKeibaGateway,
) repository.TrainerRepository {
	return &trainerRepository{
		netKeibaGateway: netKeibaGateway,
	}
}

func (t *trainerRepository) Fetch(
	ctx context.Context,
	url string,
) (*netkeiba_entity.Trainer, error) {
	trainer, err := t.netKeibaGateway.FetchTrainer(ctx, url)
	if err != nil {
		return nil, err
	}

	return trainer, nil
}
