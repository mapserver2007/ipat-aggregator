package spreadsheet_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/shopspring/decimal"
)

type AnalysisPlaceUnHit struct {
	data  *PlaceUnHitListData
	style *PlaceUnHitListStyle
}

func NewAnalysisPlaceUnHit(
	data *PlaceUnHitListData,
	style *PlaceUnHitListStyle,
) *AnalysisPlaceUnHit {
	return &AnalysisPlaceUnHit{
		data:  data,
		style: style,
	}
}

func (a *AnalysisPlaceUnHit) Data() *PlaceUnHitListData {
	return a.data
}

func (a *AnalysisPlaceUnHit) Style() *PlaceUnHitListStyle {
	return a.style
}

type PlaceUnHitListData struct {
	raceDate               int
	raceId                 string
	raceName               string
	class                  string
	courseCategory         string
	distance               string
	trackCondition         string
	entries                int
	jockeyName             string
	horseName              string
	winOdds                string
	popularNumber          int
	yamatoMarker           string
	checkPoint1            bool
	checkPoint2            bool
	checkPoint3            bool
	checkPoint4            bool
	checkPoint5            bool
	checkPoint6            bool
	checkPoint7            bool
	checkPoint8            bool
	checkPoint9            bool
	checkPoint10           bool
	checkPoint11           bool
	checkPoint12           bool
	checkPoint13           bool
	checkPoint14           bool
	checkPoint15           bool
	tospoMarkerNum         int
	tospoFavoriteMarkerNum int
	tospoRivalMarkerNum    int
	tospoHighlyTraining    bool
	tospoTrainingComment   string
	analysisComment        string
}

func NewPlaceUnHitListData(
	raceDate types.RaceDate,
	raceId types.RaceId,
	raceName string,
	class types.GradeClass,
	courseCategory types.CourseCategory,
	distance string,
	trackCondition types.TrackCondition,
	entries int,
	jockeyName string,
	horseName string,
	winOdds decimal.Decimal,
	popularNumber int,
	yamatoMarker types.Marker,
	checkPoint1 bool,
	checkPoint2 bool,
	checkPoint3 bool,
	checkPoint4 bool,
	checkPoint5 bool,
	checkPoint6 bool,
	checkPoint7 bool,
	checkPoint8 bool,
	checkPoint9 bool,
	checkPoint10 bool,
	checkPoint11 bool,
	checkPoint12 bool,
	checkPoint13 bool,
	checkPoint14 bool,
	checkPoint15 bool,
	tospoMarkerNum int,
	tospoFavoriteMarkerNum int,
	tospoRivalMarkerNum int,
	tospoHighlyTraining bool,
	tospoTrainingComment string,
	analysisComment string,
) *PlaceUnHitListData {
	return &PlaceUnHitListData{
		raceDate:               raceDate.Value(),
		raceId:                 raceId.String(),
		raceName:               "",
		class:                  "",
		courseCategory:         "",
		distance:               "",
		trackCondition:         "",
		entries:                0,
		jockeyName:             "",
		horseName:              "",
		winOdds:                "",
		popularNumber:          0,
		yamatoMarker:           "",
		checkPoint1:            false,
		checkPoint2:            false,
		checkPoint3:            false,
		checkPoint4:            false,
		checkPoint5:            false,
		checkPoint6:            false,
		checkPoint7:            false,
		checkPoint8:            false,
		checkPoint9:            false,
		checkPoint10:           false,
		checkPoint11:           false,
		checkPoint12:           false,
		checkPoint13:           false,
		checkPoint14:           false,
		checkPoint15:           false,
		tospoMarkerNum:         0,
		tospoFavoriteMarkerNum: 0,
		tospoRivalMarkerNum:    0,
		tospoHighlyTraining:    false,
		tospoTrainingComment:   "",
		analysisComment:        "",
	}
}

type PlaceUnHitListStyle struct {
}
