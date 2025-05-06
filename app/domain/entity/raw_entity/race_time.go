package raw_entity

type RaceTimeInfo struct {
	RaceTimes []*RaceTime `json:"race_times"`
}

type RaceTime struct {
	RaceId     string   `json:"race_id"`
	RaceDate   int      `json:"race_date"`
	Time       string   `json:"time"`
	TimeIndex  int      `json:"time_index"`
	TrackIndex int      `json:"track_index"`
	RapTimes   []string `json:"rap_times"`
	First3f    string   `json:"first3f"`
	First4f    string   `json:"first4f"`
	Last3f     string   `json:"last3f"`
	Last4f     string   `json:"last4f"`
	Rap5f      string   `json:"rap5f"`
}
