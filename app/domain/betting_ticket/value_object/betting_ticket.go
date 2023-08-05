package value_object

type BettingTicket int

const (
	UnknownTicket BettingTicket = iota
	Win
	Place
	BracketQuinella
	Quinella
	Exacta
	ExactaWheelOfFirst
	QuinellaPlace
	QuinellaPlaceWheel
	Trio
	TrioFormation
	TrioWheelOfFirst
	Trifecta
	TrifectaFormation
	TrifectaWheelOfFirst
)

var bettingTicketMap = map[BettingTicket]string{
	Win:                  "単勝",
	Place:                "複勝",
	BracketQuinella:      "枠連",
	Quinella:             "馬連",
	Exacta:               "馬単",
	ExactaWheelOfFirst:   "馬単1着ながし",
	QuinellaPlace:        "ワイド",
	QuinellaPlaceWheel:   "ワイドながし",
	Trio:                 "3連複",
	TrioFormation:        "3連複フォーメーション",
	TrioWheelOfFirst:     "3連複軸1頭ながし",
	Trifecta:             "3連単",
	TrifectaFormation:    "3連単フォーメーション",
	TrifectaWheelOfFirst: "3連単1着ながし",
	UnknownTicket:        "不明",
}

func NewBettingTicket(name string) BettingTicket {
	for key, value := range bettingTicketMap {
		if value == name {
			return key
		}
	}

	return UnknownTicket
}

func (b BettingTicket) Name() string {
	return convertToBettingTicketName(b)
}

func (b BettingTicket) Value() int {
	return int(b)
}

func (b BettingTicket) ConvertToOriginBettingTicket() BettingTicket {
	switch b {
	case ExactaWheelOfFirst:
		return Exacta
	case QuinellaPlaceWheel:
		return QuinellaPlace
	case TrioFormation, TrioWheelOfFirst:
		return Trio
	case TrifectaFormation, TrifectaWheelOfFirst:
		return Trifecta
	}

	return b
}

func convertToBettingTicketName(b BettingTicket) string {
	bettingTicketName, _ := bettingTicketMap[b]
	return bettingTicketName
}
