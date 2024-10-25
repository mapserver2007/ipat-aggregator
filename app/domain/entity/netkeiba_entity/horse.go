package netkeiba_entity

type Horse struct {
	horseId       string
	horseName     string
	horseBirthDay string
	trainerId     string
	ownerId       string
	breederId     string
	horseBlood    *HorseBlood
	horseResults  []*HorseResult
}

func NewHorse(
	horseId string,
	horseName string,
	horseBirthDay string,
	trainerId string,
	ownerId string,
	breederId string,
	horseBlood *HorseBlood,
	horseResults []*HorseResult,
) *Horse {
	return &Horse{
		horseId:       horseId,
		horseName:     horseName,
		horseBirthDay: horseBirthDay,
		trainerId:     trainerId,
		ownerId:       ownerId,
		breederId:     breederId,
		horseBlood:    horseBlood,
		horseResults:  horseResults,
	}
}

func (h *Horse) HorseId() string {
	return h.horseId
}

func (h *Horse) HorseName() string {
	return h.horseName
}

func (h *Horse) HorseBirthDay() string {
	return h.horseBirthDay
}

func (h *Horse) TrainerId() string {
	return h.trainerId
}

func (h *Horse) OwnerId() string {
	return h.ownerId
}

func (h *Horse) BreederId() string {
	return h.breederId
}

func (h *Horse) HorseBlood() *HorseBlood {
	return h.horseBlood
}

func (h *Horse) HorseResults() []*HorseResult {
	return h.horseResults
}
