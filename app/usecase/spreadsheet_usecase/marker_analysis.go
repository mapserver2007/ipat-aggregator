package spreadsheet_usecase

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
)

type MarkerAnalysisUseCase struct {
	spreadSheetRepository repository.SpreadSheetMarkerAnalysisRepository
	analysisService       service.AnalysisService
	filterService         service.FilterService
}

func NewMarkerAnalysisUseCase(
	spreadSheetRepository repository.SpreadSheetMarkerAnalysisRepository,
	analysisService service.AnalysisService,
	filterService service.FilterService,
) *MarkerAnalysisUseCase {
	return &MarkerAnalysisUseCase{
		spreadSheetRepository: spreadSheetRepository,
		analysisService:       analysisService,
		filterService:         filterService,
	}
}

func (p *MarkerAnalysisUseCase) Write(
	ctx context.Context,
	analysisData *analysis_entity.Layer1,
) error {
	err := p.spreadSheetRepository.Clear(ctx)
	if err != nil {
		return err
	}
	filters := p.filterService.GetAnalysisFilters()
	spreadSheetAnalysisData := p.analysisService.CreateSpreadSheetAnalysisData(ctx, analysisData, filters)
	err = p.spreadSheetRepository.Write(ctx, spreadSheetAnalysisData, filters)
	if err != nil {
		return err
	}
	err = p.spreadSheetRepository.Style(ctx, spreadSheetAnalysisData, filters)
	if err != nil {
		return err
	}

	return nil
}
