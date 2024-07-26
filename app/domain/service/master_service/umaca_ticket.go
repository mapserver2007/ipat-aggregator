package master_service

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/umaca_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/config"
	"github.com/shopspring/decimal"
	"regexp"
	"strconv"
	"time"
)

const (
	umacaMasterFileName = "umaca_master.csv"
	umacaTicketFileName = "%d_tohyo_umaca.csv"
)

type UmacaTicket interface {
	Get(ctx context.Context, races []*data_cache_entity.Race) ([]*ticket_csv_entity.RaceTicket, error)
	CreateOrUpdate(ctx context.Context, races []*data_cache_entity.Race) error
	GetMaster(ctx context.Context) ([]*umaca_csv_entity.UmacaMaster, error)
}

type umacaTicketService struct {
	umacaTicketRepository repository.UmacaTicketRepository
	ticketRepository      repository.TicketRepository
}

func NewUmacaTicket(
	umacaTicketRepository repository.UmacaTicketRepository,
	ticketRepository repository.TicketRepository,
) UmacaTicket {
	return &umacaTicketService{
		umacaTicketRepository: umacaTicketRepository,
		ticketRepository:      ticketRepository,
	}
}

func (u *umacaTicketService) Get(
	ctx context.Context,
	races []*data_cache_entity.Race,
) ([]*ticket_csv_entity.RaceTicket, error) {
	files, err := u.umacaTicketRepository.List(ctx, config.CsvDir)
	if err != nil {
		return nil, err
	}

	raceDateMap := map[types.RaceDate][]*data_cache_entity.Race{}
	for _, race := range races {
		if _, ok := raceDateMap[race.RaceDate()]; !ok {
			raceDateMap[race.RaceDate()] = make([]*data_cache_entity.Race, 0)
		}
		raceDateMap[race.RaceDate()] = append(raceDateMap[race.RaceDate()], race)
	}

	raceTickets := make([]*ticket_csv_entity.RaceTicket, 0)
	for _, file := range files {
		tickets, err := u.ticketRepository.Read(ctx, fmt.Sprintf("%s/%s", config.CsvDir, file))
		if err != nil {
			return nil, err
		}
		for _, ticket := range tickets {
			if ticket.RaceCourse().NAR() || ticket.RaceCourse().Oversea() {
				// UMACA購入は中央のみ
				return nil, fmt.Errorf("umaca ticket is only for JRA")
			}
			raceDateRaces, ok := raceDateMap[ticket.RaceDate()]
			if !ok {
				continue
			}
			for _, race := range raceDateRaces {
				// racingNumberのように完全に紐付けることはできないので、レースNo、開催場所から特定する
				if race.RaceNumber() == ticket.RaceNo() && race.RaceCourseId() == ticket.RaceCourse() {
					raceTickets = append(raceTickets, ticket_csv_entity.NewRaceTicket(
						race.RaceId(),
						ticket,
					))
				}
			}
		}
	}

	return raceTickets, nil
}

func (u *umacaTicketService) CreateOrUpdate(
	ctx context.Context,
	races []*data_cache_entity.Race,
) error {
	raceMap := map[types.RaceId]*data_cache_entity.Race{}
	for _, race := range races {
		raceMap[race.RaceId()] = race
	}

	masters, err := u.GetMaster(ctx)

	raceDateMap := map[types.RaceDate][]*umaca_csv_entity.UmacaMaster{}
	for _, master := range masters {
		if _, ok := raceDateMap[master.RaceDate()]; !ok {
			raceDateMap[master.RaceDate()] = make([]*umaca_csv_entity.UmacaMaster, 0)
		}
		raceDateMap[master.RaceDate()] = append(raceDateMap[master.RaceDate()], master)
	}

	for raceDate, raceDateMasters := range raceDateMap {
		fileName := fmt.Sprintf(umacaTicketFileName, raceDate.Value())
		filePath := fmt.Sprintf("%s/%s", config.CsvDir, fileName)

		// ファイルが取得できない場合は処理を続行する
		tickets, _ := u.ticketRepository.Read(ctx, filePath)

		// データが取得できた場合は上書きしない
		// もし上書きしたい場合はファイルを消して対応する
		if len(tickets) > 0 {
			continue
		}

		var (
			records      [][]string
			totalPayment int
			totalPayout  int
		)

		records = append(records, []string{
			"日付",
			"受付番号",
			"通番",
			"場名",
			"曜日",
			"レース",
			"式別",
			"馬／組番",
			"購入金額",
			"的中／返還",
			"払戻単価",
			"払戻／返還金額",
		})

		for _, master := range raceDateMasters {
			var (
				decimalOdds    decimal.Decimal
				decimalPayment decimal.Decimal
			)
			ticketResult := types.TicketUnHit

			race, ok := raceMap[master.RaceId()]
			if !ok {
				return fmt.Errorf("race not found in umacaTicketService.CreateOrUpdate: %s", master.RaceId())
			}

			for _, payoutResult := range race.PayoutResults() {
				ticketType := payoutResult.TicketType()
				if master.TicketType() == ticketType {
					for idx, number := range payoutResult.Numbers() {
						if u.slicesEqual(master.BetNumber().List(), number.List()) {
							ticketResult = types.TicketHit
							decimalOdds, _ = decimal.NewFromString(payoutResult.Odds()[idx])
							v := decimal.NewFromInt(int64(master.Payment().Value()))
							decimalPayment = decimalOdds.Mul(v)
							break
						}
					}
				}
			}

			weekDayName, err := u.getWeekdayName(master.RaceDate())
			if err != nil {
				return err
			}

			hitOrReturn := "-"
			betNumberStr, err := u.paddingNumbers(master.BetNumber().String())
			if err != nil {
				return err
			}
			if ticketResult == types.TicketHit {
				hitOrReturn = fmt.Sprintf("的中%s", betNumberStr)
			}

			// patのフォーマット(全角ハイフン等)までは完全一致はさせない
			records = append(records, []string{
				strconv.Itoa(master.RaceDate().Value()),
				"0001", // dummy value
				"01",   // dummy value
				race.RaceCourseId().Name(),
				weekDayName,
				strconv.Itoa(race.RaceNumber()),
				master.TicketType().Name(),
				betNumberStr,
				strconv.Itoa(master.Payment().Value()),
				hitOrReturn,
				decimalOdds.Mul(decimal.NewFromInt32(100)).String(),
				decimalPayment.String(),
			})
			totalPayment += master.Payment().Value()
			totalPayout += int(decimalPayment.IntPart())
		}

		records = append(records, []string{
			"",
			"",
			"合計",
			"",
			"",
			"",
			"",
			"",
			strconv.Itoa(totalPayment),
			"",
			"",
			strconv.Itoa(totalPayout),
		})

		err = u.umacaTicketRepository.Write(ctx, filePath, records)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *umacaTicketService) GetMaster(ctx context.Context) ([]*umaca_csv_entity.UmacaMaster, error) {
	masters, err := u.umacaTicketRepository.GetMaster(ctx, fmt.Sprintf("%s/%s", config.CsvDir, umacaMasterFileName))
	if err != nil {
		return nil, err
	}

	return masters, nil
}

func (u *umacaTicketService) getWeekdayName(
	raceDate types.RaceDate,
) (string, error) {
	date, err := time.Parse("20060102", strconv.Itoa(raceDate.Value()))
	if err != nil {
		return "", err
	}

	weekDayName := ""
	switch date.Weekday() {
	case time.Monday:
		weekDayName = "月"
	case time.Tuesday:
		weekDayName = "火"
	case time.Wednesday:
		weekDayName = "水"
	case time.Thursday:
		weekDayName = "木"
	case time.Friday:
		weekDayName = "金"
	case time.Saturday:
		weekDayName = "土"
	case time.Sunday:
		weekDayName = "日"
	}

	return weekDayName, nil
}

func (u *umacaTicketService) slicesEqual(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func (u *umacaTicketService) paddingNumbers(input string) (string, error) {
	re := regexp.MustCompile(`\d+`)
	result := re.ReplaceAllStringFunc(input, func(numStr string) string {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return numStr
		}
		if num < 10 {
			return fmt.Sprintf("%02d", num)
		}
		return numStr
	})

	return result, nil
}
