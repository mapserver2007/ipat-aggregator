package controller

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/aggregation_usecase"
)

type AggregationInput struct {
	Master *MasterOutput
}

type Aggregation struct {
	aggregationUseCase aggregation_usecase.Aggregation
}

func NewAggregation(
	aggregationUseCase aggregation_usecase.Aggregation,
) *Aggregation {
	return &Aggregation{
		aggregationUseCase: aggregationUseCase,
	}
}

func (a *Aggregation) Execute(ctx context.Context, input *AggregationInput) error {
	return a.aggregationUseCase.Execute(ctx, &aggregation_usecase.AggregationInput{
		Tickets: input.Master.Tickets,
		Races:   input.Master.Races,
	})
}
