package raw_entity

type OddsInfo struct {
	Status      string `json:"status"`
	UpdateCount string `json:"update_count"`
	Reason      string `json:"reason"`
	Data        Data   `json:"data"`
}

type Data struct {
	OfficialDatetime string         `json:"official_datetime"`
	Odds             TicketTypeOdds `json:"odds"`
}

type TicketTypeOdds struct {
	Wins  map[string][]string `json:"1"`
	Trios map[string][]string `json:"7"`
}

type RaceOddsInfo struct {
	RaceOdds []*RaceOdds `json:"races"`
}

type RaceOdds struct {
	RaceId string  `json:"race_id"`
	Odds   []*Odds `json:"odds"`
}

type Odds struct {
	TicketType int    `json:"ticket_type"`
	Odds       string `json:"odds"`
	Popular    int    `json:"popular"`
	Number     string `json:"number"`
}
