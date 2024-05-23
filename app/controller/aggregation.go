package controller

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/aggregation_usecase"
)

type Aggregation struct {
	aggregationSummaryUseCase       aggregation_usecase.Summary
	aggregationTicketSummaryUseCase aggregation_usecase.TicketSummary
	aggregationListUseCase          aggregation_usecase.List
}

type AggregationInput struct {
	Master *MasterOutput
}

func NewAggregation(
	aggregationSummaryUseCase aggregation_usecase.Summary,
	aggregationTicketSummaryUseCase aggregation_usecase.TicketSummary,
	aggregationListUseCase aggregation_usecase.List,
) *Aggregation {
	return &Aggregation{
		aggregationSummaryUseCase:       aggregationSummaryUseCase,
		aggregationTicketSummaryUseCase: aggregationTicketSummaryUseCase,
		aggregationListUseCase:          aggregationListUseCase,
	}
}

func (a *Aggregation) Execute(ctx context.Context, input *AggregationInput) error {
	err := a.aggregationSummaryUseCase.Execute(ctx, &aggregation_usecase.SummaryInput{
		Tickets: input.Master.Tickets,
		Races:   input.Master.Races,
	})
	if err != nil {
		return err
	}

	err = a.aggregationTicketSummaryUseCase.Execute(ctx, &aggregation_usecase.TicketSummaryInput{
		Tickets: input.Master.Tickets,
	})
	if err != nil {
		return err
	}

	err = a.aggregationListUseCase.Execute(ctx, &aggregation_usecase.ListInput{
		Tickets: input.Master.Tickets,
		Races:   input.Master.Races,
		Jockeys: input.Master.Jockeys,
	})
	if err != nil {
		return err
	}

	return nil
}
