package converter

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/raw_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"strconv"
	"strings"
)

type OddsEntityConverter interface {
	DataCacheToRaw(input *data_cache_entity.Odds) *raw_entity.Odds
	RawToDataCache(input *raw_entity.Odds, raceId types.RaceId, raceDate types.RaceDate) *data_cache_entity.Odds
	NetKeibaToRaw(input *netkeiba_entity.Odds) *raw_entity.Odds
}

type oddsEntityConverter struct{}

func NewOddsEntityConverter() OddsEntityConverter {
	return &oddsEntityConverter{}
}

func (o *oddsEntityConverter) DataCacheToRaw(input *data_cache_entity.Odds) *raw_entity.Odds {
	return &raw_entity.Odds{
		TicketType: input.TicketType().Value(),
		Odds:       input.Odds(),
		Popular:    input.PopularNumber(),
		Number:     input.Number().String(),
	}
}

func (o *oddsEntityConverter) RawToDataCache(input *raw_entity.Odds, raceId types.RaceId, raceDate types.RaceDate) *data_cache_entity.Odds {
	return data_cache_entity.NewOdds(
		raceId,
		raceDate,
		types.TicketType(input.TicketType),
		types.BetNumber(input.Number),
		input.Popular,
		input.Odds,
	)
}

func (o *oddsEntityConverter) NetKeibaToRaw(input *netkeiba_entity.Odds) *raw_entity.Odds {
	numbers := input.HorseNumbers()
	strNumbers := make([]string, len(numbers))
	for i, number := range numbers {
		strNumbers[i] = strconv.Itoa(number.Value())
	}
	number := strings.Join(strNumbers, types.QuinellaSeparator)
	return &raw_entity.Odds{
		TicketType: input.TicketType().Value(),
		Odds:       input.Odds(),
		Popular:    input.PopularNumber(),
		Number:     number,
	}
}
