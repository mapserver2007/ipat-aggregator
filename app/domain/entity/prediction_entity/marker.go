package prediction_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type Marker struct {
	raceId      types.RaceId
	horseNumber types.HorseNumber
	marker      types.Marker
}

func NewMarker(
	raceId types.RaceId,
	horseNumber types.HorseNumber,
	marker types.Marker,
) *Marker {
	return &Marker{
		raceId:      raceId,
		horseNumber: horseNumber,
		marker:      marker,
	}
}

func (m *Marker) RaceId() types.RaceId {
	return m.raceId
}

func (m *Marker) HorseNumber() types.HorseNumber {
	return m.horseNumber
}

func (m *Marker) Marker() types.Marker {
	return m.marker
}
