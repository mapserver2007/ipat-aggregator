package prediction_service

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/spreadsheet_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/filter_service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/config"
	"github.com/shopspring/decimal"
)

type PlaceCandidate interface {
	GetRaceCard(ctx context.Context, raceId types.RaceId) (*prediction_entity.Race, error)
	GetRaceForecasts(ctx context.Context, raceId types.RaceId) ([]*prediction_entity.RaceForecast, error)
	GetHorse(ctx context.Context, horseId types.HorseId) (*prediction_entity.Horse, error)
	GetJockey(ctx context.Context, jockeyId types.JockeyId) (*prediction_entity.Jockey, error)
	GetTrainer(ctx context.Context, trainerId types.TrainerId) (*prediction_entity.Trainer, error)
	CreateCheckList(ctx context.Context, race *prediction_entity.Race, horse *prediction_entity.Horse, forecast *prediction_entity.RaceForecast) []bool
	Convert(ctx context.Context, race *prediction_entity.Race, horse *prediction_entity.Horse, jockey *prediction_entity.Jockey, trainer *prediction_entity.Trainer, forecast *prediction_entity.RaceForecast, calculable []*analysis_entity.PlaceCalculable, horseNumber types.HorseNumber, marker types.Marker, checkList []bool) *spreadsheet_entity.PredictionCheckList
	Write(ctx context.Context, predictionCheckList []*spreadsheet_entity.PredictionCheckList) error
}

type placeCandidateService struct {
	raceRepository         repository.RaceRepository
	raceForecastRepository repository.RaceForecastRepository
	horseRepository        repository.HorseRepository
	jockeyRepository       repository.JockeyRepository
	trainerRepository      repository.TrainerRepository
	oddRepository          repository.OddsRepository
	spreadSheetRepository  repository.SpreadSheetRepository
	raceEntityConverter    converter.RaceEntityConverter
	horseEntityConverter   converter.HorseEntityConverter
	filterService          filter_service.PredictionFilter
	placeCheckListService  PlaceCheckList
	predictionOddsService  Odds
}

func NewPlaceCandidate(
	raceRepository repository.RaceRepository,
	raceForecastRepository repository.RaceForecastRepository,
	horseRepository repository.HorseRepository,
	jockeyRepository repository.JockeyRepository,
	trainerRepository repository.TrainerRepository,
	oddRepository repository.OddsRepository,
	spreadSheetRepository repository.SpreadSheetRepository,
	raceEntityConverter converter.RaceEntityConverter,
	horseEntityConverter converter.HorseEntityConverter,
	filterService filter_service.PredictionFilter,
	placeCheckListService PlaceCheckList,
	predictionOddsService Odds,
) PlaceCandidate {
	return &placeCandidateService{
		raceRepository:         raceRepository,
		raceForecastRepository: raceForecastRepository,
		horseRepository:        horseRepository,
		trainerRepository:      trainerRepository,
		jockeyRepository:       jockeyRepository,
		oddRepository:          oddRepository,
		spreadSheetRepository:  spreadSheetRepository,
		raceEntityConverter:    raceEntityConverter,
		horseEntityConverter:   horseEntityConverter,
		filterService:          filterService,
		placeCheckListService:  placeCheckListService,
		predictionOddsService:  predictionOddsService,
	}
}

func (p *placeCandidateService) GetRaceCard(
	ctx context.Context,
	raceId types.RaceId,
) (*prediction_entity.Race, error) {
	// race_resultが取得できない状態でキャッシュさせないように制御する
	rawRace, err := p.raceRepository.FetchRaceCard(ctx, fmt.Sprintf(raceCardUrl+"&cache=false", raceId))
	if err != nil {
		return nil, err
	}

	rawOdds, err := p.oddRepository.Fetch(ctx, fmt.Sprintf(oddsUrl, raceId))
	if err != nil {
		return nil, err
	}

	filters := p.filterService.CreatePredictionOddsFilters(ctx, rawRace)

	race := p.raceEntityConverter.NetKeibaToPrediction(rawRace, rawOdds, filters)

	return race, nil
}

func (p *placeCandidateService) GetRaceForecasts(
	ctx context.Context,
	raceId types.RaceId,
) ([]*prediction_entity.RaceForecast, error) {
	rawRaceForecasts, err := p.raceForecastRepository.FetchRaceForecast(ctx, fmt.Sprintf(raceForecastUrl, raceId))
	if err != nil {
		return nil, err
	}

	rawTrainingComments, err := p.raceForecastRepository.FetchTrainingComment(ctx, fmt.Sprintf(raceTrainingCommentUrl, raceId))
	if err != nil {
		return nil, err
	}

	raceForecasts := make([]*prediction_entity.RaceForecast, len(rawRaceForecasts))
	for idx := range rawRaceForecasts {
		raceForecasts[idx] = p.raceEntityConverter.TospoToPrediction(rawRaceForecasts[idx], rawTrainingComments[idx])
	}

	return raceForecasts, nil
}

func (p *placeCandidateService) GetHorse(
	ctx context.Context,
	horseId types.HorseId,
) (*prediction_entity.Horse, error) {
	rawHorse, err := p.horseRepository.Fetch(ctx, fmt.Sprintf(horseUrl, horseId))
	if err != nil {
		return nil, err
	}

	horse, err := p.horseEntityConverter.NetKeibaToPrediction(rawHorse)
	if err != nil {
		return nil, err
	}

	return horse, nil
}

func (p *placeCandidateService) GetJockey(
	ctx context.Context,
	jockeyId types.JockeyId,
) (*prediction_entity.Jockey, error) {
	rawJockeyInfo, err := p.jockeyRepository.Read(ctx, fmt.Sprintf("%s/%s", config.CacheDir, jockeyFileName))
	if err != nil {
		return nil, err
	}
	if rawJockeyInfo == nil {
		return nil, fmt.Errorf("jockey file is empty")
	}

	var jockey *prediction_entity.Jockey
	for _, rawJockey := range rawJockeyInfo.Jockeys {
		if rawJockey.JockeyId == jockeyId.Value() {
			jockey = prediction_entity.NewJockey(
				rawJockey.JockeyId,
				rawJockey.JockeyName,
			)
		}
	}

	if jockey == nil {
		jockey = prediction_entity.NewJockey("00000", "不明")
	}

	return jockey, nil
}

func (p *placeCandidateService) GetTrainer(
	ctx context.Context,
	trainerId types.TrainerId,
) (*prediction_entity.Trainer, error) {
	rawTrainer, err := p.trainerRepository.Fetch(ctx, fmt.Sprintf(trainerUrl, trainerId))
	if err != nil {
		return nil, err
	}

	return prediction_entity.NewTrainer(
		rawTrainer.TrainerId(),
		rawTrainer.TrainerName(),
		rawTrainer.LocationName(),
	), nil
}

func (p *placeCandidateService) CreateCheckList(
	ctx context.Context,
	race *prediction_entity.Race,
	horse *prediction_entity.Horse,
	forecast *prediction_entity.RaceForecast,
) []bool {
	input := &PlaceCheckListInput{
		Race:     race,
		Horse:    horse,
		Forecast: forecast,
	}

	checkList := make([]bool, 15)
	checkList[0] = p.placeCheckListService.OkEntries(ctx, input)
	checkList[1] = p.placeCheckListService.OkWinOdds(ctx, input)
	checkList[2] = p.placeCheckListService.OkInThirdPlaceRatio(ctx, input)
	checkList[3] = p.placeCheckListService.OkNotChangeCourseCategory(ctx, input)
	checkList[4] = p.placeCheckListService.OkSameDistance(ctx, input)
	checkList[5] = p.placeCheckListService.OkSameCourseCategory(ctx, input)
	checkList[6] = p.placeCheckListService.OkInThirdPlaceRecent(ctx, input)
	checkList[7] = p.placeCheckListService.OkTrackConditionExperience(ctx, input)
	checkList[8] = p.placeCheckListService.OkNotHorseWeightUp(ctx, input)
	checkList[9] = p.placeCheckListService.OkNotClassUp(ctx, input)
	checkList[10] = p.placeCheckListService.OkContinueOrEnhancementJockey(ctx, input)
	checkList[11] = p.placeCheckListService.OkNotSlowStart(ctx, input)
	checkList[12] = p.placeCheckListService.OkFavoriteRatio(ctx, input)
	checkList[13] = p.placeCheckListService.OkOnlyFavoriteAndRival(ctx, input)
	checkList[14] = p.placeCheckListService.OkIsHighlyRecommended(ctx, input)

	return checkList
}

func (p *placeCandidateService) Convert(
	ctx context.Context,
	race *prediction_entity.Race,
	horse *prediction_entity.Horse,
	jockey *prediction_entity.Jockey,
	trainer *prediction_entity.Trainer,
	forecast *prediction_entity.RaceForecast,
	calculable []*analysis_entity.PlaceCalculable,
	horseNumber types.HorseNumber,
	marker types.Marker,
	checkList []bool,
) *spreadsheet_entity.PredictionCheckList {
	var odds decimal.Decimal
	for _, o := range race.Odds() {
		if o.HorseNumber() == horseNumber {
			odds = o.Odds()
		}
	}

	predictionPlaces := p.predictionOddsService.Convert(ctx, race, horseNumber, marker, calculable)
	inexactOdds := odds.InexactFloat64()
	var firstPlaceRate, secondPlaceRate, thirdPlaceRate string
	if inexactOdds >= 1.0 && inexactOdds <= 1.5 {
		firstPlaceRate = predictionPlaces[0].RateData().OddsRange1RateFormat()
		secondPlaceRate = predictionPlaces[1].RateData().OddsRange1RateFormat()
		thirdPlaceRate = predictionPlaces[2].RateData().OddsRange1RateFormat()
	} else if inexactOdds >= 1.6 && inexactOdds <= 2.0 {
		firstPlaceRate = predictionPlaces[0].RateData().OddsRange2RateFormat()
		secondPlaceRate = predictionPlaces[1].RateData().OddsRange2RateFormat()
		thirdPlaceRate = predictionPlaces[2].RateData().OddsRange2RateFormat()
	} else if inexactOdds >= 2.1 && inexactOdds <= 2.9 {
		firstPlaceRate = predictionPlaces[0].RateData().OddsRange3RateFormat()
		secondPlaceRate = predictionPlaces[1].RateData().OddsRange3RateFormat()
		thirdPlaceRate = predictionPlaces[2].RateData().OddsRange3RateFormat()
	} else if inexactOdds >= 3.0 && inexactOdds <= 4.9 {
		firstPlaceRate = predictionPlaces[0].RateData().OddsRange4RateFormat()
		secondPlaceRate = predictionPlaces[1].RateData().OddsRange4RateFormat()
		thirdPlaceRate = predictionPlaces[2].RateData().OddsRange4RateFormat()
	} else if inexactOdds >= 5.0 && inexactOdds <= 9.9 {
		firstPlaceRate = predictionPlaces[0].RateData().OddsRange5RateFormat()
		secondPlaceRate = predictionPlaces[1].RateData().OddsRange5RateFormat()
		thirdPlaceRate = predictionPlaces[2].RateData().OddsRange5RateFormat()
	} else if inexactOdds >= 10.0 && inexactOdds <= 19.9 {
		firstPlaceRate = predictionPlaces[0].RateData().OddsRange6RateFormat()
		secondPlaceRate = predictionPlaces[1].RateData().OddsRange6RateFormat()
		thirdPlaceRate = predictionPlaces[2].RateData().OddsRange6RateFormat()
	} else if inexactOdds >= 20.0 && inexactOdds <= 49.9 {
		firstPlaceRate = predictionPlaces[0].RateData().OddsRange7RateFormat()
		secondPlaceRate = predictionPlaces[1].RateData().OddsRange7RateFormat()
		thirdPlaceRate = predictionPlaces[2].RateData().OddsRange7RateFormat()
	} else if inexactOdds >= 50.0 {
		firstPlaceRate = predictionPlaces[0].RateData().OddsRange8RateFormat()
		secondPlaceRate = predictionPlaces[1].RateData().OddsRange8RateFormat()
		thirdPlaceRate = predictionPlaces[2].RateData().OddsRange8RateFormat()
	}

	return spreadsheet_entity.NewPredictionCheckList(
		race.RaceId(),
		race.RaceDate(),
		race.RaceName(),
		race.RaceNumber(),
		race.RaceCourse(),
		horse.HorseId(),
		horse.HorseName(),
		jockey.JockeyId(),
		jockey.JockeyName(),
		trainer.TrainerId(),
		trainer.TrainerName(),
		trainer.LocationId(),
		odds,
		marker,
		firstPlaceRate,
		secondPlaceRate,
		thirdPlaceRate,
		checkList,
		forecast.FavoriteNum(),
		forecast.RivalNum(),
		forecast.MarkerNum(),
		forecast.IsHighlyRecommended(),
		forecast.TrainingComment(),
	)
}

func (p *placeCandidateService) Write(
	ctx context.Context,
	predictionCheckList []*spreadsheet_entity.PredictionCheckList,
) error {
	return p.spreadSheetRepository.WritePredictionCheckList(ctx, predictionCheckList)
}
