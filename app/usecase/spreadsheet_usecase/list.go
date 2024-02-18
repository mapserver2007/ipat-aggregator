package spreadsheet_usecase

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/list_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
)

type listUseCase struct {
	listService           service.ListService
	spreadSheetRepository repository.SpreadsheetListRepository
}

func NewListUseCase(
	listService service.ListService,
	spreadSheetRepository repository.SpreadsheetListRepository,
) *listUseCase {
	return &listUseCase{
		listService:           listService,
		spreadSheetRepository: spreadSheetRepository,
	}
}

func (p *listUseCase) Write(
	ctx context.Context,
	listRows []*list_entity.ListRow,
	jockeys []*data_cache_entity.Jockey,
) error {
	err := p.spreadSheetRepository.Clear(ctx)
	if err != nil {
		return err
	}
	rows, styles := p.listService.Convert(ctx, listRows, jockeys)
	err = p.spreadSheetRepository.Write(ctx, rows)
	if err != nil {
		return err
	}
	err = p.spreadSheetRepository.Style(ctx, styles)
	if err != nil {
		return err
	}

	return nil
}
