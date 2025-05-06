package raw_entity

type RaceForecastInfo struct {
	RaceForecasts []*RaceForecast `json:"race_forecasts"`
}

type RaceForecast struct {
	RaceId    string      `json:"race_id"`
	RaceDate  int         `json:"race_date"`
	Forecasts []*Forecast `json:"forecasts"`
}

type Forecast struct {
	HorseNumber             int    `json:"horse_number"`
	TrainingComment         string `json:"training_comment"`
	PreviousTrainingComment string `json:"previous_training_comment"`
	HighlyRecommended       bool   `json:"highly_recommended"`
	FavoriteNum             int    `json:"favorite_num"`
	RivalNum                int    `json:"rival_num"`
	MarkerNum               int    `json:"marker_num"`
}

type ForecastInfo struct {
	Body *ForecastBody `json:"body"`
}

type ForecastBody struct {
	RaceDateInfo        *RaceDateInfo        `json:"raceInfo"`
	RaceEntries         []*RaceEntry         `json:"raceEntryList"`
	RaceForecastElement *RaceForecastElement `json:"raceForecast"`
}

type RaceDateInfo struct {
	RaceDate string `json:"raceDate"`
}

type RaceEntry struct {
	HorseNumber int    `json:"horseNumber"`
	HorseName   string `json:"horseName"`
}

type RaceForecastElement struct {
	ReporterList          []*RaceForecastReporter                  `json:"reporterList"`
	RaceForecastMarkerMap map[string]map[string]RaceForecastMarker `json:"reporterMarks"`
}

type RaceForecastReporter struct {
	ReporterId int `json:"reporterId"`
}

type RaceForecastMarker struct {
	ReporterMarkType int    `json:"reporterMarkType"`
	HorseName        string `json:"horseName"`
}

type TrainingComment struct {
	Body *TrainingCommentBody `json:"body"`
}

type TrainingCommentBody struct {
	RaceTrainingComments []*RaceTrainingComment `json:"raceCommentList"`
}

type RaceTrainingComment struct {
	HorseNumber            int                         `json:"horseNumber"`
	Prediction             string                      `json:"prediction"`
	TrainingComment        string                      `json:"interestingComment"`
	RaceHistoryCommentInfo *RaceTrainingCommentHistory `json:"raceHistoryCommentInfo"`
}

type RaceTrainingCommentHistory struct {
	TrainingComment string `json:"historyComment"`
}
