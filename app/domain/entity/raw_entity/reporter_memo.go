package raw_entity

type ReporterMemo struct {
	Body *ReporterMemoBody `json:"body"`
}

type ReporterMemoBody struct {
	ReceivedMemoList []*ReceivedMemo `json:"receivedMemoList"`
}

type ReceivedMemo struct {
	Date        string `json:"date"`
	HorseNumber int    `json:"horseNumber"`
	Comment     string `json:"memoContent"`
}
