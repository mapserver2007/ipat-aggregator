package list_usecase

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/list_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
)

type ListUseCase struct {
	listService service.ListService
}

func NewListUseCase(
	listService service.ListService,
) *ListUseCase {
	return &ListUseCase{
		listService: listService,
	}
}

func (p *ListUseCase) Read(
	ctx context.Context,
	tickets []*ticket_csv_entity.Ticket,
	racingNumbers []*data_cache_entity.RacingNumber,
	races []*data_cache_entity.Race,
	jockeys []*data_cache_entity.Jockey,
) ([]*list_entity.ListRow, error) {
	listRows, err := p.listService.Create(ctx, tickets, racingNumbers, races, jockeys)
	if err != nil {
		return nil, err
	}

	return listRows, nil
}
