package prediction_service

import (
	"context"
	"fmt"
	"time"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/marker_csv_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/filter_service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types/filter"
	"github.com/shopspring/decimal"
)

type Odds interface {
	GetRace(ctx context.Context, raceId types.RaceId) (*prediction_entity.Race, error)
	Convert(
		ctx context.Context,
		race *prediction_entity.Race,
		horseNumber types.HorseNumber,
		marker types.Marker,
		calculables []*analysis_entity.PlaceCalculable,
	) []*spreadsheet_entity.PredictionPlace
	ConvertAll(ctx context.Context,
		predictionRaces []*prediction_entity.Race,
		predictionMarkers []*marker_csv_entity.PredictionMarker,
		placeCalculables []*analysis_entity.PlaceCalculable,
		raceTimeMap map[filter.AttributeId]*spreadsheet_entity.AnalysisRaceTime,
	) (
		map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace,
		map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace,
		map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace, map[types.RaceCourse][]types.RaceId)
	Write(ctx context.Context,
		firstPlaceMap,
		secondPlaceMap,
		thirdPlaceMap map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace,
		raceCourseMap map[types.RaceCourse][]types.RaceId,
	) error
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

func (p *oddsService) GetRace(
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
			nkOdds.Odds()[0],
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
		raceCard.RaceDate(),
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
		nil,
		raceResultHorseNumbers,
		predictionOdds,
		p.filterService.CreateRaceConditionFilters(ctx, raceCard),
		p.filterService.CreateRaceTimeConditionFilters(ctx, raceCard),
	)

	return predictionRace, nil
}

func (p *oddsService) Convert(
	ctx context.Context,
	race *prediction_entity.Race,
	horseNumber types.HorseNumber,
	marker types.Marker,
	calculables []*analysis_entity.PlaceCalculable,
) []*spreadsheet_entity.PredictionPlace {
	horseNumberOddsMap := map[types.HorseNumber]decimal.Decimal{}
	for _, o := range race.Odds() {
		horseNumberOddsMap[o.HorseNumber()] = o.Odds()
	}

	var predictionFilter filter.AttributeId
	for _, f := range race.RaceConditionFilters() {
		predictionFilter |= f
	}

	predictionRace := spreadsheet_entity.NewPredictionRace(
		race.RaceId(),
		race.RaceName(),
		race.RaceNumber(),
		race.RaceCourse(),
		race.CourseCategory(),
		race.Url(),
		race.RaceConditionFilters(),
		nil, // TODO 後ほど足す
	)

	firstPlaceMap := map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace{}
	secondPlaceMap := map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace{}
	thirdPlaceMap := map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace{}
	firstPlaceMap[*predictionRace] = map[types.Marker]*spreadsheet_entity.PredictionPlace{}
	secondPlaceMap[*predictionRace] = map[types.Marker]*spreadsheet_entity.PredictionPlace{}
	thirdPlaceMap[*predictionRace] = map[types.Marker]*spreadsheet_entity.PredictionPlace{}

	raceIdMap := map[types.RaceId]bool{}
	oddsRangeHitCountSlice := make([]int, 27)
	oddsRangeUnHitCountSlice := make([]int, 27)

	onOddsRangeSlice := make([]bool, 9)
	firstPlaceHitOddsRangeSlice := make([]bool, 9)
	secondPlaceHitOddsRangeSlice := make([]bool, 9)
	thirdPlaceHitOddsRangeSlice := make([]bool, 9)

	isFirstPlaceHit := race.RaceResultHorseNumbers()[0] == horseNumber
	isSecondPlaceHit := isFirstPlaceHit || race.RaceResultHorseNumbers()[1] == horseNumber
	isThirdPlaceHit := isFirstPlaceHit || isSecondPlaceHit || race.RaceResultHorseNumbers()[2] == horseNumber
	realTimeOdds := horseNumberOddsMap[horseNumber].InexactFloat64()

	if realTimeOdds >= 1.0 && realTimeOdds <= 1.4 {
		onOddsRangeSlice[0] = true
		firstPlaceHitOddsRangeSlice[0] = isFirstPlaceHit
		secondPlaceHitOddsRangeSlice[0] = isSecondPlaceHit
		thirdPlaceHitOddsRangeSlice[0] = isThirdPlaceHit
	} else if realTimeOdds >= 1.5 && realTimeOdds <= 1.9 {
		onOddsRangeSlice[1] = true
		firstPlaceHitOddsRangeSlice[1] = isFirstPlaceHit
		secondPlaceHitOddsRangeSlice[1] = isSecondPlaceHit
		thirdPlaceHitOddsRangeSlice[1] = isThirdPlaceHit
	} else if realTimeOdds >= 2.0 && realTimeOdds <= 2.2 {
		onOddsRangeSlice[2] = true
		firstPlaceHitOddsRangeSlice[2] = isFirstPlaceHit
		secondPlaceHitOddsRangeSlice[2] = isSecondPlaceHit
		thirdPlaceHitOddsRangeSlice[2] = isThirdPlaceHit
	} else if realTimeOdds >= 2.3 && realTimeOdds <= 3.0 {
		onOddsRangeSlice[3] = true
		firstPlaceHitOddsRangeSlice[3] = isFirstPlaceHit
		secondPlaceHitOddsRangeSlice[3] = isSecondPlaceHit
		thirdPlaceHitOddsRangeSlice[3] = isThirdPlaceHit
	} else if realTimeOdds >= 3.1 && realTimeOdds <= 4.9 {
		onOddsRangeSlice[4] = true
		firstPlaceHitOddsRangeSlice[4] = isFirstPlaceHit
		secondPlaceHitOddsRangeSlice[4] = isSecondPlaceHit
		thirdPlaceHitOddsRangeSlice[4] = isThirdPlaceHit
	} else if realTimeOdds >= 5.0 && realTimeOdds <= 9.9 {
		onOddsRangeSlice[5] = true
		firstPlaceHitOddsRangeSlice[5] = isFirstPlaceHit
		secondPlaceHitOddsRangeSlice[5] = isSecondPlaceHit
		thirdPlaceHitOddsRangeSlice[5] = isThirdPlaceHit
	} else if realTimeOdds >= 10.0 && realTimeOdds <= 19.9 {
		onOddsRangeSlice[6] = true
		firstPlaceHitOddsRangeSlice[6] = isFirstPlaceHit
		secondPlaceHitOddsRangeSlice[6] = isSecondPlaceHit
		thirdPlaceHitOddsRangeSlice[6] = isThirdPlaceHit
	} else if realTimeOdds >= 20.0 && realTimeOdds <= 49.9 {
		onOddsRangeSlice[7] = true
		firstPlaceHitOddsRangeSlice[7] = isFirstPlaceHit
		secondPlaceHitOddsRangeSlice[7] = isSecondPlaceHit
		thirdPlaceHitOddsRangeSlice[7] = isThirdPlaceHit
	} else if realTimeOdds >= 50.0 {
		onOddsRangeSlice[8] = true
		firstPlaceHitOddsRangeSlice[8] = isFirstPlaceHit
		secondPlaceHitOddsRangeSlice[8] = isSecondPlaceHit
		thirdPlaceHitOddsRangeSlice[8] = isThirdPlaceHit
	}

	for _, calculable := range calculables {
		if calculable.Marker() != marker {
			continue
		}

		match := true
		for _, f := range calculable.Filters() {
			if f&predictionFilter == 0 {
				match = false
				break
			}
		}
		if match {
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

	predictionPlaces := []*spreadsheet_entity.PredictionPlace{
		spreadsheet_entity.NewPredictionPlace(firstPlaceOddsRangeRateData, firstPlaceRateStyle),
		spreadsheet_entity.NewPredictionPlace(secondPlaceOddsRangeRateData, secondPlaceRateStyle),
		spreadsheet_entity.NewPredictionPlace(thirdPlaceOddsRangeRateData, thirdPlaceRateStyle),
	}

	return predictionPlaces
}

func (p *oddsService) ConvertAll(
	ctx context.Context,
	predictionRaces []*prediction_entity.Race,
	predictionMarkers []*marker_csv_entity.PredictionMarker,
	placeCalculables []*analysis_entity.PlaceCalculable,
	raceTimeMap map[filter.AttributeId]*spreadsheet_entity.AnalysisRaceTime,
) (
	map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace,
	map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace,
	map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace,
	map[types.RaceCourse][]types.RaceId,
) {
	// TODO 内部的にConvertを呼んでループさせたい

	firstPlaceMap := map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace{}
	secondPlaceMap := map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace{}
	thirdPlaceMap := map[spreadsheet_entity.PredictionRace]map[types.Marker]*spreadsheet_entity.PredictionPlace{}
	raceCourseMap := map[types.RaceCourse][]types.RaceId{}

	predictionMarkerMap := map[types.RaceId]*marker_csv_entity.PredictionMarker{}
	for _, marker := range predictionMarkers {
		predictionMarkerMap[marker.RaceId()] = marker
	}

	for _, race := range predictionRaces {
		if _, ok := raceCourseMap[race.RaceCourse()]; !ok {
			raceCourseMap[race.RaceCourse()] = make([]types.RaceId, 0)
		}

		var (
			raceConditionFilter     filter.AttributeId
			raceTimeConditionFilter filter.AttributeId
		)
		for _, f := range race.RaceConditionFilters() {
			raceConditionFilter |= f
		}
		for _, f := range race.RaceTimeConditionFilters() {
			raceTimeConditionFilter |= f
		}

		predictionRaceTime := spreadsheet_entity.InitPredictionRaceTime()
		if raceTime, ok := raceTimeMap[raceTimeConditionFilter]; ok {
			predictionRaceTime = spreadsheet_entity.NewPredictionRaceTime(
				raceTime.AverageRaceTime(),
				raceTime.AverageFirst3f(),
				raceTime.AverageFirst4f(),
				raceTime.AverageLast3f(),
				raceTime.AverageLast4f(),
				raceTime.AverageRap5f(),
			)
		}

		predictionRace := spreadsheet_entity.NewPredictionRace(
			race.RaceId(),
			race.RaceName(),
			race.RaceNumber(),
			race.RaceCourse(),
			race.CourseCategory(),
			race.Url(),
			race.RaceConditionFilters(),
			predictionRaceTime,
		)
		horseNumberOddsMap := map[types.HorseNumber]decimal.Decimal{}
		for _, o := range race.Odds() {
			horseNumberOddsMap[o.HorseNumber()] = o.Odds()
		}

		firstPlaceMap[*predictionRace] = map[types.Marker]*spreadsheet_entity.PredictionPlace{}
		secondPlaceMap[*predictionRace] = map[types.Marker]*spreadsheet_entity.PredictionPlace{}
		thirdPlaceMap[*predictionRace] = map[types.Marker]*spreadsheet_entity.PredictionPlace{}
		raceCourseMap[race.RaceCourse()] = append(raceCourseMap[race.RaceCourse()], race.RaceId())
		predictionMarker := predictionMarkerMap[race.RaceId()]

		for _, marker := range []types.Marker{types.Favorite, types.Rival, types.BrackTriangle, types.WhiteTriangle, types.Star, types.Check} {
			raceIdMap := map[types.RaceId]bool{}
			oddsRangeHitCountSlice := make([]int, 27)
			oddsRangeUnHitCountSlice := make([]int, 27)
			onOddsRangeSlice := make([]bool, 9)
			firstPlaceHitOddsRangeSlice := make([]bool, 9)
			secondPlaceHitOddsRangeSlice := make([]bool, 9)
			thirdPlaceHitOddsRangeSlice := make([]bool, 9)

			markerHorseNumber := predictionMarker.MarkerMap()[marker]
			isFirstPlaceHit := race.RaceResultHorseNumbers()[0] == markerHorseNumber
			isSecondPlaceHit := isFirstPlaceHit || race.RaceResultHorseNumbers()[1] == markerHorseNumber
			isThirdPlaceHit := isFirstPlaceHit || isSecondPlaceHit || race.RaceResultHorseNumbers()[2] == markerHorseNumber
			realTimeOdds := horseNumberOddsMap[markerHorseNumber].InexactFloat64()
			if realTimeOdds >= 1.0 && realTimeOdds <= 1.4 {
				onOddsRangeSlice[0] = true
				firstPlaceHitOddsRangeSlice[0] = isFirstPlaceHit
				secondPlaceHitOddsRangeSlice[0] = isSecondPlaceHit
				thirdPlaceHitOddsRangeSlice[0] = isThirdPlaceHit
			} else if realTimeOdds >= 1.5 && realTimeOdds <= 1.9 {
				onOddsRangeSlice[1] = true
				firstPlaceHitOddsRangeSlice[1] = isFirstPlaceHit
				secondPlaceHitOddsRangeSlice[1] = isSecondPlaceHit
				thirdPlaceHitOddsRangeSlice[1] = isThirdPlaceHit
			} else if realTimeOdds >= 2.0 && realTimeOdds <= 2.2 {
				onOddsRangeSlice[2] = true
				firstPlaceHitOddsRangeSlice[2] = isFirstPlaceHit
				secondPlaceHitOddsRangeSlice[2] = isSecondPlaceHit
				thirdPlaceHitOddsRangeSlice[2] = isThirdPlaceHit
			} else if realTimeOdds >= 2.3 && realTimeOdds <= 3.0 {
				onOddsRangeSlice[3] = true
				firstPlaceHitOddsRangeSlice[3] = isFirstPlaceHit
				secondPlaceHitOddsRangeSlice[3] = isSecondPlaceHit
				thirdPlaceHitOddsRangeSlice[3] = isThirdPlaceHit
			} else if realTimeOdds >= 3.1 && realTimeOdds <= 4.9 {
				onOddsRangeSlice[4] = true
				firstPlaceHitOddsRangeSlice[4] = isFirstPlaceHit
				secondPlaceHitOddsRangeSlice[4] = isSecondPlaceHit
				thirdPlaceHitOddsRangeSlice[4] = isThirdPlaceHit
			} else if realTimeOdds >= 5.0 && realTimeOdds <= 9.9 {
				onOddsRangeSlice[5] = true
				firstPlaceHitOddsRangeSlice[5] = isFirstPlaceHit
				secondPlaceHitOddsRangeSlice[5] = isSecondPlaceHit
				thirdPlaceHitOddsRangeSlice[5] = isThirdPlaceHit
			} else if realTimeOdds >= 10.0 && realTimeOdds <= 19.9 {
				onOddsRangeSlice[6] = true
				firstPlaceHitOddsRangeSlice[6] = isFirstPlaceHit
				secondPlaceHitOddsRangeSlice[6] = isSecondPlaceHit
				thirdPlaceHitOddsRangeSlice[6] = isThirdPlaceHit
			} else if realTimeOdds >= 20.0 && realTimeOdds <= 49.9 {
				onOddsRangeSlice[7] = true
				firstPlaceHitOddsRangeSlice[7] = isFirstPlaceHit
				secondPlaceHitOddsRangeSlice[7] = isSecondPlaceHit
				thirdPlaceHitOddsRangeSlice[7] = isThirdPlaceHit
			} else if realTimeOdds >= 50.0 {
				onOddsRangeSlice[8] = true
				firstPlaceHitOddsRangeSlice[8] = isFirstPlaceHit
				secondPlaceHitOddsRangeSlice[8] = isSecondPlaceHit
				thirdPlaceHitOddsRangeSlice[8] = isThirdPlaceHit
			}

			for _, calculable := range placeCalculables {
				if calculable.Marker() != marker {
					continue
				}

				match := true
				for _, f := range calculable.Filters() {
					if f&raceConditionFilter == 0 {
						match = false
						break
					}
				}
				if match {
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
	return p.spreadSheetRepository.WritePredictionOdds(ctx, firstPlaceMap, secondPlaceMap, thirdPlaceMap, raceCourseMap)
}
