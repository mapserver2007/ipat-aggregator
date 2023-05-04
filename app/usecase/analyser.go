package usecase

import (
	analyse_entity "github.com/mapserver2007/ipat-aggregator/app/domain/analyse/entity"
	betting_ticket_entity "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/entity"
	race_entity "github.com/mapserver2007/ipat-aggregator/app/domain/race/entity"
	"github.com/mapserver2007/ipat-aggregator/app/service"
)

type Analyser struct {
	analyser service.Analyser
}

func NewAnalyser(
	analyser service.Analyser,
) *Analyser {
	return &Analyser{
		analyser: analyser,
	}
}

func (a *Analyser) Popular(
	records []*betting_ticket_entity.CsvEntity,
	racingNumbers []*race_entity.RacingNumber,
	races []*race_entity.Race,
) *analyse_entity.AnalyseSummary {
	popularAnalyseSummaries := a.analyser.PopularAnalyse(records, racingNumbers, races)

	return analyse_entity.NewAnalyseSummary(popularAnalyseSummaries)
}
