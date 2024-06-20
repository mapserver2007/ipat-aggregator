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
		filter.NiigataTurf1000m,
		filter.NiigataGoodToFirmTurf1000m,
		filter.NiigataGoodTurf1000m,
		filter.NiigataYieldingTurf1000m,
		filter.NiigataSoftTurf1000m,
		filter.HakodateTurf1000m,
		filter.HakodateGoodToFirmTurf1000m,
		filter.HakodateGoodTurf1000m,
		filter.HakodateYieldingTurf1000m,
		filter.HakodateSoftTurf1000m,
		filter.NakayamaTurf1200m,
		filter.NakayamaGoodToFirmTurf1200m,
		filter.NakayamaGoodTurf1200m,
		filter.NakayamaYieldingTurf1200m,
		filter.NakayamaSoftTurf1200m,
		filter.KyotoTurf1200m,
		filter.KyotoGoodToFirmTurf1200m,
		filter.KyotoGoodTurf1200m,
		filter.KyotoYieldingTurf1200m,
		filter.KyotoSoftTurf1200m,
		filter.HanshinTurf1200m,
		filter.HanshinGoodToFirmTurf1200m,
		filter.HanshinGoodTurf1200m,
		filter.HanshinYieldingTurf1200m,
		filter.HanshinSoftTurf1200m,
		filter.NiigataTurf1200m,
		filter.NiigataGoodToFirmTurf1200m,
		filter.NiigataGoodTurf1200m,
		filter.NiigataYieldingTurf1200m,
		filter.NiigataSoftTurf1200m,
		filter.ChukyoTurf1200m,
		filter.ChukyoGoodToFirmTurf1200m,
		filter.ChukyoGoodTurf1200m,
		filter.ChukyoYieldingTurf1200m,
		filter.ChukyoSoftTurf1200m,
		filter.SapporoTurf1200m,
		filter.SapporoGoodToFirmTurf1200m,
		filter.SapporoGoodTurf1200m,
		filter.SapporoYieldingTurf1200m,
		filter.SapporoSoftTurf1200m,
		filter.HakodateTurf1200m,
		filter.HakodateGoodToFirmTurf1200m,
		filter.HakodateGoodTurf1200m,
		filter.HakodateYieldingTurf1200m,
		filter.HakodateSoftTurf1200m,
		filter.FukushimaTurf1200m,
		filter.FukushimaGoodToFirmTurf1200m,
		filter.FukushimaGoodTurf1200m,
		filter.FukushimaYieldingTurf1200m,
		filter.FukushimaSoftTurf1200m,
		filter.KokuraTurf1200m,
		filter.KokuraGoodToFirmTurf1200m,
		filter.KokuraGoodTurf1200m,
		filter.KokuraYieldingTurf1200m,
		filter.KokuraSoftTurf1200m,
		filter.TokyoTurf1400m,
		filter.TokyoGoodToFirmTurf1400m,
		filter.TokyoGoodTurf1400m,
		filter.TokyoYieldingTurf1400m,
		filter.TokyoSoftTurf1400m,
		filter.KyotoTurf1400m,
		filter.KyotoGoodToFirmTurf1400m,
		filter.KyotoGoodTurf1400m,
		filter.KyotoYieldingTurf1400m,
		filter.KyotoSoftTurf1400m,
		filter.HanshinTurf1400m,
		filter.HanshinGoodToFirmTurf1400m,
		filter.HanshinGoodTurf1400m,
		filter.HanshinYieldingTurf1400m,
		filter.HanshinSoftTurf1400m,
		filter.NiigataTurf1400m,
		filter.NiigataGoodToFirmTurf1400m,
		filter.NiigataGoodTurf1400m,
		filter.NiigataYieldingTurf1400m,
		filter.NiigataSoftTurf1400m,
		filter.ChukyoTurf1400m,
		filter.ChukyoGoodToFirmTurf1400m,
		filter.ChukyoGoodTurf1400m,
		filter.ChukyoYieldingTurf1400m,
		filter.ChukyoSoftTurf1400m,
		filter.SapporoTurf1500m,
		filter.SapporoGoodToFirmTurf1500m,
		filter.SapporoGoodTurf1500m,
		filter.SapporoYieldingTurf1500m,
		filter.SapporoSoftTurf1500m,
		filter.NakayamaTurf1600m,
		filter.NakayamaGoodToFirmTurf1600m,
		filter.NakayamaGoodTurf1600m,
		filter.NakayamaYieldingTurf1600m,
		filter.NakayamaSoftTurf1600m,
		filter.TokyoTurf1600m,
		filter.TokyoGoodToFirmTurf1600m,
		filter.TokyoGoodTurf1600m,
		filter.TokyoYieldingTurf1600m,
		filter.TokyoSoftTurf1600m,
		filter.KyotoTurf1600m,
		filter.KyotoGoodToFirmTurf1600m,
		filter.KyotoGoodTurf1600m,
		filter.KyotoYieldingTurf1600m,
		filter.KyotoSoftTurf1600m,
		filter.HanshinTurf1600m,
		filter.HanshinGoodToFirmTurf1600m,
		filter.HanshinGoodTurf1600m,
		filter.HanshinYieldingTurf1600m,
		filter.HanshinSoftTurf1600m,
		filter.ChukyoTurf1600m,
		filter.ChukyoGoodToFirmTurf1600m,
		filter.ChukyoGoodTurf1600m,
		filter.ChukyoYieldingTurf1600m,
		filter.ChukyoSoftTurf1600m,
		filter.NakayamaTurf1800m,
		filter.NakayamaGoodToFirmTurf1800m,
		filter.NakayamaGoodTurf1800m,
		filter.NakayamaYieldingTurf1800m,
		filter.NakayamaSoftTurf1800m,
		filter.TokyoTurf1800m,
		filter.TokyoGoodToFirmTurf1800m,
		filter.TokyoGoodTurf1800m,
		filter.TokyoYieldingTurf1800m,
		filter.TokyoSoftTurf1800m,
		filter.KyotoTurf1800m,
		filter.KyotoGoodToFirmTurf1800m,
		filter.KyotoGoodTurf1800m,
		filter.KyotoYieldingTurf1800m,
		filter.KyotoSoftTurf1800m,
		filter.HanshinTurf1800m,
		filter.HanshinGoodToFirmTurf1800m,
		filter.HanshinGoodTurf1800m,
		filter.HanshinYieldingTurf1800m,
		filter.HanshinSoftTurf1800m,
		filter.NiigataTurf1800m,
		filter.NiigataGoodToFirmTurf1800m,
		filter.NiigataGoodTurf1800m,
		filter.NiigataYieldingTurf1800m,
		filter.NiigataSoftTurf1800m,
		filter.SapporoTurf1800m,
		filter.SapporoGoodToFirmTurf1800m,
		filter.SapporoGoodTurf1800m,
		filter.SapporoYieldingTurf1800m,
		filter.SapporoSoftTurf1800m,
		filter.HakodateTurf1800m,
		filter.HakodateGoodToFirmTurf1800m,
		filter.HakodateGoodTurf1800m,
		filter.HakodateYieldingTurf1800m,
		filter.HakodateSoftTurf1800m,
		filter.FukushimaTurf1800m,
		filter.FukushimaGoodToFirmTurf1800m,
		filter.FukushimaGoodTurf1800m,
		filter.FukushimaYieldingTurf1800m,
		filter.FukushimaSoftTurf1800m,
		filter.KokuraTurf1800m,
		filter.KokuraGoodToFirmTurf1800m,
		filter.KokuraGoodTurf1800m,
		filter.KokuraYieldingTurf1800m,
		filter.KokuraSoftTurf1800m,
		filter.NakayamaTurf2000m,
		filter.NakayamaGoodToFirmTurf2000m,
		filter.NakayamaGoodTurf2000m,
		filter.NakayamaYieldingTurf2000m,
		filter.NakayamaSoftTurf2000m,
		filter.TokyoTurf2000m,
		filter.TokyoGoodToFirmTurf2000m,
		filter.TokyoGoodTurf2000m,
		filter.TokyoYieldingTurf2000m,
		filter.TokyoSoftTurf2000m,
		filter.KyotoTurf2000m,
		filter.KyotoGoodToFirmTurf2000m,
		filter.KyotoGoodTurf2000m,
		filter.KyotoYieldingTurf2000m,
		filter.KyotoSoftTurf2000m,
		filter.NiigataTurf2000m,
		filter.NiigataGoodToFirmTurf2000m,
		filter.NiigataGoodTurf2000m,
		filter.NiigataYieldingTurf2000m,
		filter.NiigataSoftTurf2000m,
		filter.ChukyoTurf2000m,
		filter.ChukyoGoodToFirmTurf2000m,
		filter.ChukyoGoodTurf2000m,
		filter.ChukyoYieldingTurf2000m,
		filter.ChukyoSoftTurf2000m,
		filter.SapporoTurf2000m,
		filter.SapporoGoodToFirmTurf2000m,
		filter.SapporoGoodTurf2000m,
		filter.SapporoYieldingTurf2000m,
		filter.SapporoSoftTurf2000m,
		filter.HakodateTurf2000m,
		filter.HakodateGoodToFirmTurf2000m,
		filter.HakodateGoodTurf2000m,
		filter.HakodateYieldingTurf2000m,
		filter.HakodateSoftTurf2000m,
		filter.FukushimaTurf2000m,
		filter.FukushimaGoodToFirmTurf2000m,
		filter.FukushimaGoodTurf2000m,
		filter.FukushimaYieldingTurf2000m,
		filter.FukushimaSoftTurf2000m,
		filter.KokuraTurf2000m,
		filter.KokuraGoodToFirmTurf2000m,
		filter.KokuraGoodTurf2000m,
		filter.KokuraYieldingTurf2000m,
		filter.KokuraSoftTurf2000m,
		filter.NakayamaTurf2200m,
		filter.NakayamaGoodToFirmTurf2200m,
		filter.NakayamaGoodTurf2200m,
		filter.NakayamaYieldingTurf2200m,
		filter.NakayamaSoftTurf2200m,
		filter.KyotoTurf2200m,
		filter.KyotoGoodToFirmTurf2200m,
		filter.KyotoGoodTurf2200m,
		filter.KyotoYieldingTurf2200m,
		filter.KyotoSoftTurf2200m,
		filter.HanshinTurf2200m,
		filter.HanshinGoodToFirmTurf2200m,
		filter.HanshinGoodTurf2200m,
		filter.HanshinYieldingTurf2200m,
		filter.HanshinSoftTurf2200m,
		filter.NiigataTurf2200m,
		filter.NiigataGoodToFirmTurf2200m,
		filter.NiigataGoodTurf2200m,
		filter.NiigataYieldingTurf2200m,
		filter.NiigataSoftTurf2200m,
		filter.ChukyoTurf2200m,
		filter.ChukyoGoodToFirmTurf2200m,
		filter.ChukyoGoodTurf2200m,
		filter.ChukyoYieldingTurf2200m,
		filter.ChukyoSoftTurf2200m,
		filter.TokyoTurf2300m,
		filter.TokyoGoodToFirmTurf2300m,
		filter.TokyoGoodTurf2300m,
		filter.TokyoYieldingTurf2300m,
		filter.TokyoSoftTurf2300m,
		filter.TokyoTurf2400m,
		filter.TokyoGoodToFirmTurf2400m,
		filter.TokyoGoodTurf2400m,
		filter.TokyoYieldingTurf2400m,
		filter.TokyoSoftTurf2400m,
		filter.KyotoTurf2400m,
		filter.KyotoGoodToFirmTurf2400m,
		filter.KyotoGoodTurf2400m,
		filter.KyotoYieldingTurf2400m,
		filter.KyotoSoftTurf2400m,
		filter.HanshinTurf2400m,
		filter.HanshinGoodToFirmTurf2400m,
		filter.HanshinGoodTurf2400m,
		filter.HanshinYieldingTurf2400m,
		filter.HanshinSoftTurf2400m,
		filter.NiigataTurf2400m,
		filter.NiigataGoodToFirmTurf2400m,
		filter.NiigataGoodTurf2400m,
		filter.NiigataYieldingTurf2400m,
		filter.NiigataSoftTurf2400m,
		filter.NakayamaTurf2500m,
		filter.NakayamaGoodToFirmTurf2500m,
		filter.NakayamaGoodTurf2500m,
		filter.NakayamaYieldingTurf2500m,
		filter.NakayamaSoftTurf2500m,
		filter.TokyoTurf2500m,
		filter.TokyoGoodToFirmTurf2500m,
		filter.TokyoGoodTurf2500m,
		filter.TokyoYieldingTurf2500m,
		filter.TokyoSoftTurf2500m,
		filter.HanshinTurf2600m,
		filter.HanshinGoodToFirmTurf2600m,
		filter.HanshinGoodTurf2600m,
		filter.HanshinYieldingTurf2600m,
		filter.HanshinSoftTurf2600m,
		filter.SapporoTurf2600m,
		filter.SapporoGoodToFirmTurf2600m,
		filter.SapporoGoodTurf2600m,
		filter.SapporoYieldingTurf2600m,
		filter.SapporoSoftTurf2600m,
		filter.HakodateTurf2600m,
		filter.HakodateGoodToFirmTurf2600m,
		filter.HakodateGoodTurf2600m,
		filter.HakodateYieldingTurf2600m,
		filter.HakodateSoftTurf2600m,
		filter.FukushimaTurf2600m,
		filter.FukushimaGoodToFirmTurf2600m,
		filter.FukushimaGoodTurf2600m,
		filter.FukushimaYieldingTurf2600m,
		filter.FukushimaSoftTurf2600m,
		filter.KokuraTurf2600m,
		filter.KokuraGoodToFirmTurf2600m,
		filter.KokuraGoodTurf2600m,
		filter.KokuraYieldingTurf2600m,
		filter.KokuraSoftTurf2600m,
		filter.HanshinTurf3000m,
		filter.HanshinGoodToFirmTurf3000m,
		filter.HanshinGoodTurf3000m,
		filter.HanshinYieldingTurf3000m,
		filter.HanshinSoftTurf3000m,
		filter.ChukyoTurf3000m,
		filter.ChukyoGoodToFirmTurf3000m,
		filter.ChukyoGoodTurf3000m,
		filter.ChukyoYieldingTurf3000m,
		filter.ChukyoSoftTurf3000m,
		filter.KyotoTurf3200m,
		filter.KyotoGoodToFirmTurf3200m,
		filter.KyotoGoodTurf3200m,
		filter.KyotoYieldingTurf3200m,
		filter.KyotoSoftTurf3200m,
		filter.TokyoTurf3400m,
		filter.TokyoGoodToFirmTurf3400m,
		filter.TokyoGoodTurf3400m,
		filter.TokyoYieldingTurf3400m,
		filter.TokyoSoftTurf3400m,
		filter.NakayamaTurf3600m,
		filter.NakayamaGoodToFirmTurf3600m,
		filter.NakayamaGoodTurf3600m,
		filter.NakayamaYieldingTurf3600m,
		filter.NakayamaSoftTurf3600m,
		filter.SapporoDirt1000m,
		filter.SapporoGoodToFirmDirt1000m,
		filter.SapporoGoodDirt1000m,
		filter.SapporoYieldingDirt1000m,
		filter.SapporoSoftDirt1000m,
		filter.HakodateDirt1000m,
		filter.HakodateGoodToFirmDirt1000m,
		filter.HakodateGoodDirt1000m,
		filter.HakodateYieldingDirt1000m,
		filter.HakodateSoftDirt1000m,
		filter.KokuraDirt1000m,
		filter.KokuraGoodToFirmDirt1000m,
		filter.KokuraGoodDirt1000m,
		filter.KokuraYieldingDirt1000m,
		filter.KokuraSoftDirt1000m,
		filter.FukushimaDirt1150m,
		filter.FukushimaGoodToFirmDirt1150m,
		filter.FukushimaGoodDirt1150m,
		filter.FukushimaYieldingDirt1150m,
		filter.FukushimaSoftDirt1150m,
		filter.NakayamaDirt1200m,
		filter.NakayamaGoodToFirmDirt1200m,
		filter.NakayamaGoodDirt1200m,
		filter.NakayamaYieldingDirt1200m,
		filter.NakayamaSoftDirt1200m,
		filter.KyotoDirt1200m,
		filter.KyotoGoodToFirmDirt1200m,
		filter.KyotoGoodDirt1200m,
		filter.KyotoYieldingDirt1200m,
		filter.KyotoSoftDirt1200m,
		filter.NiigataDirt1200m,
		filter.NiigataGoodToFirmDirt1200m,
		filter.NiigataGoodDirt1200m,
		filter.NiigataYieldingDirt1200m,
		filter.NiigataSoftDirt1200m,
		filter.ChukyoDirt1200m,
		filter.ChukyoGoodToFirmDirt1200m,
		filter.ChukyoGoodDirt1200m,
		filter.ChukyoYieldingDirt1200m,
		filter.ChukyoSoftDirt1200m,
		filter.TokyoDirt1300m,
		filter.TokyoGoodToFirmDirt1300m,
		filter.TokyoGoodDirt1300m,
		filter.TokyoYieldingDirt1300m,
		filter.TokyoSoftDirt1300m,
		filter.TokyoDirt1400m,
		filter.TokyoGoodToFirmDirt1400m,
		filter.TokyoGoodDirt1400m,
		filter.TokyoYieldingDirt1400m,
		filter.TokyoSoftDirt1400m,
		filter.KyotoDirt1400m,
		filter.KyotoGoodToFirmDirt1400m,
		filter.KyotoGoodDirt1400m,
		filter.KyotoYieldingDirt1400m,
		filter.KyotoSoftDirt1400m,
		filter.HanshinDirt1400m,
		filter.HanshinGoodToFirmDirt1400m,
		filter.HanshinGoodDirt1400m,
		filter.HanshinYieldingDirt1400m,
		filter.HanshinSoftDirt1400m,
		filter.ChukyoDirt1400m,
		filter.ChukyoGoodToFirmDirt1400m,
		filter.ChukyoGoodDirt1400m,
		filter.ChukyoYieldingDirt1400m,
		filter.ChukyoSoftDirt1400m,
		filter.TokyoDirt1600m,
		filter.TokyoGoodToFirmDirt1600m,
		filter.TokyoGoodDirt1600m,
		filter.TokyoYieldingDirt1600m,
		filter.TokyoSoftDirt1600m,
		filter.SapporoDirt1700m,
		filter.SapporoGoodToFirmDirt1700m,
		filter.SapporoGoodDirt1700m,
		filter.SapporoYieldingDirt1700m,
		filter.SapporoSoftDirt1700m,
		filter.HakodateDirt1700m,
		filter.HakodateGoodToFirmDirt1700m,
		filter.HakodateGoodDirt1700m,
		filter.HakodateYieldingDirt1700m,
		filter.HakodateSoftDirt1700m,
		filter.FukushimaDirt1700m,
		filter.FukushimaGoodToFirmDirt1700m,
		filter.FukushimaGoodDirt1700m,
		filter.FukushimaYieldingDirt1700m,
		filter.FukushimaSoftDirt1700m,
		filter.KokuraDirt1700m,
		filter.KokuraGoodToFirmDirt1700m,
		filter.KokuraGoodDirt1700m,
		filter.KokuraYieldingDirt1700m,
		filter.KokuraSoftDirt1700m,
		filter.NakayamaDirt1800m,
		filter.NakayamaGoodToFirmDirt1800m,
		filter.NakayamaGoodDirt1800m,
		filter.NakayamaYieldingDirt1800m,
		filter.NakayamaSoftDirt1800m,
		filter.KyotoDirt1800m,
		filter.KyotoGoodToFirmDirt1800m,
		filter.KyotoGoodDirt1800m,
		filter.KyotoYieldingDirt1800m,
		filter.KyotoSoftDirt1800m,
		filter.HanshinDirt1800m,
		filter.HanshinGoodToFirmDirt1800m,
		filter.HanshinGoodDirt1800m,
		filter.HanshinYieldingDirt1800m,
		filter.HanshinSoftDirt1800m,
		filter.NiigataDirt1800m,
		filter.NiigataGoodToFirmDirt1800m,
		filter.NiigataGoodDirt1800m,
		filter.NiigataYieldingDirt1800m,
		filter.NiigataSoftDirt1800m,
		filter.ChukyoDirt1800m,
		filter.ChukyoGoodToFirmDirt1800m,
		filter.ChukyoGoodDirt1800m,
		filter.ChukyoYieldingDirt1800m,
		filter.ChukyoSoftDirt1800m,
		filter.KyotoDirt1900m,
		filter.KyotoGoodToFirmDirt1900m,
		filter.KyotoGoodDirt1900m,
		filter.KyotoYieldingDirt1900m,
		filter.KyotoSoftDirt1900m,
		filter.ChukyoDirt1900m,
		filter.ChukyoGoodToFirmDirt1900m,
		filter.ChukyoGoodDirt1900m,
		filter.ChukyoYieldingDirt1900m,
		filter.ChukyoSoftDirt1900m,
		filter.HanshinDirt2000m,
		filter.HanshinGoodToFirmDirt2000m,
		filter.HanshinGoodDirt2000m,
		filter.HanshinYieldingDirt2000m,
		filter.HanshinSoftDirt2000m,
		filter.TokyoDirt2100m,
		filter.TokyoGoodToFirmDirt2100m,
		filter.TokyoGoodDirt2100m,
		filter.TokyoYieldingDirt2100m,
		filter.TokyoSoftDirt2100m,
		filter.NakayamaDirt2400m,
		filter.NakayamaGoodToFirmDirt2400m,
		filter.NakayamaGoodDirt2400m,
		filter.NakayamaYieldingDirt2400m,
		filter.NakayamaSoftDirt2400m,
		filter.SapporoDirt2400m,
		filter.SapporoGoodToFirmDirt2400m,
		filter.SapporoGoodDirt2400m,
		filter.SapporoYieldingDirt2400m,
		filter.SapporoSoftDirt2400m,
		filter.HakodateDirt2400m,
		filter.HakodateGoodToFirmDirt2400m,
		filter.HakodateGoodDirt2400m,
		filter.HakodateYieldingDirt2400m,
		filter.HakodateSoftDirt2400m,
		filter.FukushimaDirt2400m,
		filter.FukushimaGoodToFirmDirt2400m,
		filter.FukushimaGoodDirt2400m,
		filter.FukushimaYieldingDirt2400m,
		filter.FukushimaSoftDirt2400m,
		filter.KokuraDirt2400m,
		filter.KokuraGoodToFirmDirt2400m,
		filter.KokuraGoodDirt2400m,
		filter.KokuraYieldingDirt2400m,
		filter.KokuraSoftDirt2400m,
		filter.NakayamaDirt2500m,
		filter.NakayamaGoodToFirmDirt2500m,
		filter.NakayamaGoodDirt2500m,
		filter.NakayamaYieldingDirt2500m,
		filter.NakayamaSoftDirt2500m,
		filter.NiigataDirt2500m,
		filter.NiigataGoodToFirmDirt2500m,
		filter.NiigataGoodDirt2500m,
		filter.NiigataYieldingDirt2500m,
		filter.NiigataSoftDirt2500m,
	}
}
