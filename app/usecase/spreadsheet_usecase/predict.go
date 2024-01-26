package spreadsheet_usecase

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/predict_analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
)

type predictUseCase struct {
	spreadSheetSummaryRepository repository.SpreadSheetSummaryRepository
	predictAnalysisService       service.PredictAnalysisService
}

func NewPredictUseCase(
	spreadSheetSummaryRepository repository.SpreadSheetSummaryRepository,
	predictAnalysisService service.PredictAnalysisService,
) *predictUseCase {
	return &predictUseCase{
		spreadSheetSummaryRepository: spreadSheetSummaryRepository,
		predictAnalysisService:       predictAnalysisService,
	}
}

func (p *predictUseCase) Write(
	ctx context.Context,
	analysisData *predict_analysis_entity.Layer1,
) error {
	spreadSheetAnalysisData := p.predictAnalysisService.CreateSpreadSheetAnalysisData(ctx, analysisData)

	// TODO いろいろ集計データを作る処理
	return nil
}
