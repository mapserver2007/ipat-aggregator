package service

import (
	"fmt"
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
) (*spreadsheet_entity.Summary, *spreadsheet_entity.SpreadSheetSummary, *spreadsheet_entity.SpreadSheetBettingTicketSummary, *spreadsheet_entity.SpreadSheetClassSummary, *spreadsheet_entity.SpreadSheetMonthlySummary, *spreadsheet_entity.SpreadSheetCourseCategorySummary) {

	// TODO averageがなんかあってない気がする。多分単体の馬券で計算している。averageは1レース単位。

	// TODO 移行中なのでいろいろ混在
	spreadSheetSummary := spreadsheet_entity.NewSpreadSheetSummary(
		a.summarizer.GetShortSummaryForAll(records),
		a.summarizer.GetShortSummaryForMonth(records),
		a.summarizer.GetShortSummaryForYear(records),
	)

	spreadSheetBettingTicketSummary := spreadsheet_entity.NewSpreadSheetBettingTicketSummary(
		a.summarizer.GetBettingTicketSummaryForAll(records, betting_ticket_vo.Win),
		a.summarizer.GetBettingTicketSummaryForAll(records, betting_ticket_vo.Place),
		a.summarizer.GetBettingTicketSummaryForAll(records, betting_ticket_vo.Quinella),
		a.summarizer.GetBettingTicketSummaryForAll(records, betting_ticket_vo.Exacta, betting_ticket_vo.ExactaWheelOfFirst),
		a.summarizer.GetBettingTicketSummaryForAll(records, betting_ticket_vo.QuinellaPlace, betting_ticket_vo.QuinellaPlaceWheel),
		a.summarizer.GetBettingTicketSummaryForAll(records, betting_ticket_vo.Trio, betting_ticket_vo.TrioFormation, betting_ticket_vo.TrioWheelOfFirst),
		a.summarizer.GetBettingTicketSummaryForAll(records, betting_ticket_vo.Trifecta, betting_ticket_vo.TrifectaFormation, betting_ticket_vo.TrifectaWheelOfFirst),
		a.summarizer.GetBettingTicketSummaryForAll(records, betting_ticket_vo.Win, betting_ticket_vo.Place, betting_ticket_vo.Quinella,
			betting_ticket_vo.Exacta, betting_ticket_vo.ExactaWheelOfFirst, betting_ticket_vo.QuinellaPlace, betting_ticket_vo.QuinellaPlaceWheel,
			betting_ticket_vo.Trio, betting_ticket_vo.TrioFormation, betting_ticket_vo.TrioWheelOfFirst,
			betting_ticket_vo.Trifecta, betting_ticket_vo.TrifectaFormation, betting_ticket_vo.TrifectaWheelOfFirst),
	)

	spreadSheetGradeClassSummary := spreadsheet_entity.NewSpreadSheetClassSummary(
		a.summarizer.GetGradeClassSummaryForAll(records, races, race_vo.Grade1, race_vo.Jpn1, race_vo.JumpGrade1),
		a.summarizer.GetGradeClassSummaryForAll(records, races, race_vo.Grade2, race_vo.Jpn2, race_vo.JumpGrade2),
		a.summarizer.GetGradeClassSummaryForAll(records, races, race_vo.Grade3, race_vo.Jpn3, race_vo.JumpGrade3),
		a.summarizer.GetGradeClassSummaryForAll(records, races, race_vo.OpenClass, race_vo.ListedClass, race_vo.AllowanceClass),
	)

	spreadSheetMonthlySummary := spreadsheet_entity.NewSpreadSheetMonthlySummary(a.summarizer.GetMonthlySummaryMap(records))

	spreadSheetCourseCategorySummary := spreadsheet_entity.NewSpreadSheetCourseCategorySummary(
		a.summarizer.GetCourseCategorySummaryForAll(records, racingNumbers, races, race_vo.Turf),
		a.summarizer.GetCourseCategorySummaryForAll(records, racingNumbers, races, race_vo.Dirt),
		a.summarizer.GetCourseCategorySummaryForAll(records, racingNumbers, races, race_vo.Jump),
	)

	return spreadsheet_entity.NewResult(
		spreadsheet_entity.BettingTicketSummary{},
		spreadsheet_entity.RaceClassSummary{},
		spreadsheet_entity.MonthlySummary{},
		spreadsheet_entity.YearlySummary{},
		a.getCourseCategorySummary(records, racingNumbers, races),
		a.getDistanceCategorySummary(records, racingNumbers, races),
		a.getRaceCourseSummary(records, racingNumbers, races),
	), spreadSheetSummary, spreadSheetBettingTicketSummary, spreadSheetGradeClassSummary, spreadSheetMonthlySummary, spreadSheetCourseCategorySummary
}

func (a *Aggregator) getCourseCategorySummary(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race) spreadsheet_entity.CourseCategorySummary {
	return spreadsheet_entity.NewCourseCategorySummary(a.getCourseCategoryRates(records, racingNumbers, races))
}

func (a *Aggregator) getDistanceCategorySummary(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race) spreadsheet_entity.DistanceCategorySummary {
	return spreadsheet_entity.NewDistanceCategorySummary(a.getDistanceCategoryRates(records, racingNumbers, races))
}

func (a *Aggregator) getRaceCourseSummary(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race) spreadsheet_entity.RaceCourseSummary {
	return spreadsheet_entity.NewRaceCourseSummary(a.getRaceCourseRates(records, racingNumbers, races))
}

func (a *Aggregator) getCourseCategoryRates(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race) map[race_vo.CourseCategory]spreadsheet_entity.ResultRate {
	courseCategoryRatesMap := map[race_vo.CourseCategory]spreadsheet_entity.ResultRate{}
	courseCategoryRecordsMap := map[race_vo.CourseCategory][]*betting_ticket_entity.CsvEntity{}
	raceMap := a.raceConverter.ConvertToRaceMapByRaceId(races)
	racingNumberMap := a.raceConverter.ConvertToRacingNumberMap(racingNumbers)
	for _, record := range records {
		racingNumberId := race_vo.NewRacingNumberId(record.RaceDate(), record.RaceCourse())
		racingNumber, ok := racingNumberMap[racingNumberId]
		if !ok && record.RaceCourse().Organizer() == race_vo.JRA {
			panic(fmt.Errorf("unknown racingNumberId: %s", racingNumberId))
		}
		raceId := a.raceConverter.GetRaceId(record, racingNumber)
		if race, ok := raceMap[*raceId]; ok {
			courseCategory := race.CourseCategory()
			courseCategoryRecordsMap[courseCategory] = append(courseCategoryRecordsMap[courseCategory], record)
		}
	}
	for courseCategory, records := range courseCategoryRecordsMap {
		courseCategoryRatesMap[courseCategory] = CalcSumResultRate(records)
	}

	return courseCategoryRatesMap
}

func (a *Aggregator) getDistanceCategoryRates(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race) map[race_vo.DistanceCategory]spreadsheet_entity.ResultRate {
	distanceCategoryRatesMap := map[race_vo.DistanceCategory]spreadsheet_entity.ResultRate{}
	distanceCategoryRecordsMap := map[race_vo.DistanceCategory][]*betting_ticket_entity.CsvEntity{}
	raceMap := a.raceConverter.ConvertToRaceMapByRaceId(races)
	racingNumberMap := a.raceConverter.ConvertToRacingNumberMap(racingNumbers)
	for _, record := range records {
		racingNumberId := race_vo.NewRacingNumberId(record.RaceDate(), record.RaceCourse())
		racingNumber, ok := racingNumberMap[racingNumberId]
		if !ok && record.RaceCourse().Organizer() == race_vo.JRA {
			panic(fmt.Errorf("unknown racingNumberId: %s", racingNumberId))
		}
		raceId := a.raceConverter.GetRaceId(record, racingNumber)
		if race, ok := raceMap[*raceId]; ok {
			courseCategory := race.CourseCategory()
			distanceCategory := race_vo.NewDistanceCategory(race.Distance(), courseCategory)
			distanceCategoryRecordsMap[distanceCategory] = append(distanceCategoryRecordsMap[distanceCategory], record)
		}
	}
	for distanceCategory, records := range distanceCategoryRecordsMap {
		distanceCategoryRatesMap[distanceCategory] = CalcSumResultRate(records)
	}

	return distanceCategoryRatesMap
}

func (a *Aggregator) getRaceCourseRates(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race) map[race_vo.RaceCourse]spreadsheet_entity.ResultRate {
	raceCourseRatesMap := map[race_vo.RaceCourse]spreadsheet_entity.ResultRate{}
	raceCourseRecordsMap := map[race_vo.RaceCourse][]*betting_ticket_entity.CsvEntity{}
	raceMap := a.raceConverter.ConvertToRaceMapByRaceId(races)
	racingNumberMap := a.raceConverter.ConvertToRacingNumberMap(racingNumbers)
	for _, record := range records {
		racingNumberId := race_vo.NewRacingNumberId(record.RaceDate(), record.RaceCourse())
		racingNumber, ok := racingNumberMap[racingNumberId]
		if !ok && record.RaceCourse().Organizer() == race_vo.JRA {
			panic(fmt.Errorf("unknown racingNumberId: %s", racingNumberId))
		}
		raceId := a.raceConverter.GetRaceId(record, racingNumber)
		if race, ok := raceMap[*raceId]; ok {
			raceCourse := race.RaceCourseId()
			raceCourseRecordsMap[raceCourse] = append(raceCourseRecordsMap[raceCourse], record)
		}
	}
	for raceCourse, records := range raceCourseRecordsMap {
		raceCourseRatesMap[raceCourse] = CalcSumResultRate(records)
	}

	// 開催場所をまとめる
	mergeFunc := func(o1, o2 spreadsheet_entity.ResultRate) spreadsheet_entity.ResultRate {
		o1.HitCount += o2.HitCount
		o1.VoteCount += o2.VoteCount
		o1.Payments += o2.Payments
		o1.Repayments += o2.Repayments
		return o1
	}
	newRaceCourseRatesMap := map[race_vo.RaceCourse]spreadsheet_entity.ResultRate{}
	for raceCourse := range raceCourseRatesMap {
		switch raceCourse {
		case race_vo.Longchamp, race_vo.Deauville, race_vo.Shatin, race_vo.Meydan:
			newRaceCourseRatesMap[race_vo.Overseas] = mergeFunc(newRaceCourseRatesMap[race_vo.Overseas], raceCourseRatesMap[raceCourse])
		default:
			if raceCourse == race_vo.Longchamp || raceCourse == race_vo.Deauville || raceCourse == race_vo.Shatin || raceCourse == race_vo.Meydan {
				newRaceCourseRatesMap[raceCourse] = mergeFunc(newRaceCourseRatesMap[raceCourse], raceCourseRatesMap[raceCourse])
			} else {
				newRaceCourseRatesMap[raceCourse] = raceCourseRatesMap[raceCourse]
			}
		}
	}

	return newRaceCourseRatesMap
}
