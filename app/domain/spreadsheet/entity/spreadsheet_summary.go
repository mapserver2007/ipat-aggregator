package entity

import result_summary_entity "github.com/mapserver2007/ipat-aggregator/app/domain/result/entity"

type SpreadSheetSummary struct {
	shortSummaryForAll   result_summary_entity.ShortSummary
	shortSummaryForMonth result_summary_entity.ShortSummary
	shortSummaryForYear  result_summary_entity.ShortSummary
}

func NewSpreadSheetSummary(
	shortSummaryForAll result_summary_entity.ShortSummary,
	shortSummaryForMonth result_summary_entity.ShortSummary,
	shortSummaryForYear result_summary_entity.ShortSummary,
) *SpreadSheetSummary {
	return &SpreadSheetSummary{
		shortSummaryForAll:   shortSummaryForAll,
		shortSummaryForMonth: shortSummaryForMonth,
		shortSummaryForYear:  shortSummaryForYear,
	}
}

func (s *SpreadSheetSummary) GetShortSummaryForAll() result_summary_entity.ShortSummary {
	return s.shortSummaryForAll
}

func (s *SpreadSheetSummary) GetShortSummaryForMonth() result_summary_entity.ShortSummary {
	return s.shortSummaryForMonth
}

func (s *SpreadSheetSummary) GetShortSummaryForYear() result_summary_entity.ShortSummary {
	return s.shortSummaryForYear
}
