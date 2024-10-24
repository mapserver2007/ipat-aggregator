package infrastructure

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure/gateway"
)

type horseRepository struct {
	netKeibaGateway gateway.NetKeibaGateway
}

func NewHorseRepository(
	netKeibaGateway gateway.NetKeibaGateway,
) repository.HorseRepository {
	return &horseRepository{
		netKeibaGateway: netKeibaGateway,
	}
}

func (h *horseRepository) Fetch(
	ctx context.Context,
	url string,
) (*netkeiba_entity.Horse, error) {
	horse, err := h.netKeibaGateway.FetchHorse(ctx, url)
	if err != nil {
		return nil, err
	}

	return horse, nil
}
