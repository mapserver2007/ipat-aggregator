package entity

import spreadsheet_vo "github.com/mapserver2007/ipat-aggregator/app/domain/spreadsheet/value_object"

type SpreadSheetStyle struct {
	rowIndex         int
	favoriteColor    spreadsheet_vo.PlaceColor
	rivalColor       spreadsheet_vo.PlaceColor
	firstPlaceColor  spreadsheet_vo.PopularColor
	secondPlaceColor spreadsheet_vo.PopularColor
	gradeClassColor  spreadsheet_vo.GradeClassColor
	repaymentComment spreadsheet_vo.RepaymentComments
}

func NewSpreadSheetStyle(
	rowIndex int,
	favoriteColor, rivalColor spreadsheet_vo.PlaceColor,
	firstPlaceColor, secondPlaceColor spreadsheet_vo.PopularColor,
	gradeClassColor spreadsheet_vo.GradeClassColor,
	repaymentComment spreadsheet_vo.RepaymentComments,
) *SpreadSheetStyle {
	return &SpreadSheetStyle{
		rowIndex:         rowIndex,
		favoriteColor:    favoriteColor,
		rivalColor:       rivalColor,
		firstPlaceColor:  firstPlaceColor,
		secondPlaceColor: secondPlaceColor,
		gradeClassColor:  gradeClassColor,
		repaymentComment: repaymentComment,
	}
}

func (r *SpreadSheetStyle) GetRowIndex() int {
	return r.rowIndex
}

func (r *SpreadSheetStyle) GetFavoriteColor() spreadsheet_vo.PlaceColor {
	return r.favoriteColor
}

func (r *SpreadSheetStyle) GetRivalColor() spreadsheet_vo.PlaceColor {
	return r.rivalColor
}

func (r *SpreadSheetStyle) GetFirstPlaceColor() spreadsheet_vo.PopularColor {
	return r.firstPlaceColor
}

func (r *SpreadSheetStyle) GetSecondPlaceColor() spreadsheet_vo.PopularColor {
	return r.secondPlaceColor
}

func (r *SpreadSheetStyle) GetGradeClassColor() spreadsheet_vo.GradeClassColor {
	return r.gradeClassColor
}

func (r *SpreadSheetStyle) GetRepaymentComment() string {
	return r.repaymentComment.String()
}
