package aggregation_usecase

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/aggregation_service"
)

type Summary interface {
	Execute(ctx context.Context, input *SummaryInput) error
}

type SummaryInput struct {
	Tickets []*ticket_csv_entity.RaceTicket
	Races   []*data_cache_entity.Race
}

type summary struct {
	summaryService aggregation_service.Summary
}

func NewSummary(
	summaryService aggregation_service.Summary,
) Summary {
	return &summary{
		summaryService: summaryService,
	}
}

func (a *summary) Execute(ctx context.Context, input *SummaryInput) error {
	entity := a.summaryService.Create(ctx, input.Tickets, input.Races)
	err := a.summaryService.Write(ctx, entity)
	if err != nil {
		return err
	}

	return nil
}
