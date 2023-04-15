package service

import (
	"fmt"
	betting_ticket_entity "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/entity"
	betting_ticket_vo "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/value_object"
	jockey_entity "github.com/mapserver2007/ipat-aggregator/app/domain/jockey/entity"
	predict_entity "github.com/mapserver2007/ipat-aggregator/app/domain/predict/entity"
	predict_vo "github.com/mapserver2007/ipat-aggregator/app/domain/predict/value_object"
	race_entity "github.com/mapserver2007/ipat-aggregator/app/domain/race/entity"
	race_vo "github.com/mapserver2007/ipat-aggregator/app/domain/race/value_object"
	"sort"
	"strconv"
)

const (
	WeightOfQuinella             = 1.0
	WeightOfExactaSecond         = 0.25
	WeightOfExactaThird          = 0.1
	ThresholdOfLowerLimitPayment = 0.15
)

type Predictor struct {
	raceConverter          RaceConverter
	bettingTicketConverter BettingTicketConverter
}

func NewPredictor(
	raceConverter RaceConverter,
	bettingTicketConverter BettingTicketConverter,
) *Predictor {
	return &Predictor{
		raceConverter:          raceConverter,
		bettingTicketConverter: bettingTicketConverter,
	}
}

func (p *Predictor) Predict(
	records []*betting_ticket_entity.CsvEntity,
	racingNumbers []*race_entity.RacingNumber,
	races []*race_entity.Race,
	jockeys []*jockey_entity.Jockey,
) ([]*predict_entity.PredictEntity, error) {
	raceMap := p.raceConverter.ConvertToRaceMapByRaceId(races)
	recordMap := p.getRecordMapByRaceId(records, racingNumbers)
	var entities []*predict_entity.PredictEntity

	for raceId, bettingTicketDetails := range recordMap {
		raceInfo, ok := raceMap[raceId]
		if !ok {
			return nil, fmt.Errorf("unknown raceId: %s", raceId)
		}

		var payment, repayment int
		for _, bettingTicketDetail := range bettingTicketDetails {
			payment += bettingTicketDetail.Payment()
			repayment += bettingTicketDetail.Repayment()
		}

		bettingTicketMap := p.bettingTicketConverter.ConvertToBettingTicketMap(bettingTicketDetails)
		payoutResultMap := p.bettingTicketConverter.ConvertToPayoutResultMap(raceInfo.PayoutResults())

		var (
			favorites, rivals []betting_ticket_vo.BetNumber
			favorite, rival   *betting_ticket_vo.BetNumber
			winningTickets    []*predict_entity.WinningTicketEntity
		)
		status := predict_vo.PredictUncompleted
		sortedBettingTickets := getSortedBettingTickets()

		for i := 0; i < len(sortedBettingTickets); i++ {
			if status.Matched(predict_vo.FavoriteCompleted | predict_vo.RivalCompleted) {
				break
			}

			bettingTicketValue := sortedBettingTickets[i]
			details, ok := bettingTicketMap[bettingTicketValue]
			if !ok || len(details) == 0 {
				continue
			}

			// 本命対抗確定済み
			if status.Matched(predict_vo.FavoriteCompleted | predict_vo.RivalCompleted) {
				break
			}

			// 本命候補がいる場合
			if status.Included(predict_vo.FavoriteCandidate) {
				// さらに本命が絞れない場合、次の券種
				favorites, rivals, status = p.getFavoritesAndRivals(details, predict_vo.FavoriteCandidate, favorites, nil, payment)
			}

			// 対抗候補がいる場合
			if status.Included(predict_vo.RivalCandidate) {
				_, rivals, status = p.getFavoritesAndRivals(details, predict_vo.RivalCandidate, rivals, nil, payment)
			}

			// 本命が決定済み、対抗が未決定
			if status.Included(predict_vo.FavoriteCompleted) {
				_, rivals, status = p.getFavoritesAndRivals(details, predict_vo.FavoriteCompleted, nil, favorites, payment)
			}

			// 本命、対抗が未決定
			if status.Matched(predict_vo.PredictUncompleted) {
				favorites, rivals, status = p.getFavoritesAndRivals(details, predict_vo.PredictUncompleted, nil, nil, payment)
			}
		}

		// 本命が複数の場合
		if status.Matched(predict_vo.FavoriteCandidate) {
			if len(favorites) >= 2 {
				// 人気順ソート
				sort.Slice(raceInfo.RaceResults(), func(i, j int) bool {
					return raceInfo.RaceResults()[i].PopularNumber() < raceInfo.RaceResults()[j].PopularNumber()
				})

				// 対抗が複数存在する場合、人気が高いもの本命にして残りを対抗候補にする
				for _, raceResult := range raceInfo.RaceResults() {
					betNumber := betting_ticket_vo.NewBetNumber(fmt.Sprintf("%02d", raceResult.HorseNumber()))
					if containsInSlices(favorites, betNumber) {
						for _, rivalCandidate := range favorites {
							if rivalCandidate != betNumber {
								rivals = append(rivals, rivalCandidate)
							}
						}
						favorites = []betting_ticket_vo.BetNumber{betNumber}

						// 対抗が1つに決定している場合
						if len(rivals) == 1 {
							status = predict_vo.FavoriteCompleted | predict_vo.RivalCompleted
						} else {
							status = predict_vo.FavoriteCompleted | predict_vo.RivalCandidate
						}
						break
					}
				}
			}
		}

		// 対抗が候補の場合
		if status.Matched(predict_vo.FavoriteCompleted | predict_vo.RivalCandidate) {
			if len(rivals) == 0 {
				// 対抗が存在しない場合
				status = predict_vo.FavoriteCompleted | predict_vo.RivalCompleted
			} else if len(rivals) >= 2 {
				// 人気順ソート
				sort.Slice(raceInfo.RaceResults(), func(i, j int) bool {
					return raceInfo.RaceResults()[i].PopularNumber() < raceInfo.RaceResults()[j].PopularNumber()
				})

				// 対抗が複数存在する場合、人気が高いものを採用する
				for _, raceResult := range raceInfo.RaceResults() {
					betNumber := betting_ticket_vo.NewBetNumber(fmt.Sprintf("%02d", raceResult.HorseNumber()))
					if containsInSlices(rivals, betNumber) {
						rivals = []betting_ticket_vo.BetNumber{betNumber}
						break
					}
				}
			}
		}

		if favorites != nil && len(favorites) > 0 {
			favorite = &favorites[0]
		}
		if rivals != nil && len(rivals) > 0 {
			rival = &rivals[0]
		}

		for i := 0; i < len(sortedBettingTickets); i++ {
			// 本命対抗検出ロジック内だと最後の券種に到達する前にbreakするケースがあるためもう一度ループを回す
			bettingTicketValue := sortedBettingTickets[i]
			details, ok := bettingTicketMap[bettingTicketValue]
			if !ok || len(details) == 0 {
				continue
			}

			for _, detail := range details {
				payoutResult, ok := payoutResultMap[detail.BettingTicket().ConvertToOriginBettingTicket()]
				if !ok {
					return nil, fmt.Errorf("unknown payout result in ticketType %d", detail.BettingTicket().Value())
				}
				if detail.Winning() {
					var odds string
					for idx, betNumber := range payoutResult.Numbers() {
						if betNumber == detail.BetNumber().String() {
							odds = payoutResult.Odds()[idx]
						}
					}
					winningTickets = append(winningTickets, predict_entity.NewWinningTicketEntity(
						detail.BettingTicket(),
						detail.BetNumber(),
						odds,
						detail.Repayment(),
					))
				}
			}
		}

		var (
			favoriteHorse, rivalHorse   *predict_entity.Horse
			favoriteJockey, rivalJockey *predict_entity.Jockey
		)
		for _, raceResult := range raceInfo.RaceResults() {
			if favorite != nil && raceResult.HorseNumber() == favorite.List()[0] {
				favoriteHorse = predict_entity.NewHorse(raceResult.HorseName(), raceResult.Odds(), raceResult.PopularNumber())
				favoriteJockey = predict_entity.NewJockey(raceResult.JockeyName())
			}
			if rival != nil && raceResult.HorseNumber() == rival.List()[0] {
				rivalHorse = predict_entity.NewHorse(raceResult.HorseName(), raceResult.Odds(), raceResult.PopularNumber())
				rivalJockey = predict_entity.NewJockey(raceResult.JockeyName())
			}
		}

		sort.Slice(raceInfo.RaceResults(), func(i, j int) bool {
			return raceInfo.RaceResults()[i].OrderNo() < raceInfo.RaceResults()[j].OrderNo()
		})

		race := predict_entity.NewRace(
			raceId,
			raceInfo.RaceNumber(),
			raceInfo.RaceName(),
			raceInfo.StartTime(),
			race_vo.GradeClass(raceInfo.Class()),
			raceInfo.RaceCourseId(),
			raceInfo.CourseCategory(),
			raceInfo.RaceDate(),
			raceInfo.Distance(),
			raceInfo.TrackCondition(),
			payment,
			repayment,
			raceInfo.Url(),
			raceInfo.RaceResults()[0:2],
		)
		entities = append(entities, predict_entity.NewPredictEntity(
			race, favoriteHorse, rivalHorse, favoriteJockey, rivalJockey, payment, repayment, winningTickets, status))

		if len(favorites) >= 2 || len(rivals) >= 2 {
			return nil, fmt.Errorf("failed to find favorite or rival")
		}
	}

	return entities, nil
}

func (p *Predictor) getFavoritesAndRivals(
	allDetails []*betting_ticket_entity.BettingTicketDetail,
	currentStatus predict_vo.PredictStatus,
	includeBetNumbers []betting_ticket_vo.BetNumber,
	excludeBetNumbers []betting_ticket_vo.BetNumber,
	totalPayments int,
) (
	[]betting_ticket_vo.BetNumber,
	[]betting_ticket_vo.BetNumber,
	predict_vo.PredictStatus,
) {
	var (
		favoriteBetNumbers []betting_ticket_vo.BetNumber
		rivalBetNumbers    []betting_ticket_vo.BetNumber
	)

	// 候補がすでにある場合、その馬番だけに絞って検索
	var refinedDetails []*betting_ticket_entity.BettingTicketDetail

	if currentStatus.Included(predict_vo.FavoriteCandidate | predict_vo.RivalCandidate) {
		// 本命・対抗候補が複数いる場合
		for _, detail := range allDetails {
			if containsInSlices(includeBetNumbers, detail.BetNumber()) {
				refinedDetails = append(refinedDetails, detail)
			}
		}
	} else {
		refinedDetails = allDetails
	}

	// 指定馬番の金額が最も大きいものを返す
	// 同金額で複数馬番の場合は複数返す
	if currentStatus.Matched(predict_vo.PredictUncompleted) || currentStatus.Included(predict_vo.FavoriteCandidate) {
		// 本命が決定していない場合
		favoriteBetNumbers = getMaxBetNumbers(refinedDetails, nil)

		if len(favoriteBetNumbers) == 1 {
			// 単勝の場合
			// 購入金額合計の(threshold)%未満だった場合、本命候補としないようにする
			// 単勝は優先度が高いので少額でも本命・対抗になってしまうことへの対応
			detail := getDetailByNumberForWin(favoriteBetNumbers[0], refinedDetails)
			if detail != nil {
				n := float64(detail.Payment()) / float64(totalPayments)
				if n < ThresholdOfLowerLimitPayment {
					// TODO 本来は無視ではなく本命計算割合に含めたいが、今は妥協する
					return favoriteBetNumbers, rivalBetNumbers, predict_vo.PredictUncompleted
				}
			}

			// 本命が1つに定まってる場合
			// favoriteBetNumbersの要素が1つの場合：
			// -> 本命が決定。対抗の検出のために本命馬番は除外する
			rivalBetNumbers = getMaxBetNumbers(refinedDetails, favoriteBetNumbers)
			if len(rivalBetNumbers) == 1 {
				// rivalBetNumbersの要素が１つの場合：
				// -> 対抗が決定
				return favoriteBetNumbers, rivalBetNumbers, predict_vo.FavoriteCompleted | predict_vo.RivalCompleted
			} else {
				// rivalBetNumbersの要素が複数の場合：
				// -> 対抗候補が複数のため次の券種の検査に移る
				return favoriteBetNumbers, rivalBetNumbers, predict_vo.FavoriteCompleted | predict_vo.RivalCandidate
			}
		} else {
			// 本命が複数候補いる場合：
			// -> 本命候補が複数のため次の券種の検査に移る
			return favoriteBetNumbers, nil, predict_vo.FavoriteCandidate
		}
	}

	if currentStatus.Included(predict_vo.RivalCandidate | predict_vo.FavoriteCompleted) {
		// 対抗が決定していない場合
		rivalBetNumbers = getMaxBetNumbers(refinedDetails, excludeBetNumbers)

		if len(rivalBetNumbers) == 1 {
			// 本命が1つに定まってる場合
			// favoriteBetNumbersの要素が1つの場合：
			// -> 対抗が決定
			return nil, rivalBetNumbers, predict_vo.FavoriteCompleted | predict_vo.RivalCompleted
		} else {
			// 対抗が複数候補いる場合：
			// -> 対抗候補が複数のため次の券種の検査に移る
			return nil, rivalBetNumbers, predict_vo.FavoriteCompleted | predict_vo.RivalCandidate
		}
	}

	// 現在の券種では本命・対抗、候補も決定できない
	return nil, nil, predict_vo.PredictUncompleted
}

func getMaxBetNumbers(details []*betting_ticket_entity.BettingTicketDetail, excludeBetNumbers []betting_ticket_vo.BetNumber) []betting_ticket_vo.BetNumber {
	// 馬単の買い目計算ルール
	// 1着付けの馬番の金額を1.0倍、2着付けの馬番の金額を0.25倍で計算
	// それをソートして1番目、2番目...を算出
	result := make([]betting_ticket_vo.BetNumber, 0)
	if len(details) == 0 {
		return result
	}

	isExacta := func(ticketType betting_ticket_vo.BettingTicket) bool {
		if ticketType == betting_ticket_vo.Exacta ||
			ticketType == betting_ticket_vo.Trifecta ||
			ticketType == betting_ticket_vo.TrifectaFormation ||
			ticketType == betting_ticket_vo.TrifectaWheelOfFirst {
			return true
		}

		return false
	}

	betNumberPaymentMap := map[int]int{}
	for _, detail := range details {
		nums := detail.BetNumber().List()
		size := len(detail.BetNumber().List())
		weight := WeightOfQuinella

		if size >= 1 && !containsInSlices(excludeBetNumbers, betting_ticket_vo.NewBetNumber(strconv.Itoa(nums[0]))) {
			// 1着付け
			if _, ok := betNumberPaymentMap[nums[0]]; !ok {
				betNumberPaymentMap[nums[0]] = detail.Payment()
			} else {
				betNumberPaymentMap[nums[0]] += detail.Payment()
			}
		}
		if size >= 2 && !containsInSlices(excludeBetNumbers, betting_ticket_vo.NewBetNumber(strconv.Itoa(nums[1]))) {
			// 2着付け
			if isExacta(detail.BettingTicket()) {
				weight = WeightOfExactaSecond
			}
			if _, ok := betNumberPaymentMap[nums[1]]; !ok {
				betNumberPaymentMap[nums[1]] = int(float64(detail.Payment()) * weight)
			} else {
				betNumberPaymentMap[nums[1]] += int(float64(detail.Payment()) * weight)
			}
		}
		if size >= 3 && !containsInSlices(excludeBetNumbers, betting_ticket_vo.NewBetNumber(strconv.Itoa(nums[2]))) {
			// 3着付け
			if isExacta(detail.BettingTicket()) {
				weight = WeightOfExactaThird
			}
			if _, ok := betNumberPaymentMap[nums[2]]; !ok {
				betNumberPaymentMap[nums[2]] = int(float64(detail.Payment()) * weight)
			} else {
				betNumberPaymentMap[nums[2]] += int(float64(detail.Payment()) * weight)
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
				result = append(result, betting_ticket_vo.NewBetNumber(fmt.Sprintf("%02d", number)))
			}
		}
	}

	return result
}

func getDetailByNumberForWin(number betting_ticket_vo.BetNumber, details []*betting_ticket_entity.BettingTicketDetail) *betting_ticket_entity.BettingTicketDetail {
	for _, detail := range details {
		if detail.BettingTicket() == betting_ticket_vo.Win && detail.BetNumber() == number {
			return detail
		}
	}

	return nil
}

func (p *Predictor) getRecordMapByRaceId(records []*betting_ticket_entity.CsvEntity, racingNumbers []*race_entity.RacingNumber) map[race_vo.RaceId][]*betting_ticket_entity.BettingTicketDetail {
	recordMap := map[race_vo.RaceId][]*betting_ticket_entity.BettingTicketDetail{}
	racingNumberMap := p.raceConverter.ConvertToRacingNumberMap(racingNumbers)

	for _, record := range records {
		key := race_vo.NewRacingNumberId(record.RaceDate(), record.RaceCourse())
		racingNumber, _ := racingNumberMap[key]
		raceId := p.raceConverter.GetRaceId(record, racingNumber)
		bettingTicketDetail := betting_ticket_entity.NewBettingTicketDetail(
			record.BettingTicket(),
			record.BetNumber(),
			record.Payment(),
			record.Repayment(),
			record.Winning(),
		)
		recordMap[*raceId] = append(recordMap[*raceId], bettingTicketDetail)
	}

	return recordMap
}

func getSortedBettingTickets() []betting_ticket_vo.BettingTicket {
	// 計算の優先順
	return []betting_ticket_vo.BettingTicket{
		betting_ticket_vo.Win,
		betting_ticket_vo.Exacta,
		betting_ticket_vo.Trifecta,
		betting_ticket_vo.TrifectaWheelOfFirst,
		betting_ticket_vo.TrifectaFormation,
		betting_ticket_vo.QuinellaPlaceWheel,
		betting_ticket_vo.Quinella,
		betting_ticket_vo.QuinellaPlace,
		betting_ticket_vo.TrioWheelOfFirst,
		betting_ticket_vo.Trio,
		betting_ticket_vo.TrioFormation,
		betting_ticket_vo.Place,
		betting_ticket_vo.BracketQuinella,
	}
}

func containsInSlices(betNumbers []betting_ticket_vo.BetNumber, betNumber betting_ticket_vo.BetNumber) bool {
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
