package entity

import (
	predict_vo "github.com/mapserver2007/ipat-aggregator/app/domain/predict/value_object"
)

type PredictEntity struct {
	race           *Race
	favoriteHorse  *Horse
	rivalHorse     *Horse
	favoriteJockey *Jockey
	rivalJockey    *Jockey
	payment        int
	repayment      int
	winningTickets []*WinningTicketEntity
	status         predict_vo.PredictStatus
}

func NewPredictEntity(
	race *Race,
	favoriteHorse *Horse,
	rivalHorse *Horse,
	favoriteJockey *Jockey,
	rivalJockey *Jockey,
	payment int,
	repayment int,
	winningTickets []*WinningTicketEntity,
	status predict_vo.PredictStatus,
) *PredictEntity {
	return &PredictEntity{
		race:           race,
		favoriteHorse:  favoriteHorse,
		rivalHorse:     rivalHorse,
		favoriteJockey: favoriteJockey,
		rivalJockey:    rivalJockey,
		payment:        payment,
		repayment:      repayment,
		winningTickets: winningTickets,
		status:         status,
	}
}

func (p *PredictEntity) Race() *Race {
	return p.race
}

func (p *PredictEntity) FavoriteHorse() *Horse {
	return p.favoriteHorse
}

func (p *PredictEntity) RivalHorse() *Horse {
	return p.rivalHorse
}

func (p *PredictEntity) FavoriteJockey() *Jockey {
	return p.favoriteJockey
}

func (p *PredictEntity) RivalJockey() *Jockey {
	return p.rivalJockey
}

func (p *PredictEntity) Payment() int {
	return p.payment
}

func (p *PredictEntity) Repayment() int {
	return p.repayment
}

func (p *PredictEntity) WinningTickets() []*WinningTicketEntity {
	return p.winningTickets
}

func (p *PredictEntity) Status() predict_vo.PredictStatus {
	return p.status
}
