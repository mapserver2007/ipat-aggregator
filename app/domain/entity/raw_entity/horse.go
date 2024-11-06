package raw_entity

type HorseInfo struct {
	Horses []*Horse `json:"horses"`
}

type Horse struct {
	HorseId        string         `json:"horse_id"`
	HorseName      string         `json:"horse_name"`
	HorseBirthDay  int            `json:"horse_birth_day"`
	TrainerId      string         `json:"trainer_id"`
	OwnerId        string         `json:"owner_id"`
	BreederId      string         `json:"breeder_id"`
	HorseBlood     *HorseBlood    `json:"horse_blood"`
	HorseResults   []*HorseResult `json:"horse_results"`
	LatestRaceDate int            `json:"latest_race_date"`
}

type HorseBlood struct {
	SireId          string `json:"sire_id"`
	BroodmareSireId string `json:"broodmare_sire_id"`
}

type HorseResult struct {
	RaceId           string `json:"race_id"`
	RaceDate         int    `json:"race_date"`
	RaceName         string `json:"race_name"`
	JockeyId         string `json:"jockey_id"`
	OrderNo          int    `json:"order_no"`
	PopularNumber    int    `json:"popular_number"`
	HorseNumber      int    `json:"horse_number"`
	Odds             string `json:"odds"`
	Class            int    `json:"class"`
	Entries          int    `json:"entries"`
	Distance         int    `json:"distance"`
	RaceCourseId     string `json:"race_course_id"`
	CourseCategoryId int    `json:"course_category_id"`
	TrackConditionId int    `json:"track_condition_id"`
	HorseWeight      int    `json:"horse_weight"`
	RaceWeight       string `json:"race_weight"`
	Comment          string `json:"comment"`
}
