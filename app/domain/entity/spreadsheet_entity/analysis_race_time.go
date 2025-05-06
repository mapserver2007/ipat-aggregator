package spreadsheet_entity

type AnalysisRaceTime struct {
	averageRaceTime   string
	medianRaceTime    string
	averageFirst3f    string
	medianFirst3f     string
	averageFirst4f    string
	medianFirst4f     string
	averageLast3f     string
	medianLast3f      string
	averageLast4f     string
	medianLast4f      string
	averageRap5f      string
	medianRap5f       string
	averageTrackIndex int
	maxTrackIndex     int
	minTrackIndex     int
	averageTimeIndex  int
	raceCount         int
}

func NewAnalysisRaceTime(
	averageRaceTime string,
	medianRaceTime string,
	averageFirst3f string,
	medianFirst3f string,
	averageFirst4f string,
	medianFirst4f string,
	averageLast3f string,
	medianLast3f string,
	averageLast4f string,
	medianLast4f string,
	averageRap5f string,
	medianRap5f string,
	averageTrackIndex int,
	maxTrackIndex int,
	minTrackIndex int,
	averageTimeIndex int,
	raceCount int,
) *AnalysisRaceTime {
	return &AnalysisRaceTime{
		averageRaceTime:   averageRaceTime,
		medianRaceTime:    medianRaceTime,
		averageFirst3f:    averageFirst3f,
		medianFirst3f:     medianFirst3f,
		averageFirst4f:    averageFirst4f,
		medianFirst4f:     medianFirst4f,
		averageLast3f:     averageLast3f,
		medianLast3f:      medianLast3f,
		averageLast4f:     averageLast4f,
		medianLast4f:      medianLast4f,
		averageRap5f:      averageRap5f,
		medianRap5f:       medianRap5f,
		averageTrackIndex: averageTrackIndex,
		maxTrackIndex:     maxTrackIndex,
		minTrackIndex:     minTrackIndex,
		averageTimeIndex:  averageTimeIndex,
		raceCount:         raceCount,
	}
}

func (a *AnalysisRaceTime) AverageRaceTime() string {
	return a.averageRaceTime
}

func (a *AnalysisRaceTime) MedianRaceTime() string {
	return a.medianRaceTime
}

func (a *AnalysisRaceTime) AverageFirst3f() string {
	return a.averageFirst3f
}

func (a *AnalysisRaceTime) MedianFirst3f() string {
	return a.medianFirst3f
}

func (a *AnalysisRaceTime) AverageFirst4f() string {
	return a.averageFirst4f
}

func (a *AnalysisRaceTime) MedianFirst4f() string {
	return a.medianFirst4f
}

func (a *AnalysisRaceTime) AverageLast3f() string {
	return a.averageLast3f
}

func (a *AnalysisRaceTime) MedianLast3f() string {
	return a.medianLast3f
}

func (a *AnalysisRaceTime) AverageLast4f() string {
	return a.averageLast4f
}

func (a *AnalysisRaceTime) MedianLast4f() string {
	return a.medianLast4f
}

func (a *AnalysisRaceTime) AverageRap5f() string {
	return a.averageRap5f
}

func (a *AnalysisRaceTime) MedianRap5f() string {
	return a.medianRap5f
}

func (a *AnalysisRaceTime) AverageTrackIndex() int {
	return a.averageTrackIndex
}

func (a *AnalysisRaceTime) MaxTrackIndex() int {
	return a.maxTrackIndex
}

func (a *AnalysisRaceTime) MinTrackIndex() int {
	return a.minTrackIndex
}

func (a *AnalysisRaceTime) AverageTimeIndex() int {
	return a.averageTimeIndex
}

func (a *AnalysisRaceTime) RaceCount() int {
	return a.raceCount
}
