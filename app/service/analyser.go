package service

import (
	"fmt"
	analyse_entity "github.com/mapserver2007/ipat-aggregator/app/domain/analyse/entity"
	betting_ticket_entity "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/entity"
	betting_ticket_vo "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/value_object"
	race_entity "github.com/mapserver2007/ipat-aggregator/app/domain/race/entity"
	race_vo "github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
	"math"
	"sort"
)

type Analyser struct {
	raceConverter RaceConverter
}

func NewAnalyser(
	raceConverter RaceConverter,
) Analyser {
	return Analyser{
		raceConverter: raceConverter,
	}
}

func (a *Analyser) WinAnalyse(
	records []*betting_ticket_entity.CsvEntity,
	racingNumbers []*race_entity.RacingNumber,
	races []*race_entity.Race,
) *analyse_entity.WinAnalyseSummary {
	raceMap := a.raceConverter.ConvertToRaceMapByRaceId(races)
	racingNumberMap := a.raceConverter.ConvertToRacingNumberMap(racingNumbers)

	// TODO とりあえず券種別にメソッドを割るが、後々リファクタリングするかも
	recordsByWin := a.getRecordsByWin(records)
	popularMap := map[int][]*analyse_entity.PopularAnalyse{}
	for i := 1; i <= 18; i++ {
		popularMap[i] = []*analyse_entity.PopularAnalyse{}
	}

	for _, record := range recordsByWin {
		racingNumberId := race_vo.NewRacingNumberId(record.RaceDate(), record.RaceCourse())
		racingNumber, ok := racingNumberMap[racingNumberId]
		if !ok && record.RaceCourse().Organizer() == race_vo.JRA {
			panic(fmt.Errorf("unknown racingNumberId: %s", racingNumberId))
		}
		raceId := a.raceConverter.GetRaceId(record, racingNumber)
		if race, ok := raceMap[*raceId]; ok {
			popular := a.getPopularAnalyse(record, race)
			popularMap[popular.PopularNumber()] = append(popularMap[popular.PopularNumber()], popular)
		}
	}

	allSummaries := a.convertWinPopularAnalyzeSummary(popularMap)
	grade1Summaries, grade2Summaries, grade3Summaries, allowanceClassSummaries := a.convertClassWinPopularAnalyzeSummaries(popularMap)

	return analyse_entity.NewWinAnalyseSummary(allSummaries, grade1Summaries, grade2Summaries, grade3Summaries, allowanceClassSummaries)
}

func (a *Analyser) convertWinPopularAnalyzeSummary(popularMap map[int][]*analyse_entity.PopularAnalyse) []*analyse_entity.PopularAnalyseSummary {
	popularAnalyseSummaries := make([]*analyse_entity.PopularAnalyseSummary, 0, 18)
	for popularNumber := 1; popularNumber <= 18; popularNumber++ {
		populars, ok := popularMap[popularNumber]
		if !ok {
			popularAnalyseSummaries = append(popularAnalyseSummaries, analyse_entity.DefaultPopularAnalyseSummary(popularNumber))
			continue
		}
		var (
			hitCount                                                int
			hitRate                                                 float64
			totalOddsAtVote, totalOddsAtHit                         float64
			averageOddsAtVote, averageOddsAtHit, averageOddsAtUnHit float64
			totalPayment, totalPayout                               int
			averagePayment, averagePayout                           int
			medianPayment, medianPayout                             int
			maxPayout, minPayout                                    int
			maxOddsAtHit, minOddsAtHit                              float64
			allPayments, allPayouts                                 []int
		)

		betCount := len(populars)

		for _, popular := range populars {
			totalOddsAtVote += popular.Odds()
			totalPayment += popular.Payment()
			totalPayout += popular.Payout()
			allPayments = append(allPayments, popular.Payment())
			if popular.IsHit() {
				hitCount++
				totalOddsAtHit += popular.Odds()
				allPayouts = append(allPayouts, popular.Payout())
			}
			if popular.Payout() > maxPayout {
				maxPayout = popular.Payout()
			}
			if (popular.Payout() > 0 && popular.Payout() < minPayout) || minPayout == 0 {
				minPayout = popular.Payout()
			}
			if (popular.Odds() > maxOddsAtHit) && popular.IsHit() {
				maxOddsAtHit = popular.Odds()
			}
			if ((popular.Odds() > 0 && popular.Odds() < minOddsAtHit) || minOddsAtHit == 0) && popular.IsHit() {
				minOddsAtHit = popular.Odds()
			}
		}

		if betCount > 0 {
			hitRate = math.Round((float64(hitCount)/float64(betCount))*100) / 100
			averageOddsAtVote = math.Round((totalOddsAtVote/float64(betCount))*10) / 10
			averagePayment = totalPayment / betCount
		}
		if hitCount > 0 {
			averageOddsAtHit = math.Round((totalOddsAtHit/float64(hitCount))*10) / 10
			averagePayout = totalPayout / hitCount
		}
		unHitCount := betCount - hitCount
		if unHitCount > 0 {
			averageOddsAtUnHit = math.Round(((totalOddsAtVote-totalOddsAtHit)/float64(unHitCount))*10) / 10
		}

		if len(allPayments) > 0 {
			if len(allPayments)%2 == 0 {
				medianPayment = (allPayments[len(allPayments)/2] + allPayments[len(allPayments)/2-1]) / 2
			} else {
				medianPayment = allPayments[len(allPayments)/2]
			}
		}
		if len(allPayouts) > 0 {
			if len(allPayouts) > 0 && len(allPayouts)%2 == 0 {
				medianPayout = (allPayouts[len(allPayouts)/2] + allPayouts[len(allPayouts)/2-1]) / 2
			} else {
				medianPayout = allPayouts[len(allPayouts)/2]
			}
		}

		popularAnalyseSummary := analyse_entity.NewPopularAnalyseSummary(
			popularNumber,
			betCount,
			hitCount,
			hitRate,
			averageOddsAtVote,
			averageOddsAtHit,
			averageOddsAtUnHit,
			totalPayment,
			totalPayout,
			averagePayment,
			averagePayout,
			medianPayment,
			medianPayout,
			maxPayout,
			minPayout,
			maxOddsAtHit,
			minOddsAtHit,
		)
		popularAnalyseSummaries = append(popularAnalyseSummaries, popularAnalyseSummary)
	}

	sort.Slice(popularAnalyseSummaries, func(i, j int) bool {
		return popularAnalyseSummaries[i].PopularNumber() < popularAnalyseSummaries[j].PopularNumber()
	})

	return popularAnalyseSummaries
}

func (a *Analyser) convertClassWinPopularAnalyzeSummaries(
	popularMap map[int][]*analyse_entity.PopularAnalyse,
) (
	[]*analyse_entity.PopularAnalyseSummary,
	[]*analyse_entity.PopularAnalyseSummary,
	[]*analyse_entity.PopularAnalyseSummary,
	[]*analyse_entity.PopularAnalyseSummary,
) {
	var (
		grade1Summaries, grade2Summaries, grade3Summaries, allowanceClassSummaries []*analyse_entity.PopularAnalyseSummary
	)
	grade1PopularMap := map[int][]*analyse_entity.PopularAnalyse{}
	grade2PopularMap := map[int][]*analyse_entity.PopularAnalyse{}
	grade3PopularMap := map[int][]*analyse_entity.PopularAnalyse{}
	allowanceClassPopularMap := map[int][]*analyse_entity.PopularAnalyse{}

	for popularNumber, populars := range popularMap {
		for _, popular := range populars {
			switch popular.Class() {
			case race_vo.Grade1:
				if _, ok := grade1PopularMap[popularNumber]; !ok {
					grade1PopularMap[popularNumber] = make([]*analyse_entity.PopularAnalyse, 0)
				} else {
					grade1PopularMap[popularNumber] = append(grade1PopularMap[popularNumber], popular)
				}
			case race_vo.Grade2:
				if _, ok := grade2PopularMap[popularNumber]; !ok {
					grade2PopularMap[popularNumber] = make([]*analyse_entity.PopularAnalyse, 0)
				} else {
					grade2PopularMap[popularNumber] = append(grade2PopularMap[popularNumber], popular)
				}
			case race_vo.Grade3:
				if _, ok := grade3PopularMap[popularNumber]; !ok {
					grade3PopularMap[popularNumber] = make([]*analyse_entity.PopularAnalyse, 0)
				} else {
					grade3PopularMap[popularNumber] = append(grade3PopularMap[popularNumber], popular)
				}
			case race_vo.AllowanceClass:
				if _, ok := allowanceClassPopularMap[popularNumber]; !ok {
					allowanceClassPopularMap[popularNumber] = make([]*analyse_entity.PopularAnalyse, 0)
				} else {
					allowanceClassPopularMap[popularNumber] = append(allowanceClassPopularMap[popularNumber], popular)
				}
			}
		}
	}

	grade1Summaries = a.convertWinPopularAnalyzeSummary(grade1PopularMap)
	grade2Summaries = a.convertWinPopularAnalyzeSummary(grade2PopularMap)
	grade3Summaries = a.convertWinPopularAnalyzeSummary(grade3PopularMap)
	allowanceClassSummaries = a.convertWinPopularAnalyzeSummary(allowanceClassPopularMap)

	return grade1Summaries, grade2Summaries, grade3Summaries, allowanceClassSummaries
}

func (a *Analyser) getRecordsByWin(records []*betting_ticket_entity.CsvEntity) []*betting_ticket_entity.CsvEntity {
	var recordsByWin []*betting_ticket_entity.CsvEntity
	for _, record := range records {
		if record.BettingTicket() != betting_ticket_vo.Win {
			continue
		}
		recordsByWin = append(recordsByWin, record)
	}

	return recordsByWin
}

func (a *Analyser) getRecordsByPopular(records []*betting_ticket_entity.CsvEntity) []*betting_ticket_entity.CsvEntity {
	var recordsByPopular []*betting_ticket_entity.CsvEntity
	for _, record := range records {
		record.BettingTicket()
	}

	return recordsByPopular
}

func (a *Analyser) getPopularAnalyse(record *betting_ticket_entity.CsvEntity, race *race_entity.Race) *analyse_entity.PopularAnalyse {
	for _, raceResult := range race.RaceResults() {
		betNumber := record.BetNumber().List()[0]
		if betNumber == raceResult.HorseNumber() {
			return analyse_entity.NewPopularAnalyse(
				raceResult.PopularNumber(),
				record.Payment(),
				record.Repayment(),
				raceResult.Odds(),
				record.Winning(),
				race.Class(),
			)
		}
	}
	return nil
}
