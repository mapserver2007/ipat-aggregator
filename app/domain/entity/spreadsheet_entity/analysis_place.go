package spreadsheet_entity

type AnalysisPlace struct {
	rateData       *PlaceRateData
	rateStyle      *PlaceRateStyle
	hitCountData   *PlaceHitCountData
	unHitCountData *PlaceUnHitCountData
}

func NewAnalysisPlace(
	rateData *PlaceRateData,
	rateStyle *PlaceRateStyle,
	hitCountData *PlaceHitCountData,
	unHitCountData *PlaceUnHitCountData,
) *AnalysisPlace {
	return &AnalysisPlace{
		rateData:       rateData,
		rateStyle:      rateStyle,
		hitCountData:   hitCountData,
		unHitCountData: unHitCountData,
	}
}

func (a *AnalysisPlace) RateData() *PlaceRateData {
	return a.rateData
}

func (a *AnalysisPlace) RateStyle() *PlaceRateStyle {
	return a.rateStyle
}

func (a *AnalysisPlace) HitCountData() *PlaceHitCountData {
	return a.hitCountData
}

func (a *AnalysisPlace) UnHitCountData() *PlaceUnHitCountData {
	return a.unHitCountData
}
