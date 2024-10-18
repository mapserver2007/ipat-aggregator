package prediction_service

import (
	"context"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/prediction_entity"
)

type PlaceCheckList interface {
	OkEntries(ctx context.Context, input *PlaceCheckListInput) bool
	OkWinOdds(ctx context.Context, input *PlaceCheckListInput) bool
}

type placeCheckListService struct {
}

type PlaceCheckListInput struct {
	Race  *prediction_entity.Race
	Horse *prediction_entity.Horse
}

func NewPlaceCheckList() PlaceCheckList {
	return &placeCheckListService{}
}

// OkEntries 13頭立て以下であること
func (p *placeCheckListService) OkEntries(ctx context.Context, input *PlaceCheckListInput) bool {
	return input.Race.Entries() <= 13
}

func (p *placeCheckListService) OkWinOdds(ctx context.Context, input *PlaceCheckListInput) bool {

	input.Horse.HorseId()
	input.Race.RaceEntryHorses()

	for _, odds := range input.Race.Odds() {
		odds.HorseNumber()
	}

	//TODO implement me
	panic("implement me")
}
