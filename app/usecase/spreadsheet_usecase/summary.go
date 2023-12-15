package spreadsheet_usecase

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
)

type summaryUseCase struct {
	summaryService               service.SummaryService
	spreadSheetSummaryRepository repository.SpreadSheetSummaryRepository
}

func NewSummaryUseCase(
	summaryService service.SummaryService,
	spreadSheetSummaryRepository repository.SpreadSheetSummaryRepository,
) *summaryUseCase {
	return &summaryUseCase{
		summaryService:               summaryService,
		spreadSheetSummaryRepository: spreadSheetSummaryRepository,
	}
}

func (s *summaryUseCase) Write(
	ctx context.Context,
	tickets []*ticket_csv_entity.Ticket,
) error {
	all, month, year := s.summaryService.CreateShortSummary(ctx, tickets)
	summary := spreadsheet_entity.NewSummary(all, month, year) // TODO 拡張する

	s.spreadSheetSummaryRepository.Write(ctx, summary)

	return nil
}
