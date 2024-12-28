package analysis_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type Horse struct {
	horseId        types.HorseId
	horseName      string
	horseBirthDay  types.HorseBirthDay
	trainerId      types.TrainerId
	ownerId        types.OwnerId
	breederId      types.BreederId
	horseBlood     *HorseBlood
	horseResults   []*HorseResult
	latestRaceDate types.RaceDate
}

func NewHorse(
	rawHorseId string,
	horseName string,
	horseBirthDay int,
	rawTrainerId string,
	rawOwnerId string,
	rawBreederId string,
	horseBlood *HorseBlood,
	horseResults []*HorseResult,
	latestRaceDate types.RaceDate,
) *Horse {
	return &Horse{
		horseId:        types.HorseId(rawHorseId),
		horseName:      horseName,
		horseBirthDay:  types.HorseBirthDay(horseBirthDay),
		trainerId:      types.TrainerId(rawTrainerId),
		ownerId:        types.OwnerId(rawOwnerId),
		breederId:      types.BreederId(rawBreederId),
		horseBlood:     horseBlood,
		horseResults:   horseResults,
		latestRaceDate: latestRaceDate,
	}
}

func (h *Horse) HorseId() types.HorseId {
	return h.horseId
}

func (h *Horse) HorseName() string {
	return h.horseName
}

func (h *Horse) HorseBirthDay() types.HorseBirthDay {
	return h.horseBirthDay
}

func (h *Horse) TrainerId() types.TrainerId {
	return h.trainerId
}

func (h *Horse) OwnerId() types.OwnerId {
	return h.ownerId
}

func (h *Horse) BreederId() types.BreederId {
	return h.breederId
}

func (h *Horse) HorseBlood() *HorseBlood {
	return h.horseBlood
}

func (h *Horse) HorseResults() []*HorseResult {
	return h.horseResults
}

func (h *Horse) LatestRaceDate() types.RaceDate {
	return h.latestRaceDate
}
