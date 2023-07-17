package service

import (
	"fmt"
	betting_ticket_entity "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/entity"
	betting_ticket_vo "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/value_object"
	result_entity "github.com/mapserver2007/ipat-aggregator/app/domain/result/entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/result/types"
	"strconv"
	"time"
)

type Summarizer struct {
	bettingTicketConverter BettingTicketConverter
}

func NewSummarizer(
	bettingTicketConverter BettingTicketConverter,
) Summarizer {
	return Summarizer{
		bettingTicketConverter: bettingTicketConverter,
	}
}

func (s *Summarizer) GetShortSummaryForAll(records []*betting_ticket_entity.CsvEntity) result_entity.ShortSummary {
	return result_entity.NewShortSummary(
		s.getPaymentForAll(records),
		s.getPayoutForAll(records),
		s.getRecoveryRateForAll(records),
	)
}

func (s *Summarizer) GetShortSummaryForMonth(records []*betting_ticket_entity.CsvEntity) result_entity.ShortSummary {
	return result_entity.NewShortSummary(
		s.getPaymentForMonth(records),
		s.getPayoutForMonth(records),
		s.getRecoveryRateForMonth(records),
	)
}

func (s *Summarizer) GetShortSummaryForYear(records []*betting_ticket_entity.CsvEntity) result_entity.ShortSummary {
	return result_entity.NewShortSummary(
		s.getPaymentForYear(records),
		s.getPayoutForYear(records),
		s.getRecoveryRateForYear(records),
	)
}

func (s *Summarizer) GetBettingTicketSummaryForAll(records []*betting_ticket_entity.CsvEntity) result_entity.DetailSummary {
	return result_entity.NewDetailSummary(
		s.getBettingTicketWinVoteCountForAll(records),
		s.getBettingTicketWinHitCountForAll(records),
		s.getBettingTicketWinPaymentForAll(records),
		s.getBettingTicketWinPayoutForAll(records),
		s.getBettingTicketWinAveragePayoutForAll(records),
		s.getBettingTicketWinMaxPayoutForAll(records),
		s.getBettingTicketWinMinPayoutForAll(records),
	)
}

// GetPaymentForAll 投資額の合計を取得する(全期間)
func (s *Summarizer) getPaymentForAll(records []*betting_ticket_entity.CsvEntity) types.Payment {
	payment, _ := getSumAmount(records)
	return payment
}

// getPayoutForAll 回収額の合計を取得する(全期間)
func (s *Summarizer) getPayoutForAll(records []*betting_ticket_entity.CsvEntity) types.Payout {
	_, payout := getSumAmount(records)
	return payout
}

// getRecoveryRateForAll 回収率の合計を取得する(全期間)
func (s *Summarizer) getRecoveryRateForAll(records []*betting_ticket_entity.CsvEntity) string {
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

// getBettingTicketWinPaymentForAll 券種別投資額の合計を取得する(全期間)
func (s *Summarizer) getBettingTicketWinPaymentForAll(records []*betting_ticket_entity.CsvEntity) types.Payment {
	recordsGroup := s.bettingTicketConverter.ConvertToBettingTicketRecordsMap(records)
	if recordsForWin, ok := recordsGroup[betting_ticket_vo.Win]; ok {
		payment, _ := getSumAmount(recordsForWin)
		return payment
	}

	return types.Payment(0)
}

// getBettingTicketWinPayoutForAll 券種別回収額の合計を取得する(全期間)
func (s *Summarizer) getBettingTicketWinPayoutForAll(records []*betting_ticket_entity.CsvEntity) types.Payout {
	recordsGroup := s.bettingTicketConverter.ConvertToBettingTicketRecordsMap(records)
	if recordsForWin, ok := recordsGroup[betting_ticket_vo.Win]; ok {
		_, payout := getSumAmount(recordsForWin)
		return payout
	}

	return types.Payout(0)
}

// getBettingTicketWinRecoveryRateForAll 券種別回収率の合計を取得する(全期間)
func (s *Summarizer) getBettingTicketWinRecoveryRateForAll(records []*betting_ticket_entity.CsvEntity) string {
	payment := s.getBettingTicketWinPaymentForAll(records)
	payout := s.getBettingTicketWinPayoutForAll(records)
	if payment == 0 {
		return fmt.Sprintf("%d%s", 0, "%")
	}
	return fmt.Sprintf("%s%s", strconv.FormatFloat((float64(payout)*float64(100))/float64(payment), 'f', 2, 64), "%")
}

// getBettingTicketWinVoteCountForAll 券種別投票数の合計を取得する(全期間)
func (s *Summarizer) getBettingTicketWinVoteCountForAll(records []*betting_ticket_entity.CsvEntity) types.BetCount {
	recordsGroup := s.bettingTicketConverter.ConvertToBettingTicketRecordsMap(records)
	if recordsForWin, ok := recordsGroup[betting_ticket_vo.Win]; ok {
		return types.BetCount(len(recordsForWin))
	}

	return types.BetCount(0)
}

// getBettingTicketWinHitCountForAll 券種別的中数の合計を取得する(全期間)
func (s *Summarizer) getBettingTicketWinHitCountForAll(records []*betting_ticket_entity.CsvEntity) types.HitCount {
	recordsGroup := s.bettingTicketConverter.ConvertToBettingTicketRecordsMap(records)
	if recordsForWin, ok := recordsGroup[betting_ticket_vo.Win]; ok {
		hitCount := 0
		for _, record := range recordsForWin {
			if record.BettingResult() == betting_ticket_vo.Hit {
				hitCount++
			}
		}
		return types.HitCount(hitCount)
	}

	return types.HitCount(0)
}

func (s *Summarizer) getBettingTicketWinMaxPayoutForAll(records []*betting_ticket_entity.CsvEntity) types.Payout {
	recordsGroup := s.bettingTicketConverter.ConvertToBettingTicketRecordsMap(records)
	if recordsForWin, ok := recordsGroup[betting_ticket_vo.Win]; ok {
		maxPayout := 0
		for _, record := range recordsForWin {
			if maxPayout < record.Repayment() {
				maxPayout = record.Repayment()
			}
		}
		return types.Payout(maxPayout)
	}

	return types.Payout(0)
}

func (s *Summarizer) getBettingTicketWinMinPayoutForAll(records []*betting_ticket_entity.CsvEntity) types.Payout {
	recordsGroup := s.bettingTicketConverter.ConvertToBettingTicketRecordsMap(records)
	if recordsForWin, ok := recordsGroup[betting_ticket_vo.Win]; ok {
		minPayout := 0
		for _, record := range recordsForWin {
			if record.Repayment() == 0 {
				continue
			}
			if minPayout == 0 || minPayout > record.Repayment() {
				minPayout = record.Repayment()
			}
		}
		return types.Payout(minPayout)
	}

	return types.Payout(0)
}

func (s *Summarizer) getBettingTicketWinAveragePayoutForAll(records []*betting_ticket_entity.CsvEntity) types.Payout {
	recordsGroup := s.bettingTicketConverter.ConvertToBettingTicketRecordsMap(records)
	if recordsForWin, ok := recordsGroup[betting_ticket_vo.Win]; ok {
		_, payout := getSumAmount(recordsForWin)
		return types.Payout(int(float64(payout) / float64(len(recordsForWin))))
	}

	return types.Payout(0)
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
