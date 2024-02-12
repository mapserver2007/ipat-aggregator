package spreadsheet_usecase

import (
	"context"
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
	rows []*list_entity.ListRow,
) error {

	return nil
}
