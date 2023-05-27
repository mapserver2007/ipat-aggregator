package service

import (
	"context"
	"encoding/csv"
	"github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/entity"
	betting_ticket_vo "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/value_object"
	race_vo "github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"io"
	"os"
	"strings"
)

type CsvReader struct{}

func NewCsvReader() CsvReader {
	return CsvReader{}
}

func (c *CsvReader) Read(ctx context.Context, filePath string) ([]*entity.CsvEntity, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entities []*entity.CsvEntity
	reader := csv.NewReader(transform.NewReader(f, japanese.ShiftJIS.NewDecoder()))
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if record[0] == "日付" || record[0] == "" {
			continue
		}
		if strings.Contains(record[9], "返還") {
			continue
		}

		bettingTicket := betting_ticket_vo.ConvertToBettingTicket(record[6])

		bettingResult := betting_ticket_vo.UnHit
		if record[9] != "―" {
			bettingResult = betting_ticket_vo.Hit
		}

		betNumbers := []betting_ticket_vo.BetNumber{betting_ticket_vo.NewBetNumber(record[7])}
		raceDate := ConvertToRaceDate(record[0])
		raceNo := ConvertToIntValue(record[5])
		raceCourse := race_vo.ConvertToRaceCourse(record[3])

		// 連複、連単系の流し、フォーメーションなどの場合は金額形式が異なる
		var payment int
		switch bettingTicket {
		case betting_ticket_vo.ExactaWheelOfFirst:
			betNumbers = ConvertToBetNumbersForExacta(record[7])
			payment = ConvertToPaymentForWheel(record[8])[0]
		case betting_ticket_vo.QuinellaPlaceWheel:
			betNumbers = ConvertToBetNumbersForQuinella(record[7])
			payment = ConvertToPaymentForWheel(record[8])[0]
		case betting_ticket_vo.TrioFormation:
			betNumbers = ConvertToPaymentForFoTrioFormation(record[7])
			payment = ConvertToPaymentForWheel(record[8])[0]
		case betting_ticket_vo.TrioWheelOfFirst:
			betNumbers = ConvertToBetNumbersForTrio(record[7])
			payment = ConvertToPaymentForWheel(record[8])[0]
		case betting_ticket_vo.TrifectaFormation:
			betNumbers = ConvertToBetNumbersForTrifectaFormation(record[7])
			payment = ConvertToPaymentForWheel(record[8])[0]
		case betting_ticket_vo.TrifectaWheelOfFirst:
			betNumbers = ConvertToBetNumbersForTrifecta(record[7])
			payment = ConvertToPaymentForWheel(record[8])[0]
		case betting_ticket_vo.UnknownTicket:
			panic("unknown betting ticket")
		default:
			payment = ConvertToIntValue(record[8])
		}
		var winningBetNumber betting_ticket_vo.BetNumber
		if strings.HasPrefix(record[9], "的中") {
			rawBetNumber := strings.Split(record[9], "的中")[1]
			winningBetNumber = betting_ticket_vo.NewBetNumber(rawBetNumber)
		}

		for _, betNumber := range betNumbers {
			repayment := 0
			winning := false
			if winningBetNumber == betNumber {
				repayment = ConvertToIntValue(record[11])
				winning = true
			}
			entities = append(entities, entity.NewCsvEntity(
				raceDate,
				ConvertToIntValue(record[1]),
				raceCourse,
				raceNo,
				bettingTicket,
				bettingResult,
				winning,
				betNumber,
				payment,
				repayment,
			))
		}
	}

	return entities, nil
}
