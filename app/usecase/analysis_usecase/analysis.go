package analysis_usecase

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"log"
	"os"
	"path/filepath"
)

type analysis struct {
	markerDataRepository   repository.MarkerDataRepository
	predictAnalysisService service.AnalysisService
	ticketConverter        service.TicketConverter
}

func NewAnalysis(
	markerDataRepository repository.MarkerDataRepository,
	predictAnalysisService service.AnalysisService,
	ticketConverter service.TicketConverter,
) *analysis {
	return &analysis{
		markerDataRepository:   markerDataRepository,
		predictAnalysisService: predictAnalysisService,
		ticketConverter:        ticketConverter,
	}
}

func (p *analysis) Read(ctx context.Context) ([]*marker_csv_entity.Yamato, error) {
	rootPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	dirPath, err := filepath.Abs(rootPath + "/csv/markers")
	if err != nil {
		return nil, err
	}

	filePath := fmt.Sprintf("%s/%s", dirPath, "yamato_predict.csv")
	predicts, err := p.markerDataRepository.Read(ctx, filePath)
	if err != nil {
		return nil, err
	}

	return predicts, nil
}

func (p *analysis) CreateAnalysisData(
	ctx context.Context,
	records []*marker_csv_entity.Yamato,
	races []*data_cache_entity.Race,
	tickets []*ticket_csv_entity.Ticket,
	racingNumbers []*data_cache_entity.RacingNumber,
) (*analysis_entity.Layer1, error) {
	ticketsMap := p.ticketConverter.ConvertToRaceIdMap(ctx, tickets, racingNumbers, races)
	raceMap := map[types.RaceId]*data_cache_entity.Race{}
	for _, race := range races {
		raceMap[race.RaceId()] = race
	}

	for _, record := range records {
		race, ok := raceMap[record.RaceId()]
		if !ok {
			log.Println(ctx, fmt.Sprintf("raceId: %s is not found in predict_races", record.RaceId()))
			continue
		}
		ticketsByRaceId, ok := ticketsMap[record.RaceId()]
		if !ok {
			log.Println(ctx, fmt.Sprintf("raceId: %s is not found in tickets", record.RaceId()))
			continue
		}

		for _, payoutResult := range race.PayoutResults() {
			markerCombinationIds := p.predictAnalysisService.GetMarkerCombinationIds(ctx, payoutResult, record)
			for idx, markerCombinationId := range markerCombinationIds {
				var (
					payment types.Payment
					payout  types.Payout
				)
				for _, ticket := range ticketsByRaceId {
					if ticket.TicketType() == markerCombinationId.TicketType() && ticket.BetNumber() == payoutResult.Numbers()[idx] {
						payment = ticket.Payment()
						payout = ticket.Payout()
						break
					}
				}
				numerical := analysis_entity.NewCalculable(
					payment,
					payout,
					payoutResult.Odds()[idx],
					payoutResult.Numbers()[idx],
					payoutResult.Populars()[idx],
				)
				err := p.predictAnalysisService.AddAnalysisData(ctx, markerCombinationId, race, numerical)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return p.predictAnalysisService.GetAnalysisData(), nil
}
