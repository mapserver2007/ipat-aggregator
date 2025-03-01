package marker_csv_entity

import (
	"strconv"

	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type AnalysisMarker struct {
	raceDate  types.RaceDate
	raceId    types.RaceId
	markerMap map[types.Marker]types.HorseNumber
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

	markerMap := map[types.Marker]types.HorseNumber{}

	horseNumber1, _ := strconv.Atoi(rawHorseNumber1)
	markerMap[types.Favorite] = types.HorseNumber(horseNumber1)

	horseNumber2, _ := strconv.Atoi(rawHorseNumber2)
	markerMap[types.Rival] = types.HorseNumber(horseNumber2)

	horseNumber3, _ := strconv.Atoi(rawHorseNumber3)
	markerMap[types.BrackTriangle] = types.HorseNumber(horseNumber3)

	horseNumber4, _ := strconv.Atoi(rawHorseNumber4)
	markerMap[types.WhiteTriangle] = types.HorseNumber(horseNumber4)

	horseNumber5, _ := strconv.Atoi(rawHorseNumber5)
	markerMap[types.Star] = types.HorseNumber(horseNumber5)

	horseNumber6, _ := strconv.Atoi(rawHorseNumber6)
	markerMap[types.Check] = types.HorseNumber(horseNumber6)

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

func (a *AnalysisMarker) Favorite() types.HorseNumber {
	horseNumber, ok := a.markerMap[types.Favorite]
	if !ok {
		return 0
	}
	return horseNumber
}

func (a *AnalysisMarker) Rival() types.HorseNumber {
	horseNumber, ok := a.markerMap[types.Rival]
	if !ok {
		return 0
	}
	return horseNumber
}

func (a *AnalysisMarker) BrackTriangle() types.HorseNumber {
	horseNumber, ok := a.markerMap[types.BrackTriangle]
	if !ok {
		return 0
	}
	return horseNumber
}

func (a *AnalysisMarker) WhiteTriangle() types.HorseNumber {
	horseNumber, ok := a.markerMap[types.WhiteTriangle]
	if !ok {
		return 0
	}
	return horseNumber
}

func (a *AnalysisMarker) Star() types.HorseNumber {
	horseNumber, ok := a.markerMap[types.Star]
	if !ok {
		return 0
	}
	return horseNumber
}

func (a *AnalysisMarker) Check() types.HorseNumber {
	horseNumber, ok := a.markerMap[types.Check]
	if !ok {
		return 0
	}
	return horseNumber
}

func (a *AnalysisMarker) MarkerMap() map[types.Marker]types.HorseNumber {
	return a.markerMap
}
