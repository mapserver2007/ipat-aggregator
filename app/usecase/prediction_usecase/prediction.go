package prediction_usecase

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"os"
	"path/filepath"
)

type PredictionUseCase struct {
	netKeibaService          service.NetKeibaService
	raceIdDataRepository     repository.RaceIdDataRepository
	predictionDataRepository repository.PredictionDataRepository
}

func NewPredictionUseCase(
	netKeibaService service.NetKeibaService,
	raceIdDataRepository repository.RaceIdDataRepository,
	predictionDataRepository repository.PredictionDataRepository,
) *PredictionUseCase {
	return &PredictionUseCase{
		netKeibaService:          netKeibaService,
		raceIdDataRepository:     raceIdDataRepository,
		predictionDataRepository: predictionDataRepository,
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

func (p *PredictionUseCase) Fetch(ctx context.Context, raceIds []types.RaceId) error {
	raceUrls, oddsUrls := p.netKeibaService.CreatePredictionRaceUrls(ctx, raceIds)
	for i := 0; i < len(raceUrls); i++ {
		race, odds, err := p.predictionDataRepository.Fetch(ctx, raceUrls[i], oddsUrls[i])
		if err != nil {
			return err
		}

		_ = race
		_ = odds

	}

	// TODO create urls
	return nil
}

func (p *PredictionUseCase) Predict() error {
	return nil
}
