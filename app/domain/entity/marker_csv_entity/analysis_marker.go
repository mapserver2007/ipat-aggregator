package marker_csv_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"strconv"
)

type AnalysisMarker struct {
	raceDate  types.RaceDate
	raceId    types.RaceId
	markerMap map[types.Marker]int
}

func NewAnalysisMarker(
	rawRaceDate,
	rawRaceId,
	rawHorseNumber1,
	rawHorseNumber2,
	rawHorseNumber3,
	rawHorseNumber4,
	rawHorseNumber5,
	rawHorseNumber6 string,
) (*AnalysisMarker, error) {
	raceDate, err := types.NewRaceDate(rawRaceDate)
	if err != nil {
		return nil, err
	}

	markerMap := map[types.Marker]int{}

	horseNumber1, _ := strconv.Atoi(rawHorseNumber1)
	markerMap[types.Favorite] = horseNumber1

	horseNumber2, _ := strconv.Atoi(rawHorseNumber2)
	markerMap[types.Rival] = horseNumber2

	horseNumber3, _ := strconv.Atoi(rawHorseNumber3)
	markerMap[types.BrackTriangle] = horseNumber3

	horseNumber4, _ := strconv.Atoi(rawHorseNumber4)
	markerMap[types.WhiteTriangle] = horseNumber4

	horseNumber5, _ := strconv.Atoi(rawHorseNumber5)
	markerMap[types.Star] = horseNumber5

	horseNumber6, _ := strconv.Atoi(rawHorseNumber6)
	markerMap[types.Check] = horseNumber6

	return &AnalysisMarker{
		raceDate:  raceDate,
		raceId:    types.RaceId(rawRaceId),
		markerMap: markerMap,
	}, nil
}

func (a *AnalysisMarker) RaceDate() types.RaceDate {
	return a.raceDate
}

func (a *AnalysisMarker) RaceId() types.RaceId {
	return a.raceId
}

func (a *AnalysisMarker) Favorite() int {
	horseNumber, ok := a.markerMap[types.Favorite]
	if !ok {
		return 0
	}
	return horseNumber
}

func (a *AnalysisMarker) Rival() int {
	horseNumber, ok := a.markerMap[types.Rival]
	if !ok {
		return 0
	}
	return horseNumber
}

func (a *AnalysisMarker) BrackTriangle() int {
	horseNumber, ok := a.markerMap[types.BrackTriangle]
	if !ok {
		return 0
	}
	return horseNumber
}

func (a *AnalysisMarker) WhiteTriangle() int {
	horseNumber, ok := a.markerMap[types.WhiteTriangle]
	if !ok {
		return 0
	}
	return horseNumber
}

func (a *AnalysisMarker) Star() int {
	horseNumber, ok := a.markerMap[types.Star]
	if !ok {
		return 0
	}
	return horseNumber
}

func (a *AnalysisMarker) Check() int {
	horseNumber, ok := a.markerMap[types.Check]
	if !ok {
		return 0
	}
	return horseNumber
}

func (a *AnalysisMarker) MarkerMap() map[types.Marker]int {
	return a.markerMap
}
