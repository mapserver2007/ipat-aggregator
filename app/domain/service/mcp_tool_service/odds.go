package mcp_tool_service

import (
	"fmt"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/mcp_result_entity"
)

const oddsNote = `
race_idsで指定されたレースIDでレースを絞り込み、そのレースのrace_resultsのoddsがodds_listで指定された値に一致するレースを絞り込む
`

type OddsInput struct {
	OddsList []float64
}

func (m *mcpToolResult) findByOdds(
	races []*mcp_result_entity.Race,
	input *OddsInput,
) (*mcp_result_entity.Odds, error) {
	m.logger.Info(fmt.Sprintf("find params odds: %v", input.OddsList))

	rawOdds := make([]string, len(input.OddsList))
	for i, odds := range input.OddsList {
		rawOdds[i] = fmt.Sprintf("%f", odds)
	}

	return &mcp_result_entity.Odds{
		OddsList: rawOdds,
		RaceIds:  []string{"202406010203", "202401020411"}, // TODO とりあえず固定値
		OddsNote: oddsNote,
	}, nil
}
