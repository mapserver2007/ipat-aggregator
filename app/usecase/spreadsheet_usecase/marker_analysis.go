package spreadsheet_usecase

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
)

type markerAnalysisUseCase struct {
	spreadSheetRepository repository.SpreadSheetMarkerAnalysisRepository
	analysisService       service.AnalysisService
}

func NewMarkerAnalysisUseCase(
	spreadSheetRepository repository.SpreadSheetMarkerAnalysisRepository,
	analysisService service.AnalysisService,
) *markerAnalysisUseCase {
	return &markerAnalysisUseCase{
		spreadSheetRepository: spreadSheetRepository,
		analysisService:       analysisService,
	}
}

func (p *markerAnalysisUseCase) Write(
	ctx context.Context,
	analysisData *analysis_entity.Layer1,
) error {
	err := p.spreadSheetRepository.Clear(ctx)
	if err != nil {
		return err
	}
	spreadSheetAnalysisData := p.analysisService.CreateSpreadSheetAnalysisData(ctx, analysisData)
	err = p.spreadSheetRepository.Write(ctx, spreadSheetAnalysisData)
	if err != nil {
		return err
	}
	//err = p.spreadSheetRepository.Style(ctx, spreadSheetAnalysisData)
	//if err != nil {
	//	return err
	//}

	return nil
}
