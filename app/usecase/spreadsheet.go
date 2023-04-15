package usecase

import (
	"context"
	predict_entity "github.com/mapserver2007/ipat-aggregator/app/domain/predict/entity"
	race_vo "github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
	spreadsheet_entity "github.com/mapserver2007/ipat-aggregator/app/domain/spreadsheet/entity"
	"github.com/mapserver2007/ipat-aggregator/app/repository"
	"log"
)

type SpreadSheet struct {
	spreadSheetClient     repository.SpreadSheetClient
	spreadSheetListClient repository.SpreadSheetListClient
}

func NewSpreadSheet(
	spreadSheetClient repository.SpreadSheetClient,
	spreadSheetListClient repository.SpreadSheetListClient,
) *SpreadSheet {
	return &SpreadSheet{
		spreadSheetClient:     spreadSheetClient,
		spreadSheetListClient: spreadSheetListClient,
	}
}

func (s *SpreadSheet) WriteSummary(ctx context.Context, summary *spreadsheet_entity.Summary) error {
	log.Println(ctx, "writing spreadsheet for summary")
	err := s.spreadSheetClient.WriteForTotalSummary(ctx, summary.TotalResultSummary)
	if err != nil {
		return err
	}

	err = s.spreadSheetClient.WriteForCurrentMonthSummary(ctx, summary.LatestMonthResultSummary)
	if err != nil {
		return err
	}

	err = s.spreadSheetClient.WriteForCurrentYearSummary(ctx, summary.LatestYearResultSummary)
	if err != nil {
		return err
	}

	err = s.spreadSheetClient.WriteForTotalBettingTicketRateSummary(ctx, summary.BettingTicketSummary)
	if err != nil {
		return err
	}

	err = s.spreadSheetClient.WriteForRaceClassRateSummary(ctx, summary.RaceClassSummary)
	if err != nil {
		return err
	}

	err = s.spreadSheetClient.WriteForCourseCategoryRateSummary(ctx, summary.CourseCategorySummary)
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

	err = s.spreadSheetClient.WriteForMonthlyRateSummary(ctx, summary.MonthlySummary)
	if err != nil {
		return err
	}

	return nil
}

func (s *SpreadSheet) WriteList(ctx context.Context, records []*predict_entity.PredictEntity) (map[race_vo.RaceId]*spreadsheet_entity.ResultStyle, error) {
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

func (s *SpreadSheet) WriteStyleList(ctx context.Context, records []*predict_entity.PredictEntity, styleMap map[race_vo.RaceId]*spreadsheet_entity.ResultStyle) error {
	err := s.spreadSheetListClient.WriteStyleList(ctx, records, styleMap)
	if err != nil {
		return err
	}

	return nil
}
