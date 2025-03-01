package marker_csv_entity

import (
	"strconv"

	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type PredictionMarker struct {
	raceId    types.RaceId
	markerMap map[types.Marker]types.HorseNumber
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

	return &PredictionMarker{
		raceId:    types.RaceId(rawRaceId),
		markerMap: markerMap,
	}
}

func (p *PredictionMarker) RaceId() types.RaceId {
	return p.raceId
}

func (p *PredictionMarker) Favorite() types.HorseNumber {
	horseNumber, ok := p.markerMap[types.Favorite]
	if !ok {
		return 0
	}
	return horseNumber
}

func (p *PredictionMarker) Rival() types.HorseNumber {
	horseNumber, ok := p.markerMap[types.Rival]
	if !ok {
		return 0
	}
	return horseNumber
}

func (p *PredictionMarker) BrackTriangle() types.HorseNumber {
	horseNumber, ok := p.markerMap[types.BrackTriangle]
	if !ok {
		return 0
	}
	return horseNumber
}

func (p *PredictionMarker) WhiteTriangle() types.HorseNumber {
	horseNumber, ok := p.markerMap[types.WhiteTriangle]
	if !ok {
		return 0
	}
	return horseNumber
}

func (p *PredictionMarker) Star() types.HorseNumber {
	horseNumber, ok := p.markerMap[types.Star]
	if !ok {
		return 0
	}
	return horseNumber
}

func (p *PredictionMarker) Check() types.HorseNumber {
	horseNumber, ok := p.markerMap[types.Check]
	if !ok {
		return 0
	}
	return horseNumber
}

func (p *PredictionMarker) MarkerMap() map[types.Marker]types.HorseNumber {
	return p.markerMap
}
