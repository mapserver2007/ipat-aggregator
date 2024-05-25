package analysis_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
	"github.com/shopspring/decimal"
)

type TrioCalculable struct {
	raceId              types.RaceId
	raceDate            types.RaceDate
	marker              types.Marker
	markerCombinationId types.MarkerCombinationId
	odds                decimal.Decimal
	number              types.BetNumber
	popular             int
	orderNo             int
	entries             int
	filters             []filter.Id
}
