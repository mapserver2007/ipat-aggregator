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
	"log"
	"strconv"
)

type Place interface {
	Create(ctx context.Context, markers []*marker_csv_entity.AnalysisMarker, races []*data_cache_entity.Race) ([]*analysis_entity.PlaceCalculable, error)
	Convert(ctx context.Context, calculables []*analysis_entity.PlaceCalculable) (map[types.Marker]map[filter.Id]*spreadsheet_entity.AnalysisPlace, map[types.Marker]map[filter.Id]*spreadsheet_entity.AnalysisPlace, map[types.Marker]map[filter.Id]*spreadsheet_entity.AnalysisPlace, []filter.Id)
	Write(ctx context.Context, firstPlaceMap, secondPlaceMap, thirdPlaceMap map[types.Marker]map[filter.Id]*spreadsheet_entity.AnalysisPlace, filters []filter.Id) error
}

type placeService struct {
	filterService         filter_service.AnalysisFilter
	spreadSheetRepository repository.SpreadSheetRepository
}

func NewPlace(
	filterService filter_service.AnalysisFilter,
	spreadSheetRepository repository.SpreadSheetRepository,
) Place {
	return &placeService{
		filterService:         filterService,
		spreadSheetRepository: spreadSheetRepository,
	}
}

func (p *placeService) Create(
	ctx context.Context,
	markers []*marker_csv_entity.AnalysisMarker,
	races []*data_cache_entity.Race,
) ([]*analysis_entity.PlaceCalculable, error) {
	markerMap := converter.ConvertToMap(markers, func(marker *marker_csv_entity.AnalysisMarker) types.RaceId {
		return marker.RaceId()
	})

	var calculables []*analysis_entity.PlaceCalculable
	for _, race := range races {
		raceResultMap := converter.ConvertToMap(race.RaceResults(), func(raceResult *data_cache_entity.RaceResult) types.HorseNumber {
			return raceResult.HorseNumber()
		})

		marker, ok := markerMap[race.RaceId()]
		if !ok {
			switch race.Class() {
			case types.MakeDebut, types.JumpMaiden, types.JumpGrade1, types.JumpGrade2, types.JumpGrade3:
				// 新馬・障害は分析印なしなのでスキップ
			default:
				// 印が不完全な場合がたまにあり(同じ印がついていたり、取り消しによる印6個未満の場合)、その場合はスキップ
				// log.Println(fmt.Sprintf("raceId not found in place markers: %s", race.RaceId()))
			}
			continue
		}

		filters := p.filterService.CreatePlaceFilters(ctx, race)

		// 馬番はレース結果の上位3頭から取る
		// 払い戻し結果から取ってしまうと、複勝2着払いの場合にとれなくなるため
		numbers := make([]types.HorseNumber, 0, 3)
		for _, raceResult := range race.RaceResults()[:3] {
			numbers = append(numbers, raceResult.HorseNumber())
		}
		// 的中の印
		markerCombinationIds := p.getHitMarkerCombinationIds(numbers, marker)
		// 不的中の印
		markerCombinationIds = append(markerCombinationIds, p.getUnHitMarkerCombinationIds(numbers, marker)...)

		var (
			raceCalculables []*analysis_entity.PlaceCalculable
			isRaceCanceled  bool
		)
		// 的中か不的中かは、着順から判断できるためcalculableの中でフラグ管理しない
		for _, markerCombinationId := range markerCombinationIds {
			hitMarker, err := types.NewMarker(markerCombinationId.Value() % 10)
			if err != nil {
				return nil, err
			}
			if hitMarker == types.NoMarker {
				continue
			}

			horseNumber, ok := marker.MarkerMap()[hitMarker]
			if !ok {
				return nil, fmt.Errorf("marker %s is not found in markerMap", hitMarker.String())
			}
			raceResult, ok := raceResultMap[horseNumber]
			if !ok {
				return nil, fmt.Errorf("horseNumber not found: %v", horseNumber)
			}

			// 取り消しの馬かつ、印対象だったばあいそのレースは集計対象外
			if raceResult.Odds() == "0" {
				log.Println(fmt.Sprintf("exclude analysis data for canceled, raceId: %s", race.RaceId()))
				isRaceCanceled = true
				break
			}

			raceCalculables = append(raceCalculables, analysis_entity.NewPlaceCalculable(
				race.RaceId(),
				race.RaceDate(),
				markerCombinationId,
				raceResult.Odds(),
				types.BetNumber(strconv.Itoa(raceResult.HorseNumber().Value())), // 単複のみなのでbetNumberにそのまま置き換え可能
				raceResult.PopularNumber(),
				raceResult.OrderNo(),
				race.Entries(),
				raceResult.JockeyId(),
				filters,
			))
		}

		if !isRaceCanceled {
			calculables = append(calculables, raceCalculables...)
		}
	}

	return calculables, nil
}

func (p *placeService) Convert(
	ctx context.Context,
	calculables []*analysis_entity.PlaceCalculable,
) (
	map[types.Marker]map[filter.Id]*spreadsheet_entity.AnalysisPlace,
	map[types.Marker]map[filter.Id]*spreadsheet_entity.AnalysisPlace,
	map[types.Marker]map[filter.Id]*spreadsheet_entity.AnalysisPlace,
	[]filter.Id,
) {
	firstPlaceMap := map[types.Marker]map[filter.Id]*spreadsheet_entity.AnalysisPlace{}
	secondPlaceMap := map[types.Marker]map[filter.Id]*spreadsheet_entity.AnalysisPlace{}
	thirdPlaceMap := map[types.Marker]map[filter.Id]*spreadsheet_entity.AnalysisPlace{}
	analysisFilters := p.getFilters()

	markers := []types.Marker{
		types.Favorite, types.Rival, types.BrackTriangle, types.WhiteTriangle, types.Star, types.Check,
	}

	for _, marker := range markers {
		firstPlaceMap[marker] = map[filter.Id]*spreadsheet_entity.AnalysisPlace{}
		secondPlaceMap[marker] = map[filter.Id]*spreadsheet_entity.AnalysisPlace{}
		thirdPlaceMap[marker] = map[filter.Id]*spreadsheet_entity.AnalysisPlace{}

		for _, analysisFilter := range analysisFilters {
			raceIdMap := map[types.RaceId]bool{}
			oddsRangeHitCountSlice := make([]int, 24)
			oddsRangeUnHitCountSlice := make([]int, 24)
			for _, calculable := range calculables {
				if calculable.Marker() != marker {
					continue
				}

				match := true
				for _, f := range calculable.Filters() {
					if f&analysisFilter == 0 {
						match = false
						break
					}
				}
				if match {
					if _, ok := raceIdMap[calculable.RaceId()]; !ok {
						raceIdMap[calculable.RaceId()] = true
					}

					odds := calculable.Odds().InexactFloat64()
					if odds >= 1.0 && odds <= 1.5 {
						switch calculable.OrderNo() {
						case 1:
							oddsRangeHitCountSlice[0]++
						case 2:
							oddsRangeHitCountSlice[8]++
						case 3:
							oddsRangeHitCountSlice[16]++
						}
						if calculable.OrderNo() >= 2 {
							oddsRangeUnHitCountSlice[0]++
						}
						if calculable.OrderNo() >= 3 {
							oddsRangeUnHitCountSlice[8]++
						}
						if calculable.OrderNo() >= 4 {
							oddsRangeUnHitCountSlice[16]++
						}
					} else if odds >= 1.6 && odds <= 2.0 {
						switch calculable.OrderNo() {
						case 1:
							oddsRangeHitCountSlice[1]++
						case 2:
							oddsRangeHitCountSlice[9]++
						case 3:
							oddsRangeHitCountSlice[17]++
						}
						if calculable.OrderNo() >= 2 {
							oddsRangeUnHitCountSlice[1]++
						}
						if calculable.OrderNo() >= 3 {
							oddsRangeUnHitCountSlice[9]++
						}
						if calculable.OrderNo() >= 4 {
							oddsRangeUnHitCountSlice[17]++
						}
					} else if odds >= 2.1 && odds <= 2.9 {
						switch calculable.OrderNo() {
						case 1:
							oddsRangeHitCountSlice[2]++
						case 2:
							oddsRangeHitCountSlice[10]++
						case 3:
							oddsRangeHitCountSlice[18]++
						}
						if calculable.OrderNo() >= 2 {
							oddsRangeUnHitCountSlice[2]++
						}
						if calculable.OrderNo() >= 3 {
							oddsRangeUnHitCountSlice[10]++
						}
						if calculable.OrderNo() >= 4 {
							oddsRangeUnHitCountSlice[18]++
						}
					} else if odds >= 3.0 && odds <= 4.9 {
						switch calculable.OrderNo() {
						case 1:
							oddsRangeHitCountSlice[3]++
						case 2:
							oddsRangeHitCountSlice[11]++
						case 3:
							oddsRangeHitCountSlice[19]++
						}
						if calculable.OrderNo() >= 2 {
							oddsRangeUnHitCountSlice[3]++
						}
						if calculable.OrderNo() >= 3 {
							oddsRangeUnHitCountSlice[11]++
						}
						if calculable.OrderNo() >= 4 {
							oddsRangeUnHitCountSlice[19]++
						}
					} else if odds >= 5.0 && odds <= 9.9 {
						switch calculable.OrderNo() {
						case 1:
							oddsRangeHitCountSlice[4]++
						case 2:
							oddsRangeHitCountSlice[12]++
						case 3:
							oddsRangeHitCountSlice[20]++
						}
						if calculable.OrderNo() >= 2 {
							oddsRangeUnHitCountSlice[4]++
						}
						if calculable.OrderNo() >= 3 {
							oddsRangeUnHitCountSlice[12]++
						}
						if calculable.OrderNo() >= 4 {
							oddsRangeUnHitCountSlice[20]++
						}
					} else if odds >= 10.0 && odds <= 19.9 {
						switch calculable.OrderNo() {
						case 1:
							oddsRangeHitCountSlice[5]++
						case 2:
							oddsRangeHitCountSlice[13]++
						case 3:
							oddsRangeHitCountSlice[21]++
						}
						if calculable.OrderNo() >= 2 {
							oddsRangeUnHitCountSlice[5]++
						}
						if calculable.OrderNo() >= 3 {
							oddsRangeUnHitCountSlice[13]++
						}
						if calculable.OrderNo() >= 4 {
							oddsRangeUnHitCountSlice[21]++
						}
					} else if odds >= 20.0 && odds <= 49.9 {
						switch calculable.OrderNo() {
						case 1:
							oddsRangeHitCountSlice[6]++
						case 2:
							oddsRangeHitCountSlice[14]++
						case 3:
							oddsRangeHitCountSlice[22]++
						}
						if calculable.OrderNo() >= 2 {
							oddsRangeUnHitCountSlice[6]++
						}
						if calculable.OrderNo() >= 3 {
							oddsRangeUnHitCountSlice[14]++
						}
						if calculable.OrderNo() >= 4 {
							oddsRangeUnHitCountSlice[22]++
						}
					} else if odds >= 50.0 {
						switch calculable.OrderNo() {
						case 1:
							oddsRangeHitCountSlice[7]++
						case 2:
							oddsRangeHitCountSlice[15]++
						case 3:
							oddsRangeHitCountSlice[23]++
						}
						if calculable.OrderNo() >= 2 {
							oddsRangeUnHitCountSlice[7]++
						}
						if calculable.OrderNo() >= 3 {
							oddsRangeUnHitCountSlice[15]++
						}
						if calculable.OrderNo() >= 4 {
							oddsRangeUnHitCountSlice[23]++
						}
					}
				}
			}

			firstPlaceOddsRangeHitCountSlice := make([]int, 8)
			secondPlaceOddsRangeHitCountSlice := make([]int, 8)
			thirdPlaceOddsRangeHitCountSlice := make([]int, 8)
			firstPlaceOddsRangeUnHitCountSlice := make([]int, 8)
			secondPlaceOddsRangeUnHitCountSlice := make([]int, 8)
			thirdPlaceOddsRangeUnHitCountSlice := make([]int, 8)

			for i := 0; i < 8; i++ {
				firstPlaceOddsRangeHitCountSlice[i] = oddsRangeHitCountSlice[i]
				secondPlaceOddsRangeHitCountSlice[i] = oddsRangeHitCountSlice[i] + oddsRangeHitCountSlice[i+8]
				thirdPlaceOddsRangeHitCountSlice[i] = oddsRangeHitCountSlice[i] + oddsRangeHitCountSlice[i+8] + oddsRangeHitCountSlice[i+16]
				firstPlaceOddsRangeUnHitCountSlice[i] = oddsRangeUnHitCountSlice[i]
				secondPlaceOddsRangeUnHitCountSlice[i] = oddsRangeUnHitCountSlice[i+8]
				thirdPlaceOddsRangeUnHitCountSlice[i] = oddsRangeUnHitCountSlice[i+16]
			}

			firstPlaceOddsRangeHitCountData := spreadsheet_entity.NewPlaceHitCountData(
				firstPlaceOddsRangeHitCountSlice,
				analysisFilter,
				len(raceIdMap),
			)
			secondPlaceOddsRangeHitCountData := spreadsheet_entity.NewPlaceHitCountData(
				secondPlaceOddsRangeHitCountSlice,
				analysisFilter,
				len(raceIdMap),
			)
			thirdPlaceOddsRangeHitCountData := spreadsheet_entity.NewPlaceHitCountData(
				thirdPlaceOddsRangeHitCountSlice,
				analysisFilter,
				len(raceIdMap),
			)

			firstPlaceOddsRangeUnHitCountData := spreadsheet_entity.NewPlaceUnHitCountData(
				firstPlaceOddsRangeUnHitCountSlice,
				analysisFilter,
				len(raceIdMap),
			)
			secondPlaceOddsRangeUnHitCountData := spreadsheet_entity.NewPlaceUnHitCountData(
				secondPlaceOddsRangeUnHitCountSlice,
				analysisFilter,
				len(raceIdMap),
			)
			thirdPlaceOddsRangeUnHitCountData := spreadsheet_entity.NewPlaceUnHitCountData(
				thirdPlaceOddsRangeUnHitCountSlice,
				analysisFilter,
				len(raceIdMap),
			)

			firstPlaceOddsRangeRateData := spreadsheet_entity.NewPlaceRateData(
				firstPlaceOddsRangeHitCountData,
				firstPlaceOddsRangeUnHitCountData,
			)
			firstPlaceOddsRangeRateStyle := spreadsheet_entity.NewPlaceRateStyle(firstPlaceOddsRangeRateData)

			secondPlaceOddsRangeRateData := spreadsheet_entity.NewPlaceRateData(
				secondPlaceOddsRangeHitCountData,
				secondPlaceOddsRangeUnHitCountData,
			)
			secondPlaceOddsRangeRateStyle := spreadsheet_entity.NewPlaceRateStyle(secondPlaceOddsRangeRateData)

			thirdPlaceOddsRangeRateData := spreadsheet_entity.NewPlaceRateData(
				thirdPlaceOddsRangeHitCountData,
				thirdPlaceOddsRangeUnHitCountData,
			)
			thirdPlaceOddsRangeRateStyle := spreadsheet_entity.NewPlaceRateStyle(thirdPlaceOddsRangeRateData)

			firstPlaceMap[marker][analysisFilter] = spreadsheet_entity.NewAnalysisPlace(
				firstPlaceOddsRangeRateData,
				firstPlaceOddsRangeRateStyle,
				firstPlaceOddsRangeHitCountData,
				firstPlaceOddsRangeUnHitCountData,
			)

			secondPlaceMap[marker][analysisFilter] = spreadsheet_entity.NewAnalysisPlace(
				secondPlaceOddsRangeRateData,
				secondPlaceOddsRangeRateStyle,
				secondPlaceOddsRangeHitCountData,
				secondPlaceOddsRangeUnHitCountData,
			)

			thirdPlaceMap[marker][analysisFilter] = spreadsheet_entity.NewAnalysisPlace(
				thirdPlaceOddsRangeRateData,
				thirdPlaceOddsRangeRateStyle,
				thirdPlaceOddsRangeHitCountData,
				thirdPlaceOddsRangeUnHitCountData,
			)
		}
	}

	return firstPlaceMap, secondPlaceMap, thirdPlaceMap, analysisFilters
}

func (p *placeService) Write(
	ctx context.Context,
	firstPlaceMap,
	secondPlaceMap,
	thirdPlaceMap map[types.Marker]map[filter.Id]*spreadsheet_entity.AnalysisPlace,
	filters []filter.Id,
) error {
	return p.spreadSheetRepository.WriteAnalysisPlace(ctx, firstPlaceMap, secondPlaceMap, thirdPlaceMap, filters)
}

func (p *placeService) getHitMarkerCombinationIds(
	numbers []types.HorseNumber,
	marker *marker_csv_entity.AnalysisMarker,
) []types.MarkerCombinationId {
	var markerCombinationIds []types.MarkerCombinationId
	for _, number := range numbers {
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
		markerCombinationIds = append(markerCombinationIds, markerCombinationId)
	}

	return markerCombinationIds
}

func (p *placeService) getUnHitMarkerCombinationIds(
	numbers []types.HorseNumber,
	marker *marker_csv_entity.AnalysisMarker,
) []types.MarkerCombinationId {
	unHitMarkerCombinationIdMap := map[types.MarkerCombinationId]bool{
		types.MarkerCombinationId(21): true,
		types.MarkerCombinationId(22): true,
		types.MarkerCombinationId(23): true,
		types.MarkerCombinationId(24): true,
		types.MarkerCombinationId(25): true,
		types.MarkerCombinationId(26): true,
		types.MarkerCombinationId(29): true,
	}

	for _, number := range numbers {
		switch number {
		case marker.Favorite():
			unHitMarkerCombinationIdMap[types.MarkerCombinationId(21)] = false
		case marker.Rival():
			unHitMarkerCombinationIdMap[types.MarkerCombinationId(22)] = false
		case marker.BrackTriangle():
			unHitMarkerCombinationIdMap[types.MarkerCombinationId(23)] = false
		case marker.WhiteTriangle():
			unHitMarkerCombinationIdMap[types.MarkerCombinationId(24)] = false
		case marker.Star():
			unHitMarkerCombinationIdMap[types.MarkerCombinationId(25)] = false
		case marker.Check():
			unHitMarkerCombinationIdMap[types.MarkerCombinationId(26)] = false
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

func (p *placeService) getFilters() []filter.Id {
	return []filter.Id{
		filter.All2,
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
