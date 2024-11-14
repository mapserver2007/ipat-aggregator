package analysis_service

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/netkeiba_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/tospo_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/filter_service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
	"github.com/mapserver2007/ipat-aggregator/config"
)

type PlaceUnHit interface {
	GetUnHitRaces(ctx context.Context, markers []*marker_csv_entity.AnalysisMarker, races []*data_cache_entity.Race) []*analysis_entity.Race
	GetUnHitRaceRate(ctx context.Context, race *analysis_entity.Race, calculables []*analysis_entity.PlaceCalculable) map[types.HorseId][]float64
	FetchHorse(ctx context.Context, horseId types.HorseId) (*netkeiba_entity.Horse, error)
	FetchRaceForecasts(ctx context.Context, raceId types.RaceId) ([]*tospo_entity.Forecast, error)
	FetchTrainingComments(ctx context.Context, raceId types.RaceId) ([]*tospo_entity.TrainingComment, error)
	Convert(ctx context.Context) error
}

type placeUnHitService struct {
	horseRepository        repository.HorseRepository
	raceForecastRepository repository.RaceForecastRepository
	horseEntityConverter   converter.HorseEntityConverter
	filterService          filter_service.AnalysisFilter
}

func NewPlaceUnHit(
	horseRepository repository.HorseRepository,
	raceForecastRepository repository.RaceForecastRepository,
	horseEntityConverter converter.HorseEntityConverter,
	filterService filter_service.AnalysisFilter,
) PlaceUnHit {
	return &placeUnHitService{
		horseRepository:        horseRepository,
		raceForecastRepository: raceForecastRepository,
		horseEntityConverter:   horseEntityConverter,
		filterService:          filterService,
	}
}

func (p *placeUnHitService) GetUnHitRaces(
	ctx context.Context,
	markers []*marker_csv_entity.AnalysisMarker,
	races []*data_cache_entity.Race,
) []*analysis_entity.Race {
	markerMap := converter.ConvertToMap(markers, func(marker *marker_csv_entity.AnalysisMarker) types.RaceId {
		return marker.RaceId()
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
				analysisRaceResults = append(analysisRaceResults, analysis_entity.NewRaceResult(
					raceResult.OrderNo(),
					raceResult.HorseId(),
					raceResult.HorseName(),
					raceResult.HorseNumber(),
					raceResult.JockeyId(),
					raceResult.Odds(),
					raceResult.PopularNumber(),
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
	var analysisFilter filter.Id
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

func (p *placeUnHitService) Convert(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
