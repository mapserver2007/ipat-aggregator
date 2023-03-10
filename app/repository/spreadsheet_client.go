package repository

import (
	"context"
	predict_entity "github.com/mapserver2007/tools/baken/app/domain/predict/entity"
	race_vo "github.com/mapserver2007/tools/baken/app/domain/race/value_object"
	spreadsheet_entity "github.com/mapserver2007/tools/baken/app/domain/spreadsheet/entity"
)

type SpreadSheetClient interface {
	WriteForTotalSummary(ctx context.Context, summary spreadsheet_entity.ResultSummary) error
	WriteStyleForTotalSummary(ctx context.Context) error
	WriteForCurrentMonthSummary(ctx context.Context, summary spreadsheet_entity.ResultSummary) error
	WriteStyleForCurrentMonthlySummary(ctx context.Context) error
	WriteForTotalBettingTicketRateSummary(ctx context.Context, summary spreadsheet_entity.BettingTicketSummary) error
	WriteStyleForTotalBettingTicketRateSummary(ctx context.Context) error
	WriteForRaceClassRateSummary(ctx context.Context, summary spreadsheet_entity.RaceClassSummary) error
	WriteStyleForRaceClassRateSummary(ctx context.Context) error
	WriteForCourseCategoryRateSummary(ctx context.Context, summary spreadsheet_entity.CourseCategorySummary) error
	WriteStyleForCourseCategoryRateSummary(ctx context.Context) error
	WriteForDistanceCategoryRateSummary(ctx context.Context, summary spreadsheet_entity.DistanceCategorySummary) error
	WriteStyleForDistanceCategoryRateSummary(ctx context.Context) error
	WriteForRaceCourseRateSummary(ctx context.Context, summary spreadsheet_entity.RaceCourseSummary) error
	WriteStyleForRaceCourseRateSummary(ctx context.Context) error
	WriteForMonthlyRateSummary(ctx context.Context, summary spreadsheet_entity.MonthlySummary) error
	WriteStyleForMonthlyRateSummary(ctx context.Context, summary spreadsheet_entity.MonthlySummary) error
}

type SpreadSheetListClient interface {
	WriteList(ctx context.Context, records []*predict_entity.PredictEntity) (map[race_vo.RaceId]*spreadsheet_entity.ResultStyle, error)
	WriteStyleList(ctx context.Context, records []*predict_entity.PredictEntity, styleMap map[race_vo.RaceId]*spreadsheet_entity.ResultStyle) error
	Clear(ctx context.Context) error
}
