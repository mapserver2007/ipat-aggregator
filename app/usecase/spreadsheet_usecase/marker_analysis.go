package spreadsheet_usecase

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
)

type MarkerAnalysisUseCase struct {
	spreadSheetMarkerAnalysisRepository repository.SpreadSheetMarkerAnalysisRepository
	spreadSheetTrioAnalysisRepository   repository.SpreadSheetTrioAnalysisRepository
	analysisService                     service.AnalysisService
	filterService                       service.FilterService
}

func NewMarkerAnalysisUseCase(
	spreadSheetMarkerAnalysisRepository repository.SpreadSheetMarkerAnalysisRepository,
	spreadSheetTrioAnalysisRepository repository.SpreadSheetTrioAnalysisRepository,
	analysisService service.AnalysisService,
	filterService service.FilterService,
) *MarkerAnalysisUseCase {
	return &MarkerAnalysisUseCase{
		spreadSheetMarkerAnalysisRepository: spreadSheetMarkerAnalysisRepository,
		spreadSheetTrioAnalysisRepository:   spreadSheetTrioAnalysisRepository,
		analysisService:                     analysisService,
		filterService:                       filterService,
	}
}

func (p *MarkerAnalysisUseCase) Write(
	ctx context.Context,
	analysisData *analysis_entity.Layer1,
	races []*data_cache_entity.Race,
	odds []*data_cache_entity.Odds,
) error {
	//winPlaceFilters := p.filterService.GetWinPlaceAnalysisFilters()
	trioFilters := p.filterService.GetTrioAnalysisFilters()

	//spreadSheetWinPlaceAnalysisData := p.analysisService.CreateSpreadSheetAnalysisData(ctx, analysisData, winPlaceFilters)
	spreadSheetTrioAnalysisData := p.analysisService.CreateSpreadSheetAnalysisData(ctx, analysisData, trioFilters)

	//err := p.spreadSheetMarkerAnalysisRepository.Clear(ctx)
	//if err != nil {
	//	return err
	//}
	//err = p.spreadSheetMarkerAnalysisRepository.Write(ctx, spreadSheetWinPlaceAnalysisData)
	//if err != nil {
	//	return err
	//}
	//err = p.spreadSheetMarkerAnalysisRepository.Style(ctx, spreadSheetWinPlaceAnalysisData)
	//if err != nil {
	//	return err
	//}

	err := p.spreadSheetTrioAnalysisRepository.Clear(ctx)
	if err != nil {
		return err
	}
	err = p.spreadSheetTrioAnalysisRepository.Write(ctx, spreadSheetTrioAnalysisData, races, odds)
	if err != nil {
		return err
	}
	err = p.spreadSheetTrioAnalysisRepository.Style(ctx, spreadSheetTrioAnalysisData)
	if err != nil {
		return err
	}

	return nil
}
