package analysis_service

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/filter_service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
)

type PlaceAllIn interface {
	Create(ctx context.Context, markers []*marker_csv_entity.AnalysisMarker, races []*data_cache_entity.Race, winOdds []*data_cache_entity.Odds, placeOdds []*data_cache_entity.Odds) ([]*analysis_entity.PlaceAllInCalculable, error)
	Convert(ctx context.Context, calculables []*analysis_entity.PlaceAllInCalculable) (map[filter.Id]*spreadsheet_entity.AnalysisPlaceAllIn, []filter.Id)
	Write(ctx context.Context, placeAllInMap map[filter.Id]*spreadsheet_entity.AnalysisPlaceAllIn, filters []filter.Id) error
}

type placeAllInService struct {
	filterService         filter_service.AnalysisFilter
	spreadSheetRepository repository.SpreadSheetRepository
}

func NewPlaceAllIn(
	filterService filter_service.AnalysisFilter,
	spreadSheetRepository repository.SpreadSheetRepository,
) PlaceAllIn {
	return &placeAllInService{
		filterService:         filterService,
		spreadSheetRepository: spreadSheetRepository,
	}
}

func (p *placeAllInService) Create(
	ctx context.Context,
	markers []*marker_csv_entity.AnalysisMarker,
	races []*data_cache_entity.Race,
	winOdds []*data_cache_entity.Odds,
	placeOdds []*data_cache_entity.Odds,
) ([]*analysis_entity.PlaceAllInCalculable, error) {
	// raceIdに対する、odds倍以下の複勝オッズに該当するデータを抽出
	// 抽出項目：複勝オッズ、単勝オッズ
	markerMap := converter.ConvertToMap(markers, func(marker *marker_csv_entity.AnalysisMarker) types.RaceId {
		return marker.RaceId()
	})

	raceMap := map[types.RaceId]*data_cache_entity.Race{}
	for _, race := range races {
		if _, ok := markerMap[race.RaceId()]; !ok {
			switch race.Class() {
			case types.MakeDebut, types.JumpMaiden, types.JumpGrade1, types.JumpGrade2, types.JumpGrade3:
				// 新馬・障害は分析印なしなのでスキップ
			default:
				// 印が不完全な場合がたまにあり(同じ印がついていたり、取り消しによる印6個未満の場合)、その場合はスキップ
				// log.Println(fmt.Sprintf("raceId not found in place markers: %s", race.RaceId()))
			}
			continue
		}
		raceMap[race.RaceId()] = race
	}

	placeOddsMap := map[types.RaceId]*data_cache_entity.Odds{}
	for _, o := range placeOdds {
		placeOddsMap[o.RaceId()] = o
	}

	var raceCalculables []*analysis_entity.PlaceAllInCalculable
	for _, wo := range winOdds {
		race, ok := raceMap[wo.RaceId()]
		if !ok {
			return nil, fmt.Errorf("race not found in placeAllInService.Create: %s", wo.RaceId())
		}

		po, ok := placeOddsMap[wo.RaceId()]
		if !ok {
			return nil, fmt.Errorf("placeOdds not found in placeAllInService.Create: %s", wo.RaceId())
		}

		marker, ok := markerMap[wo.RaceId()]
		if !ok {
			return nil, fmt.Errorf("marker not found in placeAllInService: %s", wo.RaceId())
		}

		horseNumber := types.HorseNumber(wo.Number().List()[0])
		markerCombinationId := p.getMarkerCombinationId(horseNumber, marker)

		var fixedPlaceOdds string
		for _, payoutResult := range race.PayoutResults() {
			if payoutResult.TicketType() == types.Place {
				for idx, number := range payoutResult.Numbers() {
					if types.HorseNumber(number.List()[0]) == horseNumber {
						fixedPlaceOdds = payoutResult.Odds()[idx]
					}
				}
			}
		}

		var (
			orderNo  int
			jockeyId types.JockeyId
		)
		for _, raceResult := range race.RaceResults() {
			if raceResult.HorseNumber() == horseNumber {
				orderNo = raceResult.OrderNo()
				jockeyId = raceResult.JockeyId()
			}
		}

		filters := p.filterService.CreatePlaceAllInFilters(ctx, race)

		raceCalculables = append(raceCalculables, analysis_entity.NewPlaceAllInCalculable(
			race.RaceId(),
			race.RaceDate(),
			markerCombinationId,
			wo.Odds()[0],
			po.Odds(),
			fixedPlaceOdds,
			wo.PopularNumber(),
			orderNo,
			race.Entries(),
			race.TrackCondition(),
			race.RaceCourseId(),
			race.Distance(),
			jockeyId,
			filters,
		))
	}

	return raceCalculables, nil
}

func (p *placeAllInService) Convert(
	ctx context.Context,
	calculables []*analysis_entity.PlaceAllInCalculable,
) (map[filter.Id]*spreadsheet_entity.AnalysisPlaceAllIn, []filter.Id) {
	isHit := func(calculable *analysis_entity.PlaceAllInCalculable) bool {
		return calculable.Entries() <= 7 && calculable.OrderNo() <= 2 || calculable.Entries() >= 8 && calculable.OrderNo() <= 3
	}

	filterPlaceAllInMap := map[filter.Id]*spreadsheet_entity.AnalysisPlaceAllIn{}
	analysisFilters := p.getFilters()
	for _, analysisFilter := range analysisFilters {
		raceIdMap := map[types.RaceId]bool{}
		winOddsHitCountSlice := make([]int, 29)
		winOddsUnHitCountSlice := make([]int, 29)
		for _, calculable := range calculables {
			match := false
			for _, f := range calculable.Filters() {
				if f&analysisFilter > 0 {
					match = true
					break
				}
			}
			if match {
				if _, ok := raceIdMap[calculable.RaceId()]; !ok {
					raceIdMap[calculable.RaceId()] = true
				}
				odds := calculable.WinOdds().InexactFloat64()
				switch odds {
				case 1.1:
					if isHit(calculable) {
						winOddsHitCountSlice[0]++
					} else {
						winOddsUnHitCountSlice[0]++
					}
				case 1.2:
					if isHit(calculable) {
						winOddsHitCountSlice[1]++
					} else {
						winOddsUnHitCountSlice[1]++
					}
				case 1.3:
					if isHit(calculable) {
						winOddsHitCountSlice[2]++
					} else {
						winOddsUnHitCountSlice[2]++
					}
				case 1.4:
					if isHit(calculable) {
						winOddsHitCountSlice[3]++
					} else {
						winOddsUnHitCountSlice[3]++
					}
				case 1.5:
					if isHit(calculable) {
						winOddsHitCountSlice[4]++
					} else {
						winOddsUnHitCountSlice[4]++
					}
				case 1.6:
					if isHit(calculable) {
						winOddsHitCountSlice[5]++
					} else {
						winOddsUnHitCountSlice[5]++
					}
				case 1.7:
					if isHit(calculable) {
						winOddsHitCountSlice[6]++
					} else {
						winOddsUnHitCountSlice[6]++
					}
				case 1.8:
					if isHit(calculable) {
						winOddsHitCountSlice[7]++
					} else {
						winOddsUnHitCountSlice[7]++
					}
				case 1.9:
					if isHit(calculable) {
						winOddsHitCountSlice[8]++
					} else {
						winOddsUnHitCountSlice[8]++
					}
				case 2.0:
					if isHit(calculable) {
						winOddsHitCountSlice[9]++
					} else {
						winOddsUnHitCountSlice[9]++
					}
				case 2.1:
					if isHit(calculable) {
						winOddsHitCountSlice[10]++
					} else {
						winOddsUnHitCountSlice[10]++
					}
				case 2.2:
					if isHit(calculable) {
						winOddsHitCountSlice[11]++
					} else {
						winOddsUnHitCountSlice[11]++
					}
				case 2.3:
					if isHit(calculable) {
						winOddsHitCountSlice[12]++
					} else {
						winOddsUnHitCountSlice[12]++
					}
				case 2.4:
					if isHit(calculable) {
						winOddsHitCountSlice[13]++
					} else {
						winOddsUnHitCountSlice[13]++
					}
				case 2.5:
					if isHit(calculable) {
						winOddsHitCountSlice[14]++
					} else {
						winOddsUnHitCountSlice[14]++
					}
				case 2.6:
					if isHit(calculable) {
						winOddsHitCountSlice[15]++
					} else {
						winOddsUnHitCountSlice[15]++
					}
				case 2.7:
					if isHit(calculable) {
						winOddsHitCountSlice[16]++
					} else {
						winOddsUnHitCountSlice[16]++
					}
				case 2.8:
					if isHit(calculable) {
						winOddsHitCountSlice[17]++
					} else {
						winOddsUnHitCountSlice[17]++
					}
				case 2.9:
					if isHit(calculable) {
						winOddsHitCountSlice[18]++
					} else {
						winOddsUnHitCountSlice[18]++
					}
				case 3.0:
					if isHit(calculable) {
						winOddsHitCountSlice[19]++
					} else {
						winOddsUnHitCountSlice[19]++
					}
				case 3.1:
					if isHit(calculable) {
						winOddsHitCountSlice[20]++
					} else {
						winOddsUnHitCountSlice[20]++
					}
				case 3.2:
					if isHit(calculable) {
						winOddsHitCountSlice[21]++
					} else {
						winOddsUnHitCountSlice[21]++
					}
				case 3.3:
					if isHit(calculable) {
						winOddsHitCountSlice[22]++
					} else {
						winOddsUnHitCountSlice[22]++
					}
				case 3.4:
					if isHit(calculable) {
						winOddsHitCountSlice[23]++
					} else {
						winOddsUnHitCountSlice[23]++
					}
				case 3.5:
					if isHit(calculable) {
						winOddsHitCountSlice[24]++
					} else {
						winOddsUnHitCountSlice[24]++
					}
				case 3.6:
					if isHit(calculable) {
						winOddsHitCountSlice[25]++
					} else {
						winOddsUnHitCountSlice[25]++
					}
				case 3.7:
					if isHit(calculable) {
						winOddsHitCountSlice[26]++
					} else {
						winOddsUnHitCountSlice[26]++
					}
				case 3.8:
					if isHit(calculable) {
						winOddsHitCountSlice[27]++
					} else {
						winOddsUnHitCountSlice[27]++
					}
				case 3.9:
					if isHit(calculable) {
						winOddsHitCountSlice[28]++
					} else {
						winOddsUnHitCountSlice[28]++
					}
				}
			}
		}

		placeAllInHitCountData := spreadsheet_entity.NewPlaceAllInHitCountData(
			winOddsHitCountSlice,
			analysisFilter,
			len(raceIdMap),
		)

		placeAllInUnHitCountData := spreadsheet_entity.NewPlaceAllInUnHitCountData(
			winOddsUnHitCountSlice,
			analysisFilter,
			len(raceIdMap),
		)

		placeAllInRateData := spreadsheet_entity.NewPlaceAllInRateData(placeAllInHitCountData, placeAllInUnHitCountData)
		placeAllInRateStyle := spreadsheet_entity.NewPlaceAllInRateStyle(placeAllInRateData)
		filterPlaceAllInMap[analysisFilter] = spreadsheet_entity.NewAnalysisPlaceAllIn(
			placeAllInRateData,
			placeAllInRateStyle,
		)
	}

	return filterPlaceAllInMap, analysisFilters
}

func (p *placeAllInService) Write(
	ctx context.Context,
	placeAllInMap map[filter.Id]*spreadsheet_entity.AnalysisPlaceAllIn,
	filters []filter.Id,
) error {
	return p.spreadSheetRepository.WriteAnalysisPlaceAllIn(ctx, placeAllInMap, filters)
}

func (p *placeAllInService) getMarkerCombinationId(
	number types.HorseNumber,
	marker *marker_csv_entity.AnalysisMarker,
) types.MarkerCombinationId {
	markerCombinationId, _ := types.NewMarkerCombinationId(29)
	switch number {
	case marker.Favorite():
		markerCombinationId, _ = types.NewMarkerCombinationId(21)
	case marker.Rival():
		markerCombinationId, _ = types.NewMarkerCombinationId(22)
	case marker.BrackTriangle():
		markerCombinationId, _ = types.NewMarkerCombinationId(23)
	case marker.WhiteTriangle():
		markerCombinationId, _ = types.NewMarkerCombinationId(24)
	case marker.Star():
		markerCombinationId, _ = types.NewMarkerCombinationId(25)
	case marker.Check():
		markerCombinationId, _ = types.NewMarkerCombinationId(26)
	}

	return markerCombinationId
}

func (p *placeAllInService) getFilters() []filter.Id {
	return []filter.Id{
		filter.All2,
		filter.Turf2,
		filter.Dirt2,
		filter.GoodToFirm,
		filter.Good,
		filter.Yielding,
		filter.Soft,
		filter.Maiden,
		filter.OneWinClass,
		filter.TwoWinClass,
		filter.ThreeWinClass,
		filter.OpenListed,
		filter.Grade3,
		filter.Grade2,
		filter.Grade1,
		filter.NiigataTurf1000m,
		filter.HakodateTurf1000m,
		filter.NakayamaTurf1200m,
		filter.KyotoTurf1200m,
		filter.HanshinTurf1200m,
		filter.NiigataTurf1200m,
		filter.ChukyoTurf1200m,
		filter.SapporoTurf1200m,
		filter.HakodateTurf1200m,
		filter.FukushimaTurf1200m,
		filter.KokuraTurf1200m,
		filter.TokyoTurf1400m,
		filter.KyotoTurf1400m,
		filter.HanshinTurf1400m,
		filter.NiigataTurf1400m,
		filter.ChukyoTurf1400m,
		filter.SapporoTurf1500m,
		filter.NakayamaTurf1600m,
		filter.TokyoTurf1600m,
		filter.KyotoTurf1600m,
		filter.HanshinTurf1600m,
		filter.ChukyoTurf1600m,
		filter.NakayamaTurf1800m,
		filter.TokyoTurf1800m,
		filter.KyotoTurf1800m,
		filter.HanshinTurf1800m,
		filter.NiigataTurf1800m,
		filter.SapporoTurf1800m,
		filter.HakodateTurf1800m,
		filter.FukushimaTurf1800m,
		filter.KokuraTurf1800m,
		filter.NakayamaTurf2000m,
		filter.TokyoTurf2000m,
		filter.KyotoTurf2000m,
		filter.NiigataTurf2000m,
		filter.ChukyoTurf2000m,
		filter.SapporoTurf2000m,
		filter.HakodateTurf2000m,
		filter.FukushimaTurf2000m,
		filter.KokuraTurf2000m,
		filter.NakayamaTurf2200m,
		filter.KyotoTurf2200m,
		filter.HanshinTurf2200m,
		filter.NiigataTurf2200m,
		filter.ChukyoTurf2200m,
		filter.TokyoTurf2300m,
		filter.TokyoTurf2400m,
		filter.KyotoTurf2400m,
		filter.HanshinTurf2400m,
		filter.NiigataTurf2400m,
		filter.NakayamaTurf2500m,
		filter.TokyoTurf2500m,
		filter.HanshinTurf2600m,
		filter.SapporoTurf2600m,
		filter.HakodateTurf2600m,
		filter.FukushimaTurf2600m,
		filter.KokuraTurf2600m,
		filter.HanshinTurf3000m,
		filter.ChukyoTurf3000m,
		filter.KyotoTurf3200m,
		filter.TokyoTurf3400m,
		filter.NakayamaTurf3600m,
		filter.SapporoDirt1000m,
		filter.HakodateDirt1000m,
		filter.KokuraDirt1000m,
		filter.FukushimaDirt1150m,
		filter.NakayamaDirt1200m,
		filter.KyotoDirt1200m,
		filter.NiigataDirt1200m,
		filter.ChukyoDirt1200m,
		filter.TokyoDirt1300m,
		filter.TokyoDirt1400m,
		filter.KyotoDirt1400m,
		filter.HanshinDirt1400m,
		filter.ChukyoDirt1400m,
		filter.TokyoDirt1600m,
		filter.SapporoDirt1700m,
		filter.HakodateDirt1700m,
		filter.FukushimaDirt1700m,
		filter.KokuraDirt1700m,
		filter.NakayamaDirt1800m,
		filter.KyotoDirt1800m,
		filter.HanshinDirt1800m,
		filter.NiigataDirt1800m,
		filter.ChukyoDirt1800m,
		filter.KyotoDirt1900m,
		filter.ChukyoDirt1900m,
		filter.HanshinDirt2000m,
		filter.TokyoDirt2100m,
		filter.NakayamaDirt2400m,
		filter.SapporoDirt2400m,
		filter.HakodateDirt2400m,
		filter.FukushimaDirt2400m,
		filter.KokuraDirt2400m,
		filter.NakayamaDirt2500m,
		filter.NiigataDirt2500m,
	}
}
