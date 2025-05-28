package controller

import (
	"context"

	"github.com/mapserver2007/ipat-aggregator/app/usecase/mcp_usecase"
	"github.com/mark3labs/mcp-go/server"
	"github.com/sirupsen/logrus"
)

const (
	serverName    = "yamato"
	serverVersion = "0.0.1"
)

type MCPServer struct {
	mcpServer   *server.MCPServer
	toolUseCase mcp_usecase.Tool
	logger      *logrus.Logger
}

func NewMCPServer(
	toolUseCase mcp_usecase.Tool,
	logger *logrus.Logger,
) *MCPServer {
	mcpServer := server.NewMCPServer(
		serverName,
		serverVersion,
		server.WithResourceCapabilities(true, true),
		server.WithLogging(),
		server.WithRecovery(),
	)

	return &MCPServer{
		mcpServer:   mcpServer,
		toolUseCase: toolUseCase,
		logger:      logger,
	}
}

func (s *MCPServer) Start() error {
	if err := server.ServeStdio(s.mcpServer); err != nil {
		s.logger.Fatalf("server start error: %v", err)
		return err
	}

	return nil
}

func (s *MCPServer) Load(ctx context.Context) error {
	s.mcpServer.AddTool(s.toolUseCase.GetRaceTimeTool(), s.toolUseCase.GetRaceTimeToolHandler)

	return nil
}
