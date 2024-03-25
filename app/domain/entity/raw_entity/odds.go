package raw_entity

type RealTimeOddsInfo struct {
	Status      string `json:"status"`
	UpdateCount string `json:"update_count"`
	Reason      string `json:"reason"`
	Data        Data   `json:"data"`
}

type Data struct {
	OfficialDatetime string       `json:"official_datetime"`
	Odds             RealTimeOdds `json:"odds"`
}

type RealTimeOdds struct {
	List map[string][]string `json:"1"`
}

type FixedOddsInfo struct {
	TicketType int `json:"ticket_type"`
}
