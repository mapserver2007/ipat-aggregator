package infrastructure

import "fmt"

func HitRateFormat(matchCount, raceCount int) string {
	if raceCount == 0 {
		return "-"
	}
	return fmt.Sprintf("%.2f%%", float64(matchCount)*100/float64(raceCount))
}
func PayoutRateFormat(payout float64, raceCount int) string {
	if raceCount == 0 {
		return "-"
	}
	return fmt.Sprintf("%.2f%%", payout*100/float64(raceCount))
}
