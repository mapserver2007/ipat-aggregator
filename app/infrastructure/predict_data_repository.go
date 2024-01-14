package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"os"
	"path/filepath"
)

type predictDataRepository struct{}

func NewPredictDataRepository() repository.PredictDataRepository {
	return &predictDataRepository{}
}

func (p predictDataRepository) Read(ctx context.Context, filePath string) ([]*raw_entity.Predict, error) {
	predicts := make([]*raw_entity.Predict, 0)
	rootPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	path, err := filepath.Abs(fmt.Sprintf("%s/cache/%s", rootPath, filePath))
	if err != nil {
		return nil, err
	}

	// ファイルが存在しない場合はエラーは返さず処理を継続する
	bytes, err := os.ReadFile(path)
	if err != nil {
		return predicts, nil
	}

	var predictInfo *raw_entity.PredictInfo
	if err := json.Unmarshal(bytes, &predictInfo); err != nil {
		return nil, err
	}
	predicts = predictInfo.Predicts

	return predicts, nil
}
