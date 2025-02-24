package aggregation_service

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/list_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/ticket_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/shopspring/decimal"
	"golang.org/x/exp/slices"
)

const (
	weightOfFirstPlace           = 1.0
	weightOfSecondPlace          = 0.25
	weightOfThirdPlace           = 0.1
	thresholdOfLowerLimitPayment = 0.15
)

var ticketSortOrders = []types.TicketType{
	types.Win,
	types.Exacta,
	types.ExactaWheelOfFirst,
	types.Trifecta,
	types.TrifectaWheelOfFirst,
	types.TrifectaWheelOfSecond,
	types.TrifectaFormation,
	types.TrifectaWheelOfFirstMulti,
	types.TrifectaWheelOfSecondMulti,
	types.QuinellaWheel,
	types.QuinellaPlaceWheel,
	types.Quinella,
	types.QuinellaPlace,
	types.QuinellaPlaceFormation,
	types.TrioWheelOfFirst,
	types.TrioWheelOfSecond,
	types.Trio,
	types.TrioFormation,
	types.TrioBox,
	types.Place,
	types.BracketQuinella,
}

type List interface {
	Create(ctx context.Context,
		tickets []*ticket_csv_entity.RaceTicket,
		races []*data_cache_entity.Race,
		jockeys []*data_cache_entity.Jockey,
	) ([]*spreadsheet_entity.ListRow, error)
	Write(ctx context.Context, listRows []*spreadsheet_entity.ListRow) error
}

type listService struct {
	raceEntityConverter   converter.RaceEntityConverter
	JockeyEntityConverter converter.JockeyEntityConverter
	spreadSheetRepository repository.SpreadSheetRepository
}

func NewList(
	raceEntityConverter converter.RaceEntityConverter,
	JockeyEntityConverter converter.JockeyEntityConverter,
	spreadSheetRepository repository.SpreadSheetRepository,
) List {
	return &listService{
		raceEntityConverter:   raceEntityConverter,
		JockeyEntityConverter: JockeyEntityConverter,
		spreadSheetRepository: spreadSheetRepository,
	}
}

func (l *listService) Create(
	ctx context.Context,
	tickets []*ticket_csv_entity.RaceTicket,
	races []*data_cache_entity.Race,
	jockeys []*data_cache_entity.Jockey,
) ([]*spreadsheet_entity.ListRow, error) {
	var listRows []*spreadsheet_entity.ListRow
	raceMap := converter.ConvertToMap(races, func(race *data_cache_entity.Race) types.RaceId {
		return race.RaceId()
	})
	raceTicketsMap := converter.ConvertToSliceMap(tickets, func(ticket *ticket_csv_entity.RaceTicket) types.RaceId {
		return ticket.RaceId()
	})
	jockeyMap := converter.ConvertToMap(jockeys, func(jockey *data_cache_entity.Jockey) types.JockeyId {
		return jockey.JockeyId()
	})

	getJockeyName := func(jockeyId types.JockeyId) string {
		jockey, ok := jockeyMap[jockeyId]
		if ok {
			return jockey.JockeyName()
		}
		return "(不明)"
	}

	for raceId, raceTickets := range raceTicketsMap {
		race, ok := raceMap[raceId]
		if !ok {
			return nil, fmt.Errorf("unknown raceId: %s", raceId)
		}
		raceResults := race.RaceResults()
		sort.Slice(raceResults, func(i, j int) bool {
			return raceResults[i].PopularNumber() < raceResults[j].PopularNumber()
		})

		var rawPayment, rawPayout int
		for _, raceTicket := range raceTickets {
			rawPayment += raceTicket.Ticket().Payment().Value()
			rawPayout += raceTicket.Ticket().Payout().Value()
		}

		var (
			favorites, rivals []types.BetNumber
			favorite, rival   types.BetNumber
			hitTickets        []*list_entity.Ticket
		)

		ticketTypeMap := converter.ConvertToSliceMap(raceTickets, func(raceTicket *ticket_csv_entity.RaceTicket) types.TicketType {
			return raceTicket.Ticket().TicketType()
		})
		payoutResultMap := converter.ConvertToSliceMap(race.PayoutResults(), func(payoutResult *data_cache_entity.PayoutResult) types.TicketType {
			return payoutResult.TicketType()
		})

		status := types.PredictUncompleted
		for _, ticketType := range ticketSortOrders {
			ticketTypeRaceTickets, ok := ticketTypeMap[ticketType]
			if !ok || len(ticketTypeRaceTickets) == 0 {
				continue
			}

			// 本命候補がいる場合
			if status.Included(types.FavoriteCandidate) {
				favorites, rivals, status = l.getFavoritesAndRivals(ctx, ticketTypeRaceTickets, types.FavoriteCandidate, favorites, nil, rawPayment)
			}

			// 対抗候補がいる場合
			if status.Included(types.RivalCandidate) {
				_, rivals, status = l.getFavoritesAndRivals(ctx, ticketTypeRaceTickets, types.RivalCandidate, rivals, nil, rawPayment)
			}

			// 本命が決定済み、対抗が未決定
			if status.Included(types.FavoriteCompleted) {
				_, rivals, status = l.getFavoritesAndRivals(ctx, ticketTypeRaceTickets, types.FavoriteCompleted, nil, favorites, rawPayment)
			}

			// 本命、対抗が未決定
			if status.Matched(types.PredictUncompleted) {
				favorites, rivals, status = l.getFavoritesAndRivals(ctx, ticketTypeRaceTickets, types.PredictUncompleted, nil, nil, rawPayment)
			}

			// 本命または対抗が決定してる場合、処理を抜ける
			if status.Matched(types.FavoriteCompleted | types.RivalCompleted) {
				break
			}
		}

		// 本命が複数の場合
		if status.Matched(types.FavoriteCandidate) && len(favorites) >= 2 {
			// 候補になっている馬番絡みの払い戻し金額が最大の馬番に絞り込み、本命候補とする
			// 同額の場合は複数返る
			// それ以外の馬番は対抗候補とする
			favoriteCandidates, rivalCandidates, isFound := l.getBetNumbersByMaxPayout(ctx, favorites, raceTickets)
			if isFound {
				favorites = favoriteCandidates
				rivals = rivalCandidates
			}
			for _, raceResult := range raceResults {
				// 本命候補が複数の場合、人気が高い方を本命とする
				betNumber := types.NewBetNumber(fmt.Sprintf("%02d", raceResult.HorseNumber()))
				if l.containsInSlices(favorites, betNumber) {
					for _, candidate := range favorites {
						// 本命馬番にマッチしないものはすべて対抗候補に回す
						if candidate != betNumber {
							rivals = append(rivals, candidate)
						}
					}
					// 本命決定
					favorites = []types.BetNumber{betNumber}

					// 対抗が1つに決定している場合
					if len(rivals) == 1 {
						status = types.FavoriteCompleted | types.RivalCompleted
					} else {
						status = types.FavoriteCompleted | types.RivalCandidate
					}
					break
				}
			}
		}

		// 対抗が候補の場合
		if status.Matched(types.FavoriteCompleted | types.RivalCandidate) {
			if len(rivals) == 0 {
				// 対抗が存在しない場合
				status = types.FavoriteCompleted | types.RivalCompleted
			} else if len(rivals) >= 2 {
				// 候補になっている馬番絡みの払い戻し金額が最大の馬番に絞り込む
				// 同額の場合は複数返る
				rivalCandidates, _, isFound := l.getBetNumbersByMaxPayout(ctx, rivals, raceTickets)
				if isFound {
					rivals = rivalCandidates
				}
				// 対抗が複数存在する場合、人気が高いものを採用する
				for _, raceResult := range raceResults {
					betNumber := types.NewBetNumber(fmt.Sprintf("%02d", raceResult.HorseNumber()))
					if l.containsInSlices(rivals, betNumber) {
						rivals = []types.BetNumber{betNumber}
						break
					}
				}
			}
		}

		if len(favorites) >= 2 || len(rivals) >= 2 {
			return nil, fmt.Errorf("failed to find favorite or rival")
		}

		if favorites != nil && len(favorites) > 0 {
			favorite = favorites[0]
		}
		if rivals != nil && len(rivals) > 0 {
			rival = rivals[0]
		}

		// 本命対抗検出ロジック内だと最後の券種に到達する前にbreakするケースがあるためもう一度ループを回す
		for _, ticketType := range ticketSortOrders {
			ticketTypeRaceTickets, ok := ticketTypeMap[ticketType]
			if !ok || len(ticketTypeRaceTickets) == 0 {
				continue
			}
			for _, raceTicket := range ticketTypeRaceTickets {
				payoutResults, ok := payoutResultMap[raceTicket.Ticket().TicketType().OriginTicketType()]
				if !ok {
					return nil, fmt.Errorf("unknown payout result in ticketType %s", raceTicket.Ticket().TicketType().OriginTicketType().Name())
				}
				if raceTicket.Ticket().TicketResult() == types.TicketHit {
					for _, payoutResult := range payoutResults {
						for idx := range payoutResult.Numbers() {
							if payoutResult.Numbers()[idx] == raceTicket.Ticket().BetNumber() {
								hitTickets = append(hitTickets, list_entity.NewTicket(
									raceTicket.Ticket(),
									payoutResult.Numbers()[idx],
									payoutResult.Odds()[idx],
									payoutResult.Populars()[idx],
								))
							}
						}
					}
				}
			}
		}

		var (
			favoriteHorse, rivalHorse   *list_entity.Horse
			favoriteJockey, rivalJockey *list_entity.Jockey
		)

		for _, raceResult := range raceResults {
			if len(favorite) > 0 && raceResult.HorseNumber().Value() == favorite.List()[0] {
				favoriteHorse = list_entity.NewHorse(raceResult.HorseName(), raceResult.Odds(), raceResult.PopularNumber())
				jockey, ok := jockeyMap[raceResult.JockeyId()]
				if !ok {
					jockey = data_cache_entity.NewJockey("00000", "(不明)")
				}
				favoriteJockey = l.JockeyEntityConverter.DataCacheToList(jockey)
			}
			if len(rival) > 0 && raceResult.HorseNumber().Value() == rival.List()[0] {
				rivalHorse = list_entity.NewHorse(raceResult.HorseName(), raceResult.Odds(), raceResult.PopularNumber())
				jockey, ok := jockeyMap[raceResult.JockeyId()]
				if !ok {
					jockey = data_cache_entity.NewJockey("00000", "(不明)")
				}
				rivalJockey = l.JockeyEntityConverter.DataCacheToList(jockey)
			}
		}

		// 単複のみなど対抗が存在しない場合
		if rivalHorse == nil {
			rivalHorse = list_entity.NewHorse("-", decimal.Zero, 0)
			rivalJockey = l.JockeyEntityConverter.DataCacheToList(data_cache_entity.NewJockey("00000", "-"))
		}

		sort.Slice(raceResults, func(i, j int) bool {
			return raceResults[i].OrderNo() < raceResults[j].OrderNo()
		})

		listRace := l.raceEntityConverter.DataCacheToList(race)

		listRows = append(listRows, spreadsheet_entity.NewListRow(
			listRace,
			favoriteHorse,
			rivalHorse,
			favoriteJockey,
			rivalJockey,
			listRace.RaceResults()[0],
			list_entity.NewJockey(
				listRace.RaceResults()[0].JockeyId(),
				getJockeyName(listRace.RaceResults()[0].JockeyId()),
			),
			listRace.RaceResults()[1],
			list_entity.NewJockey(
				listRace.RaceResults()[1].JockeyId(),
				getJockeyName(listRace.RaceResults()[1].JockeyId()),
			),
			hitTickets,
			types.Payment(rawPayment),
			types.Payout(rawPayout),
		))
	}

	sort.SliceStable(listRows, func(i, j int) bool {
		return listRows[i].Data().RaceStartTime() > listRows[j].Data().RaceStartTime()
	})
	sort.SliceStable(listRows, func(i, j int) bool {
		return listRows[i].Data().RaceDate() > listRows[j].Data().RaceDate()
	})

	return listRows, nil
}

func (l *listService) Write(
	ctx context.Context,
	listRows []*spreadsheet_entity.ListRow,
) error {
	return l.spreadSheetRepository.WriteList(ctx, listRows)
}

func (l *listService) getFavoritesAndRivals(
	ctx context.Context,
	raceTickets []*ticket_csv_entity.RaceTicket,
	status types.PredictStatus,
	includeBetNumbers []types.BetNumber,
	excludeBetNumbers []types.BetNumber,
	totalPayments int,
) (
	[]types.BetNumber,
	[]types.BetNumber,
	types.PredictStatus,
) {
	var (
		favoriteBetNumbers []types.BetNumber
		rivalBetNumbers    []types.BetNumber
	)

	// 候補がすでにある場合、その馬番だけに絞って検索
	var refinedTickets []*ticket_csv_entity.Ticket

	if status.Included(types.FavoriteCandidate | types.RivalCandidate) {
		// 本命・対抗候補が複数いる場合
		for _, raceTicket := range raceTickets {
			if l.containsInSlices(includeBetNumbers, raceTicket.Ticket().BetNumber()) {
				refinedTickets = append(refinedTickets, raceTicket.Ticket())
			}
		}
	} else {
		for _, raceTicket := range raceTickets {
			refinedTickets = append(refinedTickets, raceTicket.Ticket())
		}
	}

	// 指定馬番の金額が最も大きいものを返す
	// 同金額で複数馬番の場合は複数返す

	// 本命が決定していない場合(未決定または本命候補がいる場合)
	if status.Matched(types.PredictUncompleted) || status.Included(types.FavoriteCandidate) {
		favoriteBetNumbers = l.getMaxBetNumbers(ctx, refinedTickets, nil)
		// 本命候補が1つの場合は本命を決定する
		// そうでない場合は候補を保持したまま次の処理へ
		if len(favoriteBetNumbers) == 1 {
			// 単勝の場合
			// 購入金額合計の(threshold)%未満だった場合、本命候補としないようにする
			// 単勝は優先度が高いので少額でも本命・対抗になってしまうことへの対応
			var favoriteTicket *ticket_csv_entity.Ticket
			for _, refinedTicket := range refinedTickets {
				if refinedTicket.TicketType() == types.Win && refinedTicket.BetNumber() == favoriteBetNumbers[0] {
					favoriteTicket = refinedTicket
					break
				}
			}
			if favoriteTicket != nil {
				n := float64(favoriteTicket.Payment()) / float64(totalPayments)
				// しきい値未満の場合は本命決定しない
				if n < thresholdOfLowerLimitPayment {
					// TODO 本来は無視ではなく本命計算割合に含めたいが、今は妥協する
					return favoriteBetNumbers, rivalBetNumbers, types.PredictUncompleted
				}
			}

			// 本命が1つに定まってる場合
			// favoriteBetNumbersの要素が1つの場合：
			// -> 本命が決定。対抗の検出のために本命馬番は除外する
			rivalBetNumbers = l.getMaxBetNumbers(ctx, refinedTickets, favoriteBetNumbers)
			if len(rivalBetNumbers) == 1 {
				// rivalBetNumbersの要素が１つの場合：
				// -> 対抗が決定
				return favoriteBetNumbers, rivalBetNumbers, types.FavoriteCompleted | types.RivalCompleted
			} else {
				// rivalBetNumbersの要素が複数の場合：
				// -> 対抗候補が複数のため次の券種の検査に移る
				return favoriteBetNumbers, rivalBetNumbers, types.FavoriteCompleted | types.RivalCandidate
			}
		} else {
			// 本命が複数候補いる場合：
			// -> 本命候補が複数のため次の券種の検査に移る
			return favoriteBetNumbers, nil, types.FavoriteCandidate
		}
	}

	// 対抗が決定していない場合(未決定または対抗候補がいる場合)
	if status.Included(types.RivalCandidate | types.FavoriteCompleted) {
		rivalBetNumbers = l.getMaxBetNumbers(ctx, refinedTickets, excludeBetNumbers)
		if len(rivalBetNumbers) == 1 {
			// 本命が1つに定まってる場合
			// favoriteBetNumbersの要素が1つの場合：
			// -> 対抗が決定
			return nil, rivalBetNumbers, types.FavoriteCompleted | types.RivalCompleted
		} else {
			// 対抗が複数候補いる場合：
			// -> 対抗候補が複数のため次の券種の検査に移る
			return nil, rivalBetNumbers, types.FavoriteCompleted | types.RivalCandidate
		}
	}

	// 現在の券種では本命・対抗、候補も決定できない
	return nil, nil, types.PredictUncompleted
}

func (l *listService) getMaxBetNumbers(
	ctx context.Context,
	tickets []*ticket_csv_entity.Ticket,
	excludeBetNumbers []types.BetNumber,
) []types.BetNumber {
	// 馬単、三連単の買い目計算ルール
	// 1着付けの馬番の金額を1.0倍、2着付けの馬番の金額を0.25倍で計算
	// それをソートして1番目、2番目...を算出
	// 馬連、三連複は順番を考慮しないので着順による重み付けはない
	result := make([]types.BetNumber, 0)
	if len(tickets) == 0 {
		return result
	}

	isExactaOrTrifecta := func(ticketType types.TicketType) bool {
		return ticketType == types.Exacta ||
			ticketType == types.ExactaWheelOfFirst ||
			ticketType == types.Trifecta ||
			ticketType == types.TrifectaFormation ||
			ticketType == types.TrifectaWheelOfFirst ||
			ticketType == types.TrifectaWheelOfSecond ||
			ticketType == types.TrifectaWheelOfFirstMulti ||
			ticketType == types.TrifectaWheelOfSecondMulti
	}

	betNumberPaymentMap := map[int]int{}
	for _, ticket := range tickets {
		nums := ticket.BetNumber().List()
		size := len(nums)
		weight := weightOfFirstPlace

		if size >= 1 && !l.containsInSlices(excludeBetNumbers, types.BetNumber(strconv.Itoa(nums[0]))) {
			// 1着付け
			if _, ok := betNumberPaymentMap[nums[0]]; !ok {
				betNumberPaymentMap[nums[0]] = ticket.Payment().Value()
			} else {
				betNumberPaymentMap[nums[0]] += ticket.Payment().Value()
			}
		}
		if size >= 2 && !l.containsInSlices(excludeBetNumbers, types.BetNumber(strconv.Itoa(nums[1]))) {
			// 2着付け
			if isExactaOrTrifecta(ticket.TicketType()) {
				weight = weightOfSecondPlace
			}
			if _, ok := betNumberPaymentMap[nums[1]]; !ok {
				betNumberPaymentMap[nums[1]] = int(float64(ticket.Payment()) * weight)
			} else {
				betNumberPaymentMap[nums[1]] += int(float64(ticket.Payment()) * weight)
			}
		}
		if size >= 3 && !l.containsInSlices(excludeBetNumbers, types.BetNumber(strconv.Itoa(nums[2]))) {
			// 3着付け
			if isExactaOrTrifecta(ticket.TicketType()) {
				weight = weightOfThirdPlace
			}
			if _, ok := betNumberPaymentMap[nums[2]]; !ok {
				betNumberPaymentMap[nums[2]] = int(float64(ticket.Payment()) * weight)
			} else {
				betNumberPaymentMap[nums[2]] += int(float64(ticket.Payment()) * weight)
			}
		}
	}

	// 高い順番に馬番をソートして保持
	numbers := make([]int, 0, len(betNumberPaymentMap))
	payments := make([]int, 0, len(betNumberPaymentMap))
	for number, payment := range betNumberPaymentMap {
		numbers = append(numbers, number)
		payments = append(payments, payment)
	}
	sort.Slice(payments, func(i, j int) bool {
		return payments[i] > payments[j]
	})

	seen := map[int]bool{}
	var uniquePayments []int
	for _, payment := range payments {
		if !seen[payment] {
			uniquePayments = append(uniquePayments, payment)
			seen[payment] = true
		}
	}

	// 最も高い金額を抽出。同金額のものも含めて1つに決定
	highPayments := make([]int, 1)
	if len(uniquePayments) >= 1 {
		highPayments = uniquePayments[:1]
	}

	// 最も高い金額に対する馬番を決定
	for _, highPayment := range highPayments {
		for number, payment := range betNumberPaymentMap {
			if payment == highPayment {
				result = append(result, types.NewBetNumber(fmt.Sprintf("%02d", number)))
			}
		}
	}

	return result
}

func (l *listService) getBetNumbersByMaxPayout(
	ctx context.Context,
	candidateBetNumbers []types.BetNumber,
	raceTickets []*ticket_csv_entity.RaceTicket,
) (
	[]types.BetNumber,
	[]types.BetNumber,
	bool,
) {
	var (
		maxKeys          []types.BetNumber
		otherKeys        []types.BetNumber
		isFoundCandidate bool
	)
	totalPaymentMap := map[types.BetNumber]int{}

	for _, raceTicket := range raceTickets {
		if raceTicket.Ticket().TicketResult() != types.TicketHit {
			continue
		}
		isFoundCandidate = true
		for _, candidateBetNumber := range candidateBetNumbers {
			rawCandidateBetNumber, _ := strconv.Atoi(candidateBetNumber.String())
			if slices.Contains(raceTicket.Ticket().BetNumber().List(), rawCandidateBetNumber) {
				totalPaymentMap[candidateBetNumber] += raceTicket.Ticket().Payment().Value()
			} else {
				totalPaymentMap[candidateBetNumber] += 0
			}
		}
	}

	if !isFoundCandidate {
		return nil, nil, false
	}

	maxValue := 0

	// 最大のvalueを見つける
	for _, v := range totalPaymentMap {
		if v > maxValue {
			maxValue = v
		}
	}

	// 最大のvalueを持つkeyを取得
	for key, value := range totalPaymentMap {
		if value == maxValue {
			maxKeys = append(maxKeys, key)
		} else {
			otherKeys = append(otherKeys, key)
		}
	}

	return maxKeys, otherKeys, true
}

func (l *listService) containsInSlices(betNumbers []types.BetNumber, betNumber types.BetNumber) bool {
	var (
		slice1 []int
		slice2 []int
	)

	for _, b := range betNumbers {
		slice1 = append(slice1, b.List()...)
	}
	slice2 = betNumber.List()

	elements := make(map[int]bool)
	for _, v := range slice2 {
		elements[v] = true
	}
	for _, v := range slice1 {
		if elements[v] {
			return true
		}
	}
	return false
}
