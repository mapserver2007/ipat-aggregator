package netkeiba_entity

type Trainer struct {
	trainerId    string
	trainerName  string
	locationName string
}

func NewTrainer(
	trainerId string,
	trainerName string,
	locationName string,
) *Trainer {
	return &Trainer{
		trainerId:    trainerId,
		trainerName:  trainerName,
		locationName: locationName,
	}
}

func (t *Trainer) TrainerId() string {
	return t.trainerId
}

func (t *Trainer) TrainerName() string {
	return t.trainerName
}

func (t *Trainer) LocationName() string {
	return t.locationName
}
