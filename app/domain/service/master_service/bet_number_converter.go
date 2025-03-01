package master_service

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type BetNumberConverter interface {
	QuinellaWheelToQuinellaBetNumbers(ctx context.Context, rawBetNumber string) ([]types.BetNumber, error)
	QuinellaPlaceWheelToQuinellaBetNumbers(ctx context.Context, rawBetNumber string) ([]types.BetNumber, error)
	QuinellaPlaceFormationToQuinellaBetNumbers(ctx context.Context, rawBetNumber string) ([]types.BetNumber, error)
	ExactaWheelOfFirstToExactaBetNumbers(ctx context.Context, rawBetNumber string) ([]types.BetNumber, error)
	TrioFormationToTrioBetNumbers(ctx context.Context, rawBetNumber string) ([]types.BetNumber, error)
	TrioWheelOfFirstToTrioBetNumbers(ctx context.Context, rawBetNumber string) ([]types.BetNumber, error)
	TrioWheelOfSecondToTrioBetNumbers(ctx context.Context, rawBetNumber string) ([]types.BetNumber, error)
	TrioBoxToTrioBetNumbers(ctx context.Context, rawBetNumber string) ([]types.BetNumber, error)
	TrifectaFormationToTrifectaBetNumbers(ctx context.Context, rawBetNumber string) ([]types.BetNumber, error)
	TrifectaWheelOfFirstToTrifectaBetNumbers(ctx context.Context, rawBetNumber string) ([]types.BetNumber, error)
	TrifectaWheelOfSecondToTrifectaBetNumbers(ctx context.Context, rawBetNumber string) ([]types.BetNumber, error)
	TrifectaWheelMultiToTrifectaBetNumbers(ctx context.Context, rawBetNumber string) ([]types.BetNumber, error)
}

type betNumberConverter struct{}

func NewBetNumberConverter() BetNumberConverter {
	return &betNumberConverter{}
}

func (b *betNumberConverter) QuinellaWheelToQuinellaBetNumbers(ctx context.Context, rawBetNumber string) ([]types.BetNumber, error) {
	// 複数の買い目がまとめられてるものをバラす
	separator1 := "／"
	separator2 := "；"
	values1 := strings.Split(rawBetNumber, separator1)
	pivotalNumber, err := strconv.Atoi(values1[0]) // 軸
	if err != nil {
		return nil, err
	}
	strChallengerNumbers := strings.Split(values1[1], separator2) // 相手
	var betNumbers []types.BetNumber
	for _, strChallengerNumber := range strChallengerNumbers {
		// 馬番比較のために数値変換
		var betNumberStr string
		challengerNumber, err := strconv.Atoi(strChallengerNumber)
		if err != nil {
			return nil, err
		}
		if challengerNumber > pivotalNumber {
			betNumberStr = fmt.Sprintf("%02d%s%02d", pivotalNumber, types.QuinellaSeparator, challengerNumber)
		} else {
			betNumberStr = fmt.Sprintf("%02d%s%02d", challengerNumber, types.QuinellaSeparator, pivotalNumber)
		}
		betNumbers = append(betNumbers, types.BetNumber(betNumberStr))
	}

	return betNumbers, nil
}

// QuinellaPlaceWheelToQuinellaBetNumbers ワイドながし変換
func (b *betNumberConverter) QuinellaPlaceWheelToQuinellaBetNumbers(ctx context.Context, rawBetNumber string) ([]types.BetNumber, error) {
	// 複数の買い目がまとめられてるものをバラす
	separator1 := "／"
	separator2 := "；"
	values1 := strings.Split(rawBetNumber, separator1)
	pivotalNumber, err := strconv.Atoi(values1[0]) // 軸
	if err != nil {
		return nil, err
	}
	strChallengerNumbers := strings.Split(values1[1], separator2) // 相手
	var betNumbers []types.BetNumber
	for _, strChallengerNumber := range strChallengerNumbers {
		// 馬番比較のために数値変換
		var betNumberStr string
		challengerNumber, err := strconv.Atoi(strChallengerNumber)
		if err != nil {
			return nil, err
		}
		if challengerNumber > pivotalNumber {
			betNumberStr = fmt.Sprintf("%02d%s%02d", pivotalNumber, types.QuinellaSeparator, challengerNumber)
		} else {
			betNumberStr = fmt.Sprintf("%02d%s%02d", challengerNumber, types.QuinellaSeparator, pivotalNumber)
		}
		betNumbers = append(betNumbers, types.BetNumber(betNumberStr))
	}

	return betNumbers, nil
}

func (b *betNumberConverter) QuinellaPlaceFormationToQuinellaBetNumbers(ctx context.Context, rawBetNumber string) ([]types.BetNumber, error) {
	// 複数の買い目がまとめられてるものをバラす
	separator1 := "／"
	separator2 := "；"
	values := strings.Split(rawBetNumber, separator1)

	values1 := strings.Split(values[0], separator2)
	values2 := strings.Split(values[1], separator2)

	betNumberMap := map[string]string{}
	var betNumbers []types.BetNumber

	for i := 0; i < len(values1); i++ {
		challengerNumber1, err := strconv.Atoi(values1[i])
		if err != nil {
			return nil, err
		}
		for j := 0; j < len(values2); j++ {
			challengerNumber2, err := strconv.Atoi(values2[j])
			if err != nil {
				return nil, err
			}

			if challengerNumber1 == challengerNumber2 {
				continue
			}

			betNumberStr := fmt.Sprintf("%02d%s%02d", challengerNumber1, types.QuinellaSeparator, challengerNumber2)
			if _, ok := betNumberMap[betNumberStr]; !ok {
				betNumberMap[betNumberStr] = betNumberStr
			}
		}
	}

	for _, betNumber := range betNumberMap {
		betNumbers = append(betNumbers, types.BetNumber(betNumber))
	}

	return betNumbers, nil
}

// ExactaWheelOfFirstToExactaBetNumbers 馬単1着ながし変換
func (b *betNumberConverter) ExactaWheelOfFirstToExactaBetNumbers(ctx context.Context, rawBetNumber string) ([]types.BetNumber, error) {
	// 複数の買い目がまとめられてるものをバラす
	separator1 := "／"
	separator2 := "；"
	values1 := strings.Split(rawBetNumber, separator1)
	pivotalNumber, err := strconv.Atoi(values1[0]) // 軸
	if err != nil {
		return nil, err
	}
	strChallengerNumbers := strings.Split(values1[1], separator2) // 相手
	var betNumbers []types.BetNumber
	for _, strChallengerNumber := range strChallengerNumbers {
		var betNumberStr string
		challengerNumber, err := strconv.Atoi(strChallengerNumber)
		if err != nil {
			return nil, err
		}
		betNumberStr = fmt.Sprintf("%02d%s%02d", pivotalNumber, types.ExactaSeparator, challengerNumber)
		betNumbers = append(betNumbers, types.BetNumber(betNumberStr))
	}

	return betNumbers, nil
}

func (b *betNumberConverter) TrioFormationToTrioBetNumbers(ctx context.Context, rawBetNumber string) ([]types.BetNumber, error) {
	// 複数の買い目がまとめられてるものをバラす
	separator1 := "／"
	separator2 := "；"
	values := strings.Split(rawBetNumber, separator1)

	values1 := strings.Split(values[0], separator2)
	values2 := strings.Split(values[1], separator2)
	values3 := strings.Split(values[2], separator2)

	betNumberMap := map[string]string{}
	var betNumbers []types.BetNumber

	for i := 0; i < len(values1); i++ {
		challengerNumber1, err := strconv.Atoi(values1[i])
		if err != nil {
			return nil, err
		}
		for j := 0; j < len(values2); j++ {
			challengerNumber2, err := strconv.Atoi(values2[j])
			if err != nil {
				return nil, err
			}

			if challengerNumber1 == challengerNumber2 {
				continue
			}

			for k := 0; k < len(values3); k++ {
				challengerNumber3, err := strconv.Atoi(values3[k])
				if err != nil {
					return nil, err
				}

				if challengerNumber1 == challengerNumber3 || challengerNumber2 == challengerNumber3 {
					continue
				}

				numbers := []int{challengerNumber1, challengerNumber2, challengerNumber3}
				sort.Slice(numbers, func(k, l int) bool {
					return numbers[k] < numbers[l]
				})

				betNumberStr := fmt.Sprintf("%02d%s%02d%s%02d", numbers[0], types.QuinellaSeparator, numbers[1], types.QuinellaSeparator, numbers[2])
				if _, ok := betNumberMap[betNumberStr]; !ok {
					betNumberMap[betNumberStr] = betNumberStr
				}
			}
		}
	}

	for _, betNumber := range betNumberMap {
		betNumbers = append(betNumbers, types.BetNumber(betNumber))
	}

	return betNumbers, nil

}

// TrioWheelOfFirstToTrioBetNumbers 三連複1着ながし変換
func (b *betNumberConverter) TrioWheelOfFirstToTrioBetNumbers(ctx context.Context, rawBetNumber string) ([]types.BetNumber, error) {
	// 複数の買い目がまとめられてるものをバラす
	separator1 := "／"
	separator2 := "；"
	values1 := strings.Split(rawBetNumber, separator1)
	pivotalNumber, err := strconv.Atoi(values1[0]) // 軸
	if err != nil {
		return nil, err
	}
	strChallengerNumbers := strings.Split(values1[1], separator2) // 相手
	var betNumbers []types.BetNumber

	for i := 0; i < len(strChallengerNumbers); i++ {
		for j := i + 1; j < len(strChallengerNumbers); j++ {
			challengerNumber1, err := strconv.Atoi(strChallengerNumbers[i])
			if err != nil {
				return nil, err
			}
			challengerNumber2, err := strconv.Atoi(strChallengerNumbers[j])
			if err != nil {
				return nil, err
			}

			numbers := []int{pivotalNumber, challengerNumber1, challengerNumber2}
			sort.Slice(numbers, func(k, l int) bool {
				return numbers[k] < numbers[l]
			})

			betNumberStr := fmt.Sprintf("%02d%s%02d%s%02d", numbers[0], types.QuinellaSeparator, numbers[1], types.QuinellaSeparator, numbers[2])
			betNumbers = append(betNumbers, types.BetNumber(betNumberStr))
		}
	}

	return betNumbers, nil
}

// TrioWheelOfSecondToTrioBetNumbers 三連複2着ながし変換
func (b *betNumberConverter) TrioWheelOfSecondToTrioBetNumbers(ctx context.Context, rawBetNumber string) ([]types.BetNumber, error) {
	// 複数の買い目がまとめられてるものをバラす
	separator1 := "／"
	separator2 := "；"
	values1 := strings.Split(rawBetNumber, separator1)
	strPivotalNumbers := strings.Split(values1[0], separator2)    // 軸
	strChallengerNumbers := strings.Split(values1[1], separator2) // 相手
	var betNumbers []types.BetNumber

	strPivotalNumber1, err := strconv.Atoi(strPivotalNumbers[0])
	if err != nil {
		return nil, err
	}
	strPivotalNumber2, err := strconv.Atoi(strPivotalNumbers[1])
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(strChallengerNumbers); i++ {
		challengerNumber, err := strconv.Atoi(strChallengerNumbers[i])
		if err != nil {
			return nil, err
		}
		numbers := []int{strPivotalNumber1, strPivotalNumber2, challengerNumber}
		sort.Slice(numbers, func(k, l int) bool {
			return numbers[k] < numbers[l]
		})

		betNumberStr := fmt.Sprintf("%02d%s%02d%s%02d", numbers[0], types.QuinellaSeparator, numbers[1], types.QuinellaSeparator, numbers[2])
		betNumbers = append(betNumbers, types.BetNumber(betNumberStr))
	}

	return betNumbers, nil
}

// TrioBoxToTrioBetNumbers 三連複ボックス変換
func (b *betNumberConverter) TrioBoxToTrioBetNumbers(ctx context.Context, rawBetNumber string) ([]types.BetNumber, error) {
	// 04；07；11
	separator := "；"
	values := strings.Split(rawBetNumber, separator)
	betNumberMap := map[string]string{}
	var betNumbers []types.BetNumber

	for i := 0; i < len(values); i++ {
		challengerNumber1, err := strconv.Atoi(values[i])
		if err != nil {
			return nil, err
		}
		for j := i + 1; j < len(values); j++ {
			challengerNumber2, err := strconv.Atoi(values[j])
			if err != nil {
				return nil, err
			}

			if challengerNumber1 == challengerNumber2 {
				continue
			}

			for k := j + 1; k < len(values); k++ {
				challengerNumber3, err := strconv.Atoi(values[k])
				if err != nil {
					return nil, err
				}

				if challengerNumber1 == challengerNumber3 || challengerNumber2 == challengerNumber3 {
					continue
				}

				betNumberStr := fmt.Sprintf("%02d%s%02d%s%02d", challengerNumber1, types.QuinellaSeparator, challengerNumber2, types.QuinellaSeparator, challengerNumber3)
				if _, ok := betNumberMap[betNumberStr]; !ok {
					betNumberMap[betNumberStr] = betNumberStr
				}
			}
		}
	}

	for _, betNumber := range betNumberMap {
		betNumbers = append(betNumbers, types.BetNumber(betNumber))
	}

	return betNumbers, nil
}

// TrifectaFormationToTrifectaBetNumbers 三連単フォーメーション変換
func (b *betNumberConverter) TrifectaFormationToTrifectaBetNumbers(ctx context.Context, rawBetNumber string) ([]types.BetNumber, error) {
	// 複数の買い目がまとめられてるものをバラす
	separator1 := "／"
	separator2 := "；"
	values := strings.Split(rawBetNumber, separator1)

	values1 := strings.Split(values[0], separator2)
	values2 := strings.Split(values[1], separator2)
	values3 := strings.Split(values[2], separator2)

	betNumberMap := map[string]string{}
	var betNumbers []types.BetNumber

	for i := 0; i < len(values1); i++ {
		challengerNumber1, err := strconv.Atoi(values1[i])
		if err != nil {
			return nil, err
		}
		for j := 0; j < len(values2); j++ {
			challengerNumber2, err := strconv.Atoi(values2[j])
			if err != nil {
				return nil, err
			}

			if challengerNumber1 == challengerNumber2 {
				continue
			}

			for k := 0; k < len(values3); k++ {
				challengerNumber3, err := strconv.Atoi(values3[k])
				if err != nil {
					return nil, err
				}

				if challengerNumber1 == challengerNumber3 || challengerNumber2 == challengerNumber3 {
					continue
				}

				betNumberStr := fmt.Sprintf("%02d%s%02d%s%02d", challengerNumber1, types.ExactaSeparator, challengerNumber2, types.ExactaSeparator, challengerNumber3)
				if _, ok := betNumberMap[betNumberStr]; !ok {
					betNumberMap[betNumberStr] = betNumberStr
				}
			}
		}
	}

	for _, betNumber := range betNumberMap {
		betNumbers = append(betNumbers, types.BetNumber(betNumber))
	}

	return betNumbers, nil
}

// TrifectaWheelOfFirstToTrifectaBetNumbers 3連単1着ながし変換
func (b *betNumberConverter) TrifectaWheelOfFirstToTrifectaBetNumbers(ctx context.Context, rawBetNumber string) ([]types.BetNumber, error) {
	// 複数の買い目がまとめられてるものをバラす
	separator1 := "／"
	separator2 := "；"
	values1 := strings.Split(rawBetNumber, separator1)
	pivotalNumber, err := strconv.Atoi(values1[0]) // 軸
	if err != nil {
		return nil, err
	}
	strChallengerNumbers := strings.Split(values1[1], separator2) // 相手
	var betNumbers []types.BetNumber

	for i := 0; i < len(strChallengerNumbers); i++ {
		for j := 0; j < len(strChallengerNumbers); j++ {
			if i == j {
				continue
			}
			challengerNumber1, err := strconv.Atoi(strChallengerNumbers[i])
			if err != nil {
				return nil, err
			}
			challengerNumber2, err := strconv.Atoi(strChallengerNumbers[j])
			if err != nil {
				return nil, err
			}
			betNumberStr := fmt.Sprintf("%02d%s%02d%s%02d", pivotalNumber, types.ExactaSeparator, challengerNumber1, types.ExactaSeparator, challengerNumber2)
			betNumbers = append(betNumbers, types.BetNumber(betNumberStr))
		}
	}

	return betNumbers, nil
}

// TrifectaWheelOfSecondToTrifectaBetNumbers 3連単2着ながし変換
func (b *betNumberConverter) TrifectaWheelOfSecondToTrifectaBetNumbers(ctx context.Context, rawBetNumber string) ([]types.BetNumber, error) {
	// 複数の買い目がまとめられてるものをバラす
	separator1 := "／"
	separator2 := "；"
	values1 := strings.Split(rawBetNumber, separator1)
	pivotalNumber, err := strconv.Atoi(values1[0]) // 軸
	if err != nil {
		return nil, err
	}
	strChallengerNumbers := strings.Split(values1[1], separator2) // 相手
	var betNumbers []types.BetNumber

	for i := 0; i < len(strChallengerNumbers); i++ {
		for j := 0; j < len(strChallengerNumbers); j++ {
			if i == j {
				continue
			}
			challengerNumber1, err := strconv.Atoi(strChallengerNumbers[i])
			if err != nil {
				return nil, err
			}
			challengerNumber2, err := strconv.Atoi(strChallengerNumbers[j])
			if err != nil {
				return nil, err
			}
			betNumberStr := fmt.Sprintf("%02d%s%02d%s%02d", challengerNumber1, types.ExactaSeparator, pivotalNumber, types.ExactaSeparator, challengerNumber2)
			betNumbers = append(betNumbers, types.BetNumber(betNumberStr))
		}
	}

	return betNumbers, nil
}

// TrifectaWheelMultiToTrifectaBetNumbers 3連単軸1,2頭ながしマルチ変換
func (b *betNumberConverter) TrifectaWheelMultiToTrifectaBetNumbers(ctx context.Context, rawBetNumber string) ([]types.BetNumber, error) {
	// 複数の買い目がまとめられてるものをバラす
	separator1 := "／"
	separator2 := "；"
	values1 := strings.Split(rawBetNumber, separator1)
	strPivotalNumbers := strings.Split(values1[0], separator2)    // 軸
	strChallengerNumbers := strings.Split(values1[1], separator2) // 相手
	var betNumbers []types.BetNumber

	if len(strPivotalNumbers) == 1 {
		// 1頭軸マルチ
		var combinations [][3]int
		for i := 0; i < len(strChallengerNumbers); i++ {
			for j := 0; j < len(strChallengerNumbers); j++ {
				if i == j {
					continue
				}
				challengerNumber1, err := strconv.Atoi(strChallengerNumbers[i])
				if err != nil {
					return nil, err
				}
				challengerNumber2, err := strconv.Atoi(strChallengerNumbers[j])
				if err != nil {
					return nil, err
				}
				pivotalNumber, err := strconv.Atoi(strPivotalNumbers[0])
				if err != nil {
					return nil, err
				}

				combinations = append(combinations, [3]int{pivotalNumber, challengerNumber1, challengerNumber2})
				combinations = append(combinations, [3]int{challengerNumber1, pivotalNumber, challengerNumber2})
				combinations = append(combinations, [3]int{challengerNumber2, challengerNumber1, pivotalNumber})
			}
		}
		for _, combination := range combinations {
			betNumbers = append(betNumbers, types.BetNumber(fmt.Sprintf("%02d%s%02d%s%02d",
				combination[0],
				types.ExactaSeparator,
				combination[1],
				types.ExactaSeparator,
				combination[2])))
		}
	} else if len(strPivotalNumbers) == 2 {
		// 2頭軸マルチ
		var combinations [][3]int
		pivotalNumber1, err := strconv.Atoi(strPivotalNumbers[0])
		if err != nil {
			return nil, err
		}
		pivotalNumber2, err := strconv.Atoi(strPivotalNumbers[1])
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(strChallengerNumbers); i++ {
			challengerNumber, err := strconv.Atoi(strChallengerNumbers[i])
			if err != nil {
				return nil, err
			}
			combinations = append(combinations, [3]int{pivotalNumber1, pivotalNumber2, challengerNumber})
			combinations = append(combinations, [3]int{pivotalNumber1, challengerNumber, pivotalNumber2})
			combinations = append(combinations, [3]int{pivotalNumber2, pivotalNumber1, challengerNumber})
			combinations = append(combinations, [3]int{pivotalNumber2, challengerNumber, pivotalNumber1})
			combinations = append(combinations, [3]int{challengerNumber, pivotalNumber1, pivotalNumber2})
			combinations = append(combinations, [3]int{challengerNumber, pivotalNumber2, pivotalNumber1})
		}
		for _, combination := range combinations {
			betNumbers = append(betNumbers, types.BetNumber(fmt.Sprintf("%02d%s%02d%s%02d",
				combination[0],
				types.ExactaSeparator,
				combination[1],
				types.ExactaSeparator,
				combination[2])))
		}
	} else {
		return nil, fmt.Errorf("no support pivotal number by 3 or more")
	}

	return betNumbers, nil
}
