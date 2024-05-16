package aggregation_usecase

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/aggregation_service"
)

type List interface {
	Execute(ctx context.Context, input *ListInput) error
}

type ListInput struct {
	Tickets []*ticket_csv_entity.RaceTicket
	Races   []*data_cache_entity.Race
	Jockeys []*data_cache_entity.Jockey
}

type list struct {
	listService aggregation_service.List
}

func NewList(
	listService aggregation_service.List,
) List {
	return &list{
		listService: listService,
	}
}

func (l *list) Execute(ctx context.Context, input *ListInput) error {
	listRows, err := l.listService.Create(ctx, input.Tickets, input.Races, input.Jockeys)
	if err != nil {
		return err
	}

	err = l.listService.Write(ctx, listRows)
	if err != nil {
		return err
	}

	return nil
}
