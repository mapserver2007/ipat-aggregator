package service

import (
	"fmt"
	betting_ticket_vo "github.com/mapserver2007/tools/baken/app/domain/betting_ticket/value_object"
	race_vo "github.com/mapserver2007/tools/baken/app/domain/race/value_object"
	"sort"
	"strconv"
	"strings"
)

func ConvertToRaceDate(v string) race_vo.RaceDate {
	return race_vo.NewRaceDate(v)
}

func ConvertToIntValue(v string) int {
	i, _ := strconv.Atoi(v)
	return i
}

func ConvertToPaymentForWheel(value string) []int {
	// 3連複軸1頭ながし, 3連単1着流し:
	// (1点あたりの購入金額)／(合計金額)
	separator := "／"
	values := strings.Split(value, separator)
	var payments []int
	for _, v := range values {
		payment, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		payments = append(payments, payment)
	}
	return payments
}

// ConvertToBetNumbersForQuinella ワイド流し変換
func ConvertToBetNumbersForQuinella(value string) []betting_ticket_vo.BetNumber {
	// 複数の買い目がまとめられてるものをバラす
	separator1 := "／"
	separator2 := "；"
	values1 := strings.Split(value, separator1)
	pivotalNumber, _ := strconv.Atoi(values1[0])                  // 軸
	strChallengerNumbers := strings.Split(values1[1], separator2) // 相手
	var betNumbers []betting_ticket_vo.BetNumber
	for _, strChallengerNumber := range strChallengerNumbers {
		// 馬番比較のために数値変換
		var betNumberStr string
		challengerNumber, _ := strconv.Atoi(strChallengerNumber)
		if challengerNumber > pivotalNumber {
			betNumberStr = fmt.Sprintf("%02d%s%02d", pivotalNumber, betting_ticket_vo.QuinellaSeparator, challengerNumber)
		} else {
			betNumberStr = fmt.Sprintf("%02d%s%02d", challengerNumber, betting_ticket_vo.QuinellaSeparator, pivotalNumber)
		}
		betNumbers = append(betNumbers, betting_ticket_vo.NewBetNumber(betNumberStr))
	}

	return betNumbers
}

// ConvertToBetNumbersForTrio 3連複流し変換
func ConvertToBetNumbersForTrio(value string) []betting_ticket_vo.BetNumber {
	// 複数の買い目がまとめられてるものをバラす
	separator1 := "／"
	separator2 := "；"
	values1 := strings.Split(value, separator1)
	pivotalNumber, _ := strconv.Atoi(values1[0])                  // 軸
	strChallengerNumbers := strings.Split(values1[1], separator2) // 相手
	var betNumbers []betting_ticket_vo.BetNumber

	for i := 0; i < len(strChallengerNumbers); i++ {
		for j := i + 1; j < len(strChallengerNumbers); j++ {
			challengerNumber1, _ := strconv.Atoi(strChallengerNumbers[i])
			challengerNumber2, _ := strconv.Atoi(strChallengerNumbers[j])

			numbers := []int{pivotalNumber, challengerNumber1, challengerNumber2}
			sort.Slice(numbers, func(k, l int) bool {
				return numbers[k] < numbers[l]
			})

			betNumberStr := fmt.Sprintf("%02d%s%02d%s%02d", numbers[0], betting_ticket_vo.QuinellaSeparator, numbers[1], betting_ticket_vo.QuinellaSeparator, numbers[2])
			betNumbers = append(betNumbers, betting_ticket_vo.NewBetNumber(betNumberStr))
		}
	}

	return betNumbers
}

// ConvertToPaymentForFoTrioFormation 3連複フォーメーション変換
func ConvertToPaymentForFoTrioFormation(value string) []betting_ticket_vo.BetNumber {
	// 複数の買い目がまとめられてるものをバラす
	separator1 := "／"
	separator2 := "；"
	values := strings.Split(value, separator1)

	values1 := strings.Split(values[0], separator2)
	values2 := strings.Split(values[1], separator2)
	values3 := strings.Split(values[2], separator2)

	betNumberMap := map[string]betting_ticket_vo.BetNumber{}
	var betNumbers []betting_ticket_vo.BetNumber

	for i := 0; i < len(values1); i++ {
		challengerNumber1, _ := strconv.Atoi(values1[i])
		for j := 0; j < len(values2); j++ {
			challengerNumber2, _ := strconv.Atoi(values2[j])

			if challengerNumber1 == challengerNumber2 {
				continue
			}

			for k := 0; k < len(values3); k++ {
				challengerNumber3, _ := strconv.Atoi(values3[k])

				if challengerNumber1 == challengerNumber3 || challengerNumber2 == challengerNumber3 {
					continue
				}

				numbers := []int{challengerNumber1, challengerNumber2, challengerNumber3}
				sort.Slice(numbers, func(k, l int) bool {
					return numbers[k] < numbers[l]
				})

				betNumberStr := fmt.Sprintf("%02d%s%02d%s%02d", numbers[0], betting_ticket_vo.QuinellaSeparator, numbers[1], betting_ticket_vo.QuinellaSeparator, numbers[2])
				if _, ok := betNumberMap[betNumberStr]; !ok {
					betNumberMap[betNumberStr] = betting_ticket_vo.NewBetNumber(betNumberStr)
				}
			}
		}
	}

	for _, betNumber := range betNumberMap {
		betNumbers = append(betNumbers, betNumber)
	}

	return betNumbers
}

// ConvertToBetNumbersForTrio 3連単流し変換
func ConvertToBetNumbersForExacta(value string) []betting_ticket_vo.BetNumber {
	// 複数の買い目がまとめられてるものをバラす
	separator1 := "／"
	separator2 := "；"
	values1 := strings.Split(value, separator1)
	pivotalNumber, _ := strconv.Atoi(values1[0])                  // 軸
	strChallengerNumbers := strings.Split(values1[1], separator2) // 相手
	var betNumbers []betting_ticket_vo.BetNumber

	for i := 0; i < len(strChallengerNumbers); i++ {
		for j := 0; j < len(strChallengerNumbers); j++ {
			if i == j {
				continue
			}
			challengerNumber1, _ := strconv.Atoi(strChallengerNumbers[i])
			challengerNumber2, _ := strconv.Atoi(strChallengerNumbers[j])
			betNumberStr := fmt.Sprintf("%02d%s%02d%s%02d", pivotalNumber, betting_ticket_vo.ExactaSeparator, challengerNumber1, betting_ticket_vo.ExactaSeparator, challengerNumber2)
			betNumbers = append(betNumbers, betting_ticket_vo.NewBetNumber(betNumberStr))
		}
	}

	return betNumbers
}
