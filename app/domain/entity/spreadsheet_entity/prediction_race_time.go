package spreadsheet_entity

type PredictionRaceTime struct {
	averageRaceTime string
	averageFirst3f  string
	averageFirst4f  string
	averageLast3f   string
	averageLast4f   string
	averageRap5f    string
}

func NewPredictionRaceTime(
	averageRaceTime string,
	averageFirst3f string,
	averageFirst4f string,
	averageLast3f string,
	averageLast4f string,
	averageRap5f string,
) *PredictionRaceTime {
	return &PredictionRaceTime{
		averageRaceTime: averageRaceTime,
		averageFirst3f:  averageFirst3f,
		averageFirst4f:  averageFirst4f,
		averageLast3f:   averageLast3f,
		averageLast4f:   averageLast4f,
		averageRap5f:    averageRap5f,
	}
}

func InitPredictionRaceTime() *PredictionRaceTime {
	return NewPredictionRaceTime(
		"0:00:0",
		"0.0",
		"0.0",
		"0.0",
		"0.0",
		"0.0",
	)
}

func (p *PredictionRaceTime) AverageRaceTime() string {
	return p.averageRaceTime
}

func (p *PredictionRaceTime) AverageFirst3f() string {
	return p.averageFirst3f
}

func (p *PredictionRaceTime) AverageFirst4f() string {
	return p.averageFirst4f
}

func (p *PredictionRaceTime) AverageLast3f() string {
	return p.averageLast3f
}

func (p *PredictionRaceTime) AverageLast4f() string {
	return p.averageLast4f
}

func (p *PredictionRaceTime) AverageRap5f() string {
	return p.averageRap5f
}
