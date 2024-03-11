package prediction_usecase

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"log"
	"os"
	"path/filepath"
)

type PredictionUseCase struct {
	netKeibaService          service.NetKeibaService
	raceIdDataRepository     repository.RaceIdDataRepository
	predictionDataRepository repository.PredictionDataRepository
	raceEntityConverter      service.RaceEntityConverter
	filterService            service.FilterService
}

func NewPredictionUseCase(
	netKeibaService service.NetKeibaService,
	raceIdDataRepository repository.RaceIdDataRepository,
	predictionDataRepository repository.PredictionDataRepository,
	raceEntityConverter service.RaceEntityConverter,
	filterService service.FilterService,
) *PredictionUseCase {
	return &PredictionUseCase{
		netKeibaService:          netKeibaService,
		raceIdDataRepository:     raceIdDataRepository,
		predictionDataRepository: predictionDataRepository,
		raceEntityConverter:      raceEntityConverter,
		filterService:            filterService,
	}
}

func (p *PredictionUseCase) Read(ctx context.Context) ([]*marker_csv_entity.PredictionMarker, error) {
	rootPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	dirPath, err := filepath.Abs(rootPath + "/csv")
	if err != nil {
		return nil, err
	}
	filePath := fmt.Sprintf("%s/%s", dirPath, "prediction_marker.csv")
	return p.predictionDataRepository.Read(ctx, filePath)
}

func (p *PredictionUseCase) Fetch(ctx context.Context, raceIds []types.RaceId) ([]*prediction_entity.Race, error) {
	raceUrls, oddsUrls, raceResultUrls := p.netKeibaService.CreatePredictionRaceUrls(ctx, raceIds)
	var races []*prediction_entity.Race
	for i := 0; i < len(raceUrls); i++ {
		log.Println(ctx, "fetch prediction data from "+raceUrls[i])
		log.Println(ctx, "fetch prediction data from "+oddsUrls[i])
		log.Println(ctx, "fetch prediction data from "+raceResultUrls[i])
		race, odds, err := p.predictionDataRepository.Fetch(ctx, raceUrls[i], oddsUrls[i], raceResultUrls[i])
		if err != nil {
			return nil, err
		}
		races = append(races, p.raceEntityConverter.NetKeibaToPrediction(race, odds))
	}

	return races, nil
}
