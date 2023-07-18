package usecase

import (
	"context"
	analyze_entity "github.com/mapserver2007/ipat-aggregator/app/domain/analyze/entity"
	predict_entity "github.com/mapserver2007/ipat-aggregator/app/domain/predict/entity"
	race_vo "github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
	spreadsheet_entity "github.com/mapserver2007/ipat-aggregator/app/domain/spreadsheet/entity"
	"github.com/mapserver2007/ipat-aggregator/app/repository"
	"log"
)

type SpreadSheet struct {
	spreadSheetClient        repository.SpreadSheetClient
	spreadSheetListClient    repository.SpreadSheetListClient
	spreadSheetAnalyzeClient repository.SpreadSheetAnalyzeClient
}

func NewSpreadSheet(
	spreadSheetClient repository.SpreadSheetClient,
	spreadSheetListClient repository.SpreadSheetListClient,
	spreadSheetAnalyze repository.SpreadSheetAnalyzeClient,
) *SpreadSheet {
	return &SpreadSheet{
		spreadSheetClient:        spreadSheetClient,
		spreadSheetListClient:    spreadSheetListClient,
		spreadSheetAnalyzeClient: spreadSheetAnalyze,
	}
}

func (s *SpreadSheet) WriteSummary(
	ctx context.Context,
	summary *spreadsheet_entity.Summary,
	summary2 *spreadsheet_entity.SpreadSheetSummary,
	summary3 *spreadsheet_entity.SpreadSheetBettingTicketSummary,
	summary4 *spreadsheet_entity.SpreadSheetClassSummary,
	summary5 *spreadsheet_entity.SpreadSheetMonthlySummary,
	sumamry6 *spreadsheet_entity.SpreadSheetCourseCategorySummary,
) error {
	log.Println(ctx, "writing spreadsheet for summary")
	err := s.spreadSheetClient.WriteForTotalSummary(ctx, summary2.GetShortSummaryForAll())
	if err != nil {
		return err
	}

	err = s.spreadSheetClient.WriteForCurrentMonthSummary(ctx, summary2.GetShortSummaryForMonth())
	if err != nil {
		return err
	}

	err = s.spreadSheetClient.WriteForCurrentYearSummary(ctx, summary2.GetShortSummaryForYear())
	if err != nil {
		return err
	}

	err = s.spreadSheetClient.WriteForTotalBettingTicketRateSummary(ctx, summary3)
	if err != nil {
		return err
	}

	err = s.spreadSheetClient.WriteForRaceClassRateSummary(ctx, summary4)
	if err != nil {
		return err
	}

	err = s.spreadSheetClient.WriteForCourseCategoryRateSummary(ctx, sumamry6)
	if err != nil {
		return err
	}

	err = s.spreadSheetClient.WriteForDistanceCategoryRateSummary(ctx, summary.DistanceCategorySummary)
	if err != nil {
		return err
	}

	err = s.spreadSheetClient.WriteForRaceCourseRateSummary(ctx, summary.RaceCourseSummary)
	if err != nil {
		return err
	}

	err = s.spreadSheetClient.WriteForMonthlyRateSummary(ctx, summary5)
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheet) WriteList(ctx context.Context, records []*predict_entity.PredictEntity) (map[race_vo.RaceId]*spreadsheet_entity.SpreadSheetStyle, error) {
	err := s.spreadSheetListClient.Clear(ctx)
	if err != nil {
		return nil, err
	}
	log.Println(ctx, "writing spreadsheet for list")
	styleMap, err := s.spreadSheetListClient.WriteList(ctx, records)
	if err != nil {
		return nil, err
	}

	return styleMap, nil
}

func (s *SpreadSheet) WriteStyleSummary(ctx context.Context, summary *spreadsheet_entity.Summary) error {
	err := s.spreadSheetClient.WriteStyleForTotalSummary(ctx)
	if err != nil {
		return err
	}

	err = s.spreadSheetClient.WriteStyleForCurrentMonthlySummary(ctx)
	if err != nil {
		return err
	}

	err = s.spreadSheetClient.WriteStyleForCurrentYearSummary(ctx)
	if err != nil {
		return err
	}

	err = s.spreadSheetClient.WriteStyleForTotalBettingTicketRateSummary(ctx)
	if err != nil {
		return err
	}

	err = s.spreadSheetClient.WriteStyleForRaceClassRateSummary(ctx)
	if err != nil {
		return err
	}

	err = s.spreadSheetClient.WriteStyleForCourseCategoryRateSummary(ctx)
	if err != nil {
		return err
	}

	err = s.spreadSheetClient.WriteStyleForDistanceCategoryRateSummary(ctx)
	if err != nil {
		return err
	}

	err = s.spreadSheetClient.WriteStyleForRaceCourseRateSummary(ctx)
	if err != nil {
		return err
	}

	err = s.spreadSheetClient.WriteStyleForMonthlyRateSummary(ctx, summary.MonthlySummary)
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheet) WriteStyleList(ctx context.Context, records []*predict_entity.PredictEntity, styleMap map[race_vo.RaceId]*spreadsheet_entity.SpreadSheetStyle) error {
	err := s.spreadSheetListClient.WriteStyleList(ctx, records, styleMap)
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheet) WriteAnalyze(ctx context.Context, summary *analyze_entity.AnalyzeSummary) error {
	err := s.spreadSheetAnalyzeClient.Clear(ctx)
	if err != nil {
		return err
	}
	log.Println(ctx, "writing spreadsheet for analyze")
	err = s.spreadSheetAnalyzeClient.WriteWinPopular(ctx, summary.WinPopularitySummary())
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheet) WriteStyleAnalyze(ctx context.Context, summary *analyze_entity.AnalyzeSummary) error {
	err := s.spreadSheetAnalyzeClient.WriteStyleWinPopular(ctx, summary.WinPopularitySummary())
	if err != nil {
		return err
	}

	return nil
}
