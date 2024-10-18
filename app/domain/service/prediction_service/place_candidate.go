package prediction_service

import (
	"context"
	"fmt"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/repository"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
)

type PlaceCandidate interface {
	GetRaceCard(ctx context.Context, raceId types.RaceId) (*prediction_entity.Race, error)
	GetHorse(ctx context.Context, raceEntryHorse *prediction_entity.RaceEntryHorse) (*prediction_entity.Horse, error)
	CreateCheckList(ctx context.Context, race *prediction_entity.Race, horses []*prediction_entity.Horse) error
	Write(ctx context.Context) error
}

type placeCandidateService struct {
	raceRepository        repository.RaceRepository
	horseRepository       repository.HorseRepository
	raceEntityConverter   converter.RaceEntityConverter
	horseEntityConverter  converter.HorseEntityConverter
	placeCheckListService PlaceCheckList
}

func NewPlaceCandidate(
	raceRepository repository.RaceRepository,
	horseRepository repository.HorseRepository,
	raceEntityConverter converter.RaceEntityConverter,
	horseEntityConverter converter.HorseEntityConverter,
	placeCheckListService PlaceCheckList,
) PlaceCandidate {
	return &placeCandidateService{
		raceRepository:        raceRepository,
		horseRepository:       horseRepository,
		raceEntityConverter:   raceEntityConverter,
		horseEntityConverter:  horseEntityConverter,
		placeCheckListService: placeCheckListService,
	}
}

func (p *placeCandidateService) GetRaceCard(
	ctx context.Context,
	raceId types.RaceId,
) (*prediction_entity.Race, error) {
	rawRace, err := p.raceRepository.FetchRaceCard(ctx, fmt.Sprintf(raceCardUrl, raceId))
	if err != nil {
		return nil, err
	}

	race := p.raceEntityConverter.NetKeibaToPrediction(rawRace)

	return race, nil
}

func (p *placeCandidateService) GetHorse(
	ctx context.Context,
	raceEntryHorse *prediction_entity.RaceEntryHorse,
) (*prediction_entity.Horse, error) {
	rawHorse, err := p.horseRepository.Fetch(ctx, fmt.Sprintf(horseUrl, raceEntryHorse.HorseId()))
	if err != nil {
		return nil, err
	}

	horse, err := p.horseEntityConverter.NetKeibaToPrediction(rawHorse)
	if err != nil {
		return nil, err
	}

	return horse, nil
}

func (p *placeCandidateService) CreateCheckList(
	ctx context.Context,
	race *prediction_entity.Race,
	horses []*prediction_entity.Horse,
) error {

	p.placeCheckListService.OkEntries(ctx, &PlaceCheckListInput{
		Race: race,
	})

	for _, horse := range horses {
		input := &PlaceCheckListInput{
			Race:  race,
			Horse: horse,
		}

		p.placeCheckListService.OkWinOdds(ctx, input)
	}

	return nil
}

func (p *placeCandidateService) Write(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
