package service

import (
	"fmt"
	betting_ticket_entity "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/entity"
	betting_ticket_vo "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/value_object"
	race_entity "github.com/mapserver2007/ipat-aggregator/app/domain/race/entity"
	race_vo "github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
	result_entity "github.com/mapserver2007/ipat-aggregator/app/domain/result/entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/result/types"
	"strconv"
	"time"
)

type Summarizer struct {
	raceConverter          RaceConverter
	bettingTicketConverter BettingTicketConverter
}

func NewSummarizer(
	raceConverter RaceConverter,
	bettingTicketConverter BettingTicketConverter,
) Summarizer {
	return Summarizer{
		raceConverter:          raceConverter,
		bettingTicketConverter: bettingTicketConverter,
	}
}

func (s *Summarizer) GetShortSummaryForAll(records []*betting_ticket_entity.CsvEntity) result_entity.ShortSummary {
	return result_entity.NewShortSummary(
		s.getPayment(records),
		s.getPayout(records),
	)
}

func (s *Summarizer) GetShortSummaryForMonth(records []*betting_ticket_entity.CsvEntity) result_entity.ShortSummary {
	return result_entity.NewShortSummary(
		s.getPaymentForMonth(records),
		s.getPayoutForMonth(records),
	)
}

func (s *Summarizer) GetShortSummaryForYear(records []*betting_ticket_entity.CsvEntity) result_entity.ShortSummary {
	return result_entity.NewShortSummary(
		s.getPaymentForYear(records),
		s.getPayoutForYear(records),
	)
}

func (s *Summarizer) GetBettingTicketSummaryForAll(records []*betting_ticket_entity.CsvEntity, bettingTicketTypes ...betting_ticket_vo.BettingTicket) result_entity.DetailSummary {
	return result_entity.NewDetailSummary(
		s.getBettingTicketBetCountForAll(records, bettingTicketTypes...),
		s.getBettingTicketHitCountForAll(records, bettingTicketTypes...),
		s.getBettingTicketPaymentForAll(records, bettingTicketTypes...),
		s.getBettingTicketPayoutForAll(records, bettingTicketTypes...),
		s.getBettingTicketAveragePayoutForAll(records, bettingTicketTypes...),
		s.getBettingTicketMaxPayoutForAll(records, bettingTicketTypes...),
		s.getBettingTicketMinPayoutForAll(records, bettingTicketTypes...),
	)
}

func (s *Summarizer) GetGradeClassSummaryForAll(records []*betting_ticket_entity.CsvEntity, races []*race_entity.Race, gradeClasses ...race_vo.GradeClass) result_entity.DetailSummary {
	return result_entity.NewDetailSummary(
		s.getGradeClassBetCountForAll(records, races, gradeClasses...),
		s.getGradeClassHitCountForAll(records, races, gradeClasses...),
		s.getGradeClassPaymentForAll(records, races, gradeClasses...),
		s.getGradeClassPayoutForAll(records, races, gradeClasses...),
		s.getGradeClassAveragePayoutForAll(records, races, gradeClasses...),
		s.getGradeClassMaxPayoutForAll(records, races, gradeClasses...),
		s.getGradeClassMinPayoutForAll(records, races, gradeClasses...),
	)
}

func (s *Summarizer) GetMonthlySummaryMap(records []*betting_ticket_entity.CsvEntity) map[int]result_entity.DetailSummary {
	monthlySummaryMap := map[int]result_entity.DetailSummary{}
	for date, recordsGroup := range s.bettingTicketConverter.ConvertToMonthRecordsMap(records) {
		monthlySummaryMap[date] = result_entity.NewDetailSummary(
			s.getBetCount(recordsGroup),
			s.getHitCount(recordsGroup),
			s.getPayment(recordsGroup),
			s.getPayout(recordsGroup),
			s.getAveragePayout(recordsGroup),
			s.getMaxPayout(recordsGroup),
			s.getMinPayout(recordsGroup),
		)
	}
	return monthlySummaryMap
}

func (s *Summarizer) GetCourseCategorySummaryForAll(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, courseCategory race_vo.CourseCategory) result_entity.DetailSummary {
	return result_entity.NewDetailSummary(
		s.getCourseCategoryBetCountForAll(records, racingNumbers, races, courseCategory),
		s.getCourseCategoryHitCountForAll(records, racingNumbers, races, courseCategory),
		s.getCourseCategoryPaymentForAll(records, racingNumbers, races, courseCategory),
		s.getCourseCategoryPayoutForAll(records, racingNumbers, races, courseCategory),
		s.getCourseCategoryAveragePayoutForAll(records, racingNumbers, races, courseCategory),
		s.getCourseCategoryMaxPayoutForAll(records, racingNumbers, races, courseCategory),
		s.getCourseCategoryMinPayoutForAll(records, racingNumbers, races, courseCategory),
	)
}

func (s *Summarizer) GetDistanceSummaryForAll(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, distanceCategory race_vo.DistanceCategory) result_entity.DetailSummary {
	return result_entity.NewDetailSummary(
		s.getDistanceCategoryBetCountForAll(records, racingNumbers, races, distanceCategory),
		s.getDistanceCategoryHitCountForAll(records, racingNumbers, races, distanceCategory),
		s.getDistanceCategoryPaymentForAll(records, racingNumbers, races, distanceCategory),
		s.getDistanceCategoryPayoutForAll(records, racingNumbers, races, distanceCategory),
		s.getDistanceCategoryAveragePayoutForAll(records, racingNumbers, races, distanceCategory),
		s.getDistanceCategoryMaxPayoutForAll(records, racingNumbers, races, distanceCategory),
		s.getDistanceCategoryMinPayoutForAll(records, racingNumbers, races, distanceCategory),
	)
}

// getPayment 投資額の合計を取得する
func (s *Summarizer) getPayment(records []*betting_ticket_entity.CsvEntity) types.Payment {
	payment, _ := getSumAmount(records)
	return payment
}

// getPayout 回収額の合計を取得する
func (s *Summarizer) getPayout(records []*betting_ticket_entity.CsvEntity) types.Payout {
	_, payout := getSumAmount(records)
	return payout
}

func (s *Summarizer) getBetCount(records []*betting_ticket_entity.CsvEntity) types.BetCount {
	return types.BetCount(len(records))
}

func (s *Summarizer) getHitCount(records []*betting_ticket_entity.CsvEntity) types.HitCount {
	hitCount := 0
	for _, record := range records {
		if record.BettingResult() == betting_ticket_vo.Hit {
			hitCount++
		}
	}
	return types.HitCount(hitCount)
}

func (s *Summarizer) getAveragePayout(records []*betting_ticket_entity.CsvEntity) types.Payout {
	// 不的中を除外
	var hitRecords []*betting_ticket_entity.CsvEntity
	for _, record := range records {
		if record.BettingResult() == betting_ticket_vo.Hit {
			hitRecords = append(hitRecords, record)
		}
	}

	_, payout := getSumAmount(hitRecords)
	return types.Payout(int(float64(payout) / float64(len(hitRecords))))
}

func (s *Summarizer) getMaxPayout(records []*betting_ticket_entity.CsvEntity) types.Payout {
	maxPayout := 0
	for _, record := range records {
		if maxPayout < record.Repayment() {
			maxPayout = record.Repayment()
		}
	}
	return types.Payout(maxPayout)
}

func (s *Summarizer) getMinPayout(records []*betting_ticket_entity.CsvEntity) types.Payout {
	minPayout := 0
	for _, record := range records {
		if record.Repayment() == 0 {
			continue
		}
		if minPayout == 0 || minPayout > record.Repayment() {
			minPayout = record.Repayment()
		}
	}
	return types.Payout(minPayout)
}

// getRecoveryRate 回収率の合計を取得する
func (s *Summarizer) getRecoveryRate(records []*betting_ticket_entity.CsvEntity) string {
	payment, payout := getSumAmount(records)
	if payment == 0 {
		return fmt.Sprintf("%d%s", 0, "%")
	}
	return fmt.Sprintf("%s%s", strconv.FormatFloat((float64(payout)*float64(100))/float64(payment), 'f', 2, 64), "%")
}

// getPaymentForMonth 投資額の合計を取得する(今月)
func (s *Summarizer) getPaymentForMonth(records []*betting_ticket_entity.CsvEntity) types.Payment {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	key, _ := strconv.Atoi(fmt.Sprintf("%d%02d", year, month))
	recordsGroup := s.bettingTicketConverter.ConvertToMonthRecordsMap(records)

	if recordsForMonth, ok := recordsGroup[key]; ok {
		payment, _ := getSumAmount(recordsForMonth)
		return payment
	}

	return types.Payment(0)
}

// getPayoutForMonth 回収額の合計を取得する(今月)
func (s *Summarizer) getPayoutForMonth(records []*betting_ticket_entity.CsvEntity) types.Payout {
	now := time.Now()
	year := now.Year()
	month := int(now.Month())
	key, _ := strconv.Atoi(fmt.Sprintf("%d%02d", year, month))
	recordsGroup := s.bettingTicketConverter.ConvertToMonthRecordsMap(records)

	if recordsForMonth, ok := recordsGroup[key]; ok {
		_, payout := getSumAmount(recordsForMonth)
		return payout
	}

	return types.Payout(0)
}

// getRecoveryRateForMonth 回収率の合計を取得する(今月)
func (s *Summarizer) getRecoveryRateForMonth(records []*betting_ticket_entity.CsvEntity) string {
	payment := s.getPaymentForMonth(records)
	payout := s.getPayoutForMonth(records)
	if payment == 0 {
		return fmt.Sprintf("%d%s", 0, "%")
	}
	return fmt.Sprintf("%s%s", strconv.FormatFloat((float64(payout)*float64(100))/float64(payment), 'f', 2, 64), "%")
}

// getPaymentForYear 投資額の合計を取得する(今年)
func (s *Summarizer) getPaymentForYear(records []*betting_ticket_entity.CsvEntity) types.Payment {
	now := time.Now()
	key := now.Year()
	recordsGroup := s.bettingTicketConverter.ConvertToYearRecordsMap(records)

	if recordsForYear, ok := recordsGroup[key]; ok {
		payment, _ := getSumAmount(recordsForYear)
		return payment
	}

	return types.Payment(0)
}

// getPayoutForYear 回収額の合計を取得する(今年)
func (s *Summarizer) getPayoutForYear(records []*betting_ticket_entity.CsvEntity) types.Payout {
	now := time.Now()
	key := now.Year()
	recordsGroup := s.bettingTicketConverter.ConvertToYearRecordsMap(records)

	if recordsForYear, ok := recordsGroup[key]; ok {
		_, payout := getSumAmount(recordsForYear)
		return payout
	}

	return types.Payout(0)
}

// getRecoveryRateForYear 回収率の合計を取得する(今年)
func (s *Summarizer) getRecoveryRateForYear(records []*betting_ticket_entity.CsvEntity) string {
	payment := s.getPaymentForYear(records)
	payout := s.getPayoutForYear(records)
	if payment == 0 {
		return fmt.Sprintf("%d%s", 0, "%")
	}
	return fmt.Sprintf("%s%s", strconv.FormatFloat((float64(payout)*float64(100))/float64(payment), 'f', 2, 64), "%")
}

// getBettingTicketPaymentForAll 券種別投資額の合計を取得する(全期間)
func (s *Summarizer) getBettingTicketPaymentForAll(records []*betting_ticket_entity.CsvEntity, bettingTicketTypes ...betting_ticket_vo.BettingTicket) types.Payment {
	recordsGroup := s.bettingTicketConverter.ConvertToBettingTicketRecordsMap(records)
	var mergedRecords []*betting_ticket_entity.CsvEntity
	for _, bettingTicketType := range bettingTicketTypes {
		if recordsByBettingTicket, ok := recordsGroup[bettingTicketType]; ok {
			mergedRecords = append(mergedRecords, recordsByBettingTicket...)
		}
	}
	payment, _ := getSumAmount(mergedRecords)
	return payment
}

// getBettingTicketPayoutForAll 券種別回収額の合計を取得する(全期間)
func (s *Summarizer) getBettingTicketPayoutForAll(records []*betting_ticket_entity.CsvEntity, bettingTicketTypes ...betting_ticket_vo.BettingTicket) types.Payout {
	recordsGroup := s.bettingTicketConverter.ConvertToBettingTicketRecordsMap(records)
	var mergedRecords []*betting_ticket_entity.CsvEntity
	for _, bettingTicketType := range bettingTicketTypes {
		if recordsByBettingTicket, ok := recordsGroup[bettingTicketType]; ok {
			mergedRecords = append(mergedRecords, recordsByBettingTicket...)
		}
	}
	_, payout := getSumAmount(mergedRecords)
	return payout
}

// getBettingTicketWinRecoveryRateForAll 券種別回収率の合計を取得する(全期間)
func (s *Summarizer) getBettingTicketWinRecoveryRateForAll(records []*betting_ticket_entity.CsvEntity, bettingTicketTypes ...betting_ticket_vo.BettingTicket) string {
	payment := s.getBettingTicketPaymentForAll(records, bettingTicketTypes...)
	payout := s.getBettingTicketPayoutForAll(records, bettingTicketTypes...)
	if payment == 0 {
		return fmt.Sprintf("%d%s", 0, "%")
	}
	return fmt.Sprintf("%s%s", strconv.FormatFloat((float64(payout)*float64(100))/float64(payment), 'f', 2, 64), "%")
}

// getBettingTicketBetCountForAll 券種別投票数の合計を取得する(全期間)
func (s *Summarizer) getBettingTicketBetCountForAll(records []*betting_ticket_entity.CsvEntity, bettingTicketTypes ...betting_ticket_vo.BettingTicket) types.BetCount {
	recordsGroup := s.bettingTicketConverter.ConvertToBettingTicketRecordsMap(records)
	var mergedRecords []*betting_ticket_entity.CsvEntity
	for _, bettingTicketType := range bettingTicketTypes {
		if recordsByBettingTicket, ok := recordsGroup[bettingTicketType]; ok {
			mergedRecords = append(mergedRecords, recordsByBettingTicket...)
		}
	}
	return types.BetCount(len(mergedRecords))
}

// getBettingTicketHitCountForAll 券種別的中数の合計を取得する(全期間)
func (s *Summarizer) getBettingTicketHitCountForAll(records []*betting_ticket_entity.CsvEntity, bettingTicketTypes ...betting_ticket_vo.BettingTicket) types.HitCount {
	recordsGroup := s.bettingTicketConverter.ConvertToBettingTicketRecordsMap(records)
	var mergedRecords []*betting_ticket_entity.CsvEntity
	for _, bettingTicketType := range bettingTicketTypes {
		if recordsByBettingTicket, ok := recordsGroup[bettingTicketType]; ok {
			mergedRecords = append(mergedRecords, recordsByBettingTicket...)
		}
	}
	hitCount := 0
	for _, record := range mergedRecords {
		if record.BettingResult() == betting_ticket_vo.Hit {
			hitCount++
		}
	}
	return types.HitCount(hitCount)
}

// getBettingTicketMaxPayoutForAll 券種別最大回収額の合計を取得する(全期間)
func (s *Summarizer) getBettingTicketMaxPayoutForAll(records []*betting_ticket_entity.CsvEntity, bettingTicketTypes ...betting_ticket_vo.BettingTicket) types.Payout {
	recordsGroup := s.bettingTicketConverter.ConvertToBettingTicketRecordsMap(records)
	var mergedRecords []*betting_ticket_entity.CsvEntity
	for _, bettingTicketType := range bettingTicketTypes {
		if recordsByBettingTicket, ok := recordsGroup[bettingTicketType]; ok {
			mergedRecords = append(mergedRecords, recordsByBettingTicket...)
		}
	}
	maxPayout := 0
	for _, record := range mergedRecords {
		if maxPayout < record.Repayment() {
			maxPayout = record.Repayment()
		}
	}
	return types.Payout(maxPayout)
}

// getBettingTicketMinPayoutForAll 券種別最小回収額の合計を取得する(全期間)
func (s *Summarizer) getBettingTicketMinPayoutForAll(records []*betting_ticket_entity.CsvEntity, bettingTicketTypes ...betting_ticket_vo.BettingTicket) types.Payout {
	recordsGroup := s.bettingTicketConverter.ConvertToBettingTicketRecordsMap(records)
	var mergedRecords []*betting_ticket_entity.CsvEntity
	for _, bettingTicketType := range bettingTicketTypes {
		if recordsByBettingTicket, ok := recordsGroup[bettingTicketType]; ok {
			mergedRecords = append(mergedRecords, recordsByBettingTicket...)
		}
	}
	minPayout := 0
	for _, record := range mergedRecords {
		if record.Repayment() == 0 {
			continue
		}
		if minPayout == 0 || minPayout > record.Repayment() {
			minPayout = record.Repayment()
		}
	}
	return types.Payout(minPayout)
}

// getBettingTicketAveragePayoutForAll 券種別平均回収額の合計を取得する(全期間)
func (s *Summarizer) getBettingTicketAveragePayoutForAll(records []*betting_ticket_entity.CsvEntity, bettingTicketTypes ...betting_ticket_vo.BettingTicket) types.Payout {
	recordsGroup := s.bettingTicketConverter.ConvertToBettingTicketRecordsMap(records)
	var mergedRecords []*betting_ticket_entity.CsvEntity
	for _, bettingTicketType := range bettingTicketTypes {
		if recordsByBettingTicket, ok := recordsGroup[bettingTicketType]; ok {
			mergedRecords = append(mergedRecords, recordsByBettingTicket...)
		}
	}
	// 不的中を除外
	var hitRecords []*betting_ticket_entity.CsvEntity
	for _, record := range mergedRecords {
		if record.BettingResult() == betting_ticket_vo.Hit {
			hitRecords = append(hitRecords, record)
		}
	}

	_, payout := getSumAmount(hitRecords)
	return types.Payout(int(float64(payout) / float64(len(hitRecords))))
}

// getGradeClassBetCountForAll クラス別投票数の合計を取得する(全期間)
func (s *Summarizer) getGradeClassBetCountForAll(records []*betting_ticket_entity.CsvEntity, races []*race_entity.Race, gradeClasses ...race_vo.GradeClass) types.BetCount {
	raceMap := s.raceConverter.ConvertToRaceMapByRacingNumberId(races)
	recordsGroup := s.bettingTicketConverter.ConvertToRaceClassRecordsMap(records, raceMap)
	var mergedRecords []*betting_ticket_entity.CsvEntity
	for _, gradeClass := range gradeClasses {
		if recordsByGradeClass, ok := recordsGroup[gradeClass]; ok {
			mergedRecords = append(mergedRecords, recordsByGradeClass...)
		}
	}
	return types.BetCount(len(mergedRecords))
}

// getGradeClassHitCountForAll クラス別的中数の合計を取得する(全期間)
func (s *Summarizer) getGradeClassHitCountForAll(records []*betting_ticket_entity.CsvEntity, races []*race_entity.Race, gradeClasses ...race_vo.GradeClass) types.HitCount {
	raceMap := s.raceConverter.ConvertToRaceMapByRacingNumberId(races)
	recordsGroup := s.bettingTicketConverter.ConvertToRaceClassRecordsMap(records, raceMap)
	var mergedRecords []*betting_ticket_entity.CsvEntity
	for _, gradeClass := range gradeClasses {
		if recordsByGradeClass, ok := recordsGroup[gradeClass]; ok {
			mergedRecords = append(mergedRecords, recordsByGradeClass...)
		}
	}
	hitCount := 0
	for _, record := range mergedRecords {
		if record.BettingResult() == betting_ticket_vo.Hit {
			hitCount++
		}
	}
	return types.HitCount(hitCount)
}

// getGradeClassPaymentForAll クラス別投票金額の合計を取得する(全期間)
func (s *Summarizer) getGradeClassPaymentForAll(records []*betting_ticket_entity.CsvEntity, races []*race_entity.Race, gradeClasses ...race_vo.GradeClass) types.Payment {
	raceMap := s.raceConverter.ConvertToRaceMapByRacingNumberId(races)
	recordsGroup := s.bettingTicketConverter.ConvertToRaceClassRecordsMap(records, raceMap)
	var mergedRecords []*betting_ticket_entity.CsvEntity
	for _, gradeClass := range gradeClasses {
		if recordsByBettingTicket, ok := recordsGroup[gradeClass]; ok {
			mergedRecords = append(mergedRecords, recordsByBettingTicket...)
		}
	}
	payment, _ := getSumAmount(mergedRecords)
	return payment
}

// getGradeClassPayoutForAll クラス別回収金額の合計を取得する(全期間)
func (s *Summarizer) getGradeClassPayoutForAll(records []*betting_ticket_entity.CsvEntity, races []*race_entity.Race, gradeClasses ...race_vo.GradeClass) types.Payout {
	raceMap := s.raceConverter.ConvertToRaceMapByRacingNumberId(races)
	recordsGroup := s.bettingTicketConverter.ConvertToRaceClassRecordsMap(records, raceMap)
	var mergedRecords []*betting_ticket_entity.CsvEntity
	for _, gradeClass := range gradeClasses {
		if recordsByBettingTicket, ok := recordsGroup[gradeClass]; ok {
			mergedRecords = append(mergedRecords, recordsByBettingTicket...)
		}
	}
	_, payout := getSumAmount(mergedRecords)
	return payout
}

// getGradeClassAveragePayoutForAll クラス別平均回収額の合計を取得する(全期間)
func (s *Summarizer) getGradeClassAveragePayoutForAll(records []*betting_ticket_entity.CsvEntity, races []*race_entity.Race, gradeClasses ...race_vo.GradeClass) types.Payout {
	raceMap := s.raceConverter.ConvertToRaceMapByRacingNumberId(races)
	recordsGroup := s.bettingTicketConverter.ConvertToRaceClassRecordsMap(records, raceMap)
	var mergedRecords []*betting_ticket_entity.CsvEntity
	for _, gradeClass := range gradeClasses {
		if recordsByBettingTicket, ok := recordsGroup[gradeClass]; ok {
			mergedRecords = append(mergedRecords, recordsByBettingTicket...)
		}
	}
	// 不的中を除外
	var hitRecords []*betting_ticket_entity.CsvEntity
	for _, record := range mergedRecords {
		if record.BettingResult() == betting_ticket_vo.Hit {
			hitRecords = append(hitRecords, record)
		}
	}

	_, payout := getSumAmount(hitRecords)
	return types.Payout(int(float64(payout) / float64(len(hitRecords))))
}

// getGradeClassMaxPayoutForAll クラス別最大回収額の合計を取得する(全期間)
func (s *Summarizer) getGradeClassMaxPayoutForAll(records []*betting_ticket_entity.CsvEntity, races []*race_entity.Race, gradeClasses ...race_vo.GradeClass) types.Payout {
	raceMap := s.raceConverter.ConvertToRaceMapByRacingNumberId(races)
	recordsGroup := s.bettingTicketConverter.ConvertToRaceClassRecordsMap(records, raceMap)
	var mergedRecords []*betting_ticket_entity.CsvEntity
	for _, gradeClass := range gradeClasses {
		if recordsByBettingTicket, ok := recordsGroup[gradeClass]; ok {
			mergedRecords = append(mergedRecords, recordsByBettingTicket...)
		}
	}
	maxPayout := 0
	for _, record := range mergedRecords {
		if maxPayout < record.Repayment() {
			maxPayout = record.Repayment()
		}
	}
	return types.Payout(maxPayout)
}

// getGradeClassMinPayoutForAll クラス別最小回収額の合計を取得する(全期間)
func (s *Summarizer) getGradeClassMinPayoutForAll(records []*betting_ticket_entity.CsvEntity, races []*race_entity.Race, gradeClasses ...race_vo.GradeClass) types.Payout {
	raceMap := s.raceConverter.ConvertToRaceMapByRacingNumberId(races)
	recordsGroup := s.bettingTicketConverter.ConvertToRaceClassRecordsMap(records, raceMap)
	var mergedRecords []*betting_ticket_entity.CsvEntity
	for _, gradeClass := range gradeClasses {
		if recordsByBettingTicket, ok := recordsGroup[gradeClass]; ok {
			mergedRecords = append(mergedRecords, recordsByBettingTicket...)
		}
	}
	minPayout := 0
	for _, record := range mergedRecords {
		if record.Repayment() == 0 {
			continue
		}
		if minPayout == 0 || minPayout > record.Repayment() {
			minPayout = record.Repayment()
		}
	}
	return types.Payout(minPayout)
}

func (s *Summarizer) getCourseCategoryRecordsMap(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race) map[race_vo.CourseCategory][]*betting_ticket_entity.CsvEntity {
	raceMap := s.raceConverter.ConvertToRaceMapByRaceId(races)
	racingNumberMap := s.raceConverter.ConvertToRacingNumberMap(racingNumbers)
	courseCategoryRecordsMap := map[race_vo.CourseCategory][]*betting_ticket_entity.CsvEntity{}
	for _, record := range records {
		racingNumberId := race_vo.NewRacingNumberId(record.RaceDate(), record.RaceCourse())
		racingNumber, ok := racingNumberMap[racingNumberId]
		if !ok && record.RaceCourse().Organizer() == race_vo.JRA {
			panic(fmt.Errorf("unknown racingNumberId: %s", racingNumberId))
		}
		raceId := s.raceConverter.GetRaceId(record, racingNumber)
		if race, ok := raceMap[*raceId]; ok {
			courseCategory := race.CourseCategory()
			courseCategoryRecordsMap[courseCategory] = append(courseCategoryRecordsMap[courseCategory], record)
		}
	}

	return courseCategoryRecordsMap
}

// getCourseCategoryBetCountForAll コース別投票数の合計を取得する(全期間)
func (s *Summarizer) getCourseCategoryBetCountForAll(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, courseCategory race_vo.CourseCategory) types.BetCount {
	courseCategoryRecordsMap := s.getCourseCategoryRecordsMap(records, racingNumbers, races)
	if recordsByCourseCategory, ok := courseCategoryRecordsMap[courseCategory]; ok {
		return types.BetCount(len(recordsByCourseCategory))
	}
	return types.BetCount(0)
}

// getCourseCategoryHitCountForAll コース別的中数の合計を取得する(全期間)
func (s *Summarizer) getCourseCategoryHitCountForAll(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, courseCategory race_vo.CourseCategory) types.HitCount {
	courseCategoryRecordsMap := s.getCourseCategoryRecordsMap(records, racingNumbers, races)
	hitCount := 0
	if recordsByCourseCategory, ok := courseCategoryRecordsMap[courseCategory]; ok {
		for _, record := range recordsByCourseCategory {
			if record.BettingResult() == betting_ticket_vo.Hit {
				hitCount++
			}
		}
	}
	return types.HitCount(hitCount)
}

// getCourseCategoryPaymentForAll コース別払戻金の合計を取得する(全期間)
func (s *Summarizer) getCourseCategoryPaymentForAll(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, courseCategory race_vo.CourseCategory) types.Payment {
	courseCategoryRecordsMap := s.getCourseCategoryRecordsMap(records, racingNumbers, races)
	if recordsByCourseCategory, ok := courseCategoryRecordsMap[courseCategory]; ok {
		payment, _ := getSumAmount(recordsByCourseCategory)
		return payment
	}
	return types.Payment(0)
}

// getCourseCategoryPayoutForAll コース別払戻金の合計を取得する(全期間)
func (s *Summarizer) getCourseCategoryPayoutForAll(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, courseCategory race_vo.CourseCategory) types.Payout {
	courseCategoryRecordsMap := s.getCourseCategoryRecordsMap(records, racingNumbers, races)
	if recordsByCourseCategory, ok := courseCategoryRecordsMap[courseCategory]; ok {
		_, payout := getSumAmount(recordsByCourseCategory)
		return payout
	}
	return types.Payout(0)
}

// getCourseCategoryAveragePayoutForAll コース別平均払戻金を取得する(全期間)
func (s *Summarizer) getCourseCategoryAveragePayoutForAll(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, courseCategory race_vo.CourseCategory) types.Payout {
	courseCategoryRecordsMap := s.getCourseCategoryRecordsMap(records, racingNumbers, races)
	var hitRecords []*betting_ticket_entity.CsvEntity
	if recordsByCourseCategory, ok := courseCategoryRecordsMap[courseCategory]; ok {
		for _, record := range recordsByCourseCategory {
			if record.BettingResult() == betting_ticket_vo.Hit {
				hitRecords = append(hitRecords, record)
			}
		}
	}
	_, payout := getSumAmount(hitRecords)
	return types.Payout(int(float64(payout) / float64(len(hitRecords))))
}

// getCourseCategoryMaxPayoutForAll コース別最大払戻金を取得する(全期間)
func (s *Summarizer) getCourseCategoryMaxPayoutForAll(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, courseCategory race_vo.CourseCategory) types.Payout {
	courseCategoryRecordsMap := s.getCourseCategoryRecordsMap(records, racingNumbers, races)
	maxPayout := 0
	if recordsByCourseCategory, ok := courseCategoryRecordsMap[courseCategory]; ok {
		for _, record := range recordsByCourseCategory {
			if maxPayout < record.Repayment() {
				maxPayout = record.Repayment()
			}
		}
	}
	return types.Payout(maxPayout)
}

// getCourseCategoryMinPayoutForAll コース別最小払戻金を取得する(全期間)
func (s *Summarizer) getCourseCategoryMinPayoutForAll(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, courseCategory race_vo.CourseCategory) types.Payout {
	courseCategoryRecordsMap := s.getCourseCategoryRecordsMap(records, racingNumbers, races)
	minPayout := 0
	if recordsByCourseCategory, ok := courseCategoryRecordsMap[courseCategory]; ok {
		for _, record := range recordsByCourseCategory {
			if record.Repayment() == 0 {
				continue
			}
			if minPayout == 0 || minPayout > record.Repayment() {
				minPayout = record.Repayment()
			}
		}
	}
	return types.Payout(minPayout)
}

func (s *Summarizer) getDistanceCategoryRecordsMap(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race) map[race_vo.DistanceCategory][]*betting_ticket_entity.CsvEntity {
	raceMap := s.raceConverter.ConvertToRaceMapByRaceId(races)
	racingNumberMap := s.raceConverter.ConvertToRacingNumberMap(racingNumbers)
	distanceCategoryRecordsMap := map[race_vo.DistanceCategory][]*betting_ticket_entity.CsvEntity{}
	for _, record := range records {
		racingNumberId := race_vo.NewRacingNumberId(record.RaceDate(), record.RaceCourse())
		racingNumber, ok := racingNumberMap[racingNumberId]
		if !ok && record.RaceCourse().Organizer() == race_vo.JRA {
			panic(fmt.Errorf("unknown racingNumberId: %s", racingNumberId))
		}
		raceId := s.raceConverter.GetRaceId(record, racingNumber)
		if race, ok := raceMap[*raceId]; ok {
			courseCategory := race.CourseCategory()
			distanceCategory := race_vo.NewDistanceCategory(race.Distance(), courseCategory)
			distanceCategoryRecordsMap[distanceCategory] = append(distanceCategoryRecordsMap[distanceCategory], record)
		}
	}

	return distanceCategoryRecordsMap
}

// getDistanceCategoryBetCountForAll 距離別投票数を取得する(全期間)
func (s *Summarizer) getDistanceCategoryBetCountForAll(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, distanceCategory race_vo.DistanceCategory) types.BetCount {
	distanceCategoryRecordsMap := s.getDistanceCategoryRecordsMap(records, racingNumbers, races)
	if recordsByDistanceCategory, ok := distanceCategoryRecordsMap[distanceCategory]; ok {
		return types.BetCount(len(recordsByDistanceCategory))
	}
	return types.BetCount(0)
}

// getDistanceCategoryHitCountForAll 距離別的中数を取得する(全期間)
func (s *Summarizer) getDistanceCategoryHitCountForAll(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, distanceCategory race_vo.DistanceCategory) types.HitCount {
	distanceCategoryRecordsMap := s.getDistanceCategoryRecordsMap(records, racingNumbers, races)
	hitCount := 0
	if recordsByDistanceCategory, ok := distanceCategoryRecordsMap[distanceCategory]; ok {
		for _, record := range recordsByDistanceCategory {
			if record.BettingResult() == betting_ticket_vo.Hit {
				hitCount++
			}
		}
	}
	return types.HitCount(hitCount)
}

// getDistanceCategoryPaymentForAll 距離別投票金額を取得する(全期間)
func (s *Summarizer) getDistanceCategoryPaymentForAll(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, distanceCategory race_vo.DistanceCategory) types.Payment {
	distanceCategoryRecordsMap := s.getDistanceCategoryRecordsMap(records, racingNumbers, races)
	if recordsByDistanceCategory, ok := distanceCategoryRecordsMap[distanceCategory]; ok {
		payment, _ := getSumAmount(recordsByDistanceCategory)
		return payment
	}
	return types.Payment(0)
}

// getDistanceCategoryPayoutForAll 距離別払戻金を取得する(全期間)
func (s *Summarizer) getDistanceCategoryPayoutForAll(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, distanceCategory race_vo.DistanceCategory) types.Payout {
	distanceCategoryRecordsMap := s.getDistanceCategoryRecordsMap(records, racingNumbers, races)
	if recordsByDistanceCategory, ok := distanceCategoryRecordsMap[distanceCategory]; ok {
		_, payout := getSumAmount(recordsByDistanceCategory)
		return payout
	}
	return types.Payout(0)
}

// getDistanceCategoryAveragePayoutForAll 距離別平均払戻金を取得する(全期間)
func (s *Summarizer) getDistanceCategoryAveragePayoutForAll(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, distanceCategory race_vo.DistanceCategory) types.Payout {
	distanceCategoryRecordsMap := s.getDistanceCategoryRecordsMap(records, racingNumbers, races)
	var hitRecords []*betting_ticket_entity.CsvEntity
	if recordsByDistanceCategory, ok := distanceCategoryRecordsMap[distanceCategory]; ok {
		for _, record := range recordsByDistanceCategory {
			if record.BettingResult() == betting_ticket_vo.Hit {
				hitRecords = append(hitRecords, record)
			}
		}
	}
	_, payout := getSumAmount(hitRecords)
	return types.Payout(int(float64(payout) / float64(len(hitRecords))))
}

// getDistanceCategoryMaxPayoutForAll 距離別最高払戻金を取得する(全期間)
func (s *Summarizer) getDistanceCategoryMaxPayoutForAll(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, distanceCategory race_vo.DistanceCategory) types.Payout {
	distanceCategoryRecordsMap := s.getDistanceCategoryRecordsMap(records, racingNumbers, races)
	maxPayout := 0
	if recordsByDistanceCategory, ok := distanceCategoryRecordsMap[distanceCategory]; ok {
		for _, record := range recordsByDistanceCategory {
			if maxPayout < record.Repayment() {
				maxPayout = record.Repayment()
			}
		}
	}
	return types.Payout(maxPayout)
}

// getDistanceCategoryMinPayoutForAll 距離別最低払戻金を取得する(全期間)
func (s *Summarizer) getDistanceCategoryMinPayoutForAll(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, distanceCategory race_vo.DistanceCategory) types.Payout {
	distanceCategoryRecordsMap := s.getDistanceCategoryRecordsMap(records, racingNumbers, races)
	minPayout := 0
	if recordsByDistanceCategory, ok := distanceCategoryRecordsMap[distanceCategory]; ok {
		for _, record := range recordsByDistanceCategory {
			if record.Repayment() == 0 {
				continue
			}
			if minPayout == 0 || minPayout > record.Repayment() {
				minPayout = record.Repayment()
			}
		}
	}
	return types.Payout(minPayout)
}

func getSumAmount(records []*betting_ticket_entity.CsvEntity) (types.Payment, types.Payout) {
	var (
		sumPayment int
		sumPayout  int
	)
	for _, record := range records {
		sumPayment += record.Payment()
		sumPayout += record.Repayment()
	}

	return types.Payment(sumPayment), types.Payout(sumPayout)
}
