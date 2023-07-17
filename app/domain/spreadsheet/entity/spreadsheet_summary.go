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
