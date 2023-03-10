package entity

import (
	predict_vo "github.com/mapserver2007/tools/baken/app/domain/predict/value_object"
)

type PredictEntity struct {
	Race           *RaceEntity
	FavoriteHorse  *HorseEntity
	RivalHorse     *HorseEntity
	Payment        int
	Repayment      int
	WinningTickets []*WinningTicketEntity
	Status         predict_vo.PredictStatus
}

func NewPredictEntity(
	race *RaceEntity,
	favoriteHorse *HorseEntity,
	rivalHorse *HorseEntity,
	payment int,
	repayment int,
	winningTickets []*WinningTicketEntity,
	status predict_vo.PredictStatus,
) *PredictEntity {
	return &PredictEntity{
		Race:           race,
		FavoriteHorse:  favoriteHorse,
		RivalHorse:     rivalHorse,
		Payment:        payment,
		Repayment:      repayment,
		WinningTickets: winningTickets,
		Status:         status,
	}
}
