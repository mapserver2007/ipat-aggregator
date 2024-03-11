package prediction_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type OddsRange struct {
	oddsRangeType types.OddsRangeType
	inOrder       types.InOrder
}

func NewOddsRange(
	oddsRangeType types.OddsRangeType,
	inOrder types.InOrder,
) *OddsRange {
	return &OddsRange{
		oddsRangeType: oddsRangeType,
		inOrder:       inOrder,
	}
}

func (o *OddsRange) OddsRangeType() types.OddsRangeType {
	return o.oddsRangeType
}

func (o *OddsRange) InOrder() types.InOrder {
	return o.inOrder
}
