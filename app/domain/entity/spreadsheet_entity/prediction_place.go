package spreadsheet_entity

type PredictionPlace struct {
	rateData  *PredictionRateData
	rateStyle *PredictionRateStyle
}

func NewPredictionPlace(
	rateData *PredictionRateData,
	rateStyle *PredictionRateStyle,
) *PredictionPlace {
	return &PredictionPlace{
		rateData:  rateData,
		rateStyle: rateStyle,
	}
}

func (a *PredictionPlace) RateData() *PredictionRateData {
	return a.rateData
}

func (a *PredictionPlace) RateStyle() *PredictionRateStyle {
	return a.rateStyle
}
