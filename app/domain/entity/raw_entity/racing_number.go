package raw_entity

type RacingNumberInfo struct {
	RacingNumbers []*RacingNumber `json:"racing_numbers"`
}

type RacingNumber struct {
	Date         int `json:"date"`
	Round        int `json:"round"`
	Day          int `json:"day"`
	RaceCourseId int `json:"race_course_id"`
}
