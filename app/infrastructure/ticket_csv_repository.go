package infrastructure

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"io"
	"os"
	"strings"
)

type ticketCsvRepository struct {
	betNumberConverter service.BetNumberConverter
}

func NewTicketCsvRepository(
	betNumberConverter service.BetNumberConverter,
) repository.TicketCsvRepository {
	return &ticketCsvRepository{
		betNumberConverter: betNumberConverter,
	}
}

func (t *ticketCsvRepository) Read(ctx context.Context, path string) ([]*ticket_csv_entity.Ticket, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var tickets []*ticket_csv_entity.Ticket
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

		rawRaceDate := record[0]
		rawRaceEntryNo := record[1]
		rawRaceCourse := record[3]
		rawRaceNo := record[5]
		rawTicketType := record[6]
		rawPayment := t.extractPayment(rawTicketType, record[8])

		rawBetNumbers, err := t.convertToSubTicketTypeBetNumbers(ctx, rawTicketType, record[7])
		if err != nil {
			return nil, err
		}

		var hitBetNumber string
		if strings.HasPrefix(record[9], "的中") {
			hitBetNumber = strings.Split(record[9], "的中")[1]
		}

		for _, rawBetNumber := range rawBetNumbers {
			var (
				rawPayout       string
				rawTicketResult bool
			)
			if hitBetNumber == rawBetNumber {
				rawPayout = record[11]
				rawTicketResult = true
			}
			ticket, err := ticket_csv_entity.NewTicket(
				rawRaceDate,
				rawRaceEntryNo,
				rawRaceCourse,
				rawRaceNo,
				rawBetNumber,
				rawTicketType,
				rawTicketResult,
				rawPayment,
				rawPayout,
			)
			if err != nil {
				return nil, err
			}

			tickets = append(tickets, ticket)
		}
	}

	return tickets, nil
}

func (t *ticketCsvRepository) convertToSubTicketTypeBetNumbers(
	ctx context.Context,
	rawTicketType,
	rawBetNumber string,
) ([]string, error) {
	ticketType := types.NewTicketType(rawTicketType)
	switch ticketType {
	case types.ExactaWheelOfFirst:
		return t.betNumberConverter.ExactaWheelOfFirstToExactaBetNumbers(ctx, rawBetNumber)
	case types.QuinellaPlaceWheel:
		return t.betNumberConverter.QuinellaPlaceWheelToQuinellaBetNumbers(ctx, rawBetNumber)
	case types.TrioFormation:
		return t.betNumberConverter.TrioFormationToTrioBetNumbers(ctx, rawBetNumber)
	case types.TrioWheelOfFirst:
		return t.betNumberConverter.TrioWheelOfFirstToTrioBetNumbers(ctx, rawBetNumber)
	case types.TrifectaFormation:
		return t.betNumberConverter.TrifectaFormationToTrifectaBetNumbers(ctx, rawBetNumber)
	case types.TrifectaWheelOfFirst:
		return t.betNumberConverter.TrifectaWheelOfFirstToTrifectaBetNumbers(ctx, rawBetNumber)
	case types.TrifectaWheelOfSecondMulti:
		return t.betNumberConverter.TrifectaWheelMultiToTrifectaBetNumbers(ctx, rawBetNumber)
	case types.UnknownTicketType:
		return nil, fmt.Errorf("unknown betting ticket")
	}

	return []string{rawBetNumber}, nil
}

func (t *ticketCsvRepository) extractPayment(rawTicketType, rawPayment string) string {
	ticketType := types.NewTicketType(rawTicketType)
	if types.ExactaWheelOfFirst == ticketType ||
		types.QuinellaPlaceWheel == ticketType ||
		types.TrioFormation == ticketType ||
		types.TrioWheelOfFirst == ticketType ||
		types.TrifectaFormation == ticketType ||
		types.TrifectaWheelOfFirst == ticketType ||
		types.TrifectaWheelOfSecondMulti == ticketType {
		// 3連複軸1頭ながし, 馬単、3連単1着流し:
		// (1点あたりの購入金額)／(合計金額)
		separator := "／"
		values := strings.Split(rawPayment, separator)
		return values[0]
	}

	return rawPayment
}
