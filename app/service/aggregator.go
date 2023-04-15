package service

import (
	"fmt"
	betting_ticket_entity "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/entity"
	betting_ticket_vo "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/value_object"
	race_entity "github.com/mapserver2007/ipat-aggregator/app/domain/race/entity"
	race_vo "github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
	spreadsheet_entity "github.com/mapserver2007/ipat-aggregator/app/domain/spreadsheet/entity"
	"sort"
)

type Aggregator struct {
	raceConverter          RaceConverter
	bettingTicketConverter BettingTicketConverter
}

func NewAggregator(
	raceConverter RaceConverter,
	bettingTicketConverter BettingTicketConverter,
) *Aggregator {
	return &Aggregator{
		raceConverter:          raceConverter,
		bettingTicketConverter: bettingTicketConverter,
	}
}

func (a *Aggregator) GetSummary(
	records []*betting_ticket_entity.CsvEntity,
	racingNumbers []*race_entity.RacingNumber,
	races []*race_entity.Race,
) *spreadsheet_entity.Summary {
	return spreadsheet_entity.NewResult(
		a.getTotalResultSummary(records),
		a.getLatestMonthlyResultSummary(records),
		a.getLatestYearResultSummary(records),
		a.getBettingTicketSummary(records),
		a.getRaceClassSummary(records, races),
		a.getMonthlySummary(records),
		a.getYearlySummary(records),
		a.getCourseCategorySummary(records, racingNumbers, races),
		a.getDistanceCategorySummary(records, racingNumbers, races),
		a.getRaceCourseSummary(records, racingNumbers, races),
	)
}

func (a *Aggregator) getTotalResultSummary(records []*betting_ticket_entity.CsvEntity) spreadsheet_entity.ResultSummary {
	resultRate := a.getTotalBettingTicketRate(records)
	return spreadsheet_entity.NewResultSummary(resultRate.Payments, resultRate.Repayments)
}

func (a *Aggregator) getLatestMonthlyResultSummary(records []*betting_ticket_entity.CsvEntity) spreadsheet_entity.ResultSummary {
	resultRate := a.getLatestMonthlyBettingTicketRate(records)
	return spreadsheet_entity.NewResultSummary(resultRate.Payments, resultRate.Repayments)
}

func (a *Aggregator) getLatestYearResultSummary(records []*betting_ticket_entity.CsvEntity) spreadsheet_entity.ResultSummary {
	resultRate := a.getLatestYearBettingTicketRate(records)
	return spreadsheet_entity.NewResultSummary(resultRate.Payments, resultRate.Repayments)
}

func (a *Aggregator) getBettingTicketSummary(records []*betting_ticket_entity.CsvEntity) spreadsheet_entity.BettingTicketSummary {
	return spreadsheet_entity.NewBettingTicketSummary(a.getBettingTicketResultRate(records))
}

func (a *Aggregator) getRaceClassSummary(records []*betting_ticket_entity.CsvEntity, races []*race_entity.Race) spreadsheet_entity.RaceClassSummary {
	return spreadsheet_entity.NewRaceClassSummary(a.getRaceClassResultRate(records, races))
}

func (a *Aggregator) getMonthlySummary(records []*betting_ticket_entity.CsvEntity) spreadsheet_entity.MonthlySummary {
	return spreadsheet_entity.NewMonthlySummary(a.getMonthlyResultRate(records))
}

func (a *Aggregator) getYearlySummary(records []*betting_ticket_entity.CsvEntity) spreadsheet_entity.YearlySummary {
	return spreadsheet_entity.NewYearlySummary(a.getYearlyResultRate(records))
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

func (a *Aggregator) getTotalBettingTicketRate(records []*betting_ticket_entity.CsvEntity) spreadsheet_entity.ResultRate {
	return CalcSumResultRate(records)
}

func (a *Aggregator) getLatestMonthlyBettingTicketRate(records []*betting_ticket_entity.CsvEntity) spreadsheet_entity.ResultRate {
	monthlyRatesMap := a.getMonthlyResultRate(records)

	var dateList []int
	for date := range monthlyRatesMap {
		dateList = append(dateList, date)
	}
	sort.Slice(dateList, func(i, j int) bool {
		return dateList[i] > dateList[j]
	})
	latestDate := dateList[0]

	resultRate, _ := monthlyRatesMap[latestDate]

	return resultRate
}

func (a *Aggregator) getLatestYearBettingTicketRate(records []*betting_ticket_entity.CsvEntity) spreadsheet_entity.ResultRate {
	yearlyRatesMap := a.getYearlyResultRate(records)

	var dateList []int
	for date := range yearlyRatesMap {
		dateList = append(dateList, date)
	}
	sort.Slice(dateList, func(i, j int) bool {
		return dateList[i] > dateList[j]
	})
	latestDate := dateList[0]

	resultRate, _ := yearlyRatesMap[latestDate]

	return resultRate
}

func (a *Aggregator) getBettingTicketResultRate(records []*betting_ticket_entity.CsvEntity) map[betting_ticket_vo.BettingTicket]spreadsheet_entity.ResultRate {
	bettingTicketRatesMap := map[betting_ticket_vo.BettingTicket]spreadsheet_entity.ResultRate{}

	for bettingTicket, recordsGroup := range a.bettingTicketConverter.ConvertToBettingTicketRecordsMap(records) {
		bettingTicketRatesMap[bettingTicket] = CalcSumResultRate(recordsGroup)
	}

	// 同一券種をまとめる
	mergeFunc := func(o1, o2 spreadsheet_entity.ResultRate) spreadsheet_entity.ResultRate {
		o1.HitCount += o2.HitCount
		o1.VoteCount += o2.VoteCount
		o1.Payments += o2.Payments
		o1.Repayments += o2.Repayments
		return o1
	}
	newBettingTicketRatesMap := map[betting_ticket_vo.BettingTicket]spreadsheet_entity.ResultRate{}
	for kind := range bettingTicketRatesMap {
		switch kind {
		case betting_ticket_vo.QuinellaPlaceWheel:
			newBettingTicketRatesMap[betting_ticket_vo.Quinella] = mergeFunc(newBettingTicketRatesMap[betting_ticket_vo.Quinella], bettingTicketRatesMap[kind])
		case betting_ticket_vo.TrioFormation, betting_ticket_vo.TrioWheelOfFirst:
			newBettingTicketRatesMap[betting_ticket_vo.Trio] = mergeFunc(newBettingTicketRatesMap[betting_ticket_vo.Trio], bettingTicketRatesMap[kind])
		case betting_ticket_vo.TrifectaFormation, betting_ticket_vo.TrifectaWheelOfFirst:
			newBettingTicketRatesMap[betting_ticket_vo.Trifecta] = mergeFunc(newBettingTicketRatesMap[betting_ticket_vo.Trifecta], bettingTicketRatesMap[kind])
		default:
			if kind == betting_ticket_vo.Quinella || kind == betting_ticket_vo.Trio || kind == betting_ticket_vo.Trifecta {
				newBettingTicketRatesMap[kind] = mergeFunc(newBettingTicketRatesMap[kind], bettingTicketRatesMap[kind])
			} else {
				newBettingTicketRatesMap[kind] = bettingTicketRatesMap[kind]
			}
		}
	}

	return newBettingTicketRatesMap
}

func (a *Aggregator) getMonthlyResultRate(records []*betting_ticket_entity.CsvEntity) map[int]spreadsheet_entity.ResultRate {
	monthlyRatesMap := map[int]spreadsheet_entity.ResultRate{}
	for date, recordsGroup := range a.bettingTicketConverter.ConvertToMonthRecordsMap(records) {
		monthlyRatesMap[date] = CalcSumResultRate(recordsGroup)
	}

	return monthlyRatesMap
}

func (a *Aggregator) getYearlyResultRate(records []*betting_ticket_entity.CsvEntity) map[int]spreadsheet_entity.ResultRate {
	yearlyRatesMap := map[int]spreadsheet_entity.ResultRate{}
	for date, recordsGroup := range a.bettingTicketConverter.ConvertToYearRecordsMap(records) {
		yearlyRatesMap[date] = CalcSumResultRate(recordsGroup)
	}

	return yearlyRatesMap
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

func (a *Aggregator) getRaceClassResultRate(records []*betting_ticket_entity.CsvEntity, races []*race_entity.Race) map[race_vo.GradeClass]spreadsheet_entity.ResultRate {
	raceClassRatesMap := map[race_vo.GradeClass]spreadsheet_entity.ResultRate{}
	raceMap := a.raceConverter.ConvertToRaceMapByRacingNumberId(races)
	raceClassMap := a.bettingTicketConverter.ConvertToRaceClassRecordsMap(records, raceMap)
	for raceClass, records := range raceClassMap {
		raceClassRatesMap[raceClass] = CalcSumResultRate(records)
	}

	// クラスをまとめる
	mergeFunc := func(o1, o2 spreadsheet_entity.ResultRate) spreadsheet_entity.ResultRate {
		o1.HitCount += o2.HitCount
		o1.VoteCount += o2.VoteCount
		o1.Payments += o2.Payments
		o1.Repayments += o2.Repayments
		return o1
	}
	newRaceClassRatesMap := map[race_vo.GradeClass]spreadsheet_entity.ResultRate{}
	for raceClass := range raceClassRatesMap {
		switch raceClass {
		case race_vo.Jpn1, race_vo.JumpGrade1:
			newRaceClassRatesMap[race_vo.Grade1] = mergeFunc(newRaceClassRatesMap[race_vo.Grade1], raceClassRatesMap[raceClass])
		case race_vo.Jpn2, race_vo.JumpGrade2:
			newRaceClassRatesMap[race_vo.Grade2] = mergeFunc(newRaceClassRatesMap[race_vo.Grade2], raceClassRatesMap[raceClass])
		case race_vo.Jpn3, race_vo.JumpGrade3:
			newRaceClassRatesMap[race_vo.Grade3] = mergeFunc(newRaceClassRatesMap[race_vo.Grade3], raceClassRatesMap[raceClass])
		case race_vo.ListedClass, race_vo.OpenClass, race_vo.AllowanceClass:
			newRaceClassRatesMap[race_vo.NonGradeClass] = mergeFunc(newRaceClassRatesMap[race_vo.NonGradeClass], raceClassRatesMap[raceClass])
		default:
			newRaceClassRatesMap[raceClass] = raceClassRatesMap[raceClass]
		}
	}

	return newRaceClassRatesMap
}
