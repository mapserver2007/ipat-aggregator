package analysis_service

import (
	"context"
	"fmt"
	"slices"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/tospo_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/filter_service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
	"github.com/mapserver2007/ipat-aggregator/config"
	"github.com/shopspring/decimal"
)

type PlaceUnHit interface {
	GetUnHitRaces(ctx context.Context,
		markers []*marker_csv_entity.AnalysisMarker,
		races []*data_cache_entity.Race,
		jockeys []*data_cache_entity.Jockey,
	) []*analysis_entity.Race
	GetUnHitRaceRate(ctx context.Context,
		race *analysis_entity.Race,
		calculables []*analysis_entity.PlaceCalculable,
	) map[types.HorseId][]float64
	GetCheckList(ctx context.Context,
		race *analysis_entity.Race,
		horse *data_cache_entity.Horse,
		raceForecast *data_cache_entity.RaceForecast,
	) error
	GetWinRedOdds(ctx context.Context,
		oddsList []*data_cache_entity.Odds,
		thresholdOdds decimal.Decimal,
		raceMap map[types.RaceId]*analysis_entity.Race,
	) ([]*analysis_entity.Odds, error)
	GetWinOddsFaults(ctx context.Context,
		oddsList []*data_cache_entity.Odds,
		raceMap map[types.RaceId]*analysis_entity.Race,
	) ([]*analysis_entity.OddsFault, error)
	GetTrioOdds(ctx context.Context,
		oddsList []*data_cache_entity.Odds,
		popularNumber int,
		raceMap map[types.RaceId]*analysis_entity.Race,
	) ([]*analysis_entity.Odds, error)
	GetQuinellaOddsWheelCombinations(ctx context.Context,
		oddsList []*data_cache_entity.Odds,
		winOddsMap map[types.RaceId][]*analysis_entity.Odds,
		raceMap map[types.RaceId]*analysis_entity.Race,
	) ([]*analysis_entity.Odds, error)
	FetchHorse(ctx context.Context, horseId types.HorseId) (*netkeiba_entity.Horse, error)
	FetchRaceForecasts(ctx context.Context, raceId types.RaceId) ([]*tospo_entity.Forecast, error)
	FetchTrainingComments(ctx context.Context, raceId types.RaceId) ([]*tospo_entity.TrainingComment, error)
	CreateUnhitRaces(ctx context.Context,
		races []*analysis_entity.Race,
		raceRateMap map[types.RaceId]map[types.HorseId][]float64,
		raceForecastMap map[types.RaceId]*data_cache_entity.RaceForecast,
		horseMap map[types.HorseId]*data_cache_entity.Horse,
		winMultiOddsMap map[types.RaceId][]*analysis_entity.Odds,
		winOddsFaultMap map[types.RaceId][]*analysis_entity.OddsFault,
		trioOddsMap map[types.RaceId]*analysis_entity.Odds,
		quinellaConsecutiveNumberMap map[types.RaceId]int,
		quinellaCombinationTotalOddsMap map[types.RaceId]decimal.Decimal,
	) ([]*spreadsheet_entity.AnalysisPlaceUnhit, error)
	CreateCheckPoints(ctx context.Context, race *analysis_entity.Race) error
	Write(ctx context.Context, analysisPlaceUnhits []*spreadsheet_entity.AnalysisPlaceUnhit) error
}

type placeUnHitService struct {
	horseRepository        repository.HorseRepository
	raceForecastRepository repository.RaceForecastRepository
	spreadSheetRepository  repository.SpreadSheetRepository
	horseEntityConverter   converter.HorseEntityConverter
	filterService          filter_service.AnalysisFilter
	placeCheckListService  PlaceCheckList
	placeCheckPointService PlaceCheckPoint
}

func NewPlaceUnHit(
	horseRepository repository.HorseRepository,
	raceForecastRepository repository.RaceForecastRepository,
	spreadSheetRepository repository.SpreadSheetRepository,
	horseEntityConverter converter.HorseEntityConverter,
	filterService filter_service.AnalysisFilter,
	placeCheckListService PlaceCheckList,
	placeCheckPointService PlaceCheckPoint,
) PlaceUnHit {
	return &placeUnHitService{
		horseRepository:        horseRepository,
		raceForecastRepository: raceForecastRepository,
		spreadSheetRepository:  spreadSheetRepository,
		horseEntityConverter:   horseEntityConverter,
		filterService:          filterService,
		placeCheckListService:  placeCheckListService,
		placeCheckPointService: placeCheckPointService,
	}
}

func (p *placeUnHitService) GetUnHitRaces(
	ctx context.Context,
	markers []*marker_csv_entity.AnalysisMarker,
	races []*data_cache_entity.Race,
	jockeys []*data_cache_entity.Jockey,
) []*analysis_entity.Race {
	markerMap := converter.ConvertToMap(markers, func(marker *marker_csv_entity.AnalysisMarker) types.RaceId {
		return marker.RaceId()
	})

	jockeyMap := converter.ConvertToMap(jockeys, func(jockey *data_cache_entity.Jockey) types.JockeyId {
		return jockey.JockeyId()
	})

	var unHitRaces []*analysis_entity.Race
	for _, race := range races {
		analysisMarker, ok := markerMap[race.RaceId()]
		if !ok {
			continue
		}

		var (
			analysisRaceResults []*analysis_entity.RaceResult
			analysisMarkers     []*analysis_entity.Marker
		)
		for _, raceResult := range race.RaceResults() {
			// 除外になった馬は0倍
			if raceResult.Odds().IsZero() {
				continue
			}

			if raceResult.Odds().InexactFloat64() < config.AnalysisUnHitWinLowerOdds && raceResult.OrderNo() > 3 {
				jockeyName := ""
				jockey, ok := jockeyMap[raceResult.JockeyId()]
				if ok {
					jockeyName = jockey.JockeyName()
				}

				analysisRaceResults = append(analysisRaceResults, analysis_entity.NewRaceResult(
					raceResult.OrderNo(),
					raceResult.HorseId(),
					raceResult.HorseName(),
					raceResult.HorseNumber(),
					raceResult.JockeyId(),
					jockeyName,
					raceResult.Odds(),
					raceResult.PopularNumber(),
					raceResult.JockeyWeight(),
					raceResult.HorseWeight(),
					raceResult.HorseWeightAdd(),
				))

				marker := types.NoMarker
				switch raceResult.HorseNumber() {
				case analysisMarker.Favorite():
					marker = types.Favorite
				case analysisMarker.Rival():
					marker = types.Rival
				case analysisMarker.BrackTriangle():
					marker = types.BrackTriangle
				case analysisMarker.WhiteTriangle():
					marker = types.WhiteTriangle
				case analysisMarker.Star():
					marker = types.Star
				case analysisMarker.Check():
					marker = types.Check
				}

				analysisMarkers = append(analysisMarkers, analysis_entity.NewMarker(
					marker, raceResult.HorseNumber(),
				))
			}
		}

		if len(analysisRaceResults) > 0 {
			unHitRaces = append(unHitRaces, analysis_entity.NewRace(
				race.RaceId(),
				race.RaceDate(),
				race.RaceNumber(),
				race.RaceCourseId(),
				race.RaceName(),
				race.Url(),
				race.Entries(),
				race.Distance(),
				race.Class(),
				race.CourseCategory(),
				race.TrackCondition(),
				race.RaceWeightCondition(),
				analysisRaceResults,
				analysisMarkers,
				p.filterService.CreatePlaceFilters(ctx, race),
			))
		}
	}

	return unHitRaces
}

func (p *placeUnHitService) GetUnHitRaceRate(
	ctx context.Context,
	race *analysis_entity.Race,
	calculables []*analysis_entity.PlaceCalculable,
) map[types.HorseId][]float64 {
	var analysisFilter filter.AttributeId
	for _, f := range race.AnalysisFilters() {
		analysisFilter |= f
	}

	placeRateMap := map[types.HorseId][]float64{}
	var firstPlaceRate, secondPlaceRate, thirdPlaceRate float64

	// 対象レース日付以降のデータは計算に含める
	// 過去時点の率と異なるが初期の頃のレースのデータ数が少なすぎるため
	// 分析用途として見るので当時と現在の率が異なっても良いこととする
	for idx, raceResult := range race.RaceResults() {
		placeRateMap[raceResult.HorseId()] = make([]float64, 3)

		oddsRangeHitCountSlice := make([]int, 27)
		oddsRangeUnHitCountSlice := make([]int, 27)
		marker := race.Markers()[idx].Marker()

		for _, calculable := range calculables {
			if calculable.Marker() != marker {
				continue
			}
			calculable.Odds()

			match := true
			for _, f := range calculable.Filters() {
				if f&analysisFilter == 0 {
					match = false
					break
				}
			}

			// 対象レース日付以降のデータは計算に含めたい場合は解除する
			//if race.RaceDate() >= calculable.RaceDate() {
			//	continue
			//}

			if match {
				inexactOdds := calculable.Odds().InexactFloat64()
				if inexactOdds >= 1.0 && inexactOdds <= 1.4 {
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
				} else if inexactOdds >= 1.5 && inexactOdds <= 1.9 {
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
				} else if inexactOdds >= 2.0 && inexactOdds <= 2.2 {
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
				} else if inexactOdds >= 2.3 && inexactOdds <= 3.0 {
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
				} else if inexactOdds >= 3.1 && inexactOdds <= 4.9 {
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
				} else if inexactOdds >= 5.0 && inexactOdds <= 9.9 {
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
				} else if inexactOdds >= 10.0 && inexactOdds <= 19.9 {
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
				} else if inexactOdds >= 20.0 && inexactOdds <= 49.9 {
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
				} else if inexactOdds >= 50.0 {
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

		inexactOdds := raceResult.Odds().InexactFloat64()
		if inexactOdds >= 1.0 && inexactOdds <= 1.4 {
			firstPlaceRateCount := oddsRangeHitCountSlice[0]
			secondPlaceRateCount := oddsRangeHitCountSlice[0] + oddsRangeHitCountSlice[9]
			thirdPlaceRateCount := oddsRangeHitCountSlice[0] + oddsRangeHitCountSlice[9] + oddsRangeHitCountSlice[18]
			firstPlaceRate = float64(firstPlaceRateCount) * 100 / float64(firstPlaceRateCount+oddsRangeUnHitCountSlice[0])
			secondPlaceRate = float64(secondPlaceRateCount) * 100 / float64(secondPlaceRateCount+oddsRangeUnHitCountSlice[9])
			thirdPlaceRate = float64(thirdPlaceRateCount) * 100 / float64(thirdPlaceRateCount+oddsRangeUnHitCountSlice[18])
		} else if inexactOdds >= 1.5 && inexactOdds <= 1.9 {
			firstPlaceRateCount := oddsRangeHitCountSlice[1]
			secondPlaceRateCount := oddsRangeHitCountSlice[1] + oddsRangeHitCountSlice[10]
			thirdPlaceRateCount := oddsRangeHitCountSlice[1] + oddsRangeHitCountSlice[10] + oddsRangeHitCountSlice[19]
			firstPlaceRate = float64(firstPlaceRateCount) * 100 / float64(firstPlaceRateCount+oddsRangeUnHitCountSlice[1])
			secondPlaceRate = float64(secondPlaceRateCount) * 100 / float64(secondPlaceRateCount+oddsRangeUnHitCountSlice[10])
			thirdPlaceRate = float64(thirdPlaceRateCount) * 100 / float64(thirdPlaceRateCount+oddsRangeUnHitCountSlice[19])
		} else if inexactOdds >= 2.0 && inexactOdds <= 2.2 {
			firstPlaceRateCount := oddsRangeHitCountSlice[2]
			secondPlaceRateCount := oddsRangeHitCountSlice[2] + oddsRangeHitCountSlice[11]
			thirdPlaceRateCount := oddsRangeHitCountSlice[2] + oddsRangeHitCountSlice[11] + oddsRangeHitCountSlice[20]
			firstPlaceRate = float64(firstPlaceRateCount) * 100 / float64(firstPlaceRateCount+oddsRangeUnHitCountSlice[2])
			secondPlaceRate = float64(secondPlaceRateCount) * 100 / float64(secondPlaceRateCount+oddsRangeUnHitCountSlice[11])
			thirdPlaceRate = float64(thirdPlaceRateCount) * 100 / float64(thirdPlaceRateCount+oddsRangeUnHitCountSlice[20])
		} else if inexactOdds >= 2.3 && inexactOdds <= 3.0 {
			firstPlaceRateCount := oddsRangeHitCountSlice[3]
			secondPlaceRateCount := oddsRangeHitCountSlice[3] + oddsRangeHitCountSlice[12]
			thirdPlaceRateCount := oddsRangeHitCountSlice[3] + oddsRangeHitCountSlice[12] + oddsRangeHitCountSlice[21]
			firstPlaceRate = float64(firstPlaceRateCount) * 100 / float64(firstPlaceRateCount+oddsRangeUnHitCountSlice[3])
			secondPlaceRate = float64(secondPlaceRateCount) * 100 / float64(secondPlaceRateCount+oddsRangeUnHitCountSlice[12])
			thirdPlaceRate = float64(thirdPlaceRateCount) * 100 / float64(thirdPlaceRateCount+oddsRangeUnHitCountSlice[21])
		} else if inexactOdds >= 3.1 && inexactOdds <= 4.9 {
			firstPlaceRateCount := oddsRangeHitCountSlice[4]
			secondPlaceRateCount := oddsRangeHitCountSlice[4] + oddsRangeHitCountSlice[13]
			thirdPlaceRateCount := oddsRangeHitCountSlice[4] + oddsRangeHitCountSlice[13] + oddsRangeHitCountSlice[22]
			firstPlaceRate = float64(firstPlaceRateCount) * 100 / float64(firstPlaceRateCount+oddsRangeUnHitCountSlice[4])
			secondPlaceRate = float64(secondPlaceRateCount) * 100 / float64(secondPlaceRateCount+oddsRangeUnHitCountSlice[13])
			thirdPlaceRate = float64(thirdPlaceRateCount) * 100 / float64(thirdPlaceRateCount+oddsRangeUnHitCountSlice[22])
		} else if inexactOdds >= 5.0 && inexactOdds <= 9.9 {
			firstPlaceRateCount := oddsRangeHitCountSlice[5]
			secondPlaceRateCount := oddsRangeHitCountSlice[5] + oddsRangeHitCountSlice[14]
			thirdPlaceRateCount := oddsRangeHitCountSlice[5] + oddsRangeHitCountSlice[14] + oddsRangeHitCountSlice[23]
			firstPlaceRate = float64(firstPlaceRateCount) * 100 / float64(firstPlaceRateCount+oddsRangeUnHitCountSlice[5])
			secondPlaceRate = float64(secondPlaceRateCount) * 100 / float64(secondPlaceRateCount+oddsRangeUnHitCountSlice[14])
			thirdPlaceRate = float64(thirdPlaceRateCount) * 100 / float64(thirdPlaceRateCount+oddsRangeUnHitCountSlice[23])
		} else if inexactOdds >= 10.0 && inexactOdds <= 19.9 {
			firstPlaceRateCount := oddsRangeHitCountSlice[6]
			secondPlaceRateCount := oddsRangeHitCountSlice[6] + oddsRangeHitCountSlice[15]
			thirdPlaceRateCount := oddsRangeHitCountSlice[6] + oddsRangeHitCountSlice[15] + oddsRangeHitCountSlice[24]
			firstPlaceRate = float64(firstPlaceRateCount) * 100 / float64(firstPlaceRateCount+oddsRangeUnHitCountSlice[6])
			secondPlaceRate = float64(secondPlaceRateCount) * 100 / float64(secondPlaceRateCount+oddsRangeUnHitCountSlice[15])
			thirdPlaceRate = float64(thirdPlaceRateCount) * 100 / float64(thirdPlaceRateCount+oddsRangeUnHitCountSlice[24])
		} else if inexactOdds >= 20.0 && inexactOdds <= 49.9 {
			firstPlaceRateCount := oddsRangeHitCountSlice[7]
			secondPlaceRateCount := oddsRangeHitCountSlice[7] + oddsRangeHitCountSlice[16]
			thirdPlaceRateCount := oddsRangeHitCountSlice[7] + oddsRangeHitCountSlice[16] + oddsRangeHitCountSlice[25]
			firstPlaceRate = float64(firstPlaceRateCount) * 100 / float64(firstPlaceRateCount+oddsRangeUnHitCountSlice[7])
			secondPlaceRate = float64(secondPlaceRateCount) * 100 / float64(secondPlaceRateCount+oddsRangeUnHitCountSlice[16])
			thirdPlaceRate = float64(thirdPlaceRateCount) * 100 / float64(thirdPlaceRateCount+oddsRangeUnHitCountSlice[25])
		} else if inexactOdds >= 50.0 {
			firstPlaceRateCount := oddsRangeHitCountSlice[8]
			secondPlaceRateCount := oddsRangeHitCountSlice[8] + oddsRangeHitCountSlice[17]
			thirdPlaceRateCount := oddsRangeHitCountSlice[8] + oddsRangeHitCountSlice[17] + oddsRangeHitCountSlice[26]
			firstPlaceRate = float64(firstPlaceRateCount) * 100 / float64(firstPlaceRateCount+oddsRangeUnHitCountSlice[8])
			secondPlaceRate = float64(secondPlaceRateCount) * 100 / float64(secondPlaceRateCount+oddsRangeUnHitCountSlice[17])
			thirdPlaceRate = float64(thirdPlaceRateCount) * 100 / float64(thirdPlaceRateCount+oddsRangeUnHitCountSlice[26])
		}

		placeRateMap[raceResult.HorseId()][0] = firstPlaceRate
		placeRateMap[raceResult.HorseId()][1] = secondPlaceRate
		placeRateMap[raceResult.HorseId()][2] = thirdPlaceRate
	}

	return placeRateMap
}

func (p *placeUnHitService) GetCheckList(
	ctx context.Context,
	race *analysis_entity.Race,
	horse *data_cache_entity.Horse,
	raceForecast *data_cache_entity.RaceForecast,
) error {
	//TODO implement me
	panic("implement me")
}

func (p *placeUnHitService) GetWinRedOdds(
	ctx context.Context,
	oddsList []*data_cache_entity.Odds,
	thresholdOdds decimal.Decimal,
	raceMap map[types.RaceId]*analysis_entity.Race,
) ([]*analysis_entity.Odds, error) {
	winOdds := make([]*analysis_entity.Odds, 0)
	for _, odds := range oddsList {
		if _, ok := raceMap[odds.RaceId()]; !ok {
			continue
		}

		newWinOdds := make([]string, 0)
		for _, oddsStr := range odds.Odds() {
			decimalOdds, err := decimal.NewFromString(oddsStr)
			if err != nil {
				return nil, err
			}
			if decimalOdds.LessThan(thresholdOdds) {
				newWinOdds = append(newWinOdds, oddsStr)
			}
		}
		if len(newWinOdds) > 0 {
			newOdds, err := analysis_entity.NewOdds(
				odds.RaceId(),
				odds.RaceDate(),
				odds.TicketType(),
				odds.Number(),
				odds.PopularNumber(),
				newWinOdds,
			)
			if err != nil {
				return nil, err
			}
			winOdds = append(winOdds, newOdds)
		}
	}

	return winOdds, nil
}

func (p *placeUnHitService) GetWinOddsFaults(
	ctx context.Context,
	oddsList []*data_cache_entity.Odds,
	raceMap map[types.RaceId]*analysis_entity.Race,
) ([]*analysis_entity.OddsFault, error) {
	oddsFaults := make([]*analysis_entity.OddsFault, 0)
	if len(oddsList) < 2 {
		return nil, fmt.Errorf("oddsList must contain at least two elements")
	}

	for i := range oddsList {
		if _, ok := raceMap[oddsList[i].RaceId()]; !ok {
			continue
		}

		if i == len(oddsList)-1 {
			break
		}
		oddsFault, err := analysis_entity.NewOddsFault(
			oddsList[i].RaceId(),
			oddsList[i].Odds()[0],
			oddsList[i+1].Odds()[0],
			oddsList[i].PopularNumber(),
		)
		if err != nil {
			return nil, err
		}
		oddsFaults = append(oddsFaults, oddsFault)
	}

	return oddsFaults, nil
}

func (p *placeUnHitService) GetTrioOdds(
	ctx context.Context,
	oddsList []*data_cache_entity.Odds,
	popularNumber int,
	raceMap map[types.RaceId]*analysis_entity.Race,
) ([]*analysis_entity.Odds, error) {
	trioOdds := make([]*analysis_entity.Odds, 0)
	for _, odds := range oddsList {
		if _, ok := raceMap[odds.RaceId()]; !ok {
			continue
		}
		if odds.PopularNumber() == popularNumber {
			newTrioOdds, err := analysis_entity.NewOdds(
				odds.RaceId(),
				odds.RaceDate(),
				odds.TicketType(),
				odds.Number(),
				odds.PopularNumber(),
				odds.Odds(),
			)
			if err != nil {
				return nil, err
			}
			trioOdds = append(trioOdds, newTrioOdds)
		}
	}

	return trioOdds, nil
}

func (p *placeUnHitService) GetQuinellaOddsWheelCombinations(
	ctx context.Context,
	oddsList []*data_cache_entity.Odds,
	winOddsMap map[types.RaceId][]*analysis_entity.Odds,
	raceMap map[types.RaceId]*analysis_entity.Race,
) ([]*analysis_entity.Odds, error) {
	quinellaWheelCombinationsMap := make(map[types.RaceId][]*analysis_entity.Odds)
	for _, odds := range oddsList {
		if _, ok := raceMap[odds.RaceId()]; !ok {
			continue
		}
		if _, ok := quinellaWheelCombinationsMap[odds.RaceId()]; !ok {
			quinellaWheelCombinationsMap[odds.RaceId()] = make([]*analysis_entity.Odds, 0)
		}
		newOdds, err := analysis_entity.NewOdds(
			odds.RaceId(),
			odds.RaceDate(),
			odds.TicketType(),
			odds.Number(),
			odds.PopularNumber(),
			odds.Odds(),
		)
		if err != nil {
			return nil, err
		}
		quinellaWheelCombinationsMap[odds.RaceId()] = append(quinellaWheelCombinationsMap[odds.RaceId()], newOdds)
	}

	quinellaWheelCombinations := make([]*analysis_entity.Odds, 0)
	for raceId, quinellaOddsList := range quinellaWheelCombinationsMap {
		winOddsList, ok := winOddsMap[raceId]
		if !ok {
			return nil, fmt.Errorf("winOddsMap not found: %s", raceId)
		}
		if len(winOddsList) == 0 {
			return nil, fmt.Errorf("winOddsList is empty: %s", raceId)
		}

		var firstWinOdds *analysis_entity.Odds
		for _, winOdds := range winOddsList {
			if winOdds.PopularNumber() == 1 {
				firstWinOdds = winOdds
				break
			}
		}

		if firstWinOdds == nil {
			return nil, fmt.Errorf("firstWinOdds not found: %s", raceId)
		}

		for _, odds := range quinellaOddsList {
			horseNumber := firstWinOdds.Number().List()[0]
			if slices.Contains(odds.Number().List(), horseNumber) {
				quinellaWheelCombinations = append(quinellaWheelCombinations, odds)
			}
		}
	}

	return quinellaWheelCombinations, nil
}

func (p *placeUnHitService) FetchHorse(
	ctx context.Context,
	horseId types.HorseId,
) (*netkeiba_entity.Horse, error) {
	horse, err := p.horseRepository.Fetch(ctx, fmt.Sprintf(horseUrl, horseId))
	if err != nil {
		return nil, err
	}

	return horse, nil
}

func (p *placeUnHitService) FetchRaceForecasts(
	ctx context.Context,
	raceId types.RaceId,
) ([]*tospo_entity.Forecast, error) {
	forecasts, err := p.raceForecastRepository.FetchRaceForecast(ctx, fmt.Sprintf(raceForecastUrl, raceId))
	if err != nil {
		return nil, err
	}

	return forecasts, nil
}

func (p *placeUnHitService) FetchTrainingComments(
	ctx context.Context,
	raceId types.RaceId,
) ([]*tospo_entity.TrainingComment, error) {
	trainingComments, err := p.raceForecastRepository.FetchTrainingComment(ctx, fmt.Sprintf(raceTrainingCommentUrl, raceId))
	if err != nil {
		return nil, err
	}

	return trainingComments, nil
}

func (p *placeUnHitService) CreateUnhitRaces(ctx context.Context,
	races []*analysis_entity.Race,
	raceRateMap map[types.RaceId]map[types.HorseId][]float64,
	raceForecastMap map[types.RaceId]*data_cache_entity.RaceForecast,
	horseMap map[types.HorseId]*data_cache_entity.Horse,
	winMultiOddsMap map[types.RaceId][]*analysis_entity.Odds,
	winOddsFaultMap map[types.RaceId][]*analysis_entity.OddsFault,
	trioOddsMap map[types.RaceId]*analysis_entity.Odds,
	quinellaConsecutiveNumberMap map[types.RaceId]int,
	quinellaCombinationTotalOddsMap map[types.RaceId]decimal.Decimal,
) ([]*spreadsheet_entity.AnalysisPlaceUnhit, error) {
	placeUnHitEntites := make([]*spreadsheet_entity.AnalysisPlaceUnhit, 0, len(races))
	for _, race := range races {
		raceForecast, ok := raceForecastMap[race.RaceId()]
		if !ok {
			return nil, fmt.Errorf("raceForecast not found: %s", race.RaceId())
		}
		raceForecastMap := converter.ConvertToMap(raceForecast.Forecasts(), func(raceForecast *data_cache_entity.Forecast) types.HorseNumber {
			return raceForecast.HorseNumber()
		})

		winOddsFaults, ok := winOddsFaultMap[race.RaceId()]
		if !ok {
			return nil, fmt.Errorf("winOddsFaultMap not found: %s", race.RaceId())
		}

		for _, raceResult := range race.RaceResults() {
			raceForecast, ok := raceForecastMap[raceResult.HorseNumber()]
			if !ok {
				return nil, fmt.Errorf("raceForecast not found: raceId: %s, horseNumber: %d", race.RaceId(), raceResult.HorseNumber())
			}

			winMultiOdds := winMultiOddsMap[race.RaceId()]
			trioOdds := trioOddsMap[race.RaceId()]
			quinellaConsecutiveNumber := quinellaConsecutiveNumberMap[race.RaceId()]
			quinellaCombinationTotalOdds := quinellaCombinationTotalOddsMap[race.RaceId()]

			placeUnHitEntites = append(placeUnHitEntites, spreadsheet_entity.NewAnalysisPlaceUnhit(
				race.RaceId(),
				race.RaceDate(),
				race.RaceNumber(),
				race.RaceCourse(),
				race.RaceName(),
				race.Class(),
				race.CourseCategory(),
				race.Distance(),
				race.RaceWeightCondition(),
				race.TrackCondition(),
				race.Entries(),
				raceResult.HorseNumber(),
				raceResult.HorseId(),
				raceResult.HorseName(),
				raceResult.JockeyId(),
				raceResult.JockeyName(),
				raceResult.PopularNumber(),
				raceResult.Odds(),
				raceResult.OrderNo(),
				raceResult.JockeyWeight(),
				raceResult.HorseWeight(),
				raceResult.HorseWeightAdd(),
				func() *decimal.Decimal {
					if trioOdds != nil {
						odds := trioOdds.Odds()
						return &odds
					}
					return nil
				}(),
				len(winMultiOdds),
				winOddsFaults[0].OddsFault(),
				winOddsFaults[1].OddsFault(),
				quinellaConsecutiveNumber,
				quinellaCombinationTotalOdds,
				raceForecast.TrainingComment(),
			))
		}
	}

	return placeUnHitEntites, nil
}

func (p *placeUnHitService) CreateCheckPoints(
	ctx context.Context,
	race *analysis_entity.Race,
) error {

	p.placeCheckPointService.GetNegativePoint(ctx, race)

	// 実装を追加
	return nil
}

func (p *placeUnHitService) Write(
	ctx context.Context,
	analysisPlaceUnhits []*spreadsheet_entity.AnalysisPlaceUnhit,
) error {
	return p.spreadSheetRepository.WriteAnalysisPlaceUnhit(ctx, analysisPlaceUnhits)
}
