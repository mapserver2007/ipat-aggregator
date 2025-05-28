package mcp_tool_service

import (
	"fmt"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/data_cache_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/mcp_result_entity"
)

type OddsInput struct {
	OddsList []int
}

func (m *mcpToolResult) findByOdds(
	races []*data_cache_entity.Race,
	input *OddsInput,
) ([]*mcp_result_entity.Race, error) {
	m.logger.Info(fmt.Sprintf("find params odds: %v", input.OddsList))
	return nil, nil
}
