package raw_entity

type JockeyInfo struct {
	Jockeys          []*Jockey `json:"jockeys"`
	ExcludeJockeyIds []int     `json:"exclude_jockey_ids"`
}

type Jockey struct {
	JockeyId   int    `json:"jockey_id"`
	JockeyName string `json:"jockey_name"`
}
