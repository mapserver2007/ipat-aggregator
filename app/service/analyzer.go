package service

import (
	"fmt"
	analyze_entity "github.com/mapserver2007/ipat-aggregator/app/domain/analyze/entity"
	betting_ticket_entity "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/entity"
	betting_ticket_vo "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/value_object"
	race_entity "github.com/mapserver2007/ipat-aggregator/app/domain/race/entity"
	race_vo "github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
	"github.com/mapserver2007/ipat-aggregator/app/service/factory"
	"math"
	"sort"
)

type Analyzer struct {
	raceConverter RaceConverter
}

func NewAnalyzer(
	raceConverter RaceConverter,
) Analyzer {
	return Analyzer{
		raceConverter: raceConverter,
	}
}

func (a *Analyzer) WinPopularAnalyze(
	records []*betting_ticket_entity.CsvEntity,
	racingNumbers []*race_entity.RacingNumber,
	races []*race_entity.Race,
) *analyze_entity.WinAnalyzeSummary {
	entities := a.getWinAnalyzeEntities(records, racingNumbers, races)
	popularMap := map[int][]*analyze_entity.WinAnalyze{}
	for _, entity := range entities {
		popularMap[entity.PopularNumber()] = append(popularMap[entity.PopularNumber()], entity)
	}
	allSummaries := a.convertWinPopularAnalyzeSummary(popularMap)
	grade1Summaries, grade2Summaries, grade3Summaries, allowanceClassSummaries := a.convertClassWinPopularAnalyzeSummaries(popularMap)

	return analyze_entity.NewWinAnalyzeSummary(allSummaries, grade1Summaries, grade2Summaries, grade3Summaries, allowanceClassSummaries)
}

func (a *Analyzer) WinOddsAnalyzer(
	records []*betting_ticket_entity.CsvEntity,
	racingNumbers []*race_entity.RacingNumber,
	races []*race_entity.Race,
) {
	entities := a.getWinAnalyzeEntities(records, racingNumbers, races)
	oddsMap := map[string][]*analyze_entity.WinAnalyze{}
	for _, entity := range entities {
		oddsMap[entity.Odds().OddsRange()] = append(oddsMap[entity.Odds().OddsRange()], entity)
	}

	allSummaries := a.convertWinOddsAnalyzeSummary(oddsMap)

	fmt.Println(allSummaries)
}

func (a *Analyzer) convertWinPopularAnalyzeSummary(popularMap map[int][]*analyze_entity.WinAnalyze) []*analyze_entity.WinPopularAnalyzeSummary {
	winPopularAnalyzeSummaries := factory.DefaultWinPopularAnalyzeSummarySlice()
	size := len(winPopularAnalyzeSummaries)
	for idx := 0; idx < size; idx++ {
		popularNumber := idx + 1
		popularEntities, ok := popularMap[popularNumber]
		if !ok {
			winPopularAnalyzeSummaries = append(winPopularAnalyzeSummaries, analyze_entity.DefaultWinPopularAnalyzeSummary(popularNumber))
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

		betCount := len(popularEntities)

		for _, popularEntity := range popularEntities {
			totalOddsAtVote += popularEntity.Odds().Value()
			totalPayment += popularEntity.Payment()
			totalPayout += popularEntity.Payout()
			allPayments = append(allPayments, popularEntity.Payment())
			if popularEntity.IsHit() {
				hitCount++
				totalOddsAtHit += popularEntity.Odds().Value()
				allPayouts = append(allPayouts, popularEntity.Payout())
			}
			if popularEntity.Payout() > maxPayout {
				maxPayout = popularEntity.Payout()
			}
			if (popularEntity.Payout() > 0 && popularEntity.Payout() < minPayout) || minPayout == 0 {
				minPayout = popularEntity.Payout()
			}
			if (popularEntity.Odds().Value() > maxOddsAtHit) && popularEntity.IsHit() {
				maxOddsAtHit = popularEntity.Odds().Value()
			}
			if ((popularEntity.Odds() > 0 && popularEntity.Odds().Value() < minOddsAtHit) || minOddsAtHit == 0) && popularEntity.IsHit() {
				minOddsAtHit = popularEntity.Odds().Value()
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

		popularAnalyzeSummary := analyze_entity.NewWinPopularAnalyzeSummary(
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
		winPopularAnalyzeSummaries = append(winPopularAnalyzeSummaries, popularAnalyzeSummary)
	}

	sort.Slice(winPopularAnalyzeSummaries, func(i, j int) bool {
		return winPopularAnalyzeSummaries[i].PopularNumber() < winPopularAnalyzeSummaries[j].PopularNumber()
	})

	return winPopularAnalyzeSummaries
}

func (a *Analyzer) convertWinOddsAnalyzeSummary(oddsMap map[string][]*analyze_entity.WinAnalyze) map[string]*analyze_entity.WinOddsAnalyzeSummary {
	winOddsAnalyzeSummaryMap := factory.DefaultWinOddsAnalyzeSummaryMap()
	for oddRange := range winOddsAnalyzeSummaryMap {
		oddsEntities, ok := oddsMap[oddRange]
		if !ok {
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

		betCount := len(oddsEntities)
		for _, oddsEntity := range oddsEntities {
			totalOddsAtVote += oddsEntity.Odds().Value()
			totalPayment += oddsEntity.Payment()
			totalPayout += oddsEntity.Payout()
			allPayments = append(allPayments, oddsEntity.Payment())
			if oddsEntity.IsHit() {
				hitCount++
				totalOddsAtHit += oddsEntity.Odds().Value()
				allPayouts = append(allPayouts, oddsEntity.Payout())
			}
			if oddsEntity.Payout() > maxPayout {
				maxPayout = oddsEntity.Payout()
			}
			if (oddsEntity.Payout() > 0 && oddsEntity.Payout() < minPayout) || minPayout == 0 {
				minPayout = oddsEntity.Payout()
			}
			if (oddsEntity.Odds().Value() > maxOddsAtHit) && oddsEntity.IsHit() {
				maxOddsAtHit = oddsEntity.Odds().Value()
			}
			if ((oddsEntity.Odds() > 0 && oddsEntity.Odds().Value() < minOddsAtHit) || minOddsAtHit == 0) && oddsEntity.IsHit() {
				minOddsAtHit = oddsEntity.Odds().Value()
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

		oddsAnalyzeSummary := analyze_entity.NewWinOddsAnalyzeSummary(
			oddRange,
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

		winOddsAnalyzeSummaryMap[oddRange] = oddsAnalyzeSummary
	}

	return winOddsAnalyzeSummaryMap
}

func (a *Analyzer) convertClassWinPopularAnalyzeSummaries(
	popularMap map[int][]*analyze_entity.WinAnalyze,
) (
	[]*analyze_entity.WinPopularAnalyzeSummary,
	[]*analyze_entity.WinPopularAnalyzeSummary,
	[]*analyze_entity.WinPopularAnalyzeSummary,
	[]*analyze_entity.WinPopularAnalyzeSummary,
) {
	var (
		grade1Summaries, grade2Summaries, grade3Summaries, allowanceClassSummaries []*analyze_entity.WinPopularAnalyzeSummary
	)
	grade1PopularMap := map[int][]*analyze_entity.WinAnalyze{}
	grade2PopularMap := map[int][]*analyze_entity.WinAnalyze{}
	grade3PopularMap := map[int][]*analyze_entity.WinAnalyze{}
	allowanceClassPopularMap := map[int][]*analyze_entity.WinAnalyze{}

	for popularNumber, populars := range popularMap {
		for _, popular := range populars {
			switch popular.Class() {
			case race_vo.Grade1:
				if _, ok := grade1PopularMap[popularNumber]; !ok {
					grade1PopularMap[popularNumber] = make([]*analyze_entity.WinAnalyze, 0)
				} else {
					grade1PopularMap[popularNumber] = append(grade1PopularMap[popularNumber], popular)
				}
			case race_vo.Grade2:
				if _, ok := grade2PopularMap[popularNumber]; !ok {
					grade2PopularMap[popularNumber] = make([]*analyze_entity.WinAnalyze, 0)
				} else {
					grade2PopularMap[popularNumber] = append(grade2PopularMap[popularNumber], popular)
				}
			case race_vo.Grade3:
				if _, ok := grade3PopularMap[popularNumber]; !ok {
					grade3PopularMap[popularNumber] = make([]*analyze_entity.WinAnalyze, 0)
				} else {
					grade3PopularMap[popularNumber] = append(grade3PopularMap[popularNumber], popular)
				}
			case race_vo.AllowanceClass:
				if _, ok := allowanceClassPopularMap[popularNumber]; !ok {
					allowanceClassPopularMap[popularNumber] = make([]*analyze_entity.WinAnalyze, 0)
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

func (a *Analyzer) getWinAnalyzeEntities(
	records []*betting_ticket_entity.CsvEntity,
	racingNumbers []*race_entity.RacingNumber,
	races []*race_entity.Race,
) []*analyze_entity.WinAnalyze {
	raceMap := a.raceConverter.ConvertToRaceMapByRaceId(races)
	racingNumberMap := a.raceConverter.ConvertToRacingNumberMap(racingNumbers)
	recordsByWin := a.getRecordsByWin(records)

	var winAnalyzeEntities []*analyze_entity.WinAnalyze
	for _, record := range recordsByWin {
		racingNumberId := race_vo.NewRacingNumberId(record.RaceDate(), record.RaceCourse())
		racingNumber, ok := racingNumberMap[racingNumberId]
		if !ok {
			panic(fmt.Errorf("unknown racingNumberId: %s", racingNumberId))
		}
		raceId := a.raceConverter.GetRaceId(record, racingNumber)
		if race, ok := raceMap[*raceId]; ok {
			winAnalyze := a.getWinAnalyze(record, race)
			winAnalyzeEntities = append(winAnalyzeEntities, winAnalyze)
		}
	}

	return winAnalyzeEntities
}

func (a *Analyzer) getRecordsByWin(records []*betting_ticket_entity.CsvEntity) []*betting_ticket_entity.CsvEntity {
	var recordsByWin []*betting_ticket_entity.CsvEntity
	for _, record := range records {
		if record.BettingTicket() != betting_ticket_vo.Win || record.RaceCourse().Organizer() != race_vo.JRA {
			continue
		}
		recordsByWin = append(recordsByWin, record)
	}

	return recordsByWin
}

func (a *Analyzer) getRecordsByPopular(records []*betting_ticket_entity.CsvEntity) []*betting_ticket_entity.CsvEntity {
	var recordsByPopular []*betting_ticket_entity.CsvEntity
	for _, record := range records {
		record.BettingTicket()
	}

	return recordsByPopular
}

func (a *Analyzer) getWinAnalyze(record *betting_ticket_entity.CsvEntity, race *race_entity.Race) *analyze_entity.WinAnalyze {
	for _, raceResult := range race.RaceResults() {
		betNumber := record.BetNumber().List()[0]
		if betNumber == raceResult.HorseNumber() {
			return analyze_entity.NewWinAnalyze(
				raceResult.PopularNumber(),
				record.Payment().Value(),
				record.Payout().Value(),
				raceResult.Odds(),
				record.BettingResult() == betting_ticket_vo.Hit,
				race.Class(),
			)
		}
	}
	return nil
}
