package netkeiba_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"strconv"
)

type Marker struct {
	horseNumber types.HorseNumber
	marker      types.Marker
}

func NewMarker(
	rawHorseNumber string,
	rawMarker string,
) (*Marker, error) {
	var marker types.Marker
	switch rawMarker {
	case "0":
		marker = types.NoMarker
	case "1":
		marker = types.Favorite
	case "2":
		marker = types.Rival
	case "3":
		marker = types.BrackTriangle
	case "4":
		marker = types.WhiteTriangle
	case "5":
		marker = types.Star
	case "98":
		marker = types.Check
	}

	horseNumber, err := strconv.Atoi(rawHorseNumber)
	if err != nil {
		return nil, err
	}

	return &Marker{
		horseNumber: types.HorseNumber(horseNumber),
		marker:      marker,
	}, nil
}

func (m *Marker) HorseNumber() types.HorseNumber {
	return m.horseNumber
}

func (m *Marker) Marker() types.Marker {
	return m.marker
}
