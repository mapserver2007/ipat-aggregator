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
	spreadSheetClient                     repository.SpreadSheetClient
	spreadSheetMonthlyBettingTicketClient repository.SpreadSheetMonthlyBettingTicketClient
	spreadSheetListClient                 repository.SpreadSheetListClient
	spreadSheetAnalyzeClient              repository.SpreadSheetAnalyzeClient
}

func NewSpreadSheet(
	spreadSheetClient repository.SpreadSheetClient,
	spreadSheetMonthlyBettingTicketClient repository.SpreadSheetMonthlyBettingTicketClient,
	spreadSheetListClient repository.SpreadSheetListClient,
	spreadSheetAnalyze repository.SpreadSheetAnalyzeClient,
) *SpreadSheet {
	return &SpreadSheet{
		spreadSheetClient:                     spreadSheetClient,
		spreadSheetMonthlyBettingTicketClient: spreadSheetMonthlyBettingTicketClient,
		spreadSheetListClient:                 spreadSheetListClient,
		spreadSheetAnalyzeClient:              spreadSheetAnalyze,
	}
}

func (s *SpreadSheet) WriteSummary(
	ctx context.Context,
	summary *spreadsheet_entity.SpreadSheetSummary,
) error {
	log.Println(ctx, "writing spreadsheet for summary")
	err := s.spreadSheetClient.Clear(ctx)
	if err != nil {
		return err
	}
	err = s.spreadSheetClient.WriteForTotalSummary(ctx, summary.GetShortSummary().GetShortSummaryForAll())
	if err != nil {
		return err
	}

	err = s.spreadSheetClient.WriteForCurrentMonthSummary(ctx, summary.GetShortSummary().GetShortSummaryForMonth())
	if err != nil {
		return err
	}

	err = s.spreadSheetClient.WriteForCurrentYearSummary(ctx, summary.GetShortSummary().GetShortSummaryForYear())
	if err != nil {
		return err
	}

	err = s.spreadSheetClient.WriteForTotalBettingTicketRateSummary(ctx, summary.GetBettingTicketSummary())
	if err != nil {
		return err
	}

	err = s.spreadSheetClient.WriteForRaceClassRateSummary(ctx, summary.GetClassSummary())
	if err != nil {
		return err
	}

	err = s.spreadSheetClient.WriteForCourseCategoryRateSummary(ctx, summary.GetCourseCategorySummary())
	if err != nil {
		return err
	}

	err = s.spreadSheetClient.WriteForDistanceCategoryRateSummary(ctx, summary.GetDistanceCategorySummary())
	if err != nil {
		return err
	}

	err = s.spreadSheetClient.WriteForRaceCourseRateSummary(ctx, summary.GetRaceCourseSummary())
	if err != nil {
		return err
	}

	err = s.spreadSheetClient.WriteForMonthlyRateSummary(ctx, summary.GetMonthlySummary())
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheet) WriteStyleSummary(ctx context.Context, summary *spreadsheet_entity.SpreadSheetSummary) error {
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

	rowCount := len(summary.GetMonthlySummary().GetMonthlySummaryMap())
	err = s.spreadSheetClient.WriteStyleForMonthlyRateSummary(ctx, rowCount)
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

func (s *SpreadSheet) WriteStyleList(ctx context.Context, records []*predict_entity.PredictEntity, styleMap map[race_vo.RaceId]*spreadsheet_entity.SpreadSheetStyle) error {
	err := s.spreadSheetListClient.WriteStyleList(ctx, records, styleMap)
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheet) WriteMonthlyBettingTicketSummary(
	ctx context.Context,
	summary *spreadsheet_entity.SpreadSheetMonthlyBettingTicketSummary,
) error {
	err := s.spreadSheetMonthlyBettingTicketClient.Clear(ctx)
	if err != nil {
		return err
	}
	log.Println(ctx, "writing spreadsheet for monthly betting ticket summary")
	err = s.spreadSheetMonthlyBettingTicketClient.Write(ctx, summary)
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheet) WriteStyleMonthlyBettingTicketSummary(
	ctx context.Context,
	summary *spreadsheet_entity.SpreadSheetMonthlyBettingTicketSummary,
) error {
	rowCount := len(summary.GetMonthlyBettingTicketSummaryMap())
	err := s.spreadSheetMonthlyBettingTicketClient.WriteStyle(ctx, rowCount)
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
