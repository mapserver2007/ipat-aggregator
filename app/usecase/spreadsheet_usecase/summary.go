package spreadsheet_usecase

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
)

type summaryUseCase struct {
	summaryService         service.SummaryService
	spreadSheetSRepository repository.SpreadSheetSummaryRepository
}

func NewSummaryUseCase(
	summaryService service.SummaryService,
	spreadSheetSRepository repository.SpreadSheetSummaryRepository,
) *summaryUseCase {
	return &summaryUseCase{
		summaryService:         summaryService,
		spreadSheetSRepository: spreadSheetSRepository,
	}
}

func (s *summaryUseCase) Write(
	ctx context.Context,
	tickets []*ticket_csv_entity.Ticket,
	racingNumbers []*data_cache_entity.RacingNumber,
	races []*data_cache_entity.Race,
) error {
	summary := s.summaryService.CreateSummary(ctx, tickets, racingNumbers, races)
	err := s.spreadSheetSRepository.Write(ctx, summary)
	if err != nil {
		return err
	}
	err = s.spreadSheetSRepository.Style(ctx, summary)
	if err != nil {
		return err
	}

	return nil
}
