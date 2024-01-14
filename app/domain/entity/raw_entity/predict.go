package raw_entity

type PredictInfo struct {
	Predicts []*Predict `json:"predicts"`
}

type Predict struct {
	RaceId  string    `json:"race_id"`
	Markers []*Marker `json:"markers"`
}

type Marker struct {
	Marker      int `json:"marker"`
	HorseNumber int `json:"horse_number"`
}
