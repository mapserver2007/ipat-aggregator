package aggregation_usecase

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/aggregation_service"
)

type Aggregation interface {
	Execute(ctx context.Context, input *AggregationInput) error
}

type AggregationInput struct {
	Tickets []*ticket_csv_entity.RaceTicket
	Races   []*data_cache_entity.Race
}

type aggregation struct {
	summaryService aggregation_service.Summary
}

func NewAggregation(
	summaryService aggregation_service.Summary,
) Aggregation {
	return &aggregation{
		summaryService: summaryService,
	}
}

func (a *aggregation) Execute(ctx context.Context, input *AggregationInput) error {
	summary := a.summaryService.Create(ctx, input.Tickets, input.Races)
	err := a.summaryService.Write(ctx, summary)
	if err != nil {
		return err
	}

	return nil
}
