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
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type analysis struct {
	markerDataRepository repository.MarkerDataRepository
	analysisService      service.AnalysisService
	ticketConverter      service.TicketConverter
}

func NewAnalysis(
	markerDataRepository repository.MarkerDataRepository,
	analysisService service.AnalysisService,
	ticketConverter service.TicketConverter,
) *analysis {
	return &analysis{
		markerDataRepository: markerDataRepository,
		analysisService:      analysisService,
		ticketConverter:      ticketConverter,
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
	markers []*marker_csv_entity.Yamato,
	races []*data_cache_entity.Race,
	tickets []*ticket_csv_entity.Ticket,
	racingNumbers []*data_cache_entity.RacingNumber,
) (*analysis_entity.Layer1, []filter.Id, error) {
	ticketsMap := p.ticketConverter.ConvertToRaceIdMap(ctx, tickets, racingNumbers, races)
	raceMap := map[types.RaceId]*data_cache_entity.Race{}
	for _, race := range races {
		raceMap[race.RaceId()] = race
	}

	for _, marker := range markers {
		race, ok := raceMap[marker.RaceId()]
		if !ok {
			log.Println(ctx, fmt.Sprintf("raceId: %s is not found in predict_races", marker.RaceId()))
			continue
		}
		// 馬券購入がない場合はnilになる
		ticketsByRaceId := ticketsMap[marker.RaceId()]

		raceResultMap := map[int]*data_cache_entity.RaceResult{}
		for _, raceResult := range race.RaceResults() {
			raceResultMap[raceResult.HorseNumber()] = raceResult
		}

		filters := p.analysisService.CreateAnalysisFilters(ctx, race)
		for _, payoutResult := range race.PayoutResults() {
			hitMarkerCombinationIds := p.analysisService.GetHitMarkerCombinationIds(ctx, payoutResult, marker)
			for idx, markerCombinationId := range hitMarkerCombinationIds {
				var (
					payment types.Payment
					payout  types.Payout
				)
				if ticketsByRaceId != nil {
					for _, ticket := range ticketsByRaceId {
						if ticket.TicketType() == markerCombinationId.TicketType() && ticket.BetNumber() == payoutResult.Numbers()[idx] {
							payment = ticket.Payment()
							payout = ticket.Payout()
							break
						}
					}
				}
				calculable := analysis_entity.NewCalculable(
					payment,
					payout,
					payoutResult.Odds()[idx],
					payoutResult.Numbers()[idx],
					payoutResult.Populars()[idx],
					1,
					race.Entries(),
					filters,
				)
				err := p.analysisService.AddAnalysisData(ctx, markerCombinationId, race, calculable, true)
				if err != nil {
					return nil, nil, err
				}
			}

			unHitMarkerCombinationIds := p.analysisService.GetUnHitMarkerCombinationIds(ctx, payoutResult, marker)
			for _, markerCombinationId := range unHitMarkerCombinationIds {
				// 不的中の集計については単複のみ(他の券種は組合せのオッズの取得ができないため)
				if markerCombinationId.TicketType() == types.Win || markerCombinationId.TicketType() == types.Place {
					unHitMarker, err := types.NewMarker(markerCombinationId.Value() % 10)
					if err != nil {
						return nil, nil, err
					}
					horseNumber, ok := marker.MarkerMap()[unHitMarker]
					if !ok && unHitMarker != types.NoMarker {
						return nil, nil, fmt.Errorf("marker %s is not found in markerMap", unHitMarker.String())
					}
					if raceResult, ok := raceResultMap[horseNumber]; ok {
						var (
							payment types.Payment
							payout  types.Payout
						)
						if ticketsByRaceId != nil {
							for _, ticket := range ticketsByRaceId {
								betNumber := ticket.BetNumber().List()[0] // 単複のみ
								if ticket.TicketType() == markerCombinationId.TicketType() && betNumber == raceResult.HorseNumber() {
									payment = ticket.Payment()
									payout = ticket.Payout()
									break
								}
							}
						}
						calculable := analysis_entity.NewCalculable(
							payment,
							payout,
							raceResult.Odds(),
							types.BetNumber(strconv.Itoa(raceResult.HorseNumber())), // 単複のみなのでbetNumberにそのまま置き換え可能
							raceResult.PopularNumber(),
							raceResult.OrderNo(),
							race.Entries(),
							filters,
						)

						err := p.analysisService.AddAnalysisData(ctx, markerCombinationId, race, calculable, false)
						if err != nil {
							return nil, nil, err
						}
					}
				}
			}
		}
	}

	return p.analysisService.GetAnalysisData(), p.analysisService.GetSearchFilters(), nil
}
