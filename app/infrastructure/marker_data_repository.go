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

func (p *markerDataRepository) Read(ctx context.Context, filePath string) ([]*marker_csv_entity.Yamato, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var predicts []*marker_csv_entity.Yamato
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

		predict, err := marker_csv_entity.NewYamato(
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

		predicts = append(predicts, predict)
		rowNum++
	}

	return predicts, nil
}