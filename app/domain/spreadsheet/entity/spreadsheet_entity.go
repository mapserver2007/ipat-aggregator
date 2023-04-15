package entity

import (
	"fmt"
	betting_ticket_vo "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/value_object"
	race_entity "github.com/mapserver2007/ipat-aggregator/app/domain/race/entity"
	race_vo "github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
	spreadsheet_vo "github.com/mapserver2007/ipat-aggregator/app/domain/spreadsheet/value_object"
	"strconv"
)

type Summary struct {
	TotalResultSummary       ResultSummary
	LatestMonthResultSummary ResultSummary
	LatestYearResultSummary  ResultSummary
	BettingTicketSummary     BettingTicketSummary
	RaceClassSummary         RaceClassSummary
	MonthlySummary           MonthlySummary
	YearlySummary            YearlySummary
	CourseCategorySummary    CourseCategorySummary
	RaceCourseSummary        RaceCourseSummary
	DistanceCategorySummary  DistanceCategorySummary
}

type ResultSummary struct {
	Payments   int
	Repayments int
}

func (r *ResultSummary) CalcReturnOnInvestment() string {
	if r.Payments == 0 {
		return fmt.Sprintf("%d%s", 0, "%")
	}
	return fmt.Sprintf("%s%s", strconv.FormatFloat((float64(r.Repayments)*float64(100))/float64(r.Payments), 'f', 2, 64), "%")
}

type BettingTicketSummary struct {
	BettingTicketRates map[betting_ticket_vo.BettingTicket]ResultRate
}

type RaceClassSummary struct {
	RaceClassRates map[race_vo.GradeClass]ResultRate
}

type MonthlySummary struct {
	MonthlyRates map[int]ResultRate
}

type YearlySummary struct {
	YearlyRates map[int]ResultRate
}

type CourseCategorySummary struct {
	CourseCategoryRates map[race_vo.CourseCategory]ResultRate
}

type RaceCourseSummary struct {
	RaceCourseRates map[race_vo.RaceCourse]ResultRate
}

type DistanceCategorySummary struct {
	DistanceCategoryRates map[race_vo.DistanceCategory]ResultRate
}

type RaceResultSummary struct {
	ResultSummary           ResultSummary
	Race                    race_entity.Race
	RaceHandicappingSummary RaceHandicappingSummary
}

type RaceHandicappingSummary struct {
	Favorite  HorseInfo
	Contender HorseInfo
}

type HorseInfo struct {
	name          string
	orderNo       int
	popularNumber int
	odds          string
}

type ResultRate struct {
	VoteCount  int
	HitCount   int
	Payments   int
	Repayments int
}

func (m *ResultRate) GetHitRate() float64 {
	if m.VoteCount == 0 {
		return 0
	}
	return (float64(m.HitCount) * float64(100)) / float64(m.VoteCount)
}

func (m *ResultRate) HitRateFormat() string {
	return fmt.Sprintf("%s%s", strconv.FormatFloat(m.GetHitRate(), 'f', 1, 64), "%")
}

func (m *ResultRate) ReturnOnInvestmentFormat() string {
	if m.Payments == 0 {
		return fmt.Sprintf("%d%s", 0, "%")
	}
	return fmt.Sprintf("%s%s", strconv.FormatFloat((float64(m.Repayments)*float64(100))/float64(m.Payments), 'f', 2, 64), "%")
}

type ResultStyle struct {
	RowIndex         int
	FavoriteColor    spreadsheet_vo.PlaceColor
	RivalColor       spreadsheet_vo.PlaceColor
	GradeClassColor  spreadsheet_vo.GradeClassColor
	RepaymentComment string
}

func NewResultStyle(
	rowIndex int,
	favoriteColor, rivalColor spreadsheet_vo.PlaceColor,
	gradeClassColor spreadsheet_vo.GradeClassColor,
	repaymentComment string,
) *ResultStyle {
	return &ResultStyle{
		RowIndex:         rowIndex,
		FavoriteColor:    favoriteColor,
		RivalColor:       rivalColor,
		GradeClassColor:  gradeClassColor,
		RepaymentComment: repaymentComment,
	}
}
