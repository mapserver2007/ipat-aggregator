package analysis_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type Marker struct {
	marker      types.Marker
	horseNumber types.HorseNumber
}

func NewMarker(
	marker types.Marker,
	horseNumber types.HorseNumber,
) *Marker {
	return &Marker{
		marker:      marker,
		horseNumber: horseNumber,
	}
}

func (m *Marker) Marker() types.Marker {
	return m.marker
}

func (m *Marker) HorseNumber() types.HorseNumber {
	return m.horseNumber
}
