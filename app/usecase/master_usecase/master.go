package master_usecase

import (
	"context"
	"strconv"

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
	WinOdds           []*data_cache_entity.Odds
	PlaceOdds         []*data_cache_entity.Odds
	TrioOdds          []*data_cache_entity.Odds
	AnalysisMarkers   []*marker_csv_entity.AnalysisMarker
	PredictionMarkers []*marker_csv_entity.PredictionMarker
}

type master struct {
	ticketService           master_service.Ticket
	raceIdService           master_service.RaceId
	raceService             master_service.Race
	raceForecastService     master_service.RaceForecast
	jockeyService           master_service.Jockey
	winOddsService          master_service.WinOdds
	placeOddsService        master_service.PlaceOdds
	trioOddsService         master_service.TrioOdds
	analysisMarkerService   master_service.AnalysisMarker
	predictionMarkerService master_service.PredictionMarker
	umacaTicketService      master_service.UmacaTicket
}

func NewMaster(
	ticketService master_service.Ticket,
	raceIdService master_service.RaceId,
	raceService master_service.Race,
	raceForecastService master_service.RaceForecast,
	jockeyService master_service.Jockey,
	winOddsService master_service.WinOdds,
	placeOddsService master_service.PlaceOdds,
	trioOddsService master_service.TrioOdds,
	analysisMarkerService master_service.AnalysisMarker,
	predictionMarkerService master_service.PredictionMarker,
	umacaTicketService master_service.UmacaTicket,
) Master {
	return &master{
		ticketService:           ticketService,
		raceIdService:           raceIdService,
		raceService:             raceService,
		raceForecastService:     raceForecastService,
		jockeyService:           jockeyService,
		winOddsService:          winOddsService,
		placeOddsService:        placeOddsService,
		trioOddsService:         trioOddsService,
		analysisMarkerService:   analysisMarkerService,
		predictionMarkerService: predictionMarkerService,
		umacaTicketService:      umacaTicketService,
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

	winOdds, err := m.winOddsService.Get(ctx)
	if err != nil {
		return nil, err
	}

	placeOdds, err := m.placeOddsService.Get(ctx)
	if err != nil {
		return nil, err
	}

	trioOdds, err := m.trioOddsService.Get(ctx)
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
		WinOdds:           winOdds,
		PlaceOdds:         placeOdds,
		TrioOdds:          trioOdds,
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

	umacaMasters, err := m.umacaTicketService.GetMaster(ctx)
	if err != nil {
		return err
	}
	// umaca投票データに追加されたraceIdは未キャッシュなので更新する
	latestEndDate, err := types.NewRaceDate(input.EndDate)
	if err != nil {
		return err
	}
	for _, umacaMaster := range umacaMasters {
		if latestEndDate.Value() < umacaMaster.RaceDate().Value() {
			latestEndDate = umacaMaster.RaceDate()
		}
	}

	err = m.raceIdService.CreateOrUpdate(ctx, input.StartDate, strconv.Itoa(latestEndDate.Value()))
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

	races, err = m.raceService.Get(ctx)
	if err != nil {
		return err
	}

	err = m.umacaTicketService.CreateOrUpdate(ctx, races)
	if err != nil {
		return err
	}

	umacaRaceTickets, err := m.umacaTicketService.Get(ctx, races)
	if err != nil {
		return err
	}

	// 地方・海外のレースデータを取得するために馬券情報を取得する
	// 中央は期間から自動計算するが、地方・海外は馬券情報からRaceIdを割り出して取得する
	raceTickets, err := m.ticketService.Get(ctx, races)
	if err != nil {
		return err
	}

	// patとumacaのデータを統合
	raceTickets = append(raceTickets, umacaRaceTickets...)

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

	// tospoデータ差分更新
	races, err = m.raceService.Get(ctx)
	if err != nil {
		return err
	}

	//err = m.raceForecastService.CreateOrUpdate(ctx, races)
	//if err != nil {
	//	return err
	//}

	jockeys, excludeJockeyIds, err := m.jockeyService.Get(ctx)
	if err != nil {
		return err
	}

	err = m.jockeyService.CreateOrUpdate(ctx, jockeys, excludeJockeyIds)
	if err != nil {
		return err
	}

	winOdds, err := m.winOddsService.Get(ctx)
	if err != nil {
		return err
	}

	placeOdds, err := m.placeOddsService.Get(ctx)
	if err != nil {
		return err
	}

	trioOdds, err := m.trioOddsService.Get(ctx)
	if err != nil {
		return err
	}

	markers, err := m.analysisMarkerService.Get(ctx)
	if err != nil {
		return err
	}

	err = m.winOddsService.CreateOrUpdate(ctx, winOdds, markers)
	if err != nil {
		return err
	}

	err = m.placeOddsService.CreateOrUpdate(ctx, placeOdds, markers)
	if err != nil {
		return err
	}

	err = m.trioOddsService.CreateOrUpdate(ctx, trioOdds, markers)
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
