package tospo_entity

import "github.com/mapserver2007/ipat-aggregator/app/domain/types"

type TrainingComment struct {
	horseNumber             types.HorseNumber
	trainingComment         string
	previousTrainingComment string
	isHighlyRecommended     bool
}

const highlyRecommended = "イチ押し"

func NewTrainingComment(
	horseNumber int,
	trainingComment string,
	previousTrainingComment string,
	prediction string,
) *TrainingComment {
	return &TrainingComment{
		horseNumber:             types.HorseNumber(horseNumber),
		trainingComment:         trainingComment,
		previousTrainingComment: previousTrainingComment,
		isHighlyRecommended:     prediction == highlyRecommended,
	}
}

func (t *TrainingComment) HorseNumber() types.HorseNumber {
	return t.horseNumber
}

func (t *TrainingComment) TrainingComment() string {
	return t.trainingComment
}

func (t *TrainingComment) PreviousTrainingComment() string {
	return t.previousTrainingComment
}

func (t *TrainingComment) IsHighlyRecommended() bool {
	return t.isHighlyRecommended
}
