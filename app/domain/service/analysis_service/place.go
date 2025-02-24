package analysis_service

import (
	"context"
	"fmt"
	"strconv"

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

type Place interface {
	Create(ctx context.Context, markers []*marker_csv_entity.AnalysisMarker, races []*data_cache_entity.Race) ([]*analysis_entity.PlaceCalculable, error)
	Convert(ctx context.Context, calculables []*analysis_entity.PlaceCalculable) (map[types.Marker]map[filter.AttributeId]*spreadsheet_entity.AnalysisPlace, map[types.Marker]map[filter.AttributeId]*spreadsheet_entity.AnalysisPlace, map[types.Marker]map[filter.AttributeId]*spreadsheet_entity.AnalysisPlace, []filter.AttributeId)
	Write(ctx context.Context, firstPlaceMap, secondPlaceMap, thirdPlaceMap map[types.Marker]map[filter.AttributeId]*spreadsheet_entity.AnalysisPlace, filters []filter.AttributeId) error
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
				return nil, fmt.Errorf("horseNumber %v not found in raceId %v", horseNumber, race.RaceId())
			}

			// 取り消しの馬かつ、印対象だったばあいそのレースは集計対象外
			if raceResult.Odds().IsZero() {
				//log.Println(fmt.Sprintf("exclude analysis data for canceled, raceId: %s", race.RaceId()))
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
	map[types.Marker]map[filter.AttributeId]*spreadsheet_entity.AnalysisPlace,
	map[types.Marker]map[filter.AttributeId]*spreadsheet_entity.AnalysisPlace,
	map[types.Marker]map[filter.AttributeId]*spreadsheet_entity.AnalysisPlace,
	[]filter.AttributeId,
) {
	firstPlaceMap := map[types.Marker]map[filter.AttributeId]*spreadsheet_entity.AnalysisPlace{}
	secondPlaceMap := map[types.Marker]map[filter.AttributeId]*spreadsheet_entity.AnalysisPlace{}
	thirdPlaceMap := map[types.Marker]map[filter.AttributeId]*spreadsheet_entity.AnalysisPlace{}
	analysisFilters := p.getFilters()

	markers := []types.Marker{
		types.Favorite, types.Rival, types.BrackTriangle, types.WhiteTriangle, types.Star, types.Check,
	}

	for _, marker := range markers {
		firstPlaceMap[marker] = map[filter.AttributeId]*spreadsheet_entity.AnalysisPlace{}
		secondPlaceMap[marker] = map[filter.AttributeId]*spreadsheet_entity.AnalysisPlace{}
		thirdPlaceMap[marker] = map[filter.AttributeId]*spreadsheet_entity.AnalysisPlace{}

		for _, analysisFilter := range analysisFilters {
			raceIdMap := map[types.RaceId]bool{}
			oddsRangeHitCountSlice := make([]int, 27)
			oddsRangeUnHitCountSlice := make([]int, 27)

			for _, calculable := range calculables {
				if calculable.Marker() != marker {
					continue
				}

				var calcFilter filter.AttributeId
				for _, f := range calculable.Filters() {
					calcFilter |= f
				}

				if analysisFilter == filter.All || analysisFilter&calcFilter == analysisFilter {
					if _, ok := raceIdMap[calculable.RaceId()]; !ok {
						raceIdMap[calculable.RaceId()] = true
					}

					odds := calculable.Odds().InexactFloat64()
					if odds >= 1.0 && odds <= 1.4 {
						switch calculable.OrderNo() {
						case 1:
							oddsRangeHitCountSlice[0]++
						case 2:
							oddsRangeHitCountSlice[9]++
						case 3:
							oddsRangeHitCountSlice[18]++
						}
						if calculable.OrderNo() >= 2 {
							oddsRangeUnHitCountSlice[0]++
						}
						if calculable.OrderNo() >= 3 {
							oddsRangeUnHitCountSlice[9]++
						}
						if calculable.OrderNo() >= 4 {
							oddsRangeUnHitCountSlice[18]++
						}
					} else if odds >= 1.5 && odds <= 1.9 {
						switch calculable.OrderNo() {
						case 1:
							oddsRangeHitCountSlice[1]++
						case 2:
							oddsRangeHitCountSlice[10]++
						case 3:
							oddsRangeHitCountSlice[19]++
						}
						if calculable.OrderNo() >= 2 {
							oddsRangeUnHitCountSlice[1]++
						}
						if calculable.OrderNo() >= 3 {
							oddsRangeUnHitCountSlice[10]++
						}
						if calculable.OrderNo() >= 4 {
							oddsRangeUnHitCountSlice[19]++
						}
					} else if odds >= 2.0 && odds <= 2.2 {
						switch calculable.OrderNo() {
						case 1:
							oddsRangeHitCountSlice[2]++
						case 2:
							oddsRangeHitCountSlice[11]++
						case 3:
							oddsRangeHitCountSlice[20]++
						}
						if calculable.OrderNo() >= 2 {
							oddsRangeUnHitCountSlice[2]++
						}
						if calculable.OrderNo() >= 3 {
							oddsRangeUnHitCountSlice[11]++
						}
						if calculable.OrderNo() >= 4 {
							oddsRangeUnHitCountSlice[20]++
						}
					} else if odds >= 2.3 && odds <= 3.0 {
						switch calculable.OrderNo() {
						case 1:
							oddsRangeHitCountSlice[3]++
						case 2:
							oddsRangeHitCountSlice[12]++
						case 3:
							oddsRangeHitCountSlice[21]++
						}
						if calculable.OrderNo() >= 2 {
							oddsRangeUnHitCountSlice[3]++
						}
						if calculable.OrderNo() >= 3 {
							oddsRangeUnHitCountSlice[12]++
						}
						if calculable.OrderNo() >= 4 {
							oddsRangeUnHitCountSlice[21]++
						}
					} else if odds >= 3.1 && odds <= 4.9 {
						switch calculable.OrderNo() {
						case 1:
							oddsRangeHitCountSlice[4]++
						case 2:
							oddsRangeHitCountSlice[13]++
						case 3:
							oddsRangeHitCountSlice[22]++
						}
						if calculable.OrderNo() >= 2 {
							oddsRangeUnHitCountSlice[4]++
						}
						if calculable.OrderNo() >= 3 {
							oddsRangeUnHitCountSlice[13]++
						}
						if calculable.OrderNo() >= 4 {
							oddsRangeUnHitCountSlice[22]++
						}
					} else if odds >= 5.0 && odds <= 9.9 {
						switch calculable.OrderNo() {
						case 1:
							oddsRangeHitCountSlice[5]++
						case 2:
							oddsRangeHitCountSlice[14]++
						case 3:
							oddsRangeHitCountSlice[23]++
						}
						if calculable.OrderNo() >= 2 {
							oddsRangeUnHitCountSlice[5]++
						}
						if calculable.OrderNo() >= 3 {
							oddsRangeUnHitCountSlice[14]++
						}
						if calculable.OrderNo() >= 4 {
							oddsRangeUnHitCountSlice[23]++
						}
					} else if odds >= 10.0 && odds <= 19.9 {
						switch calculable.OrderNo() {
						case 1:
							oddsRangeHitCountSlice[6]++
						case 2:
							oddsRangeHitCountSlice[15]++
						case 3:
							oddsRangeHitCountSlice[24]++
						}
						if calculable.OrderNo() >= 2 {
							oddsRangeUnHitCountSlice[6]++
						}
						if calculable.OrderNo() >= 3 {
							oddsRangeUnHitCountSlice[15]++
						}
						if calculable.OrderNo() >= 4 {
							oddsRangeUnHitCountSlice[24]++
						}
					} else if odds >= 20.0 && odds <= 49.9 {
						switch calculable.OrderNo() {
						case 1:
							oddsRangeHitCountSlice[7]++
						case 2:
							oddsRangeHitCountSlice[16]++
						case 3:
							oddsRangeHitCountSlice[25]++
						}
						if calculable.OrderNo() >= 2 {
							oddsRangeUnHitCountSlice[7]++
						}
						if calculable.OrderNo() >= 3 {
							oddsRangeUnHitCountSlice[16]++
						}
						if calculable.OrderNo() >= 4 {
							oddsRangeUnHitCountSlice[25]++
						}
					} else if odds >= 50.0 {
						switch calculable.OrderNo() {
						case 1:
							oddsRangeHitCountSlice[8]++
						case 2:
							oddsRangeHitCountSlice[17]++
						case 3:
							oddsRangeHitCountSlice[26]++
						}
						if calculable.OrderNo() >= 2 {
							oddsRangeUnHitCountSlice[8]++
						}
						if calculable.OrderNo() >= 3 {
							oddsRangeUnHitCountSlice[17]++
						}
						if calculable.OrderNo() >= 4 {
							oddsRangeUnHitCountSlice[26]++
						}
					}
				}
			}

			firstPlaceOddsRangeHitCountSlice := make([]int, 9)
			secondPlaceOddsRangeHitCountSlice := make([]int, 9)
			thirdPlaceOddsRangeHitCountSlice := make([]int, 9)
			firstPlaceOddsRangeUnHitCountSlice := make([]int, 9)
			secondPlaceOddsRangeUnHitCountSlice := make([]int, 9)
			thirdPlaceOddsRangeUnHitCountSlice := make([]int, 9)

			for i := 0; i < 9; i++ {
				firstPlaceOddsRangeHitCountSlice[i] = oddsRangeHitCountSlice[i]
				secondPlaceOddsRangeHitCountSlice[i] = oddsRangeHitCountSlice[i] + oddsRangeHitCountSlice[i+9]
				thirdPlaceOddsRangeHitCountSlice[i] = oddsRangeHitCountSlice[i] + oddsRangeHitCountSlice[i+9] + oddsRangeHitCountSlice[i+18]
				firstPlaceOddsRangeUnHitCountSlice[i] = oddsRangeUnHitCountSlice[i]
				secondPlaceOddsRangeUnHitCountSlice[i] = oddsRangeUnHitCountSlice[i+9]
				thirdPlaceOddsRangeUnHitCountSlice[i] = oddsRangeUnHitCountSlice[i+18]
			}

			firstPlaceOddsRangeHitCountData := spreadsheet_entity.NewPlaceHitCountData(
				firstPlaceOddsRangeHitCountSlice,
				len(raceIdMap),
			)
			secondPlaceOddsRangeHitCountData := spreadsheet_entity.NewPlaceHitCountData(
				secondPlaceOddsRangeHitCountSlice,
				len(raceIdMap),
			)
			thirdPlaceOddsRangeHitCountData := spreadsheet_entity.NewPlaceHitCountData(
				thirdPlaceOddsRangeHitCountSlice,
				len(raceIdMap),
			)

			firstPlaceOddsRangeUnHitCountData := spreadsheet_entity.NewPlaceUnHitCountData(
				firstPlaceOddsRangeUnHitCountSlice,
				len(raceIdMap),
			)
			secondPlaceOddsRangeUnHitCountData := spreadsheet_entity.NewPlaceUnHitCountData(
				secondPlaceOddsRangeUnHitCountSlice,
				len(raceIdMap),
			)
			thirdPlaceOddsRangeUnHitCountData := spreadsheet_entity.NewPlaceUnHitCountData(
				thirdPlaceOddsRangeUnHitCountSlice,
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
	thirdPlaceMap map[types.Marker]map[filter.AttributeId]*spreadsheet_entity.AnalysisPlace,
	filters []filter.AttributeId,
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

func (p *placeService) getFilters() []filter.AttributeId {
	return []filter.AttributeId{
		filter.All,
		filter.Turf | filter.Niigata | filter.Distance1000m,
		//filter.Turf | filter.Hakodate | filter.Distance1000m, // 現在は使われていない
		filter.Turf | filter.Nakayama | filter.Distance1200m,
		filter.Turf | filter.Kyoto | filter.Distance1200m,
		filter.Turf | filter.Hanshin | filter.Distance1200m,
		filter.Turf | filter.Niigata | filter.Distance1200m,
		filter.Turf | filter.Chukyo | filter.Distance1200m,
		filter.Turf | filter.Sapporo | filter.Distance1200m,
		filter.Turf | filter.Hakodate | filter.Distance1200m,
		filter.Turf | filter.Fukushima | filter.Distance1200m,
		filter.Turf | filter.Kokura | filter.Distance1200m,
		filter.Turf | filter.Tokyo | filter.Distance1400m,
		filter.Turf | filter.Kyoto | filter.Distance1400m,
		filter.Turf | filter.Hanshin | filter.Distance1400m,
		filter.Turf | filter.Niigata | filter.Distance1400m,
		filter.Turf | filter.Chukyo | filter.Distance1400m,
		filter.Turf | filter.Sapporo | filter.Distance1500m,
		filter.Turf | filter.Nakayama | filter.Distance1600m,
		filter.Turf | filter.Tokyo | filter.Distance1600m,
		filter.Turf | filter.Kyoto | filter.Distance1600m,
		filter.Turf | filter.Hanshin | filter.Distance1600m,
		filter.Turf | filter.Niigata | filter.Distance1600m,
		filter.Turf | filter.Chukyo | filter.Distance1600m,
		filter.Turf | filter.Nakayama | filter.Distance1800m,
		filter.Turf | filter.Tokyo | filter.Distance1800m,
		filter.Turf | filter.Kyoto | filter.Distance1800m,
		filter.Turf | filter.Hanshin | filter.Distance1800m,
		filter.Turf | filter.Niigata | filter.Distance1800m,
		filter.Turf | filter.Sapporo | filter.Distance1800m,
		filter.Turf | filter.Hakodate | filter.Distance1800m,
		filter.Turf | filter.Fukushima | filter.Distance1800m,
		filter.Turf | filter.Kokura | filter.Distance1800m,
		filter.Turf | filter.Nakayama | filter.Distance2000m,
		filter.Turf | filter.Tokyo | filter.Distance2000m,
		filter.Turf | filter.Kyoto | filter.Distance2000m,
		filter.Turf | filter.Hanshin | filter.Distance2000m,
		filter.Turf | filter.Niigata | filter.Distance2000m,
		filter.Turf | filter.Chukyo | filter.Distance2000m,
		filter.Turf | filter.Sapporo | filter.Distance2000m,
		filter.Turf | filter.Hakodate | filter.Distance2000m,
		filter.Turf | filter.Fukushima | filter.Distance2000m,
		filter.Turf | filter.Kokura | filter.Distance2000m,
		filter.Turf | filter.Nakayama | filter.Distance2200m,
		filter.Turf | filter.Kyoto | filter.Distance2200m,
		filter.Turf | filter.Hanshin | filter.Distance2200m,
		filter.Turf | filter.Niigata | filter.Distance2200m,
		filter.Turf | filter.Chukyo | filter.Distance2200m,
		filter.Turf | filter.Tokyo | filter.Distance2300m,
		filter.Turf | filter.Tokyo | filter.Distance2400m,
		filter.Turf | filter.Kyoto | filter.Distance2400m,
		filter.Turf | filter.Hanshin | filter.Distance2400m,
		filter.Turf | filter.Niigata | filter.Distance2400m,
		filter.Turf | filter.Nakayama | filter.Distance2500m,
		filter.Turf | filter.Tokyo | filter.Distance2500m,
		filter.Turf | filter.Hanshin | filter.Distance2600m,
		filter.Turf | filter.Sapporo | filter.Distance2600m,
		filter.Turf | filter.Hakodate | filter.Distance2600m,
		filter.Turf | filter.Fukushima | filter.Distance2600m,
		filter.Turf | filter.Kokura | filter.Distance2600m,
		filter.Turf | filter.Kyoto | filter.Distance3000m,
		filter.Turf | filter.Hanshin | filter.Distance3000m,
		//filter.Turf | filter.Chukyo | filter.Distance3000m, // 現在は使われていない
		filter.Turf | filter.Kyoto | filter.Distance3200m,
		filter.Turf | filter.Tokyo | filter.Distance3400m,
		filter.Turf | filter.Nakayama | filter.Distance3600m,
		filter.Dirt | filter.Sapporo | filter.Distance1000m,
		filter.Dirt | filter.Hakodate | filter.Distance1000m,
		filter.Dirt | filter.Kokura | filter.Distance1000m,
		filter.Dirt | filter.Fukushima | filter.Distance1150m,
		filter.Dirt | filter.Nakayama | filter.Distance1200m,
		filter.Dirt | filter.Kyoto | filter.Distance1200m,
		filter.Dirt | filter.Hanshin | filter.Distance1200m,
		filter.Dirt | filter.Niigata | filter.Distance1200m,
		filter.Dirt | filter.Chukyo | filter.Distance1200m,
		filter.Dirt | filter.Tokyo | filter.Distance1300m,
		filter.Dirt | filter.Tokyo | filter.Distance1400m,
		filter.Dirt | filter.Kyoto | filter.Distance1400m,
		filter.Dirt | filter.Hanshin | filter.Distance1400m,
		filter.Dirt | filter.Chukyo | filter.Distance1400m,
		filter.Dirt | filter.Tokyo | filter.Distance1600m,
		filter.Dirt | filter.Sapporo | filter.Distance1700m,
		filter.Dirt | filter.Hakodate | filter.Distance1700m,
		filter.Dirt | filter.Fukushima | filter.Distance1700m,
		filter.Dirt | filter.Kokura | filter.Distance1700m,
		filter.Dirt | filter.Nakayama | filter.Distance1800m,
		filter.Dirt | filter.Kyoto | filter.Distance1800m,
		filter.Dirt | filter.Hanshin | filter.Distance1800m,
		filter.Dirt | filter.Niigata | filter.Distance1800m,
		filter.Dirt | filter.Chukyo | filter.Distance1800m,
		filter.Dirt | filter.Kyoto | filter.Distance1900m,
		filter.Dirt | filter.Chukyo | filter.Distance1900m,
		filter.Dirt | filter.Hanshin | filter.Distance2000m,
		filter.Dirt | filter.Tokyo | filter.Distance2100m,
		filter.Dirt | filter.Nakayama | filter.Distance2400m,
		filter.Dirt | filter.Sapporo | filter.Distance2400m,
		filter.Dirt | filter.Hakodate | filter.Distance2400m,
		filter.Dirt | filter.Fukushima | filter.Distance2400m,
		filter.Dirt | filter.Kokura | filter.Distance2400m,
		filter.Dirt | filter.Nakayama | filter.Distance2500m,
		filter.Dirt | filter.Niigata | filter.Distance2500m,
	}
}
