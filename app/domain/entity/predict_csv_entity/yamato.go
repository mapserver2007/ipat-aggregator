package predict_csv_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"strconv"
)

type Yamato struct {
	raceDate  types.RaceDate
	raceId    types.RaceId
	markerMap map[types.Marker]int
}

func NewYamato(
	rawRaceDate,
	rawRaceId,
	rawHorseNumber1,
	rawHorseNumber2,
	rawHorseNumber3,
	rawHorseNumber4,
	rawHorseNumber5,
	rawHorseNumber6 string,
) (*Yamato, error) {
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

	return &Yamato{
		raceDate:  raceDate,
		raceId:    types.RaceId(rawRaceId),
		markerMap: markerMap,
	}, nil
}

func (y *Yamato) RaceDate() types.RaceDate {
	return y.raceDate
}

func (y *Yamato) RaceId() types.RaceId {
	return y.raceId
}

func (y *Yamato) Favorite() int {
	horseNumber, ok := y.markerMap[types.Favorite]
	if !ok {
		return 0
	}
	return horseNumber
}

func (y *Yamato) Rival() int {
	horseNumber, ok := y.markerMap[types.Rival]
	if !ok {
		return 0
	}
	return horseNumber
}

func (y *Yamato) BrackTriangle() int {
	horseNumber, ok := y.markerMap[types.BrackTriangle]
	if !ok {
		return 0
	}
	return horseNumber
}

func (y *Yamato) WhiteTriangle() int {
	horseNumber, ok := y.markerMap[types.WhiteTriangle]
	if !ok {
		return 0
	}
	return horseNumber
}

func (y *Yamato) Star() int {
	horseNumber, ok := y.markerMap[types.Star]
	if !ok {
		return 0
	}
	return horseNumber
}

func (y *Yamato) Check() int {
	horseNumber, ok := y.markerMap[types.Check]
	if !ok {
		return 0
	}
	return horseNumber
}

func (y *Yamato) MarkerMap() map[types.Marker]int {
	return y.markerMap
}
