package entity

import betting_ticket_entity "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/entity"

type RaceRecord struct {
	recordDetails []*RecordDetail
}

func NewRaceRecordEntity(recordDetails []*RecordDetail) *RaceRecord {
	return &RaceRecord{recordDetails: recordDetails}
}

func (r *RaceRecord) RecordDetails() []*RecordDetail {
	return r.recordDetails
}

// RecordDetail レース単位のレース情報、購入馬券、本命対抗
type RecordDetail struct {
	race                 Race
	bettingTicketDetails []*betting_ticket_entity.BettingTicketDetail
	predictionForHorse   betting_ticket_entity.PredictionForHorse
}

func NewRecordDetail(
	race Race,
	bettingTicketDetails []*betting_ticket_entity.BettingTicketDetail,
	predictionForHorse betting_ticket_entity.PredictionForHorse,
) *RecordDetail {
	return &RecordDetail{
		race:                 race,
		bettingTicketDetails: bettingTicketDetails,
		predictionForHorse:   predictionForHorse,
	}
}

func (r *RecordDetail) Race() Race {
	return r.race
}

func (r *RecordDetail) BettingTicketDetails() []*betting_ticket_entity.BettingTicketDetail {
	return r.bettingTicketDetails
}

func (r *RecordDetail) PredictionForHorse() betting_ticket_entity.PredictionForHorse {
	return r.predictionForHorse
}
