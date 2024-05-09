package controller

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/master_usecase"
)

type MasterInput struct {
	StartDate string
	EndDate   string
}

type MasterOutput struct {
	Tickets           []*ticket_csv_entity.RaceTicket
	Races             []*data_cache_entity.Race
	Jockeys           []*data_cache_entity.Jockey
	Odds              []*data_cache_entity.Odds
	AnalysisMarkers   []*marker_csv_entity.AnalysisMarker
	PredictionMarkers []*marker_csv_entity.PredictionMarker
}

type Master struct {
	masterUseCase master_usecase.Master
}

func NewMaster(masterUseCase master_usecase.Master) *Master {
	return &Master{
		masterUseCase: masterUseCase,
	}
}

func (m *Master) Execute(ctx context.Context, input *MasterInput) (*MasterOutput, error) {
	err := m.masterUseCase.CreateOrUpdate(ctx, &master_usecase.MasterInput{
		StartDate: input.StartDate,
		EndDate:   input.EndDate,
	})
	if err != nil {
		return nil, err
	}

	output, err := m.masterUseCase.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &MasterOutput{
		Tickets:           output.Tickets,
		Races:             output.Races,
		Jockeys:           output.Jockeys,
		Odds:              output.Odds,
		AnalysisMarkers:   output.AnalysisMarkers,
		PredictionMarkers: output.PredictionMarkers,
	}, nil
}
