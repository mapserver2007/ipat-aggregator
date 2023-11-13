package repository

import (
	"context"
	analyze_entity "github.com/mapserver2007/ipat-aggregator/app/domain/analyze/entity"
	predict_entity "github.com/mapserver2007/ipat-aggregator/app/domain/predict/entity"
	race_vo "github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
	result_summary_entity "github.com/mapserver2007/ipat-aggregator/app/domain/result/entity"
	spreadsheet_entity "github.com/mapserver2007/ipat-aggregator/app/domain/spreadsheet/entity"
)

type SpreadSheetClient interface {
	WriteForTotalSummary(ctx context.Context, summary result_summary_entity.ShortSummary) error
	WriteStyleForTotalSummary(ctx context.Context) error
	WriteForCurrentMonthSummary(ctx context.Context, summary result_summary_entity.ShortSummary) error
	WriteStyleForCurrentMonthlySummary(ctx context.Context) error
	WriteForCurrentYearSummary(ctx context.Context, summary result_summary_entity.ShortSummary) error
	WriteStyleForCurrentYearSummary(ctx context.Context) error
	WriteForTotalBettingTicketRateSummary(ctx context.Context, summary *spreadsheet_entity.SpreadSheetBettingTicketSummary) error
	WriteStyleForTotalBettingTicketRateSummary(ctx context.Context) error
	WriteForRaceClassRateSummary(ctx context.Context, summary *spreadsheet_entity.SpreadSheetClassSummary) error
	WriteStyleForRaceClassRateSummary(ctx context.Context) error
	WriteForCourseCategoryRateSummary(ctx context.Context, summary *spreadsheet_entity.SpreadSheetCourseCategorySummary) error
	WriteStyleForCourseCategoryRateSummary(ctx context.Context) error
	WriteForDistanceCategoryRateSummary(ctx context.Context, summary *spreadsheet_entity.SpreadSheetDistanceCategorySummary) error
	WriteStyleForDistanceCategoryRateSummary(ctx context.Context) error
	WriteForRaceCourseRateSummary(ctx context.Context, summary *spreadsheet_entity.SpreadSheetRaceCourseSummary) error
	WriteStyleForRaceCourseRateSummary(ctx context.Context) error
	WriteForMonthlyRateSummary(ctx context.Context, summary *spreadsheet_entity.SpreadSheetMonthlySummary) error
	WriteStyleForMonthlyRateSummary(ctx context.Context, rowCount int) error
}

type SpreadSheetMonthlyBettingTicketClient interface {
	Write(ctx context.Context, summary *spreadsheet_entity.SpreadSheetMonthlyBettingTicketSummary) error
	WriteStyle(ctx context.Context, rowCount int) error
	Clear(ctx context.Context) error
}

type SpreadSheetListClient interface {
	WriteList(ctx context.Context, records []*predict_entity.PredictEntity) (map[race_vo.RaceId]*spreadsheet_entity.SpreadSheetStyle, error)
	WriteStyleList(ctx context.Context, records []*predict_entity.PredictEntity, styleMap map[race_vo.RaceId]*spreadsheet_entity.SpreadSheetStyle) error
	Clear(ctx context.Context) error
}

type SpreadSheetAnalyzeClient interface {
	WriteWinPopular(ctx context.Context, summary *analyze_entity.WinAnalyzeSummary) error
	WriteStyleWinPopular(ctx context.Context, summary *analyze_entity.WinAnalyzeSummary) error
	//WriteWinOdds(ctx context.Context) error
	Clear(ctx context.Context) error
}
