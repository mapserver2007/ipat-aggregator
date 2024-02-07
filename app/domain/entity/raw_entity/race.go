package raw_entity

type RaceInfo struct {
	Races []*Race `json:"races"`
}

type Race struct {
	RaceId              string          `json:"race_id"`
	RaceDate            int             `json:"race_date"`
	RaceNumber          int             `json:"race_number"`
	RaceCourseId        string          `json:"race_course_id"`
	RaceName            string          `json:"race_name"`
	Organizer           int             `json:"organizer"`
	Url                 string          `json:"url"`
	Time                string          `json:"time"`
	StartTime           string          `json:"start_time"`
	Entries             int             `json:"entries"`
	Distance            int             `json:"distance"`
	Class               int             `json:"class"`
	CourseCategory      int             `json:"course_category"`
	TrackCondition      string          `json:"track_condition"`
	RaceSexCondition    int             `json:"race_sex_condition"`
	RaceWeightCondition int             `json:"race_weight_condition"`
	RaceResults         []*RaceResult   `json:"race_results"`
	PayoutResults       []*PayoutResult `json:"payout_results"`
}

type RaceResult struct {
	OrderNo       int    `json:"order_no"`
	HorseName     string `json:"horse_name"`
	BracketNumber int    `json:"bracket_number"`
	HorseNumber   int    `json:"horse_number"`
	JockeyId      int    `json:"jockey_id"`
	Odds          string `json:"odds"`
	PopularNumber int    `json:"popular_number"`
}

type PayoutResult struct {
	TicketType int      `json:"ticket_type"`
	Numbers    []string `json:"numbers"`
	Odds       []string `json:"odds"`
	Populars   []int    `json:"populars"`
}
