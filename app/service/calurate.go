package service

import (
	"github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/entity"
	betting_ticket_vo "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/value_object"
	spreadsheet_entity "github.com/mapserver2007/ipat-aggregator/app/domain/spreadsheet/entity"
)

// TODO コンストラクタ化してDIしたい
// というかいらなくなるかも

func CalcSumResultRate(records []*entity.CsvEntity) spreadsheet_entity.ResultRate {
	var voteCount, hitCount, repayments, payments int
	for _, record := range records {
		voteCount += 1
		if record.BettingResult() == betting_ticket_vo.Hit {
			hitCount += 1
			repayments += record.Repayment()
		}
		payments += record.Payment()
	}

	return spreadsheet_entity.NewResultRate(voteCount, hitCount, payments, repayments)
}
