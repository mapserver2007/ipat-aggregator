package service

import (
	betting_ticket_entity "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/entity"
	betting_ticket_vo "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/value_object"
	race_entity "github.com/mapserver2007/ipat-aggregator/app/domain/race/entity"
	race_vo "github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
	spreadsheet_entity "github.com/mapserver2007/ipat-aggregator/app/domain/spreadsheet/entity"
)

type Aggregator struct {
	raceConverter          RaceConverter
	bettingTicketConverter BettingTicketConverter
	summarizer             Summarizer
}

func NewAggregator(
	raceConverter RaceConverter,
	bettingTicketConverter BettingTicketConverter,
	summarizer Summarizer,
) *Aggregator {
	return &Aggregator{
		raceConverter:          raceConverter,
		bettingTicketConverter: bettingTicketConverter,
		summarizer:             summarizer,
	}
}

func (a *Aggregator) GetSummary(
	records []*betting_ticket_entity.CsvEntity,
	racingNumbers []*race_entity.RacingNumber,
	races []*race_entity.Race,
) *spreadsheet_entity.SpreadSheetSummary {
	spreadSheetShortSummary := spreadsheet_entity.NewSpreadSheetShortSummary(
		a.summarizer.GetShortSummary(records),
		a.summarizer.GetShortSummaryForMonth(records),
		a.summarizer.GetShortSummaryForYear(records),
	)

	spreadSheetBettingTicketSummary := spreadsheet_entity.NewSpreadSheetBettingTicketSummary(
		a.summarizer.GetBettingTicketSummary(records, racingNumbers, races, betting_ticket_vo.Win),
		a.summarizer.GetBettingTicketSummary(records, racingNumbers, races, betting_ticket_vo.Place),
		a.summarizer.GetBettingTicketSummary(records, racingNumbers, races, betting_ticket_vo.Quinella),
		a.summarizer.GetBettingTicketSummary(records, racingNumbers, races, betting_ticket_vo.Exacta, betting_ticket_vo.ExactaWheelOfFirst),
		a.summarizer.GetBettingTicketSummary(records, racingNumbers, races, betting_ticket_vo.QuinellaPlace, betting_ticket_vo.QuinellaPlaceWheel),
		a.summarizer.GetBettingTicketSummary(records, racingNumbers, races, betting_ticket_vo.Trio, betting_ticket_vo.TrioFormation, betting_ticket_vo.TrioWheelOfFirst),
		a.summarizer.GetBettingTicketSummary(records, racingNumbers, races, betting_ticket_vo.Trifecta, betting_ticket_vo.TrifectaFormation, betting_ticket_vo.TrifectaWheelOfFirst, betting_ticket_vo.TrifectaWheelOfSecondMulti),
		a.summarizer.GetBettingTicketSummary(records, racingNumbers, races, betting_ticket_vo.Win, betting_ticket_vo.Place, betting_ticket_vo.Quinella,
			betting_ticket_vo.Exacta, betting_ticket_vo.ExactaWheelOfFirst, betting_ticket_vo.QuinellaPlace, betting_ticket_vo.QuinellaPlaceWheel,
			betting_ticket_vo.Trio, betting_ticket_vo.TrioFormation, betting_ticket_vo.TrioWheelOfFirst,
			betting_ticket_vo.Trifecta, betting_ticket_vo.TrifectaFormation, betting_ticket_vo.TrifectaWheelOfFirst, betting_ticket_vo.TrifectaWheelOfSecondMulti),
	)

	spreadSheetGradeClassSummary := spreadsheet_entity.NewSpreadSheetClassSummary(
		a.summarizer.GetGradeClassSummary(records, racingNumbers, races, race_vo.Grade1, race_vo.Jpn1, race_vo.JumpGrade1),
		a.summarizer.GetGradeClassSummary(records, racingNumbers, races, race_vo.Grade2, race_vo.Jpn2, race_vo.JumpGrade2),
		a.summarizer.GetGradeClassSummary(records, racingNumbers, races, race_vo.Grade3, race_vo.Jpn3, race_vo.JumpGrade3),
		a.summarizer.GetGradeClassSummary(records, racingNumbers, races, race_vo.OpenClass, race_vo.ListedClass),
		a.summarizer.GetGradeClassSummary(records, racingNumbers, races, race_vo.ThreeWinClass),
		a.summarizer.GetGradeClassSummary(records, racingNumbers, races, race_vo.TwoWinClass),
		a.summarizer.GetGradeClassSummary(records, racingNumbers, races, race_vo.OneWinClass),
		a.summarizer.GetGradeClassSummary(records, racingNumbers, races, race_vo.Maiden, race_vo.JumpMaiden),
	)

	spreadSheetMonthlySummary := spreadsheet_entity.NewSpreadSheetMonthlySummary(a.summarizer.GetMonthlySummaryMap(records))

	spreadSheetCourseCategorySummary := spreadsheet_entity.NewSpreadSheetCourseCategorySummary(
		a.summarizer.GetCourseCategorySummaryMap(records, racingNumbers, races),
	)

	spreadSheetDistanceCategorySummary := spreadsheet_entity.NewSpreadSheetDistanceCategorySummary(
		a.summarizer.GetDistanceCategorySummaryMap(records, racingNumbers, races),
	)

	spreadSheetRaceCourseSummary := spreadsheet_entity.NewSpreadSheetRaceCourseSummary(
		a.summarizer.GetRaceCourseSummaryMap(records, racingNumbers, races),
	)

	return spreadsheet_entity.NewSpreadSheetSummary(
		spreadSheetShortSummary,
		spreadSheetBettingTicketSummary,
		spreadSheetGradeClassSummary,
		spreadSheetMonthlySummary,
		spreadSheetCourseCategorySummary,
		spreadSheetDistanceCategorySummary,
		spreadSheetRaceCourseSummary,
	)
}

func (a *Aggregator) GetyMonthlyBettingTicketSummary(
	records []*betting_ticket_entity.CsvEntity,
	racingNumbers []*race_entity.RacingNumber,
	races []*race_entity.Race,
) *spreadsheet_entity.SpreadSheetMonthlyBettingTicketSummary {
	return spreadsheet_entity.NewSpreadSheetMonthlyBettingTicketSummary(a.summarizer.GetMonthlyBettingTicketSummary(records, racingNumbers, races))
}
