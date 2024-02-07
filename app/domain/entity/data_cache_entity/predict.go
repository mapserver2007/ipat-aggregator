package data_cache_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type Predict struct {
	raceId  types.RaceId
	markers []*Marker
}

type Marker struct {
	marker      types.Marker
	horseNumber int
}

func NewPredict(
	raceId string,
	markers []*Marker,
) *Predict {
	return &Predict{
		raceId:  types.RaceId(raceId),
		markers: markers,
	}
}

func (p *Predict) RaceId() types.RaceId {
	return p.raceId
}

func (p *Predict) Markers() []*Marker {
	return p.markers
}

//func NewMarker(
//	rawMarker int,
//	horseNumber int,
//) (*Marker, error) {
//	marker, err := types.NewMarker(rawMarker)
//	if err != nil {
//		return nil, err
//	}
//
//	return &Marker{
//		marker:      marker,
//		horseNumber: horseNumber,
//	}, nil
//}
//
//func (m *Marker) Marker() types.Marker {
//	return m.marker
//}
//
//func (m *Marker) HorseNumber() int {
//	return m.horseNumber
//}
