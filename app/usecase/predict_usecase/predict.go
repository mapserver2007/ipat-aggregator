package predict_usecase

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/predict_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"os"
	"path/filepath"
)

const (
	startDate = "20240101"
	endDate   = "20240110"
)

type predict struct {
	predictDataRepository repository.PredictDataRepository
}

func NewPredict(
	predictDataRepository repository.PredictDataRepository,
) *predict {
	return &predict{
		predictDataRepository: predictDataRepository,
	}
}

func (p *predict) Read(ctx context.Context) ([]*predict_csv_entity.Yamato, error) {
	rootPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	dirPath, err := filepath.Abs(rootPath + "/csv/markers")
	if err != nil {
		return nil, err
	}

	filePath := fmt.Sprintf("%s/%s", dirPath, "yamato_predict.csv")
	predicts, err := p.predictDataRepository.Read(ctx, filePath)
	if err != nil {
		return nil, err
	}

	return predicts, nil
}

func (p *predict) Predict(ctx context.Context) error {
	// TODO いろいろ集計データを作る処理
	return nil
}
