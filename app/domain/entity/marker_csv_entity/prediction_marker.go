package marker_csv_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type PredictionMarker struct {
	raceId    types.RaceId
	markerMap map[types.Marker]int
}

func NewPredictionMarker(
	rawRaceId,
	rawHorseNumber1,
	rawHorseNumber2,
	rawHorseNumber3,
	rawHorseNumber4,
	rawHorseNumber5,
	rawHorseNumber6 string,
) *PredictionMarker {
	markerMap := map[types.Marker]int{}

	horseNumber1 := 0
	if rawHorseNumber1 != "" {
		horseNumber1 = 1
	}
	markerMap[types.Favorite] = horseNumber1

	horseNumber2 := 0
	if rawHorseNumber2 != "" {
		horseNumber2 = 1
	}
	markerMap[types.Rival] = horseNumber2

	horseNumber3 := 0
	if rawHorseNumber3 != "" {
		horseNumber3 = 1
	}
	markerMap[types.BrackTriangle] = horseNumber3

	horseNumber4 := 0
	if rawHorseNumber4 != "" {
		horseNumber4 = 1
	}
	markerMap[types.WhiteTriangle] = horseNumber4

	horseNumber5 := 0
	if rawHorseNumber5 != "" {
		horseNumber5 = 1
	}
	markerMap[types.Star] = horseNumber5

	horseNumber6 := 0
	if rawHorseNumber6 != "" {
		horseNumber6 = 1
	}
	markerMap[types.Check] = horseNumber6

	return &PredictionMarker{
		raceId:    types.RaceId(rawRaceId),
		markerMap: markerMap,
	}
}

func (p *PredictionMarker) RaceId() types.RaceId {
	return p.raceId
}

func (p *PredictionMarker) Favorite() int {
	horseNumber, ok := p.markerMap[types.Favorite]
	if !ok {
		return 0
	}
	return horseNumber
}

func (p *PredictionMarker) Rival() int {
	horseNumber, ok := p.markerMap[types.Rival]
	if !ok {
		return 0
	}
	return horseNumber
}

func (p *PredictionMarker) BrackTriangle() int {
	horseNumber, ok := p.markerMap[types.BrackTriangle]
	if !ok {
		return 0
	}
	return horseNumber
}

func (p *PredictionMarker) WhiteTriangle() int {
	horseNumber, ok := p.markerMap[types.WhiteTriangle]
	if !ok {
		return 0
	}
	return horseNumber
}

func (p *PredictionMarker) Star() int {
	horseNumber, ok := p.markerMap[types.Star]
	if !ok {
		return 0
	}
	return horseNumber
}

func (p *PredictionMarker) Check() int {
	horseNumber, ok := p.markerMap[types.Check]
	if !ok {
		return 0
	}
	return horseNumber
}

func (p *PredictionMarker) MarkerMap() map[types.Marker]int {
	return p.markerMap
}
