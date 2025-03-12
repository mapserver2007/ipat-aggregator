package spreadsheet_entity

type AnalysisPlaceCheckPoint struct {
	itemName string
	point    int
}

func NewAnalysisPlaceCheckPoint(
	itemName string,
	point int,
) *AnalysisPlaceCheckPoint {
	return &AnalysisPlaceCheckPoint{
		itemName: itemName,
		point:    point,
	}
}

func (a *AnalysisPlaceCheckPoint) ItemName() string {
	return a.itemName
}

func (a *AnalysisPlaceCheckPoint) Point() int {
	return a.point
}
