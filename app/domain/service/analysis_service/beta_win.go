package analysis_service

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/filter_service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
	"github.com/shopspring/decimal"
)

type BetaWin interface {
	Create(ctx context.Context, markers []*marker_csv_entity.AnalysisMarker, races []*data_cache_entity.Race) ([]*analysis_entity.BetaCalculable, error)
	Convert(ctx context.Context, calculables []*analysis_entity.BetaCalculable) error
}

type betaWinService struct {
	filterService filter_service.AnalysisFilter
}

func NewBetaWin(
	filterService filter_service.AnalysisFilter,
) BetaWin {
	return &betaWinService{
		filterService: filterService,
	}
}

func (b *betaWinService) Create(
	ctx context.Context,
	markers []*marker_csv_entity.AnalysisMarker,
	races []*data_cache_entity.Race,
) ([]*analysis_entity.BetaCalculable, error) {
	markerMap := converter.ConvertToMap(markers, func(marker *marker_csv_entity.AnalysisMarker) types.RaceId {
		return marker.RaceId()
	})

	var calculables []*analysis_entity.BetaCalculable
	for _, race := range races {
		raceResultMap := converter.ConvertToMap(race.RaceResults(), func(raceResult *data_cache_entity.RaceResult) types.HorseNumber {
			return raceResult.HorseNumber()
		})
		payoutResultMap := converter.ConvertToMap(race.PayoutResults(), func(payoutResult *data_cache_entity.PayoutResult) types.TicketType {
			return payoutResult.TicketType().OriginTicketType()
		})
		payoutResult := payoutResultMap[types.Win]

		marker, ok := markerMap[race.RaceId()]
		if !ok {
			continue
		}

		markerHorseNumberMap := map[types.HorseNumber]struct{}{}
		for _, number := range marker.MarkerMap() {
			markerHorseNumberMap[number] = struct{}{}
		}

		isRaceCanceled := false
		for _, raceResult := range race.RaceResults() {
			if _, ok = markerHorseNumberMap[raceResult.HorseNumber()]; ok && raceResult.Odds().IsZero() {
				isRaceCanceled = true
				break
			}
		}

		// 取り消しの馬かつ、印対象だった場合そのレースは集計対象外
		if isRaceCanceled {
			continue
		}

		// 着順が1着のものを抽出(同着を考慮して複数保持)
		numbers := make([]types.HorseNumber, 0)
		for _, raceResult := range race.RaceResults() {
			if raceResult.OrderNo() > 1 {
				break
			}
			numbers = append(numbers, raceResult.HorseNumber())
		}

		markerCombinationIds := b.getHitMarkerCombinationIds(numbers, marker)
		filters := b.filterService.CreateBetaFilters(ctx, race, markerCombinationIds)

		raceCalculables := make([]*analysis_entity.BetaCalculable, 0, len(markerCombinationIds))
		for i, number := range numbers {
			raceResult, ok := raceResultMap[number]
			if !ok {
				return nil, fmt.Errorf("horseNumber %v not found in raceId %v", number, race.RaceId())
			}

			for j, hitNumber := range payoutResult.Numbers() {
				// 同着の場合、馬番が一致しないのでスキップ
				if hitNumber.List()[0] != number.Value() {
					continue
				}

				rawOdds := payoutResult.Odds()[j]
				odds, err := decimal.NewFromString(rawOdds)
				if err != nil {
					return nil, err
				}

				raceCalculables = append(raceCalculables, analysis_entity.NewBetaCalculable(
					race.RaceId(),
					race.RaceDate(),
					markerCombinationIds[i],
					odds,
					hitNumber,
					raceResult.PopularNumber(),
					race.Entries(),
					raceResult.JockeyId(),
					filters,
				))
			}
		}

		calculables = append(calculables, raceCalculables...)
	}

	return calculables, nil
}

func (b *betaWinService) Convert(
	ctx context.Context,
	calculables []*analysis_entity.BetaCalculable,
) error {
	filterBetaWinMap := map[filter.AttributeId]map[types.MarkerCombinationId]*spreadsheet_entity.AnalysisBetaRate{}
	analysisFilters := b.getFilters()
	for _, analysisFilter := range analysisFilters {
		filterBetaWinMap[analysisFilter] = map[types.MarkerCombinationId]*spreadsheet_entity.AnalysisBetaRate{}
		markerFilteredOddsMap := map[types.MarkerCombinationId][]decimal.Decimal{}

		raceCount := 0
		for _, calculable := range calculables {
			var calcFilter filter.AttributeId
			for _, f := range calculable.Filters() {
				calcFilter |= f
			}
			if analysisFilter == filter.All || analysisFilter&calcFilter == analysisFilter {
				if _, ok := markerFilteredOddsMap[calculable.MarkerCombinationId()]; !ok {
					markerFilteredOddsMap[calculable.MarkerCombinationId()] = []decimal.Decimal{}
				}
				markerFilteredOddsMap[calculable.MarkerCombinationId()] = append(markerFilteredOddsMap[calculable.MarkerCombinationId()], calculable.Odds())
				raceCount++
			}
		}

		for markerCombinationId, odds := range markerFilteredOddsMap {
			filterBetaWinMap[analysisFilter][markerCombinationId] = spreadsheet_entity.NewAnalysisBetaRate(odds, raceCount)
		}

		_ = markerFilteredOddsMap
	}

	return nil
}

func (b *betaWinService) getHitMarkerCombinationIds(
	numbers []types.HorseNumber,
	marker *marker_csv_entity.AnalysisMarker,
) []types.MarkerCombinationId {
	var markerCombinationIds []types.MarkerCombinationId
	for _, number := range numbers {
		markerCombinationId, _ := types.NewMarkerCombinationId(19)
		switch number {
		case marker.Favorite():
			markerCombinationId, _ = types.NewMarkerCombinationId(11)
		case marker.Rival():
			markerCombinationId, _ = types.NewMarkerCombinationId(12)
		case marker.BrackTriangle():
			markerCombinationId, _ = types.NewMarkerCombinationId(13)
		case marker.WhiteTriangle():
			markerCombinationId, _ = types.NewMarkerCombinationId(14)
		case marker.Star():
			markerCombinationId, _ = types.NewMarkerCombinationId(15)
		case marker.Check():
			markerCombinationId, _ = types.NewMarkerCombinationId(16)
		}
		markerCombinationIds = append(markerCombinationIds, markerCombinationId)
	}

	return markerCombinationIds
}

func (b *betaWinService) getUnHitMarkerCombinationIds(
	numbers []types.HorseNumber,
	marker *marker_csv_entity.AnalysisMarker,
) []types.MarkerCombinationId {
	unHitMarkerCombinationIdMap := map[types.MarkerCombinationId]bool{
		types.MarkerCombinationId(11): true,
		types.MarkerCombinationId(12): true,
		types.MarkerCombinationId(13): true,
		types.MarkerCombinationId(14): true,
		types.MarkerCombinationId(15): true,
		types.MarkerCombinationId(16): true,
		types.MarkerCombinationId(19): true,
	}

	for _, number := range numbers {
		switch number {
		case marker.Favorite():
			unHitMarkerCombinationIdMap[types.MarkerCombinationId(11)] = false
		case marker.Rival():
			unHitMarkerCombinationIdMap[types.MarkerCombinationId(12)] = false
		case marker.BrackTriangle():
			unHitMarkerCombinationIdMap[types.MarkerCombinationId(13)] = false
		case marker.WhiteTriangle():
			unHitMarkerCombinationIdMap[types.MarkerCombinationId(14)] = false
		case marker.Star():
			unHitMarkerCombinationIdMap[types.MarkerCombinationId(15)] = false
		case marker.Check():
			unHitMarkerCombinationIdMap[types.MarkerCombinationId(16)] = false
		}
	}

	var unHitMarkerCombinationIds []types.MarkerCombinationId
	for markerCombinationId, unHit := range unHitMarkerCombinationIdMap {
		if unHit {
			unHitMarkerCombinationIds = append(unHitMarkerCombinationIds, markerCombinationId)
		}
	}

	return unHitMarkerCombinationIds
}

func (b *betaWinService) getFilters() []filter.AttributeId {
	return []filter.AttributeId{
		filter.All,
		filter.Turf,
		filter.Dirt,
	}
}
