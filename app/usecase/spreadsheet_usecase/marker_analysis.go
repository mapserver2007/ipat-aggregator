package spreadsheet_usecase

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
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
	filters := []filter.Id{filter.All, filter.Turf, filter.Dirt}
	spreadSheetAnalysisData := p.analysisService.CreateSpreadSheetAnalysisData(ctx, analysisData)
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
