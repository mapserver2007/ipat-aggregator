package mcp_usecase

import (
	"context"

	"github.com/mapserver2007/ipat-aggregator/config"

	"github.com/mapserver2007/ipat-aggregator/app/domain/service/mcp_tool_service"
	"github.com/mapserver2007/ipat-aggregator/app/domain/types"
	"github.com/mapserver2007/ipat-aggregator/app/usecase/master_usecase"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/sirupsen/logrus"
)

const (
	toolNameRaceName       = "raceName"
	toolNameRaceCourse     = "raceCourse"
	toolNameCourseCategory = "courseCategory"
	toolNameDistance       = "distance"
)

type Tool interface {
	GetRaceTimeTool() mcp.Tool
	GetRaceTimeToolHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
}

type tool struct {
	masterUseCase master_usecase.Master
	mcpToolResult mcp_tool_service.MCPToolResult
	logger        *logrus.Logger
}

func NewTool(
	masterUseCase master_usecase.Master,
	mcpToolResult mcp_tool_service.MCPToolResult,
	logger *logrus.Logger,
) Tool {
	return &tool{
		masterUseCase: masterUseCase,
		mcpToolResult: mcpToolResult,
		logger:        logger,
	}
}

func (t *tool) GetRaceTimeTool() mcp.Tool {
	return mcp.NewTool(
		"race_time",
		mcp_tool_service.GetToolOps()...,
	)
}

func (t *tool) GetRaceTimeToolHandler(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	startDate, err := types.NewRaceDate(config.RaceStartDate)
	if err != nil {
		t.logger.Errorf("failed to create race date: %v", err)
		return nil, err
	}

	endDate, err := types.NewRaceDate(config.RaceEndDate)
	if err != nil {
		t.logger.Errorf("failed to create race date: %v", err)
		return nil, err
	}

	err = t.masterUseCase.CreateOrUpdate(ctx, &master_usecase.MasterInput{
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		return nil, err
	}

	master, err := t.masterUseCase.Get(ctx)
	if err != nil {
		t.logger.Errorf("failed to get master: %v", err)
		return nil, err
	}

	mcpArgument, err := mcp_tool_service.LoadMCPArgument(request)
	t.logger.Infof("mcp argument: %v", mcpArgument)
	if err != nil {
		t.logger.Errorf("failed to load mcp argument: %v", err)
		return nil, err
	}

	mcpToolResult, err := t.mcpToolResult.Create(ctx, mcpArgument, master)
	if err != nil {
		t.logger.Errorf("failed to create mcp tool result: %v", err)
		return nil, err
	}

	// raceTimeResult, err := t.mcpRaceTimeService.Find(ctx, mcpArgument)
	// if err != nil {
	// 	return nil, err
	// }

	return mcpToolResult, nil
}
