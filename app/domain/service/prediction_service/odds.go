package prediction_service

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/filter_service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/shopspring/decimal"
	"time"
)

const (
	raceCardUrl   = "https://race.netkeiba.com/race/shutuba.html?race_id=%s"
	oddsUrl       = "https://race.netkeiba.com/api/api_get_jra_odds.html?race_id=%s&type=1&action=update"
	raceResultUrl = "https://race.netkeiba.com/race/result.html?race_id=%s&organizer=1&race_date=%s"
)

type Odds interface {
	Get(ctx context.Context, raceId types.RaceId) (*prediction_entity.Race, error)
	Convert(ctx context.Context, predictionRaces []*prediction_entity.Race, predictionMarkers []*marker_csv_entity.PredictionMarker, calculables []*analysis_entity.PlaceCalculable) (map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace, map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace, map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace, map[types.RaceCourse][]types.RaceId)
	Write(ctx context.Context, firstPlaceMap, secondPlaceMap, thirdPlaceMap map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace, raceCourseMap map[types.RaceCourse][]types.RaceId) error
}

type oddsService struct {
	oddRepository         repository.OddsRepository
	raceRepository        repository.RaceRepository
	spreadSheetRepository repository.SpreadSheetRepository
	filterService         filter_service.PredictionFilter
}

func NewOdds(
	oddRepository repository.OddsRepository,
	raceRepository repository.RaceRepository,
	spreadSheetRepository repository.SpreadSheetRepository,
	filterService filter_service.PredictionFilter,
) Odds {
	return &oddsService{
		oddRepository:         oddRepository,
		raceRepository:        raceRepository,
		spreadSheetRepository: spreadSheetRepository,
		filterService:         filterService,
	}
}

func (p *oddsService) Get(
	ctx context.Context,
	raceId types.RaceId,
) (*prediction_entity.Race, error) {
	odds, err := p.oddRepository.Fetch(ctx, fmt.Sprintf(oddsUrl, raceId))
	if err != nil {
		return nil, err
	}

	raceCard, err := p.raceRepository.FetchRaceCard(ctx, fmt.Sprintf(raceCardUrl, raceId))
	if err != nil {
		return nil, err
	}

	raceDate := time.Now().Format("20060102")
	race, err := p.raceRepository.FetchRace(ctx, fmt.Sprintf(raceResultUrl, raceId, raceDate))
	if err != nil {
		return nil, err
	}

	var predictionOdds []*prediction_entity.Odds
	for _, nkOdds := range odds {
		predictionOdds = append(predictionOdds, prediction_entity.NewOdds(
			nkOdds.Odds(),
			nkOdds.PopularNumber(),
			nkOdds.HorseNumbers()[0],
		))
	}

	// レース結果のうち、必要なのは着順に対する馬番のみ
	raceResultHorseNumbers := make([]int, 3)
	if race.RaceResults() != nil && len(race.RaceResults()) >= 3 {
		for idx, raceResult := range race.RaceResults()[:3] {
			raceResultHorseNumbers[idx] = raceResult.HorseNumber()
		}
	}

	predictionRace := prediction_entity.NewRace(
		raceCard.RaceId(),
		raceCard.RaceName(),
		raceCard.RaceNumber(),
		raceCard.Entries(),
		raceCard.Distance(),
		raceCard.Class(),
		raceCard.CourseCategory(),
		raceCard.TrackCondition(),
		raceCard.RaceSexCondition(),
		raceCard.RaceWeightCondition(),
		raceCard.RaceCourseId(),
		raceCard.Url(),
		raceResultHorseNumbers,
		predictionOdds,
		p.filterService.Create(ctx, raceCard),
	)

	return predictionRace, nil
}

func (p *oddsService) Convert(
	ctx context.Context,
	predictionRaces []*prediction_entity.Race,
	predictionMarkers []*marker_csv_entity.PredictionMarker,
	calculables []*analysis_entity.PlaceCalculable,
) (
	map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace,
	map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace,
	map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace,
	map[types.RaceCourse][]types.RaceId,
) {
	firstPlaceMap := map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace{}
	secondPlaceMap := map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace{}
	thirdPlaceMap := map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace{}
	raceCourseMap := map[types.RaceCourse][]types.RaceId{}

	predictionMarkerMap := map[types.RaceId]*marker_csv_entity.PredictionMarker{}
	for _, marker := range predictionMarkers {
		predictionMarkerMap[marker.RaceId()] = marker
	}

	for _, race := range predictionRaces {
		if _, ok := raceCourseMap[race.RaceCourseId()]; !ok {
			raceCourseMap[race.RaceCourseId()] = make([]types.RaceId, 0)
		}
		predictionRace := spreadsheet_entity.NewPredictionRace(
			race.RaceId(),
			race.RaceName(),
			race.RaceNumber(),
			race.RaceCourseId(),
			race.CourseCategory(),
			race.Url(),
			race.PredictionFilter(),
		)
		horseNumberOddsMap := map[types.HorseNumber]decimal.Decimal{}
		for _, o := range race.Odds() {
			horseNumberOddsMap[o.HorseNumber()] = o.Odds()
		}

		firstPlaceMap[*predictionRace] = map[types.Marker]*spreadsheet_entity.PredictionPlace{}
		secondPlaceMap[*predictionRace] = map[types.Marker]*spreadsheet_entity.PredictionPlace{}
		thirdPlaceMap[*predictionRace] = map[types.Marker]*spreadsheet_entity.PredictionPlace{}
		raceCourseMap[race.RaceCourseId()] = append(raceCourseMap[race.RaceCourseId()], race.RaceId())
		predictionMarker := predictionMarkerMap[race.RaceId()]

		for _, marker := range []types.Marker{types.Favorite, types.Rival, types.BrackTriangle, types.WhiteTriangle, types.Star, types.Check} {
			raceIdMap := map[types.RaceId]bool{}
			oddsRangeHitCountSlice := make([]int, 24)
			oddsRangeUnHitCountSlice := make([]int, 24)
			onOddsRangeSlice := make([]bool, 8)
			firstPlaceHitOddsRangeSlice := make([]bool, 8)
			secondPlaceHitOddsRangeSlice := make([]bool, 8)
			thirdPlaceHitOddsRangeSlice := make([]bool, 8)

			markerHorseNumber := predictionMarker.MarkerMap()[marker]
			isFirstPlaceHit := race.RaceResultHorseNumbers()[0] == markerHorseNumber
			isSecondPlaceHit := isFirstPlaceHit || race.RaceResultHorseNumbers()[1] == markerHorseNumber
			isThirdPlaceHit := isFirstPlaceHit || isSecondPlaceHit || race.RaceResultHorseNumbers()[2] == markerHorseNumber
			realTimeOdds := horseNumberOddsMap[markerHorseNumber].InexactFloat64()
			if realTimeOdds >= 1.0 && realTimeOdds <= 1.5 {
				onOddsRangeSlice[0] = true
				firstPlaceHitOddsRangeSlice[0] = isFirstPlaceHit
				secondPlaceHitOddsRangeSlice[0] = isSecondPlaceHit
				thirdPlaceHitOddsRangeSlice[0] = isThirdPlaceHit
			} else if realTimeOdds >= 1.6 && realTimeOdds <= 2.0 {
				onOddsRangeSlice[1] = true
				firstPlaceHitOddsRangeSlice[1] = isFirstPlaceHit
				secondPlaceHitOddsRangeSlice[1] = isSecondPlaceHit
				thirdPlaceHitOddsRangeSlice[1] = isThirdPlaceHit
			} else if realTimeOdds >= 2.1 && realTimeOdds <= 2.9 {
				onOddsRangeSlice[2] = true
				firstPlaceHitOddsRangeSlice[2] = isFirstPlaceHit
				secondPlaceHitOddsRangeSlice[2] = isSecondPlaceHit
				thirdPlaceHitOddsRangeSlice[2] = isThirdPlaceHit
			} else if realTimeOdds >= 3.0 && realTimeOdds <= 4.9 {
				onOddsRangeSlice[3] = true
				firstPlaceHitOddsRangeSlice[3] = isFirstPlaceHit
				secondPlaceHitOddsRangeSlice[3] = isSecondPlaceHit
				thirdPlaceHitOddsRangeSlice[3] = isThirdPlaceHit
			} else if realTimeOdds >= 5.0 && realTimeOdds <= 9.9 {
				onOddsRangeSlice[4] = true
				firstPlaceHitOddsRangeSlice[4] = isFirstPlaceHit
				secondPlaceHitOddsRangeSlice[4] = isSecondPlaceHit
				thirdPlaceHitOddsRangeSlice[4] = isThirdPlaceHit
			} else if realTimeOdds >= 10.0 && realTimeOdds <= 19.9 {
				onOddsRangeSlice[5] = true
				firstPlaceHitOddsRangeSlice[5] = isFirstPlaceHit
				secondPlaceHitOddsRangeSlice[5] = isSecondPlaceHit
				thirdPlaceHitOddsRangeSlice[5] = isThirdPlaceHit
			} else if realTimeOdds >= 20.0 && realTimeOdds <= 49.9 {
				onOddsRangeSlice[6] = true
				firstPlaceHitOddsRangeSlice[6] = isFirstPlaceHit
				secondPlaceHitOddsRangeSlice[6] = isSecondPlaceHit
				thirdPlaceHitOddsRangeSlice[6] = isThirdPlaceHit
			} else if realTimeOdds >= 50.0 {
				onOddsRangeSlice[7] = true
				firstPlaceHitOddsRangeSlice[7] = isFirstPlaceHit
				secondPlaceHitOddsRangeSlice[7] = isSecondPlaceHit
				thirdPlaceHitOddsRangeSlice[7] = isThirdPlaceHit
			}

			for _, calculable := range calculables {
				if calculable.Marker() != marker {
					continue
				}

				match := true
				for _, f := range calculable.Filters() {
					if f&race.PredictionFilter() == 0 {
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
				race.PredictionFilter(),
				len(raceIdMap),
			)
			secondPlaceOddsRangeHitCountData := spreadsheet_entity.NewPlaceHitCountData(
				secondPlaceOddsRangeHitCountSlice,
				race.PredictionFilter(),
				len(raceIdMap),
			)
			thirdPlaceOddsRangeHitCountData := spreadsheet_entity.NewPlaceHitCountData(
				thirdPlaceOddsRangeHitCountSlice,
				race.PredictionFilter(),
				len(raceIdMap),
			)

			firstPlaceOddsRangeUnHitCountData := spreadsheet_entity.NewPlaceUnHitCountData(
				firstPlaceOddsRangeUnHitCountSlice,
				race.PredictionFilter(),
				len(raceIdMap),
			)
			secondPlaceOddsRangeUnHitCountData := spreadsheet_entity.NewPlaceUnHitCountData(
				secondPlaceOddsRangeUnHitCountSlice,
				race.PredictionFilter(),
				len(raceIdMap),
			)
			thirdPlaceOddsRangeUnHitCountData := spreadsheet_entity.NewPlaceUnHitCountData(
				thirdPlaceOddsRangeUnHitCountSlice,
				race.PredictionFilter(),
				len(raceIdMap),
			)

			firstPlaceOddsRangeRateData := spreadsheet_entity.NewPredictionRateData(
				firstPlaceOddsRangeHitCountData,
				firstPlaceOddsRangeUnHitCountData,
				firstPlaceHitOddsRangeSlice,
			)
			firstPlaceRateStyle := spreadsheet_entity.NewPredictionRateStyle(onOddsRangeSlice)

			secondPlaceOddsRangeRateData := spreadsheet_entity.NewPredictionRateData(
				secondPlaceOddsRangeHitCountData,
				secondPlaceOddsRangeUnHitCountData,
				secondPlaceHitOddsRangeSlice,
			)
			secondPlaceRateStyle := spreadsheet_entity.NewPredictionRateStyle(onOddsRangeSlice)

			thirdPlaceOddsRangeRateData := spreadsheet_entity.NewPredictionRateData(
				thirdPlaceOddsRangeHitCountData,
				thirdPlaceOddsRangeUnHitCountData,
				thirdPlaceHitOddsRangeSlice,
			)
			thirdPlaceRateStyle := spreadsheet_entity.NewPredictionRateStyle(onOddsRangeSlice)

			firstPlaceMap[*predictionRace][marker] = spreadsheet_entity.NewPredictionPlace(
				firstPlaceOddsRangeRateData,
				firstPlaceRateStyle,
			)
			secondPlaceMap[*predictionRace][marker] = spreadsheet_entity.NewPredictionPlace(
				secondPlaceOddsRangeRateData,
				secondPlaceRateStyle,
			)
			thirdPlaceMap[*predictionRace][marker] = spreadsheet_entity.NewPredictionPlace(
				thirdPlaceOddsRangeRateData,
				thirdPlaceRateStyle,
			)
		}
	}

	return firstPlaceMap, secondPlaceMap, thirdPlaceMap, raceCourseMap
}

func (p *oddsService) Write(
	ctx context.Context,
	firstPlaceMap,
	secondPlaceMap,
	thirdPlaceMap map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace,
	raceCourseMap map[types.RaceCourse][]types.RaceId,
) error {
	return p.spreadSheetRepository.WritePrediction(ctx, firstPlaceMap, secondPlaceMap, thirdPlaceMap, raceCourseMap)
}
