package service

import (
	"fmt"
	betting_ticket_entity "github.com/mapserver2007/tools/baken/app/domain/betting_ticket/entity"
	betting_ticket_vo "github.com/mapserver2007/tools/baken/app/domain/betting_ticket/value_object"
	race_entity "github.com/mapserver2007/tools/baken/app/domain/race/entity"
	race_vo "github.com/mapserver2007/tools/baken/app/domain/race/value_object"
	spreadsheet_entity "github.com/mapserver2007/tools/baken/app/domain/spreadsheet/entity"
	"sort"
)

type Aggregator struct {
	raceConverter          RaceConverter
	bettingTicketConverter BettingTicketConverter
	entities               []*betting_ticket_entity.CsvEntity
	racingNumberInfo       *race_entity.RacingNumberInfo
	raceInfo               *race_entity.RaceInfo
}

func NewAggregator(
	raceConverter RaceConverter,
	bettingTicketConverter BettingTicketConverter,
	entities []*betting_ticket_entity.CsvEntity,
	racingNumberInfo *race_entity.RacingNumberInfo,
	raceInfo *race_entity.RaceInfo,
) Aggregator {
	return Aggregator{
		raceConverter:          raceConverter,
		bettingTicketConverter: bettingTicketConverter,
		entities:               entities,
		racingNumberInfo:       racingNumberInfo,
		raceInfo:               raceInfo,
	}
}

func (a *Aggregator) GetSummary() *spreadsheet_entity.Summary {
	return spreadsheet_entity.NewResult(
		a.getTotalResultSummary(),
		a.getLatestMonthlyResultSummary(),
		a.getBettingTicketSummary(),
		a.getRaceClassSummary(),
		a.getMonthlySummary(),
		a.getCourseCategorySummary(),
		a.getDistanceCategorySummary(),
		a.getRaceCourseSummary(),
	)
}

func (a *Aggregator) getTotalResultSummary() spreadsheet_entity.ResultSummary {
	resultRate := a.getTotalBettingTicketRate()
	return spreadsheet_entity.NewResultSummary(resultRate.Payments, resultRate.Repayments)
}

func (a *Aggregator) getLatestMonthlyResultSummary() spreadsheet_entity.ResultSummary {
	resultRate := a.getLatestMonthlyBettingTicketRate()
	return spreadsheet_entity.NewResultSummary(resultRate.Payments, resultRate.Repayments)
}

func (a *Aggregator) getBettingTicketSummary() spreadsheet_entity.BettingTicketSummary {
	return spreadsheet_entity.NewBettingTicketSummary(a.getBettingTicketResultRate())
}

func (a *Aggregator) getRaceClassSummary() spreadsheet_entity.RaceClassSummary {
	return spreadsheet_entity.NewRaceClassSummary(a.getRaceClassResultRate())
}

func (a *Aggregator) getMonthlySummary() spreadsheet_entity.MonthlySummary {
	return spreadsheet_entity.NewMonthlySummary(a.getMonthlyResultRate())
}

func (a *Aggregator) getCourseCategorySummary() spreadsheet_entity.CourseCategorySummary {
	return spreadsheet_entity.NewCourseCategorySummary(a.getCourseCategoryRates())
}

func (a *Aggregator) getDistanceCategorySummary() spreadsheet_entity.DistanceCategorySummary {
	return spreadsheet_entity.NewDistanceCategorySummary(a.getDistanceCategoryRates())
}

func (a *Aggregator) getRaceCourseSummary() spreadsheet_entity.RaceCourseSummary {
	return spreadsheet_entity.NewRaceCourseSummary(a.getRaceCourseRates())
}

func (a *Aggregator) getTotalBettingTicketRate() spreadsheet_entity.ResultRate {
	return CalcSumResultRate(a.entities)
}

func (a *Aggregator) getLatestMonthlyBettingTicketRate() spreadsheet_entity.ResultRate {
	monthlyRatesMap := a.getMonthlyResultRate()

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

func (a *Aggregator) getBettingTicketResultRate() map[betting_ticket_vo.BettingTicket]spreadsheet_entity.ResultRate {
	bettingTicketRatesMap := map[betting_ticket_vo.BettingTicket]spreadsheet_entity.ResultRate{}

	for bettingTicket, records := range a.bettingTicketConverter.ConvertToBettingTicketRecordsMap(a.entities) {
		bettingTicketRatesMap[bettingTicket] = CalcSumResultRate(records)
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

func (a *Aggregator) getMonthlyResultRate() map[int]spreadsheet_entity.ResultRate {
	monthlyRatesMap := map[int]spreadsheet_entity.ResultRate{}
	for date, records := range a.bettingTicketConverter.ConvertToMonthRecordsMap(a.entities) {
		monthlyRatesMap[date] = CalcSumResultRate(records)
	}

	return monthlyRatesMap
}

func (a *Aggregator) getCourseCategoryRates() map[race_vo.CourseCategory]spreadsheet_entity.ResultRate {
	courseCategoryRatesMap := map[race_vo.CourseCategory]spreadsheet_entity.ResultRate{}
	courseCategoryRecordsMap := map[race_vo.CourseCategory][]*betting_ticket_entity.CsvEntity{}
	raceMap := a.raceConverter.ConvertToRaceMapByRaceId(a.raceInfo.Races())
	racingNumberMap := a.raceConverter.ConvertToRacingNumberMap(a.racingNumberInfo.RacingNumbers())
	for _, entity := range a.entities {
		racingNumberId := race_vo.NewRacingNumberId(entity.RaceDate(), entity.RaceCourse())
		racingNumber, ok := racingNumberMap[racingNumberId]
		if !ok {
			panic(fmt.Errorf("unknown racingNumberId: %s", racingNumberId))
		}
		raceId := a.raceConverter.GetRaceId(entity, racingNumber)
		if race, ok := raceMap[*raceId]; ok {
			courseCategory := race.CourseCategory()
			courseCategoryRecordsMap[courseCategory] = append(courseCategoryRecordsMap[courseCategory], entity)
		}
	}
	for courseCategory, records := range courseCategoryRecordsMap {
		courseCategoryRatesMap[courseCategory] = CalcSumResultRate(records)
	}

	return courseCategoryRatesMap
}

func (a *Aggregator) getDistanceCategoryRates() map[race_vo.DistanceCategory]spreadsheet_entity.ResultRate {
	distanceCategoryRatesMap := map[race_vo.DistanceCategory]spreadsheet_entity.ResultRate{}
	distanceCategoryRecordsMap := map[race_vo.DistanceCategory][]*betting_ticket_entity.CsvEntity{}
	raceMap := a.raceConverter.ConvertToRaceMapByRaceId(a.raceInfo.Races())
	racingNumberMap := a.raceConverter.ConvertToRacingNumberMap(a.racingNumberInfo.RacingNumbers())
	for _, entity := range a.entities {
		racingNumberId := race_vo.NewRacingNumberId(entity.RaceDate(), entity.RaceCourse())
		racingNumber, ok := racingNumberMap[racingNumberId]
		if !ok {
			panic(fmt.Errorf("unknown racingNumberId: %s", racingNumberId))
		}
		raceId := a.raceConverter.GetRaceId(entity, racingNumber)
		if race, ok := raceMap[*raceId]; ok {
			courseCategory := race.CourseCategory()
			distanceCategory := race_vo.NewDistanceCategory(race.Distance(), courseCategory)
			distanceCategoryRecordsMap[distanceCategory] = append(distanceCategoryRecordsMap[distanceCategory], entity)
		}
	}
	for distanceCategory, records := range distanceCategoryRecordsMap {
		distanceCategoryRatesMap[distanceCategory] = CalcSumResultRate(records)
	}

	return distanceCategoryRatesMap
}

func (a *Aggregator) getRaceCourseRates() map[race_vo.RaceCourse]spreadsheet_entity.ResultRate {
	raceCourseRatesMap := map[race_vo.RaceCourse]spreadsheet_entity.ResultRate{}
	raceCourseRecordsMap := map[race_vo.RaceCourse][]*betting_ticket_entity.CsvEntity{}
	raceMap := a.raceConverter.ConvertToRaceMapByRaceId(a.raceInfo.Races())
	racingNumberMap := a.raceConverter.ConvertToRacingNumberMap(a.racingNumberInfo.RacingNumbers())
	for _, entity := range a.entities {
		racingNumberId := race_vo.NewRacingNumberId(entity.RaceDate(), entity.RaceCourse())
		racingNumber, ok := racingNumberMap[racingNumberId]
		if !ok {
			panic(fmt.Errorf("unknown racingNumberId: %s", racingNumberId))
		}
		raceId := a.raceConverter.GetRaceId(entity, racingNumber)
		if race, ok := raceMap[*raceId]; ok {
			raceCourse := race.RaceCourseId()
			raceCourseRecordsMap[raceCourse] = append(raceCourseRecordsMap[raceCourse], entity)
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
		case race_vo.Longchamp, race_vo.Deauville, race_vo.Shatin:
			newRaceCourseRatesMap[race_vo.Overseas] = mergeFunc(newRaceCourseRatesMap[race_vo.Overseas], raceCourseRatesMap[raceCourse])
		default:
			if raceCourse == race_vo.Longchamp || raceCourse == race_vo.Deauville || raceCourse == race_vo.Shatin {
				newRaceCourseRatesMap[raceCourse] = mergeFunc(newRaceCourseRatesMap[raceCourse], raceCourseRatesMap[raceCourse])
			} else {
				newRaceCourseRatesMap[raceCourse] = raceCourseRatesMap[raceCourse]
			}
		}
	}

	return newRaceCourseRatesMap
}

func (a *Aggregator) getRaceClassResultRate() map[race_vo.GradeClass]spreadsheet_entity.ResultRate {
	raceClassRatesMap := map[race_vo.GradeClass]spreadsheet_entity.ResultRate{}
	raceMap := a.raceConverter.ConvertToRaceMapByRacingNumberId(a.raceInfo.Races())
	raceClassMap := a.bettingTicketConverter.ConvertToRaceClassRecordsMap(a.entities, raceMap)
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
