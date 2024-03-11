package analysis_usecase

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type AnalysisUseCase struct {
	markerDataRepository repository.MarkerDataRepository
	analysisService      service.AnalysisService
	filterService        service.FilterService
	ticketConverter      service.TicketConverter
}

func NewAnalysisUseCase(
	markerDataRepository repository.MarkerDataRepository,
	analysisService service.AnalysisService,
	filterService service.FilterService,
	ticketConverter service.TicketConverter,
) *AnalysisUseCase {
	return &AnalysisUseCase{
		markerDataRepository: markerDataRepository,
		analysisService:      analysisService,
		filterService:        filterService,
		ticketConverter:      ticketConverter,
	}
}

func (p *AnalysisUseCase) Read(ctx context.Context) ([]*marker_csv_entity.AnalysisMarker, error) {
	rootPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	dirPath, err := filepath.Abs(rootPath + "/csv")
	if err != nil {
		return nil, err
	}
	filePath := fmt.Sprintf("%s/%s", dirPath, "analysis_marker.csv")
	return p.markerDataRepository.Read(ctx, filePath)
}

func (p *AnalysisUseCase) CreateAnalysisData(
	ctx context.Context,
	markers []*marker_csv_entity.AnalysisMarker,
	races []*data_cache_entity.Race,
) (*analysis_entity.Layer1, error) {
	raceMap := map[types.RaceId]*data_cache_entity.Race{}
	for _, race := range races {
		raceMap[race.RaceId()] = race
	}

	for _, marker := range markers {
		race, ok := raceMap[marker.RaceId()]
		if !ok {
			log.Println(ctx, fmt.Sprintf("raceId: %s is not found in races", marker.RaceId()))
			continue
		}

		raceResultMap := map[int]*data_cache_entity.RaceResult{}
		for _, raceResult := range race.RaceResults() {
			raceResultMap[raceResult.HorseNumber()] = raceResult
		}

		for _, payoutResult := range race.PayoutResults() {
			hitMarkerCombinationIds := p.analysisService.GetHitMarkerCombinationIds(ctx, payoutResult, marker)
			for _, markerCombinationId := range hitMarkerCombinationIds {
				// 集計については単複のみ(他の券種は組合せのオッズの取得ができないため)
				if markerCombinationId.TicketType() == types.Win || markerCombinationId.TicketType() == types.Place {
					hitMarker, err := types.NewMarker(markerCombinationId.Value() % 10)
					if err != nil {
						return nil, err
					}
					horseNumber, ok := marker.MarkerMap()[hitMarker]
					if !ok && hitMarker != types.NoMarker {
						return nil, fmt.Errorf("marker %s is not found in markerMap", hitMarker.String())
					}
					if raceResult, ok := raceResultMap[horseNumber]; ok {
						calculable := analysis_entity.NewCalculable(
							raceResult.Odds(),
							types.BetNumber(strconv.Itoa(raceResult.HorseNumber())), // 単複のみなのでbetNumberにそのまま置き換え可能
							raceResult.PopularNumber(),
							raceResult.OrderNo(),
							race.Entries(),
							p.filterService.CreateAnalysisFilters(ctx, race, raceResult),
						)

						err = p.analysisService.AddAnalysisData(ctx, markerCombinationId, race, calculable)
						if err != nil {
							return nil, err
						}
					}
				}
			}

			unHitMarkerCombinationIds := p.analysisService.GetUnHitMarkerCombinationIds(ctx, payoutResult, marker)
			for _, markerCombinationId := range unHitMarkerCombinationIds {
				// 集計については単複のみ(他の券種は組合せのオッズの取得ができないため)
				if markerCombinationId.TicketType() == types.Win || markerCombinationId.TicketType() == types.Place {
					unHitMarker, err := types.NewMarker(markerCombinationId.Value() % 10)
					if err != nil {
						return nil, err
					}
					horseNumber, ok := marker.MarkerMap()[unHitMarker]
					if !ok && unHitMarker != types.NoMarker {
						return nil, fmt.Errorf("marker %s is not found in markerMap", unHitMarker.String())
					}
					if raceResult, ok := raceResultMap[horseNumber]; ok {
						calculable := analysis_entity.NewCalculable(
							raceResult.Odds(),
							types.BetNumber(strconv.Itoa(raceResult.HorseNumber())), // 単複のみなのでbetNumberにそのまま置き換え可能
							raceResult.PopularNumber(),
							raceResult.OrderNo(),
							race.Entries(),
							p.filterService.CreateAnalysisFilters(ctx, race, raceResult),
						)

						err = p.analysisService.AddAnalysisData(ctx, markerCombinationId, race, calculable)
						if err != nil {
							return nil, err
						}
					}
				}
			}
		}
	}

	return p.analysisService.GetAnalysisData(), nil
}
