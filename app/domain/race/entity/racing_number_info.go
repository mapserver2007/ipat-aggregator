package entity

type RacingNumberInfo struct {
	racingNumbers []*RacingNumber
}

func NewRacingNumberInfo(
	racingNumbers []*RacingNumber,
) *RacingNumberInfo {
	return &RacingNumberInfo{racingNumbers: racingNumbers}
}

func (r *RacingNumberInfo) RacingNumbers() []*RacingNumber {
	return r.racingNumbers
}

func (r *RacingNumberInfo) Get(idx int) *RacingNumber {
	return r.racingNumbers[idx]
}
