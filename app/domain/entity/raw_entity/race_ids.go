package raw_entity

type RaceIdInfo struct {
	RaceDates    []*RaceDate `json:"race_dates"`
	ExcludeDates []int       `json:"exclude_dates"`
}

type RaceDate struct {
	RaceDate int      `json:"race_date"`
	RaceIds  []string `json:"race_ids"`
}
