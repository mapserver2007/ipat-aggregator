package mcp_tool_service

import (
	"fmt"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/mcp_result_entity"
)

type RaceNameInput struct {
	RaceNames []string
}

func (m *mcpToolResult) findByRaceName(
	races []*data_cache_entity.Race,
	input *RaceNameInput,
) ([]*mcp_result_entity.Race, error) {
	m.logger.Info(fmt.Sprintf("find params raceNames: %v", input.RaceNames))

	var mcpRaces []*mcp_result_entity.Race
	for _, race := range races {
		for _, raceName := range input.RaceNames {
			if race.RaceName() == raceName {
				mcpRaces = append(mcpRaces, m.raceEntityConverter.DataCacheToMCPResult(race))
			}
		}
	}

	return mcpRaces, nil
}
