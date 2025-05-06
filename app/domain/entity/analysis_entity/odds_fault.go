package analysis_entity

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/shopspring/decimal"
)

type OddsFault struct {
	raceId      types.RaceId
	oddsFaultNo int
	oddsFault   decimal.Decimal
}

func NewOddsFault(
	raceId types.RaceId,
	odds1 string,
	odds2 string,
	oddsFaultNo int,
) (*OddsFault, error) {
	odds1Decimal, err := decimal.NewFromString(odds1)
	if err != nil {
		return nil, err
	}
	odds2Decimal, err := decimal.NewFromString(odds2)
	if err != nil {
		return nil, err
	}

	oddsFault := odds2Decimal.Sub(odds1Decimal)

	return &OddsFault{
		raceId:      raceId,
		oddsFaultNo: oddsFaultNo,
		oddsFault:   oddsFault,
	}, nil
}

func (o *OddsFault) RaceId() types.RaceId {
	return o.raceId
}

func (o *OddsFault) OddsFaultNo() int {
	return o.oddsFaultNo
}

func (o *OddsFault) OddsFault() decimal.Decimal {
	return o.oddsFault
}
