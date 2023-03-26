package service

import (
	betting_ticket_entity "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/entity"
)

type Analyser struct {
	records []*betting_ticket_entity.CsvEntity
}

func NewAnalyser(
	records []*betting_ticket_entity.CsvEntity,
) *Analyser {
	return &Analyser{
		records: records,
	}
}

func (a *Analyser) Analyse() error {

	// 投票回数の計算
	//for _, record := range a.records {
	//	fmt.Println(record)
	//}

	return nil
}
