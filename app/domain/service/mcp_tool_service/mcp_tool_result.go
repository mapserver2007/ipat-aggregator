package mcp_tool_service

import (
	"context"
	"encoding/json"

	"github.com/mapserver2007/ipat-aggregator/app/domain/entity/mcp_result_entity"
	"github.com/mapserver2007/ipat-aggregator/app/domain/service/converter"
	"github.com/mapserver2007/ipat-aggregator/app/domain/vo"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/master_usecase"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
)

type MCPToolResult interface {
	Create(ctx context.Context, arg *vo.MCPArgument, master *master_usecase.MasterOutput) (*mcp.CallToolResult, error)
}

type mcpToolResult struct {
	raceEntityConverter   converter.RaceEntityConverter
	jockeyEntityConverter converter.JockeyEntityConverter
	logger                *logrus.Logger
}

func NewMCPToolResult(
	raceEntityConverter converter.RaceEntityConverter,
	jockeyEntityConverter converter.JockeyEntityConverter,
	logger *logrus.Logger,
) MCPToolResult {
	return &mcpToolResult{
		raceEntityConverter:   raceEntityConverter,
		jockeyEntityConverter: jockeyEntityConverter,
		logger:                logger,
	}
}

func (m *mcpToolResult) Create(
	ctx context.Context,
	arg *vo.MCPArgument,
	master *master_usecase.MasterOutput,
) (*mcp.CallToolResult, error) {
	// TODO package配下の部品で各絞り込みを行い、最終的にunion allして返す
	var jockeys []*mcp_result_entity.Jockey
	for _, jockey := range master.Jockeys {
		jockeys = append(jockeys, m.jockeyEntityConverter.DataCacheToMCPResult(jockey))
	}

	races, err := m.findByRaceName(master.Races, &RaceNameInput{
		RaceNames: arg.RaceNames(),
	})
	if err != nil {
		return nil, err
	}
	m.logger.Infof("findByRaceName len: %v", len(races))

	// races, err = m.findByOdds(master.Races, &OddsInput{
	// 	OddsList: arg.Odds(),
	// })
	// if err != nil {
	// 	return nil, err
	// }

	result := mcp_result_entity.MCPResult{
		Races:    races,
		Jockeys:  jockeys,
		RaceNote: raceNote,
	}

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return nil, err
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}
