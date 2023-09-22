package value_object

import (
	"strconv"
	"strings"
)

type WinOdds float64

// https://saisyu-corner.com/odds_ninki_tokei_00
var oddsWinRateMap = map[string]float64{
	"1.1":     0.786,
	"1.2":     0.646,
	"1.3":     0.662,
	"1.4":     0.603,
	"1.5":     0.520,
	"1.6":     0.509,
	"1.7":     0.448,
	"1.8":     0.435,
	"1.9":     0.392,
	"2.0-2.1": 0.357,
	"2.2-2.3": 0.330,
	"2.4-2.5": 0.324,
	"2.6-2.7": 0.300,
	"2.8-2.9": 0.276,
	"3.0-3.4": 0.247,
	"3.5-3.9": 0.224,
	"4.0-4.9": 0.176,
	"5.0-6.9": 0.147,
	"7.0-9.9": 0.100,
	"10-14.9": 0.070,
	"15-19.9": 0.049,
	"20-29.9": 0.034,
	"30-49.9": 0.021,
	"50-99.9": 0.011,
	"100-":    0.003,
}

//var oddsWinPayoutMap = map[string]float64{
//	"1.7":     0.76,
//	"1.8":     0.78,
//	"1.9":     0.75,
//	"2.0-2.1": 0.73,
//	"2.2-2.3": 0.74,
//	"2.4-2.5": 0.79,
//	"2.6-2.7": 0.80,
//	"2.8-2.9": 0.79,
//	"3.0-3.4": 0.79,
//	"3.5-3.9": 0.83,
//	"4.0-4.9": 0.78,
//	"5.0-6.9": 0.79,
//	"7.0-9.9": 0.83,
//	"10-14.9": 0.86,
//	"15-19.9": 0.84,
//	"20-29.9": 0.83,
//	"30-49.9": 0.80,
//	"50-99.9": 0.77,
//	"100-":    0.47,
//}

func (w WinOdds) Value() float64 {
	return float64(w)
}

func (w WinOdds) WinRate() float64 {
	for k, v := range oddsWinRateMap {
		var n1, n2 float64
		nums := strings.Split(k, "-")
		if len(nums) == 1 {
			n1, _ = strconv.ParseFloat(nums[0], 64)
			n2 = n1
		} else {
			n1, _ = strconv.ParseFloat(nums[0], 64)
			n2, _ = strconv.ParseFloat(nums[1], 64)
			if n2 == 0 {
				n2 = 999
			}
		}
		if n1 <= float64(w) && float64(w) <= n2 {
			return v
		}
	}
	return 0
}

func (w WinOdds) OddsRange() string {
	for k := range oddsWinRateMap {
		var n1, n2 float64
		nums := strings.Split(k, "-")
		if len(nums) == 1 {
			n1, _ = strconv.ParseFloat(nums[0], 64)
			n2 = n1
		} else {
			n1, _ = strconv.ParseFloat(nums[0], 64)
			n2, _ = strconv.ParseFloat(nums[1], 64)
			if n2 == 0 {
				n2 = 999
			}
		}
		if n1 <= float64(w) && float64(w) <= n2 {
			return k
		}
	}
	return ""
}
