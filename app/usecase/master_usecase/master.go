package master_usecase

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/master_service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type Master interface {
	Get(ctx context.Context) (*MasterOutput, error)
	CreateOrUpdate(ctx context.Context, input *MasterInput) error
}

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

type master struct {
	ticketService           master_service.Ticket
	raceIdService           master_service.RaceId
	raceService             master_service.Race
	jockeyService           master_service.Jockey
	trioOddsService         master_service.TrioOdds
	analysisMarkerService   master_service.AnalysisMarker
	predictionMarkerService master_service.PredictionMarker
}

func NewMaster(
	ticketService master_service.Ticket,
	raceIdService master_service.RaceId,
	raceService master_service.Race,
	jockeyService master_service.Jockey,
	trioOddsService master_service.TrioOdds,
	analysisMarkerService master_service.AnalysisMarker,
	predictionMarkerService master_service.PredictionMarker,
) Master {
	return &master{
		ticketService:           ticketService,
		raceIdService:           raceIdService,
		raceService:             raceService,
		jockeyService:           jockeyService,
		trioOddsService:         trioOddsService,
		analysisMarkerService:   analysisMarkerService,
		predictionMarkerService: predictionMarkerService,
	}
}

func (m *master) Get(ctx context.Context) (*MasterOutput, error) {
	races, err := m.raceService.Get(ctx)
	if err != nil {
		return nil, err
	}

	raceTickets, err := m.ticketService.Get(ctx, races)
	if err != nil {
		return nil, err
	}

	jockeys, _, err := m.jockeyService.Get(ctx)
	if err != nil {
		return nil, err
	}

	odds, err := m.trioOddsService.Get(ctx)
	if err != nil {
		return nil, err
	}

	analysisMarkers, err := m.analysisMarkerService.Get(ctx)
	if err != nil {
		return nil, err
	}

	predictionMarkers, err := m.predictionMarkerService.Get(ctx)
	if err != nil {
		return nil, err
	}

	return &MasterOutput{
		Tickets:           raceTickets,
		Races:             races,
		Jockeys:           jockeys,
		Odds:              odds,
		AnalysisMarkers:   analysisMarkers,
		PredictionMarkers: predictionMarkers,
	}, nil
}

func (m *master) CreateOrUpdate(ctx context.Context, input *MasterInput) error {
	err := m.raceIdService.CreateOrUpdate(ctx, input.StartDate, input.EndDate)
	if err != nil {
		return err
	}

	raceDateMap, _, err := m.raceIdService.Get(ctx)
	if err != nil {
		return err
	}

	races, err := m.raceService.Get(ctx)
	if err != nil {
		return err
	}

	err = m.raceService.CreateOrUpdate(ctx, races, raceDateMap)
	if err != nil {
		return err
	}

	// 地方・海外のレースデータを取得するために馬券情報を取得する
	// 中央は期間から自動計算するが、地方・海外は馬券情報からRaceIdを割り出して取得する
	raceTickets, err := m.ticketService.Get(ctx, races)
	if err != nil {
		return err
	}

	raceDateMapForNAROrOversea := map[types.RaceDate][]types.RaceId{}
	for _, raceTicket := range raceTickets {
		if raceTicket.Ticket().RaceCourse().NAR() || raceTicket.Ticket().RaceCourse().Oversea() {
			if _, ok := raceDateMapForNAROrOversea[raceTicket.Ticket().RaceDate()]; !ok {
				raceDateMapForNAROrOversea[raceTicket.Ticket().RaceDate()] = make([]types.RaceId, 0)
			}
			// 全馬券をもとに処理しているので同一raceIdが入る
			raceDateMapForNAROrOversea[raceTicket.Ticket().RaceDate()] = append(raceDateMapForNAROrOversea[raceTicket.Ticket().RaceDate()], raceTicket.RaceId())
		}
	}

	for raceDate, raceIds := range raceDateMapForNAROrOversea {
		m.uniqueSlice(&raceIds)
		raceDateMapForNAROrOversea[raceDate] = raceIds
	}

	err = m.raceIdService.Update(ctx, raceDateMapForNAROrOversea)
	if err != nil {
		return err
	}

	raceDateMap, _, err = m.raceIdService.Get(ctx)
	if err != nil {
		return err
	}

	races, err = m.raceService.Get(ctx)
	if err != nil {
		return err
	}

	err = m.raceService.CreateOrUpdate(ctx, races, raceDateMap)
	if err != nil {
		return err
	}

	jockeys, excludeJockeyIds, err := m.jockeyService.Get(ctx)
	if err != nil {
		return err
	}

	err = m.jockeyService.CreateOrUpdate(ctx, jockeys, excludeJockeyIds)
	if err != nil {
		return err
	}

	odds, err := m.trioOddsService.Get(ctx)
	if err != nil {
		return err
	}

	markers, err := m.analysisMarkerService.Get(ctx)
	if err != nil {
		return err
	}

	err = m.trioOddsService.CreateOrUpdate(ctx, odds, markers)
	if err != nil {
		return err
	}

	return nil
}

func (m *master) uniqueSlice(slice *[]types.RaceId) {
	seen := make(map[types.RaceId]bool)
	j := 0
	for i := 0; i < len(*slice); i++ {
		if !seen[(*slice)[i]] {
			seen[(*slice)[i]] = true
			(*slice)[j] = (*slice)[i]
			j++
		}
	}
	*slice = (*slice)[:j]
}
