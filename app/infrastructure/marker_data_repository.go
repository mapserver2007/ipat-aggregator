package infrastructure

import (
	"context"
	"encoding/csv"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"io"
	"os"
)

type markerDataRepository struct{}

func NewMarkerDataRepository() repository.MarkerDataRepository {
	return &markerDataRepository{}
}

func (p *markerDataRepository) Read(ctx context.Context, filePath string) ([]*marker_csv_entity.AnalysisMarker, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var markers []*marker_csv_entity.AnalysisMarker
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

		marker, err := marker_csv_entity.NewAnalysisMarker(
			record[0],
			record[1],
			record[2],
			record[3],
			record[4],
			record[5],
			record[6],
			record[7],
		)
		if err != nil {
			return nil, err
		}

		markers = append(markers, marker)
		rowNum++
	}

	return markers, nil
}
