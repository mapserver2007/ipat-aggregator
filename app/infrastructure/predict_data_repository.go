package infrastructure

import (
	"context"
	"encoding/csv"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/predict_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"io"
	"os"
)

type predictDataRepository struct{}

func NewPredictDataRepository() repository.PredictDataRepository {
	return &predictDataRepository{}
}

func (p predictDataRepository) Read(ctx context.Context, filePath string) ([]*predict_csv_entity.Yamato, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var predicts []*predict_csv_entity.Yamato
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

		predict, err := predict_csv_entity.NewYamato(
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
