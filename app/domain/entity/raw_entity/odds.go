package raw_entity

type OddsInfo struct {
	Status      string `json:"status"`
	UpdateCount string `json:"update_count"`
	Reason      string `json:"reason"`
	Data        Data   `json:"data"`
}

type Data struct {
	OfficialDatetime string `json:"official_datetime"`
	Odds             Odds   `json:"odds"`
}

type Odds struct {
	List map[string][]string `json:"1"`
}
