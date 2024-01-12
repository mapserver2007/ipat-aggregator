package raw_entity

type RaceIdInfo struct {
	RaceDates    []*RaceDate `json:"race_dates"`
	ExcludeDates []string    `json:"exclude_dates"`
}

type RaceDate struct {
	RaceDate string   `json:"race_date"`
	RaceIds  []string `json:"race_ids"`
}
