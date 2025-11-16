package mcp_result_entity

type MCPResult struct {
	Races     []*Race    `json:"races"`
	RaceNote  string     `json:"race_note"`
	Condition *Condition `json:"condition"`
}

type Race struct {
	RaceId              string        `json:"race_id"`
	RaceDate            string        `json:"race_date"`
	RaceName            string        `json:"race_name"`
	RaceNumber          int           `json:"race_number"`
	RaceCourse          string        `json:"race_course"`
	RaceUrl             string        `json:"race_url"`
	Entries             int           `json:"entries"`
	Distance            int           `json:"distance"`
	Class               string        `json:"class"`
	CourseCategory      string        `json:"course_category"`
	TrackCondition      string        `json:"track_condition"`
	RaceSexCondition    string        `json:"race_sex_condition"`
	RaceWeightCondition string        `json:"race_weight_condition"`
	RaceAgeCondition    string        `json:"race_age_condition"`
	RaceResults         []*RaceResult `json:"race_results"`
}

type RaceResult struct {
	OrderNo        int    `json:"order_no"`
	HorseId        string `json:"horse_id"`
	HorseName      string `json:"horse_name"`
	BracketNumber  int    `json:"bracket_number"`
	HorseNumber    int    `json:"horse_number"`
	JockeyId       string `json:"jockey_id"`
	Odds           string `json:"odds"`
	PopularNumber  int    `json:"popular_number"`
	JockeyWeight   string `json:"jockey_weight"`
	HorseWeight    int    `json:"horse_weight"`
	HorseWeightAdd int    `json:"horse_weight_add"`
}

type Condition struct {
	Odds          *Odds  `json:"odds"`
	ConditionNote string `json:"condition_note"`
}

// type Jockey struct {
// 	JockeyId   string `json:"jockey_id"`
// 	JockeyName string `json:"jockey_name"`
// }

type Odds struct {
	OddsList []string `json:"odds_list"`
	RaceIds  []string `json:"race_ids"`
	OddsNote string   `json:"odds_note"`
}

// type RaceUrl struct {
// 	Url  string `json:"url"`
// 	Note string `json:"note"`
// }

// type RaceTime struct {
// }
