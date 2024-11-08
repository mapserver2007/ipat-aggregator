package types

type TrackCondition int

// 厳密には芝とダートでは表記が違うが芝の表記に統一
const (
	UnknownTrackCondition TrackCondition = iota
	GoodToFirm
	Good
	Yielding
	Soft
)

func NewTrackCondition(name string) TrackCondition {
	var trackCondition TrackCondition
	switch name {
	case "良":
		trackCondition = GoodToFirm
	case "稍":
		trackCondition = Good
	case "重":
		trackCondition = Yielding
	case "不":
		trackCondition = Soft
	}

	return trackCondition
}

var trackConditionMap = map[TrackCondition]string{
	UnknownTrackCondition: "不明",
	GoodToFirm:            "良",
	Good:                  "稍",
	Yielding:              "重",
	Soft:                  "不",
}

func (t TrackCondition) Value() int {
	return int(t)
}

func (t TrackCondition) String() string {
	trackConditionName, _ := trackConditionMap[t]
	return trackConditionName
}
