package raw_entity

type JockeyInfo struct {
	Jockeys          []*Jockey `json:"jockeys"`
	ExcludeJockeyIds []string  `json:"exclude_jockey_ids"`
}

type Jockey struct {
	JockeyId   string `json:"jockey_id"`
	JockeyName string `json:"jockey_name"`
}
