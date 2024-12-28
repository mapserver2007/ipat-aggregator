package types

type TicketType int

const (
	UnknownTicketType TicketType = iota
	Win
	Place
	BracketQuinella
	Quinella
	QuinellaWheel
	Exacta
	ExactaWheelOfFirst
	QuinellaPlace
	QuinellaPlaceWheel
	QuinellaPlaceFormation
	Trio
	TrioFormation
	TrioWheelOfFirst
	TrioWheelOfSecond
	TrioBox
	Trifecta
	TrifectaFormation
	TrifectaWheelOfFirst
	TrifectaWheelOfSecond
	TrifectaWheelOfFirstMulti
	TrifectaWheelOfSecondMulti
	AllTicketType
)

var ticketTypeMap = map[TicketType]string{
	Win:                        "単勝",
	Place:                      "複勝",
	BracketQuinella:            "枠連",
	Quinella:                   "馬連",
	QuinellaWheel:              "馬連ながし",
	Exacta:                     "馬単",
	ExactaWheelOfFirst:         "馬単1着ながし",
	QuinellaPlace:              "ワイド",
	QuinellaPlaceWheel:         "ワイドながし",
	QuinellaPlaceFormation:     "ワイドフォーメーション",
	Trio:                       "3連複",
	TrioFormation:              "3連複フォーメーション",
	TrioWheelOfFirst:           "3連複軸1頭ながし",
	TrioWheelOfSecond:          "3連複軸2頭ながし",
	TrioBox:                    "3連複ＢＯＸ",
	Trifecta:                   "3連単",
	TrifectaFormation:          "3連単フォーメーション",
	TrifectaWheelOfFirst:       "3連単1着ながし",
	TrifectaWheelOfSecond:      "3連単2着ながし",
	TrifectaWheelOfFirstMulti:  "3連単軸1頭ながしマルチ",
	TrifectaWheelOfSecondMulti: "3連単軸2頭ながしマルチ",
	AllTicketType:              "全券種合計",
	UnknownTicketType:          "不明",
}

func NewTicketType(name string) TicketType {
	for key, value := range ticketTypeMap {
		if value == name {
			return key
		}
	}

	return UnknownTicketType
}

func (b TicketType) Name() string {
	name, _ := ticketTypeMap[b]
	return name
}

func (b TicketType) OriginTicketType() TicketType {
	switch b {
	case QuinellaWheel:
		return Quinella
	case ExactaWheelOfFirst:
		return Exacta
	case QuinellaPlaceWheel, QuinellaPlaceFormation:
		return QuinellaPlace
	case TrioFormation, TrioWheelOfFirst, TrioWheelOfSecond, TrioBox:
		return Trio
	case TrifectaFormation, TrifectaWheelOfFirst, TrifectaWheelOfSecond, TrifectaWheelOfFirstMulti, TrifectaWheelOfSecondMulti:
		return Trifecta
	}
	return b
}

func (b TicketType) Value() int {
	return int(b)
}
