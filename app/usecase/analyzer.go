package usecase

import (
	analyze_entity "github.com/mapserver2007/ipat-aggregator/app/domain/analyze/entity"
	betting_ticket_entity "github.com/mapserver2007/ipat-aggregator/app/domain/betting_ticket/entity"
	race_entity "github.com/mapserver2007/ipat-aggregator/app/domain/race/entity"
	"github.com/mapserver2007/ipat-aggregator/app/service"
)

type Analyzer struct {
	analyzer service.Analyzer
}

func NewAnalyzer(
	analyzer service.Analyzer,
) *Analyzer {
	return &Analyzer{
		analyzer: analyzer,
	}
}

func (a *Analyzer) Popular(
	records []*betting_ticket_entity.CsvEntity,
	racingNumbers []*race_entity.RacingNumber,
	races []*race_entity.Race,
) *analyze_entity.AnalyzeSummary {
	winAnalyzeSummary := a.analyzer.WinAnalyze(records, racingNumbers, races)

	return analyze_entity.NewAnalyzeSummary(winAnalyzeSummary)
}
