package service

import (
	"fmt"
	betting_ticket_entity "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/entity"
	betting_ticket_vo "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/value_object"
	race_entity "github.com/mapserver2007/ipat-aggregator/app/domain/race/entity"
	race_vo "github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
	result_entity "github.com/mapserver2007/ipat-aggregator/app/domain/result/entity"
	spreadsheet_entity "github.com/mapserver2007/ipat-aggregator/app/domain/spreadsheet/entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
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

func (s *Summarizer) GetShortSummary(records []*betting_ticket_entity.CsvEntity) result_entity.ShortSummary {
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

func (s *Summarizer) GetBettingTicketSummary(records []*betting_ticket_entity.CsvEntity, bettingTicketTypes ...betting_ticket_vo.BettingTicket) result_entity.DetailSummary {
	return result_entity.NewDetailSummary(
		s.getBettingTicketBetCount(records, bettingTicketTypes...),
		s.getBettingTicketHitCount(records, bettingTicketTypes...),
		s.getBettingTicketPayment(records, bettingTicketTypes...),
		s.getBettingTicketPayout(records, bettingTicketTypes...),
		s.getBettingTicketAveragePayout(records, bettingTicketTypes...),
		s.getBettingTicketMaxPayout(records, bettingTicketTypes...),
		s.getBettingTicketMinPayout(records, bettingTicketTypes...),
	)
}

func (s *Summarizer) GetGradeClassSummary(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, gradeClasses ...race_vo.GradeClass) result_entity.DetailSummary {
	return result_entity.NewDetailSummary(
		s.getGradeClassBetCount(records, racingNumbers, races, gradeClasses...),
		s.getGradeClassHitCount(records, racingNumbers, races, gradeClasses...),
		s.getGradeClassPayment(records, racingNumbers, races, gradeClasses...),
		s.getGradeClassPayout(records, racingNumbers, races, gradeClasses...),
		s.getGradeClassAveragePayout(records, racingNumbers, races, gradeClasses...),
		s.getGradeClassMaxPayout(records, racingNumbers, races, gradeClasses...),
		s.getGradeClassMinPayout(records, racingNumbers, races, gradeClasses...),
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

func (s *Summarizer) GetCourseCategorySummaryMap(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race) map[race_vo.CourseCategory]result_entity.DetailSummary {
	courseCategories := []race_vo.CourseCategory{
		race_vo.Turf, race_vo.Dirt, race_vo.Jump,
	}
	courseCategoryMap := map[race_vo.CourseCategory]result_entity.DetailSummary{}
	for _, courseCategory := range courseCategories {
		courseCategoryMap[courseCategory] = result_entity.NewDetailSummary(
			s.getCourseCategoryBetCount(records, racingNumbers, races, courseCategory),
			s.getCourseCategoryHitCount(records, racingNumbers, races, courseCategory),
			s.getCourseCategoryPayment(records, racingNumbers, races, courseCategory),
			s.getCourseCategoryPayout(records, racingNumbers, races, courseCategory),
			s.getCourseCategoryAveragePayout(records, racingNumbers, races, courseCategory),
			s.getCourseCategoryMaxPayout(records, racingNumbers, races, courseCategory),
			s.getCourseCategoryMinPayout(records, racingNumbers, races, courseCategory),
		)
	}
	return courseCategoryMap
}

func (s *Summarizer) GetDistanceCategorySummaryMap(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race) map[race_vo.DistanceCategory]result_entity.DetailSummary {
	distanceCategories := []race_vo.DistanceCategory{
		race_vo.TurfSprint,
		race_vo.TurfMile,
		race_vo.TurfIntermediate,
		race_vo.TurfLong,
		race_vo.TurfExtended,
		race_vo.DirtSprint,
		race_vo.DirtMile,
		race_vo.DirtIntermediate,
		race_vo.DirtLong,
		race_vo.DirtExtended,
		race_vo.JumpAllDistance,
	}
	distanceCategoryMap := map[race_vo.DistanceCategory]result_entity.DetailSummary{}
	for _, distanceCategory := range distanceCategories {
		distanceCategoryMap[distanceCategory] = result_entity.NewDetailSummary(
			s.getDistanceCategoryBetCount(records, racingNumbers, races, distanceCategory),
			s.getDistanceCategoryHitCount(records, racingNumbers, races, distanceCategory),
			s.getDistanceCategoryPayment(records, racingNumbers, races, distanceCategory),
			s.getDistanceCategoryPayout(records, racingNumbers, races, distanceCategory),
			s.getDistanceCategoryAveragePayout(records, racingNumbers, races, distanceCategory),
			s.getDistanceCategoryMaxPayout(records, racingNumbers, races, distanceCategory),
			s.getDistanceCategoryMinPayout(records, racingNumbers, races, distanceCategory),
		)
	}
	return distanceCategoryMap
}

func (s *Summarizer) GetRaceCourseSummaryMap(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race) map[race_vo.RaceCourse]result_entity.DetailSummary {
	raceCourses := []race_vo.RaceCourse{
		race_vo.Sapporo, race_vo.Hakodate, race_vo.Fukushima, race_vo.Niigata, race_vo.Tokyo, race_vo.Nakayama, race_vo.Chukyo, race_vo.Kyoto, race_vo.Hanshin, race_vo.Kokura,
		race_vo.Monbetsu, race_vo.Morioka, race_vo.Urawa, race_vo.Hunabashi, race_vo.Ooi, race_vo.Kawasaki, race_vo.Nagoya, race_vo.Sonoda, race_vo.Kouchi, race_vo.Saga,
	}
	raceCourseMap := map[race_vo.RaceCourse]result_entity.DetailSummary{}
	for _, raceCourse := range raceCourses {
		raceCourseMap[raceCourse] = result_entity.NewDetailSummary(
			s.getRaceCourseBetCount(records, racingNumbers, races, raceCourse),
			s.getRaceCourseHitCount(records, racingNumbers, races, raceCourse),
			s.getRaceCoursePayment(records, racingNumbers, races, raceCourse),
			s.getRaceCoursePayout(records, racingNumbers, races, raceCourse),
			s.getRaceCourseAveragePayout(records, racingNumbers, races, raceCourse),
			s.getRaceCourseMaxPayout(records, racingNumbers, races, raceCourse),
			s.getRaceCourseMinPayout(records, racingNumbers, races, raceCourse),
		)
	}
	raceCourseMap[race_vo.Overseas] = result_entity.NewDetailSummary(
		s.getRaceCourseBetCount(records, racingNumbers, races, race_vo.Longchamp, race_vo.Deauville, race_vo.Shatin, race_vo.Meydan),
		s.getRaceCourseHitCount(records, racingNumbers, races, race_vo.Longchamp, race_vo.Deauville, race_vo.Shatin, race_vo.Meydan),
		s.getRaceCoursePayment(records, racingNumbers, races, race_vo.Longchamp, race_vo.Deauville, race_vo.Shatin, race_vo.Meydan),
		s.getRaceCoursePayout(records, racingNumbers, races, race_vo.Longchamp, race_vo.Deauville, race_vo.Shatin, race_vo.Meydan),
		s.getRaceCourseAveragePayout(records, racingNumbers, races, race_vo.Longchamp, race_vo.Deauville, race_vo.Shatin, race_vo.Meydan),
		s.getRaceCourseMaxPayout(records, racingNumbers, races, race_vo.Longchamp, race_vo.Deauville, race_vo.Shatin, race_vo.Meydan),
		s.getRaceCourseMinPayout(records, racingNumbers, races, race_vo.Longchamp, race_vo.Deauville, race_vo.Shatin, race_vo.Meydan),
	)

	return raceCourseMap
}

func (s *Summarizer) GetMonthlyBettingTicketSummary(records []*betting_ticket_entity.CsvEntity) map[int]*spreadsheet_entity.SpreadSheetBettingTicketSummary {
	bettingTicketMonthlyMap := map[int]*spreadsheet_entity.SpreadSheetBettingTicketSummary{}
	for date, recordsGroup := range s.bettingTicketConverter.ConvertToMonthRecordsMap(records) {
		spreadSheetBettingTicketSummary := spreadsheet_entity.NewSpreadSheetBettingTicketSummary(
			s.GetBettingTicketSummary(recordsGroup, betting_ticket_vo.Win),
			s.GetBettingTicketSummary(recordsGroup, betting_ticket_vo.Place),
			s.GetBettingTicketSummary(recordsGroup, betting_ticket_vo.Quinella),
			s.GetBettingTicketSummary(recordsGroup, betting_ticket_vo.Exacta, betting_ticket_vo.ExactaWheelOfFirst),
			s.GetBettingTicketSummary(recordsGroup, betting_ticket_vo.QuinellaPlace, betting_ticket_vo.QuinellaPlaceWheel),
			s.GetBettingTicketSummary(recordsGroup, betting_ticket_vo.Trio, betting_ticket_vo.TrioFormation, betting_ticket_vo.TrioWheelOfFirst),
			s.GetBettingTicketSummary(recordsGroup, betting_ticket_vo.Trifecta, betting_ticket_vo.TrifectaFormation, betting_ticket_vo.TrifectaWheelOfFirst),
			s.GetBettingTicketSummary(recordsGroup, betting_ticket_vo.Win, betting_ticket_vo.Place, betting_ticket_vo.Quinella,
				betting_ticket_vo.Exacta, betting_ticket_vo.ExactaWheelOfFirst, betting_ticket_vo.QuinellaPlace, betting_ticket_vo.QuinellaPlaceWheel,
				betting_ticket_vo.Trio, betting_ticket_vo.TrioFormation, betting_ticket_vo.TrioWheelOfFirst,
				betting_ticket_vo.Trifecta, betting_ticket_vo.TrifectaFormation, betting_ticket_vo.TrifectaWheelOfFirst),
		)
		bettingTicketMonthlyMap[date] = spreadSheetBettingTicketSummary
	}

	return bettingTicketMonthlyMap
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
	if payout > 0 {
		return types.Payout(int(float64(payout) / float64(len(hitRecords))))
	}

	return types.Payout(0)
}

func (s *Summarizer) getMaxPayout(records []*betting_ticket_entity.CsvEntity) types.Payout {
	maxPayout := 0
	for _, record := range records {
		if maxPayout < record.Payout().Value() {
			maxPayout = record.Payout().Value()
		}
	}
	return types.Payout(maxPayout)
}

func (s *Summarizer) getMinPayout(records []*betting_ticket_entity.CsvEntity) types.Payout {
	minPayout := 0
	for _, record := range records {
		if record.Payout() == 0 {
			continue
		}
		if minPayout == 0 || minPayout > record.Payout().Value() {
			minPayout = record.Payout().Value()
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

// getBettingTicketPayment 券種別投資額の合計を取得する(全期間)
func (s *Summarizer) getBettingTicketPayment(records []*betting_ticket_entity.CsvEntity, bettingTicketTypes ...betting_ticket_vo.BettingTicket) types.Payment {
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

// getBettingTicketPayout 券種別回収額の合計を取得する(全期間)
func (s *Summarizer) getBettingTicketPayout(records []*betting_ticket_entity.CsvEntity, bettingTicketTypes ...betting_ticket_vo.BettingTicket) types.Payout {
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

// getBettingTicketWinRecoveryRate 券種別回収率の合計を取得する(全期間)
func (s *Summarizer) getBettingTicketWinRecoveryRate(records []*betting_ticket_entity.CsvEntity, bettingTicketTypes ...betting_ticket_vo.BettingTicket) string {
	payment := s.getBettingTicketPayment(records, bettingTicketTypes...)
	payout := s.getBettingTicketPayout(records, bettingTicketTypes...)
	if payment == 0 {
		return fmt.Sprintf("%d%s", 0, "%")
	}
	return fmt.Sprintf("%s%s", strconv.FormatFloat((float64(payout)*float64(100))/float64(payment), 'f', 2, 64), "%")
}

// getBettingTicketBetCount 券種別投票数の合計を取得する(全期間)
func (s *Summarizer) getBettingTicketBetCount(records []*betting_ticket_entity.CsvEntity, bettingTicketTypes ...betting_ticket_vo.BettingTicket) types.BetCount {
	recordsGroup := s.bettingTicketConverter.ConvertToBettingTicketRecordsMap(records)
	var mergedRecords []*betting_ticket_entity.CsvEntity
	for _, bettingTicketType := range bettingTicketTypes {
		if recordsByBettingTicket, ok := recordsGroup[bettingTicketType]; ok {
			mergedRecords = append(mergedRecords, recordsByBettingTicket...)
		}
	}
	return types.BetCount(len(mergedRecords))
}

// getBettingTicketHitCount 券種別的中数の合計を取得する(全期間)
func (s *Summarizer) getBettingTicketHitCount(records []*betting_ticket_entity.CsvEntity, bettingTicketTypes ...betting_ticket_vo.BettingTicket) types.HitCount {
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

// getBettingTicketMaxPayout 券種別最大回収額の合計を取得する(全期間)
func (s *Summarizer) getBettingTicketMaxPayout(records []*betting_ticket_entity.CsvEntity, bettingTicketTypes ...betting_ticket_vo.BettingTicket) types.Payout {
	recordsGroup := s.bettingTicketConverter.ConvertToBettingTicketRecordsMap(records)
	var mergedRecords []*betting_ticket_entity.CsvEntity
	for _, bettingTicketType := range bettingTicketTypes {
		if recordsByBettingTicket, ok := recordsGroup[bettingTicketType]; ok {
			mergedRecords = append(mergedRecords, recordsByBettingTicket...)
		}
	}
	maxPayout := 0
	for _, record := range mergedRecords {
		if maxPayout < record.Payout().Value() {
			maxPayout = record.Payout().Value()
		}
	}
	return types.Payout(maxPayout)
}

// getBettingTicketMinPayout 券種別最小回収額の合計を取得する(全期間)
func (s *Summarizer) getBettingTicketMinPayout(records []*betting_ticket_entity.CsvEntity, bettingTicketTypes ...betting_ticket_vo.BettingTicket) types.Payout {
	recordsGroup := s.bettingTicketConverter.ConvertToBettingTicketRecordsMap(records)
	var mergedRecords []*betting_ticket_entity.CsvEntity
	for _, bettingTicketType := range bettingTicketTypes {
		if recordsByBettingTicket, ok := recordsGroup[bettingTicketType]; ok {
			mergedRecords = append(mergedRecords, recordsByBettingTicket...)
		}
	}
	minPayout := 0
	for _, record := range mergedRecords {
		if record.Payout() == 0 {
			continue
		}
		if minPayout == 0 || minPayout > record.Payout().Value() {
			minPayout = record.Payout().Value()
		}
	}
	return types.Payout(minPayout)
}

// getBettingTicketAveragePayout 券種別平均回収額の合計を取得する(全期間)
func (s *Summarizer) getBettingTicketAveragePayout(records []*betting_ticket_entity.CsvEntity, bettingTicketTypes ...betting_ticket_vo.BettingTicket) types.Payout {
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
	if payout > 0 {
		return types.Payout(int(float64(payout) / float64(len(hitRecords))))
	}

	return types.Payout(0)
}

func (s *Summarizer) getBettingTicketMinOdds() {

}

// getGradeClassBetCount クラス別投票数の合計を取得する(全期間)
func (s *Summarizer) getGradeClassBetCount(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, gradeClasses ...race_vo.GradeClass) types.BetCount {
	recordsGroup := s.bettingTicketConverter.ConvertToRaceClassRecordsMap(records, racingNumbers, races)
	var mergedRecords []*betting_ticket_entity.CsvEntity
	for _, gradeClass := range gradeClasses {
		if recordsByGradeClass, ok := recordsGroup[gradeClass]; ok {
			mergedRecords = append(mergedRecords, recordsByGradeClass...)
		}
	}
	return types.BetCount(len(mergedRecords))
}

// getGradeClassHitCount クラス別的中数の合計を取得する(全期間)
func (s *Summarizer) getGradeClassHitCount(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, gradeClasses ...race_vo.GradeClass) types.HitCount {
	recordsGroup := s.bettingTicketConverter.ConvertToRaceClassRecordsMap(records, racingNumbers, races)
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

// getGradeClassPayment クラス別投票金額の合計を取得する(全期間)
func (s *Summarizer) getGradeClassPayment(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, gradeClasses ...race_vo.GradeClass) types.Payment {
	recordsGroup := s.bettingTicketConverter.ConvertToRaceClassRecordsMap(records, racingNumbers, races)
	var mergedRecords []*betting_ticket_entity.CsvEntity
	for _, gradeClass := range gradeClasses {
		if recordsByBettingTicket, ok := recordsGroup[gradeClass]; ok {
			mergedRecords = append(mergedRecords, recordsByBettingTicket...)
		}
	}
	payment, _ := getSumAmount(mergedRecords)
	return payment
}

// getGradeClassPayout クラス別回収金額の合計を取得する(全期間)
func (s *Summarizer) getGradeClassPayout(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, gradeClasses ...race_vo.GradeClass) types.Payout {
	recordsGroup := s.bettingTicketConverter.ConvertToRaceClassRecordsMap(records, racingNumbers, races)
	var mergedRecords []*betting_ticket_entity.CsvEntity
	for _, gradeClass := range gradeClasses {
		if recordsByBettingTicket, ok := recordsGroup[gradeClass]; ok {
			mergedRecords = append(mergedRecords, recordsByBettingTicket...)
		}
	}
	_, payout := getSumAmount(mergedRecords)
	return payout
}

// getGradeClassAveragePayout クラス別平均回収額の合計を取得する(全期間)
func (s *Summarizer) getGradeClassAveragePayout(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, gradeClasses ...race_vo.GradeClass) types.Payout {
	recordsGroup := s.bettingTicketConverter.ConvertToRaceClassRecordsMap(records, racingNumbers, races)
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
	if payout > 0 {
		return types.Payout(int(float64(payout) / float64(len(hitRecords))))
	}

	return types.Payout(0)
}

// getGradeClassMaxPayout クラス別最大回収額の合計を取得する(全期間)
func (s *Summarizer) getGradeClassMaxPayout(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, gradeClasses ...race_vo.GradeClass) types.Payout {
	recordsGroup := s.bettingTicketConverter.ConvertToRaceClassRecordsMap(records, racingNumbers, races)
	var mergedRecords []*betting_ticket_entity.CsvEntity
	for _, gradeClass := range gradeClasses {
		if recordsByBettingTicket, ok := recordsGroup[gradeClass]; ok {
			mergedRecords = append(mergedRecords, recordsByBettingTicket...)
		}
	}
	maxPayout := 0
	for _, record := range mergedRecords {
		if maxPayout < record.Payout().Value() {
			maxPayout = record.Payout().Value()
		}
	}
	return types.Payout(maxPayout)
}

// getGradeClassMinPayout クラス別最小回収額の合計を取得する(全期間)
func (s *Summarizer) getGradeClassMinPayout(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, gradeClasses ...race_vo.GradeClass) types.Payout {
	recordsGroup := s.bettingTicketConverter.ConvertToRaceClassRecordsMap(records, racingNumbers, races)
	var mergedRecords []*betting_ticket_entity.CsvEntity
	for _, gradeClass := range gradeClasses {
		if recordsByBettingTicket, ok := recordsGroup[gradeClass]; ok {
			mergedRecords = append(mergedRecords, recordsByBettingTicket...)
		}
	}
	minPayout := 0
	for _, record := range mergedRecords {
		if record.Payout() == 0 {
			continue
		}
		if minPayout == 0 || minPayout > record.Payout().Value() {
			minPayout = record.Payout().Value()
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

// getCourseCategoryBetCount コース別投票数の合計を取得する(全期間)
func (s *Summarizer) getCourseCategoryBetCount(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, courseCategory race_vo.CourseCategory) types.BetCount {
	courseCategoryRecordsMap := s.getCourseCategoryRecordsMap(records, racingNumbers, races)
	if recordsByCourseCategory, ok := courseCategoryRecordsMap[courseCategory]; ok {
		return types.BetCount(len(recordsByCourseCategory))
	}
	return types.BetCount(0)
}

// getCourseCategoryHitCount コース別的中数の合計を取得する(全期間)
func (s *Summarizer) getCourseCategoryHitCount(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, courseCategory race_vo.CourseCategory) types.HitCount {
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

// getCourseCategoryPayment コース別払戻金の合計を取得する(全期間)
func (s *Summarizer) getCourseCategoryPayment(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, courseCategory race_vo.CourseCategory) types.Payment {
	courseCategoryRecordsMap := s.getCourseCategoryRecordsMap(records, racingNumbers, races)
	if recordsByCourseCategory, ok := courseCategoryRecordsMap[courseCategory]; ok {
		payment, _ := getSumAmount(recordsByCourseCategory)
		return payment
	}
	return types.Payment(0)
}

// getCourseCategoryPayout コース別払戻金の合計を取得する(全期間)
func (s *Summarizer) getCourseCategoryPayout(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, courseCategory race_vo.CourseCategory) types.Payout {
	courseCategoryRecordsMap := s.getCourseCategoryRecordsMap(records, racingNumbers, races)
	if recordsByCourseCategory, ok := courseCategoryRecordsMap[courseCategory]; ok {
		_, payout := getSumAmount(recordsByCourseCategory)
		return payout
	}
	return types.Payout(0)
}

// getCourseCategoryAveragePayout コース別平均払戻金を取得する(全期間)
func (s *Summarizer) getCourseCategoryAveragePayout(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, courseCategory race_vo.CourseCategory) types.Payout {
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
	if payout > 0 {
		return types.Payout(int(float64(payout) / float64(len(hitRecords))))
	}

	return types.Payout(0)
}

// getCourseCategoryMaxPayout コース別最大払戻金を取得する(全期間)
func (s *Summarizer) getCourseCategoryMaxPayout(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, courseCategory race_vo.CourseCategory) types.Payout {
	courseCategoryRecordsMap := s.getCourseCategoryRecordsMap(records, racingNumbers, races)
	maxPayout := 0
	if recordsByCourseCategory, ok := courseCategoryRecordsMap[courseCategory]; ok {
		for _, record := range recordsByCourseCategory {
			if maxPayout < record.Payout().Value() {
				maxPayout = record.Payout().Value()
			}
		}
	}
	return types.Payout(maxPayout)
}

// getCourseCategoryMinPayout コース別最小払戻金を取得する(全期間)
func (s *Summarizer) getCourseCategoryMinPayout(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, courseCategory race_vo.CourseCategory) types.Payout {
	courseCategoryRecordsMap := s.getCourseCategoryRecordsMap(records, racingNumbers, races)
	minPayout := 0
	if recordsByCourseCategory, ok := courseCategoryRecordsMap[courseCategory]; ok {
		for _, record := range recordsByCourseCategory {
			if record.Payout() == 0 {
				continue
			}
			if minPayout == 0 || minPayout > record.Payout().Value() {
				minPayout = record.Payout().Value()
			}
		}
	}
	return types.Payout(minPayout)
}

func (s *Summarizer) getDistanceCategoryRecordsMap(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race) map[race_vo.DistanceCategory][]*betting_ticket_entity.CsvEntity {
	distanceCategoryRecordsMap := map[race_vo.DistanceCategory][]*betting_ticket_entity.CsvEntity{}
	for distanceCategory, recordsGroup := range s.bettingTicketConverter.ConvertToDistanceCategoryRecordsMap(records, racingNumbers, races) {
		distanceCategoryRecordsMap[distanceCategory] = recordsGroup
	}
	return distanceCategoryRecordsMap
}

// getDistanceCategoryBetCount 距離別投票数を取得する(全期間)
func (s *Summarizer) getDistanceCategoryBetCount(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, distanceCategory race_vo.DistanceCategory) types.BetCount {
	distanceCategoryRecordsMap := s.getDistanceCategoryRecordsMap(records, racingNumbers, races)
	if recordsByDistanceCategory, ok := distanceCategoryRecordsMap[distanceCategory]; ok {
		return types.BetCount(len(recordsByDistanceCategory))
	}
	return types.BetCount(0)
}

// getDistanceCategoryHitCount 距離別的中数を取得する(全期間)
func (s *Summarizer) getDistanceCategoryHitCount(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, distanceCategory race_vo.DistanceCategory) types.HitCount {
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

// getDistanceCategoryPayment 距離別投票金額を取得する(全期間)
func (s *Summarizer) getDistanceCategoryPayment(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, distanceCategory race_vo.DistanceCategory) types.Payment {
	distanceCategoryRecordsMap := s.getDistanceCategoryRecordsMap(records, racingNumbers, races)
	if recordsByDistanceCategory, ok := distanceCategoryRecordsMap[distanceCategory]; ok {
		payment, _ := getSumAmount(recordsByDistanceCategory)
		return payment
	}
	return types.Payment(0)
}

// getDistanceCategoryPayout 距離別払戻金を取得する(全期間)
func (s *Summarizer) getDistanceCategoryPayout(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, distanceCategory race_vo.DistanceCategory) types.Payout {
	distanceCategoryRecordsMap := s.getDistanceCategoryRecordsMap(records, racingNumbers, races)
	if recordsByDistanceCategory, ok := distanceCategoryRecordsMap[distanceCategory]; ok {
		_, payout := getSumAmount(recordsByDistanceCategory)
		return payout
	}
	return types.Payout(0)
}

// getDistanceCategoryAveragePayout 距離別平均払戻金を取得する(全期間)
func (s *Summarizer) getDistanceCategoryAveragePayout(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, distanceCategory race_vo.DistanceCategory) types.Payout {
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
	if payout > 0 {
		return types.Payout(int(float64(payout) / float64(len(hitRecords))))
	}

	return types.Payout(0)
}

// getDistanceCategoryMaxPayout 距離別最高払戻金を取得する(全期間)
func (s *Summarizer) getDistanceCategoryMaxPayout(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, distanceCategory race_vo.DistanceCategory) types.Payout {
	distanceCategoryRecordsMap := s.getDistanceCategoryRecordsMap(records, racingNumbers, races)
	maxPayout := 0
	if recordsByDistanceCategory, ok := distanceCategoryRecordsMap[distanceCategory]; ok {
		for _, record := range recordsByDistanceCategory {
			if maxPayout < record.Payout().Value() {
				maxPayout = record.Payout().Value()
			}
		}
	}
	return types.Payout(maxPayout)
}

// getDistanceCategoryMinPayout 距離別最低払戻金を取得する(全期間)
func (s *Summarizer) getDistanceCategoryMinPayout(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, distanceCategory race_vo.DistanceCategory) types.Payout {
	distanceCategoryRecordsMap := s.getDistanceCategoryRecordsMap(records, racingNumbers, races)
	minPayout := 0
	if recordsByDistanceCategory, ok := distanceCategoryRecordsMap[distanceCategory]; ok {
		for _, record := range recordsByDistanceCategory {
			if record.Payout() == 0 {
				continue
			}
			if minPayout == 0 || minPayout > record.Payout().Value() {
				minPayout = record.Payout().Value()
			}
		}
	}
	return types.Payout(minPayout)
}

// getRaceCourseBetCount 競馬場別投票数を取得する(全期間)
func (s *Summarizer) getRaceCourseBetCount(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, raceCourses ...race_vo.RaceCourse) types.BetCount {
	raceCourseRecordsMap := s.getRaceCourseRecordsMap(records, racingNumbers, races)
	var mergedRecords []*betting_ticket_entity.CsvEntity
	for _, raceCourse := range raceCourses {
		if recordsByBettingTicket, ok := raceCourseRecordsMap[raceCourse]; ok {
			mergedRecords = append(mergedRecords, recordsByBettingTicket...)
		}
	}
	return types.BetCount(len(mergedRecords))
}

// getRaceCourseHitCount 競馬場別的中数を取得する(全期間)
func (s *Summarizer) getRaceCourseHitCount(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, raceCourses ...race_vo.RaceCourse) types.HitCount {
	raceCourseRecordsMap := s.getRaceCourseRecordsMap(records, racingNumbers, races)
	var mergedRecords []*betting_ticket_entity.CsvEntity
	for _, raceCourse := range raceCourses {
		if recordsByBettingTicket, ok := raceCourseRecordsMap[raceCourse]; ok {
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

// getRaceCoursePayment 競馬場別投票金額を取得する(全期間)
func (s *Summarizer) getRaceCoursePayment(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, raceCourses ...race_vo.RaceCourse) types.Payment {
	raceCourseRecordsMap := s.getRaceCourseRecordsMap(records, racingNumbers, races)
	var mergedRecords []*betting_ticket_entity.CsvEntity
	for _, raceCourse := range raceCourses {
		if recordsByBettingTicket, ok := raceCourseRecordsMap[raceCourse]; ok {
			mergedRecords = append(mergedRecords, recordsByBettingTicket...)
		}
	}
	payment, _ := getSumAmount(mergedRecords)
	return payment
}

// getRaceCoursePayout 競馬場別払戻金を取得する(全期間)
func (s *Summarizer) getRaceCoursePayout(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, raceCourses ...race_vo.RaceCourse) types.Payout {
	raceCourseRecordsMap := s.getRaceCourseRecordsMap(records, racingNumbers, races)
	var mergedRecords []*betting_ticket_entity.CsvEntity
	for _, raceCourse := range raceCourses {
		if recordsByBettingTicket, ok := raceCourseRecordsMap[raceCourse]; ok {
			mergedRecords = append(mergedRecords, recordsByBettingTicket...)
		}
	}
	_, payout := getSumAmount(mergedRecords)
	return payout
}

// getRaceCourseAveragePayout 競馬場別平均払戻金を取得する(全期間)
func (s *Summarizer) getRaceCourseAveragePayout(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, raceCourses ...race_vo.RaceCourse) types.Payout {
	raceCourseRecordsMap := s.getRaceCourseRecordsMap(records, racingNumbers, races)
	var mergedRecords []*betting_ticket_entity.CsvEntity
	for _, raceCourse := range raceCourses {
		if recordsByBettingTicket, ok := raceCourseRecordsMap[raceCourse]; ok {
			mergedRecords = append(mergedRecords, recordsByBettingTicket...)
		}
	}
	var hitRecords []*betting_ticket_entity.CsvEntity
	for _, record := range mergedRecords {
		if record.BettingResult() == betting_ticket_vo.Hit {
			hitRecords = append(hitRecords, record)
		}
	}
	_, payout := getSumAmount(hitRecords)
	if payout > 0 {
		return types.Payout(int(float64(payout) / float64(len(hitRecords))))
	}

	return types.Payout(0)
}

// getRaceCourseMaxPayout 競馬場別最高払戻金を取得する(全期間)
func (s *Summarizer) getRaceCourseMaxPayout(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, raceCourses ...race_vo.RaceCourse) types.Payout {
	raceCourseRecordsMap := s.getRaceCourseRecordsMap(records, racingNumbers, races)
	var mergedRecords []*betting_ticket_entity.CsvEntity
	for _, raceCourse := range raceCourses {
		if recordsByBettingTicket, ok := raceCourseRecordsMap[raceCourse]; ok {
			mergedRecords = append(mergedRecords, recordsByBettingTicket...)
		}
	}
	maxPayout := 0
	for _, record := range mergedRecords {
		if maxPayout < record.Payout().Value() {
			maxPayout = record.Payout().Value()
		}
	}
	return types.Payout(maxPayout)
}

// getRaceCourseMinPayout 競馬場別最低払戻金を取得する(全期間)
func (s *Summarizer) getRaceCourseMinPayout(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race, raceCourses ...race_vo.RaceCourse) types.Payout {
	raceCourseRecordsMap := s.getRaceCourseRecordsMap(records, racingNumbers, races)
	var mergedRecords []*betting_ticket_entity.CsvEntity
	for _, raceCourse := range raceCourses {
		if recordsByBettingTicket, ok := raceCourseRecordsMap[raceCourse]; ok {
			mergedRecords = append(mergedRecords, recordsByBettingTicket...)
		}
	}
	minPayout := 0
	for _, record := range mergedRecords {
		if record.Payout() == 0 {
			continue
		}
		if minPayout == 0 || minPayout > record.Payout().Value() {
			minPayout = record.Payout().Value()
		}
	}
	return types.Payout(minPayout)
}

func (s *Summarizer) getRaceCourseRecordsMap(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber, races []*race_entity.Race) map[race_vo.RaceCourse][]*betting_ticket_entity.CsvEntity {
	raceCourseRecordsMap := map[race_vo.RaceCourse][]*betting_ticket_entity.CsvEntity{}
	for raceCourse, recordsGroup := range s.bettingTicketConverter.ConvertToRaceCourseRecordsMap(records, racingNumbers, races) {
		raceCourseRecordsMap[raceCourse] = recordsGroup
	}
	return raceCourseRecordsMap
}

func (s *Summarizer) getBettingTicketRecordsMapForMonthly(records []*betting_ticket_entity.CsvEntity, bettingTicketTypes ...betting_ticket_vo.BettingTicket) (map[int][]*betting_ticket_entity.CsvEntity, []int) {
	recordsGroup := s.bettingTicketConverter.ConvertToBettingTicketRecordsMap(records)
	mergedMonthlyRecordsMap := map[int][]*betting_ticket_entity.CsvEntity{}
	var dateList []int
	for _, bettingTicketType := range bettingTicketTypes {
		if recordsByBettingTicket, ok := recordsGroup[bettingTicketType]; ok {
			for date, monthlyRecordsByBettingTicket := range s.bettingTicketConverter.ConvertToMonthRecordsMap(recordsByBettingTicket) {
				mergedMonthlyRecordsMap[date] = append(mergedMonthlyRecordsMap[date], monthlyRecordsByBettingTicket...)
				dateList = append(dateList, date)
			}
		}
	}
	return mergedMonthlyRecordsMap, dateList
}

func getSumAmount(records []*betting_ticket_entity.CsvEntity) (types.Payment, types.Payout) {
	var (
		sumPayment int
		sumPayout  int
	)
	for _, record := range records {
		sumPayment += record.Payment().Value()
		sumPayout += record.Payout().Value()
	}

	return types.Payment(sumPayment), types.Payout(sumPayout)
}
