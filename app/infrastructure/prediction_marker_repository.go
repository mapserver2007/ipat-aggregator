package infrastructure

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure/file_gateway"
	"github.com/mapserver2007/ipat-aggregator/app/infrastructure/gateway"
)

type predictionMarkerRepository struct {
	netKeibaGateway gateway.NetKeibaGateway
	pathOptimizer   file_gateway.PathOptimizer
}

func NewPredictionMarkerRepository(
	netKeibaGateway gateway.NetKeibaGateway,
	pathOptimizer file_gateway.PathOptimizer,
) repository.PredictionMarkerRepository {
	return &predictionMarkerRepository{
		netKeibaGateway: netKeibaGateway,
		pathOptimizer:   pathOptimizer,
	}
}

func (p *predictionMarkerRepository) Read(
	ctx context.Context,
	path string,
) ([]*marker_csv_entity.PredictionMarker, error) {
	rootPath, err := p.pathOptimizer.GetProjectRoot()
	if err != nil {
		return nil, err
	}
	absPath, err := filepath.Abs(fmt.Sprintf("%s/%s", rootPath, path))
	if err != nil {
		return nil, err
	}

	f, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var markers []*marker_csv_entity.PredictionMarker
	reader := csv.NewReader(f)
	rowNum := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if rowNum == 0 {
			rowNum++
			continue
		}

		marker := marker_csv_entity.NewPredictionMarker(
			record[0],
			record[1],
			record[2],
			record[3],
			record[4],
			record[5],
			record[6],
		)
		if err != nil {
			return nil, err
		}

		markers = append(markers, marker)
		rowNum++
	}

	return markers, nil
}
