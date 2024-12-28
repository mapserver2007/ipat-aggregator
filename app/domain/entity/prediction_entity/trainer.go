package prediction_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type Trainer struct {
	trainerId   types.TrainerId
	trainerName string
	locationId  types.LocationId
}

func NewTrainer(
	rawTrainerId string,
	trainerName string,
	locationName string,
) *Trainer {
	return &Trainer{
		trainerId:   types.TrainerId(rawTrainerId),
		trainerName: trainerName,
		locationId:  types.NewLocationId(locationName),
	}
}

func (t *Trainer) TrainerId() types.TrainerId {
	return t.trainerId
}

func (t *Trainer) TrainerName() string {
	return t.trainerName
}

func (t *Trainer) LocationId() types.LocationId {
	return t.locationId
}
