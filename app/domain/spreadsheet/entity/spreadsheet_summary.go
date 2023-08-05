package entity

import (
	race_vo "github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
	result_summary_entity "github.com/mapserver2007/ipat-aggregator/app/domain/result/entity"
)

type SpreadSheetSummary struct {
	shortSummary            *SpreadSheetShortSummary
	bettingTicketSummary    *SpreadSheetBettingTicketSummary
	classSummary            *SpreadSheetClassSummary
	monthlySummary          *SpreadSheetMonthlySummary
	courseCategorySummary   *SpreadSheetCourseCategorySummary
	distanceCategorySummary *SpreadSheetDistanceCategorySummary
	raceCourseSummary       *SpreadSheetRaceCourseSummary
}

func NewSpreadSheetSummary(
	shortSummary *SpreadSheetShortSummary,
	bettingTicketSummary *SpreadSheetBettingTicketSummary,
	classSummary *SpreadSheetClassSummary,
	monthlySummary *SpreadSheetMonthlySummary,
	courseCategorySummary *SpreadSheetCourseCategorySummary,
	distanceCategorySummary *SpreadSheetDistanceCategorySummary,
	raceCourseSummary *SpreadSheetRaceCourseSummary,
) *SpreadSheetSummary {
	return &SpreadSheetSummary{
		shortSummary:            shortSummary,
		bettingTicketSummary:    bettingTicketSummary,
		classSummary:            classSummary,
		monthlySummary:          monthlySummary,
		courseCategorySummary:   courseCategorySummary,
		distanceCategorySummary: distanceCategorySummary,
		raceCourseSummary:       raceCourseSummary,
	}
}

func (s *SpreadSheetSummary) GetShortSummary() *SpreadSheetShortSummary {
	return s.shortSummary
}

func (s *SpreadSheetSummary) GetBettingTicketSummary() *SpreadSheetBettingTicketSummary {
	return s.bettingTicketSummary
}

func (s *SpreadSheetSummary) GetClassSummary() *SpreadSheetClassSummary {
	return s.classSummary

}

func (s *SpreadSheetSummary) GetMonthlySummary() *SpreadSheetMonthlySummary {
	return s.monthlySummary
}

func (s *SpreadSheetSummary) GetCourseCategorySummary() *SpreadSheetCourseCategorySummary {
	return s.courseCategorySummary
}

func (s *SpreadSheetSummary) GetDistanceCategorySummary() *SpreadSheetDistanceCategorySummary {
	return s.distanceCategorySummary
}

func (s *SpreadSheetSummary) GetRaceCourseSummary() *SpreadSheetRaceCourseSummary {
	return s.raceCourseSummary
}

type SpreadSheetShortSummary struct {
	shortSummaryForAll   result_summary_entity.ShortSummary
	shortSummaryForMonth result_summary_entity.ShortSummary
	shortSummaryForYear  result_summary_entity.ShortSummary
}

func NewSpreadSheetShortSummary(
	shortSummaryForAll result_summary_entity.ShortSummary,
	shortSummaryForMonth result_summary_entity.ShortSummary,
	shortSummaryForYear result_summary_entity.ShortSummary,
) *SpreadSheetShortSummary {
	return &SpreadSheetShortSummary{
		shortSummaryForAll:   shortSummaryForAll,
		shortSummaryForMonth: shortSummaryForMonth,
		shortSummaryForYear:  shortSummaryForYear,
	}
}

func (s *SpreadSheetShortSummary) GetShortSummaryForAll() result_summary_entity.ShortSummary {
	return s.shortSummaryForAll
}

func (s *SpreadSheetShortSummary) GetShortSummaryForMonth() result_summary_entity.ShortSummary {
	return s.shortSummaryForMonth
}

func (s *SpreadSheetShortSummary) GetShortSummaryForYear() result_summary_entity.ShortSummary {
	return s.shortSummaryForYear
}

type SpreadSheetBettingTicketSummary struct {
	winSummary           result_summary_entity.DetailSummary
	placeSummary         result_summary_entity.DetailSummary
	quinellaSummary      result_summary_entity.DetailSummary
	exactaSummary        result_summary_entity.DetailSummary
	quinellaPlaceSummary result_summary_entity.DetailSummary
	trioSummary          result_summary_entity.DetailSummary
	trifectaSummary      result_summary_entity.DetailSummary
	totalSummary         result_summary_entity.DetailSummary
}

func NewSpreadSheetBettingTicketSummary(
	winSummary result_summary_entity.DetailSummary,
	placeSummary result_summary_entity.DetailSummary,
	quinellaSummary result_summary_entity.DetailSummary,
	exactaSummary result_summary_entity.DetailSummary,
	quinellaPlaceSummary result_summary_entity.DetailSummary,
	trioSummary result_summary_entity.DetailSummary,
	trifectaSummary result_summary_entity.DetailSummary,
	totalSummary result_summary_entity.DetailSummary,
) *SpreadSheetBettingTicketSummary {
	return &SpreadSheetBettingTicketSummary{
		winSummary:           winSummary,
		placeSummary:         placeSummary,
		quinellaSummary:      quinellaSummary,
		exactaSummary:        exactaSummary,
		quinellaPlaceSummary: quinellaPlaceSummary,
		trioSummary:          trioSummary,
		trifectaSummary:      trifectaSummary,
		totalSummary:         totalSummary,
	}
}

func (s *SpreadSheetBettingTicketSummary) GetWinSummary() result_summary_entity.DetailSummary {
	return s.winSummary
}

func (s *SpreadSheetBettingTicketSummary) GetPlaceSummary() result_summary_entity.DetailSummary {
	return s.placeSummary
}

func (s *SpreadSheetBettingTicketSummary) GetQuinellaSummary() result_summary_entity.DetailSummary {
	return s.quinellaSummary
}

func (s *SpreadSheetBettingTicketSummary) GetExactaSummary() result_summary_entity.DetailSummary {
	return s.exactaSummary
}

func (s *SpreadSheetBettingTicketSummary) GetQuinellaPlaceSummary() result_summary_entity.DetailSummary {
	return s.quinellaPlaceSummary
}

func (s *SpreadSheetBettingTicketSummary) GetTrioSummary() result_summary_entity.DetailSummary {
	return s.trioSummary
}

func (s *SpreadSheetBettingTicketSummary) GetTrifectaSummary() result_summary_entity.DetailSummary {
	return s.trifectaSummary
}

func (s *SpreadSheetBettingTicketSummary) GetTotalSummary() result_summary_entity.DetailSummary {
	return s.totalSummary
}

type SpreadSheetClassSummary struct {
	grade1Summary        result_summary_entity.DetailSummary
	grade2Summary        result_summary_entity.DetailSummary
	grade3Summary        result_summary_entity.DetailSummary
	openClassSummary     result_summary_entity.DetailSummary
	threeWinClassSummary result_summary_entity.DetailSummary
	twoWinClassSummary   result_summary_entity.DetailSummary
	oneWinClassSummary   result_summary_entity.DetailSummary
	maidenClassSummary   result_summary_entity.DetailSummary
}

func NewSpreadSheetClassSummary(
	grade1Summary result_summary_entity.DetailSummary,
	grade2Summary result_summary_entity.DetailSummary,
	grade3Summary result_summary_entity.DetailSummary,
	openClassSummary result_summary_entity.DetailSummary,
	threeWinClassSummary result_summary_entity.DetailSummary,
	twoWinClassSummary result_summary_entity.DetailSummary,
	oneWinClassSummary result_summary_entity.DetailSummary,
	maidenClassSummary result_summary_entity.DetailSummary,
) *SpreadSheetClassSummary {
	return &SpreadSheetClassSummary{
		grade1Summary:        grade1Summary,
		grade2Summary:        grade2Summary,
		grade3Summary:        grade3Summary,
		openClassSummary:     openClassSummary,
		threeWinClassSummary: threeWinClassSummary,
		twoWinClassSummary:   twoWinClassSummary,
		oneWinClassSummary:   oneWinClassSummary,
		maidenClassSummary:   maidenClassSummary,
	}
}

func (s *SpreadSheetClassSummary) GetGrade1Summary() result_summary_entity.DetailSummary {
	return s.grade1Summary
}

func (s *SpreadSheetClassSummary) GetGrade2Summary() result_summary_entity.DetailSummary {
	return s.grade2Summary
}

func (s *SpreadSheetClassSummary) GetGrade3Summary() result_summary_entity.DetailSummary {
	return s.grade3Summary
}

func (s *SpreadSheetClassSummary) GetOpenClassSummary() result_summary_entity.DetailSummary {
	return s.openClassSummary
}

func (s *SpreadSheetClassSummary) GetThreeWinClassSummary() result_summary_entity.DetailSummary {
	return s.threeWinClassSummary
}

func (s *SpreadSheetClassSummary) GetTwoWinClassSummary() result_summary_entity.DetailSummary {
	return s.twoWinClassSummary
}

func (s *SpreadSheetClassSummary) GetOneWinClassSummary() result_summary_entity.DetailSummary {
	return s.oneWinClassSummary
}

func (s *SpreadSheetClassSummary) GetMaidenClassSummary() result_summary_entity.DetailSummary {
	return s.maidenClassSummary
}

type SpreadSheetMonthlySummary struct {
	monthlySummaryMap map[int]result_summary_entity.DetailSummary
}

func NewSpreadSheetMonthlySummary(monthlySummaryMap map[int]result_summary_entity.DetailSummary) *SpreadSheetMonthlySummary {
	return &SpreadSheetMonthlySummary{
		monthlySummaryMap: monthlySummaryMap,
	}
}

func (s *SpreadSheetMonthlySummary) GetMonthlySummaryMap() map[int]result_summary_entity.DetailSummary {
	return s.monthlySummaryMap
}

type SpreadSheetCourseCategorySummary struct {
	courseCategorySummaryMap map[race_vo.CourseCategory]result_summary_entity.DetailSummary
}

func NewSpreadSheetCourseCategorySummary(
	courseCategorySummaryMap map[race_vo.CourseCategory]result_summary_entity.DetailSummary,
) *SpreadSheetCourseCategorySummary {
	return &SpreadSheetCourseCategorySummary{
		courseCategorySummaryMap: courseCategorySummaryMap,
	}
}

func (s *SpreadSheetCourseCategorySummary) GetCourseCategorySummary(courseCategory race_vo.CourseCategory) result_summary_entity.DetailSummary {
	if courseCategorySummary, ok := s.courseCategorySummaryMap[courseCategory]; ok {
		return courseCategorySummary
	}
	return result_summary_entity.DetailSummary{}
}

type SpreadSheetDistanceCategorySummary struct {
	distanceCategorySummaryMap map[race_vo.DistanceCategory]result_summary_entity.DetailSummary
}

func NewSpreadSheetDistanceCategorySummary(
	distanceCategorySummaryMap map[race_vo.DistanceCategory]result_summary_entity.DetailSummary,
) *SpreadSheetDistanceCategorySummary {
	return &SpreadSheetDistanceCategorySummary{
		distanceCategorySummaryMap: distanceCategorySummaryMap,
	}
}

func (s *SpreadSheetDistanceCategorySummary) GetDistanceCategorySummary(distanceCategory race_vo.DistanceCategory) result_summary_entity.DetailSummary {
	if distanceCategorySummary, ok := s.distanceCategorySummaryMap[distanceCategory]; ok {
		return distanceCategorySummary
	}
	return result_summary_entity.DetailSummary{}
}

type SpreadSheetRaceCourseSummary struct {
	raceCourseSummaryMap map[race_vo.RaceCourse]result_summary_entity.DetailSummary
}

func NewSpreadSheetRaceCourseSummary(
	raceCourseSummaryMap map[race_vo.RaceCourse]result_summary_entity.DetailSummary,
) *SpreadSheetRaceCourseSummary {
	return &SpreadSheetRaceCourseSummary{
		raceCourseSummaryMap: raceCourseSummaryMap,
	}
}

func (s *SpreadSheetRaceCourseSummary) GetRaceCourseSummary(raceCourse race_vo.RaceCourse) result_summary_entity.DetailSummary {
	if raceCourseSummary, ok := s.raceCourseSummaryMap[raceCourse]; ok {
		return raceCourseSummary
	}
	return result_summary_entity.DetailSummary{}
}

type SpreadSheetMonthlyBettingTicketSummary struct {
	monthlyBettingTicketSummaryMap map[int]*SpreadSheetBettingTicketSummary
}

func NewSpreadSheetMonthlyBettingTicketSummary(
	monthlyBettingTicketSummaryMap map[int]*SpreadSheetBettingTicketSummary,
) *SpreadSheetMonthlyBettingTicketSummary {
	return &SpreadSheetMonthlyBettingTicketSummary{
		monthlyBettingTicketSummaryMap: monthlyBettingTicketSummaryMap,
	}
}

func (s *SpreadSheetMonthlyBettingTicketSummary) GetMonthlyBettingTicketSummaryMap() map[int]*SpreadSheetBettingTicketSummary {
	return s.monthlyBettingTicketSummaryMap
}

type SpreadSheetMonthlyBettingTicketSummary2 struct {
	winSummaryMap           map[int]*result_summary_entity.DetailSummary
	placeSummaryMap         map[int]*result_summary_entity.DetailSummary
	quinellaSummaryMap      map[int]*result_summary_entity.DetailSummary
	exactaSummaryMap        map[int]*result_summary_entity.DetailSummary
	quinellaPlaceSummaryMap map[int]*result_summary_entity.DetailSummary
	trioSummaryMap          map[int]*result_summary_entity.DetailSummary
	trifectaSummaryMap      map[int]*result_summary_entity.DetailSummary
	totalSummaryMap         map[int]*result_summary_entity.DetailSummary
}

func NewSpreadSheetMonthlyBettingTicketSummary2(
	winSummaryMap map[int]*result_summary_entity.DetailSummary,
	placeSummaryMap map[int]*result_summary_entity.DetailSummary,
	quinellaSummaryMap map[int]*result_summary_entity.DetailSummary,
	exactaSummaryMap map[int]*result_summary_entity.DetailSummary,
	quinellaPlaceSummaryMap map[int]*result_summary_entity.DetailSummary,
	trioSummaryMap map[int]*result_summary_entity.DetailSummary,
	trifectaSummaryMap map[int]*result_summary_entity.DetailSummary,
) *SpreadSheetMonthlyBettingTicketSummary2 {
	return &SpreadSheetMonthlyBettingTicketSummary2{
		winSummaryMap:           winSummaryMap,
		placeSummaryMap:         placeSummaryMap,
		quinellaSummaryMap:      quinellaSummaryMap,
		exactaSummaryMap:        exactaSummaryMap,
		quinellaPlaceSummaryMap: quinellaPlaceSummaryMap,
		trioSummaryMap:          trioSummaryMap,
		trifectaSummaryMap:      trifectaSummaryMap,
	}
}

func (s *SpreadSheetMonthlyBettingTicketSummary2) GetWinSummaryMap() map[int]*result_summary_entity.DetailSummary {
	return s.winSummaryMap
}

func (s *SpreadSheetMonthlyBettingTicketSummary2) GetPlaceSummaryMap() map[int]*result_summary_entity.DetailSummary {
	return s.placeSummaryMap
}

func (s *SpreadSheetMonthlyBettingTicketSummary2) GetQuinellaSummaryMap() map[int]*result_summary_entity.DetailSummary {
	return s.quinellaSummaryMap
}

func (s *SpreadSheetMonthlyBettingTicketSummary2) GetExactaSummaryMap() map[int]*result_summary_entity.DetailSummary {
	return s.exactaSummaryMap
}

func (s *SpreadSheetMonthlyBettingTicketSummary2) GetQuinellaPlaceSummaryMap() map[int]*result_summary_entity.DetailSummary {
	return s.quinellaPlaceSummaryMap
}

func (s *SpreadSheetMonthlyBettingTicketSummary2) GetTrioSummaryMap() map[int]*result_summary_entity.DetailSummary {
	return s.trioSummaryMap
}

func (s *SpreadSheetMonthlyBettingTicketSummary2) GetTrifectaSummaryMap() map[int]*result_summary_entity.DetailSummary {
	return s.trifectaSummaryMap
}
