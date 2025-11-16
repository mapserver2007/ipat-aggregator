package raw_entity

type JockeyResultInfo struct {
	JockeyResults     []*JockeyResult      `json:"jockey_results"`
	JockeyResultCount []*JockeyResultCount `json:"jockey_result_count"`
}

type JockeyResult struct {
	JockeyResultId   string `json:"jockey_result_id"`
	JockeyId         string `json:"jockey_id"`
	RaceId           string `json:"race_id"`
	RaceDate         int    `json:"race_date"`
	RaceCourseId     string `json:"race_course_id"`
	Odds             string `json:"odds"`
	OrderNo          int    `json:"order_no"`
	HorseId          string `json:"horse_id"`
	CourseCategoryId int    `json:"course_category_id"`
	Distance         int    `json:"distance"`
}

type JockeyResultCount struct {
	JockeyId string `json:"jockey_id"`
	Year     int    `json:"year"`
	Place    string `json:"place"`
	Count    int    `json:"count"`
}
