package types

type PredictStatus byte

const (
	PredictUncompleted = PredictStatus(0x00) // 予想未完了
	FavoriteCandidate  = PredictStatus(0x01) // 本命候補が複数ある
	FavoriteCompleted  = PredictStatus(0x02) // 本命確定
	RivalCandidate     = PredictStatus(0x04) // 対抗候補が複数ある
	RivalCompleted     = PredictStatus(0x08) // 対抗確定
)

func (p PredictStatus) Included(target PredictStatus) bool {
	return p&target != 0
}

func (p PredictStatus) Matched(target PredictStatus) bool {
	return p == target
}
