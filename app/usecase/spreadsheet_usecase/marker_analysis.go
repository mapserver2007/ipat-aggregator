package spreadsheet_usecase

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
)

type markerAnalysisUseCase struct {
	spreadSheetRepository  repository.SpreadSheetMarkerAnalysisRepository
	predictAnalysisService service.AnalysisService
}

func NewMarkerAnalysisUseCase(
	spreadSheetRepository repository.SpreadSheetMarkerAnalysisRepository,
	predictAnalysisService service.AnalysisService,
) *markerAnalysisUseCase {
	return &markerAnalysisUseCase{
		spreadSheetRepository:  spreadSheetRepository,
		predictAnalysisService: predictAnalysisService,
	}
}

func (p *markerAnalysisUseCase) Write(
	ctx context.Context,
	analysisData *analysis_entity.Layer1,
) error {
	spreadSheetAnalysisData := p.predictAnalysisService.CreateSpreadSheetAnalysisData(ctx, analysisData)
	_ = spreadSheetAnalysisData

	// TODO いろいろ集計データを作る処理
	return nil
}
