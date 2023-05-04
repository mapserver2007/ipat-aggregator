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

func (a *Analyser) PopularAnalyse(
	records []*betting_ticket_entity.CsvEntity,
	racingNumbers []*race_entity.RacingNumber,
	races []*race_entity.Race,
) []*analyse_entity.PopularAnalyseSummary {
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
			popular := a.getPopularAnalyse(record, race.RaceResults())
			popularMap[popular.PopularNumber()] = append(popularMap[popular.PopularNumber()], popular)
		}
	}

	popularAnalyseSummaries := a.convertPopularAnalyzeSummary(popularMap)

	return popularAnalyseSummaries
}

func (a *Analyser) convertPopularAnalyzeSummary(popularMap map[int][]*analyse_entity.PopularAnalyse) []*analyse_entity.PopularAnalyseSummary {
	var popularAnalyseSummaries []*analyse_entity.PopularAnalyseSummary
	for popularNumber, populars := range popularMap {
		var (
			hitCount                                                int
			hitRate                                                 float64
			totalOddsAtVote, totalOddsAtHit                         float64
			averageOddsAtVote, averageOddsAtHit, averageOddsAtUnHit float64
			totalPayment, totalPayout                               int
			averagePayment, averagePayout                           int
			medianPayment, medianPayout                             float64
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
				medianPayment = float64(allPayments[len(allPayments)/2]+allPayments[len(allPayments)/2-1]) / 2
			} else {
				medianPayment = float64(allPayments[len(allPayments)/2])
			}
		}
		if len(allPayouts) > 0 {
			if len(allPayouts) > 0 && len(allPayouts)%2 == 0 {
				medianPayout = float64(allPayouts[len(allPayouts)/2]+allPayouts[len(allPayouts)/2-1]) / 2
			} else {
				medianPayout = float64(allPayouts[len(allPayouts)/2])
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

func (a *Analyser) getPopularAnalyse(record *betting_ticket_entity.CsvEntity, raceResults []*race_entity.RaceResult) *analyse_entity.PopularAnalyse {
	for _, raceResult := range raceResults {
		betNumber := record.BetNumber().List()[0]
		if betNumber == raceResult.HorseNumber() {
			return analyse_entity.NewPopularAnalyse(
				raceResult.PopularNumber(),
				record.Payment(),
				record.Repayment(),
				raceResult.Odds(),
				record.Winning(),
			)
		}
	}
	return nil
}
