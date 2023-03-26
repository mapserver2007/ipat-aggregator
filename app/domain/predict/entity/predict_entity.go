package entity

import (
	predict_vo "github.com/mapserver2007/ipat-aggregator/app/domain/predict/value_object"
)

type PredictEntity struct {
	race           *RaceEntity
	favoriteHorse  *HorseEntity
	rivalHorse     *HorseEntity
	payment        int
	repayment      int
	winningTickets []*WinningTicketEntity
	status         predict_vo.PredictStatus
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
		race:           race,
		favoriteHorse:  favoriteHorse,
		rivalHorse:     rivalHorse,
		payment:        payment,
		repayment:      repayment,
		winningTickets: winningTickets,
		status:         status,
	}
}

func (p *PredictEntity) Race() *RaceEntity {
	return p.race
}

func (p *PredictEntity) FavoriteHorse() *HorseEntity {
	return p.favoriteHorse
}

func (p *PredictEntity) RivalHorse() *HorseEntity {
	return p.rivalHorse
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
