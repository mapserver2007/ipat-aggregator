package service

import (
	"fmt"
	betting_ticket_entity "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/entity"
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
