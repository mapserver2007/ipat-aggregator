package infrastructure

import (
	"context"
	"encoding/csv"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/master_service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	ticketPatDataSuffix = "_tohyo"
)

type ticketRepository struct {
	betNumberConverter master_service.BetNumberConverter
}

func NewTicketRepository(
	betNumberConverter master_service.BetNumberConverter,
) repository.TicketRepository {
	return &ticketRepository{
		betNumberConverter: betNumberConverter,
	}
}

func (t *ticketRepository) List(ctx context.Context, path string) ([]string, error) {
	rootPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	absPath, err := filepath.Abs(fmt.Sprintf("%s/%s", rootPath, path))
	if err != nil {
		return nil, err
	}

	pattern := filepath.Join(absPath, "*.csv")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	fileNames := make([]string, 0, len(files))
	for _, file := range files {
		if !strings.Contains(file, ticketPatDataSuffix) {
			continue
		}
		fileNames = append(fileNames, filepath.Base(file))
	}

	return fileNames, nil
}

func (t *ticketRepository) Read(ctx context.Context, path string) ([]*ticket_csv_entity.Ticket, error) {
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
		rawRaceCourse := record[3]
		rawRaceNo := record[5]
		rawTicketType := record[6]
		rawPayment := t.extractPayment(ctx, rawTicketType, record[8])

		betNumbers, err := t.convertToSubTicketTypeBetNumbers(ctx, rawTicketType, record[7])
		if err != nil {
			return nil, err
		}

		var hitBetNumber types.BetNumber
		if strings.HasPrefix(record[9], "的中") {
			hitBetNumber = types.NewBetNumber(strings.Split(record[9], "的中")[1])
		}

		for _, betNumber := range betNumbers {
			var (
				rawPayout       string
				rawTicketResult bool
			)
			if hitBetNumber == betNumber {
				rawPayout = record[11]
				rawTicketResult = true
			}
			ticket, err := ticket_csv_entity.NewTicket(
				betNumber,
				rawRaceDate,
				rawRaceCourse,
				rawRaceNo,
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

func (t *ticketRepository) convertToSubTicketTypeBetNumbers(
	ctx context.Context,
	rawTicketType,
	rawBetNumber string,
) ([]types.BetNumber, error) {
	ticketType := types.NewTicketType(rawTicketType)
	switch ticketType {
	case types.QuinellaWheel:
		return t.betNumberConverter.QuinellaWheelToQuinellaBetNumbers(ctx, rawBetNumber)
	case types.QuinellaPlaceWheel:
		return t.betNumberConverter.QuinellaPlaceWheelToQuinellaBetNumbers(ctx, rawBetNumber)
	case types.QuinellaPlaceFormation:
		return t.betNumberConverter.QuinellaPlaceFormationToQuinellaBetNumbers(ctx, rawBetNumber)
	case types.ExactaWheelOfFirst:
		return t.betNumberConverter.ExactaWheelOfFirstToExactaBetNumbers(ctx, rawBetNumber)
	case types.TrioFormation:
		return t.betNumberConverter.TrioFormationToTrioBetNumbers(ctx, rawBetNumber)
	case types.TrioWheelOfFirst:
		return t.betNumberConverter.TrioWheelOfFirstToTrioBetNumbers(ctx, rawBetNumber)
	case types.TrioWheelOfSecond:
		return t.betNumberConverter.TrioWheelOfSecondToTrioBetNumbers(ctx, rawBetNumber)
	case types.TrioBox:
		return t.betNumberConverter.TrioBoxToTrioBetNumbers(ctx, rawBetNumber)
	case types.TrifectaFormation:
		return t.betNumberConverter.TrifectaFormationToTrifectaBetNumbers(ctx, rawBetNumber)
	case types.TrifectaWheelOfFirst:
		return t.betNumberConverter.TrifectaWheelOfFirstToTrifectaBetNumbers(ctx, rawBetNumber)
	case types.TrifectaWheelOfSecond:
		return t.betNumberConverter.TrifectaWheelOfSecondToTrifectaBetNumbers(ctx, rawBetNumber)
	case types.TrifectaWheelOfFirstMulti, types.TrifectaWheelOfSecondMulti:
		return t.betNumberConverter.TrifectaWheelMultiToTrifectaBetNumbers(ctx, rawBetNumber)
	case types.UnknownTicketType:
		return nil, fmt.Errorf("unknown ticket type")
	}

	return []types.BetNumber{types.NewBetNumber(rawBetNumber)}, nil
}

func (t *ticketRepository) extractPayment(
	ctx context.Context,
	rawTicketType,
	rawPayment string,
) string {
	ticketType := types.NewTicketType(rawTicketType)
	if types.QuinellaWheel == ticketType ||
		types.QuinellaPlaceWheel == ticketType ||
		types.QuinellaPlaceFormation == ticketType ||
		types.ExactaWheelOfFirst == ticketType ||
		types.TrioFormation == ticketType ||
		types.TrioWheelOfFirst == ticketType ||
		types.TrioWheelOfSecond == ticketType ||
		types.TrioBox == ticketType ||
		types.TrifectaFormation == ticketType ||
		types.TrifectaWheelOfFirst == ticketType ||
		types.TrifectaWheelOfSecond == ticketType ||
		types.TrifectaWheelOfFirstMulti == ticketType ||
		types.TrifectaWheelOfSecondMulti == ticketType {
		// 3連複軸1頭ながし, 馬単、3連単1着流し:
		// (1点あたりの購入金額)／(合計金額)
		separator := "／"
		values := strings.Split(rawPayment, separator)
		return values[0]
	}

	return rawPayment
}
