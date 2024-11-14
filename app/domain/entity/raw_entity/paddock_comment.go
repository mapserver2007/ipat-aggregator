package raw_entity

type PaddockCommentInfo struct {
	Body *PaddockCommentBody `json:"body"`
}

type PaddockCommentBody struct {
	RaceEntryList []*PaddockComment `json:"raceEntryList"`
}

type PaddockComment struct {
	HorseNumber int    `json:"horseNumber"`
	Comment     string `json:"paddockComment"`
	Evaluation  int    `json:"evaluation"`
}
