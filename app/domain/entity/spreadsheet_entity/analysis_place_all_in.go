package spreadsheet_entity

type AnalysisPlaceAllIn struct {
	rateDate  *PlaceAllInRateData
	rateStyle *PlaceAllInRateStyle
}

func NewAnalysisPlaceAllIn(
	rateData *PlaceAllInRateData,
	rateStyle *PlaceAllInRateStyle,
) *AnalysisPlaceAllIn {
	return &AnalysisPlaceAllIn{
		rateDate:  rateData,
		rateStyle: rateStyle,
	}
}

func (a *AnalysisPlaceAllIn) RateData() *PlaceAllInRateData {
	return a.rateDate
}

func (a *AnalysisPlaceAllIn) RateStyle() *PlaceAllInRateStyle {
	return a.rateStyle
}
