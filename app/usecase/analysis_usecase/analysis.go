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
)

type AnalysisUseCase struct {
	markerDataRepository repository.MarkerDataRepository
	analysisService      service.AnalysisService
	ticketConverter      service.TicketConverter
}

func NewAnalysisUseCase(
	markerDataRepository repository.MarkerDataRepository,
	analysisService service.AnalysisService,
	ticketConverter service.TicketConverter,
) *AnalysisUseCase {
	return &AnalysisUseCase{
		markerDataRepository: markerDataRepository,
		analysisService:      analysisService,
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

		err := p.analysisService.AddAnalysisData(ctx, marker, race)
		if err != nil {
			return nil, err
		}
	}

	return p.analysisService.GetAnalysisData(), nil
}
