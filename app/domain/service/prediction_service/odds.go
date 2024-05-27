package prediction_service

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/analysis_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"time"
)

const (
	raceCardUrl   = "https://race.netkeiba.com/race/shutuba.html?race_id=%s"
	oddsUrl       = "https://race.netkeiba.com/api/api_get_jra_odds.html?race_id=%s&type=%d&action=update"
	raceResultUrl = "https://race.netkeiba.com/race/result.html?race_id=%s&organizer=1&race_date=%s"
)

type Odds interface {
	Get(ctx context.Context, raceId types.RaceId) (*prediction_entity.Race, error)
	Convert(ctx context.Context, predictionRaces []*prediction_entity.Race, calculables []*analysis_entity.PlaceCalculable) error
}

type oddsService struct {
	oddRepository  repository.OddsRepository
	raceRepository repository.RaceRepository
}

func NewOdds(
	oddRepository repository.OddsRepository,
	raceRepository repository.RaceRepository,
) Odds {
	return &oddsService{
		oddRepository:  oddRepository,
		raceRepository: raceRepository,
	}
}

func (p *oddsService) Get(
	ctx context.Context,
	raceId types.RaceId,
) (*prediction_entity.Race, error) {
	odds, err := p.oddRepository.Fetch(ctx, fmt.Sprintf(oddsUrl, raceId, 1))
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
	raceResultHorseNumbers := make([]int, 0, 3)
	if race.RaceResults() != nil && len(race.RaceResults()) >= 3 {
		for _, raceResult := range race.RaceResults()[:3] {
			raceResultHorseNumbers = append(raceResultHorseNumbers, raceResult.HorseNumber())
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
	)

	return predictionRace, nil
}

func (p *oddsService) Convert(
	ctx context.Context,
	predictionRaces []*prediction_entity.Race,
	calculables []*analysis_entity.PlaceCalculable,
) error {
	return nil
}
