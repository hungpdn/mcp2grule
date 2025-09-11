package cmd

import (
	"context"
	"os"

	"github.com/hungpdn/mcp2grule/internal/api"
	"github.com/hungpdn/mcp2grule/internal/api/handler"
	"github.com/hungpdn/mcp2grule/internal/config"
	"github.com/hungpdn/mcp2grule/internal/grule"
	"github.com/hungpdn/mcp2grule/internal/pkg/exitcode"
	"github.com/hungpdn/mcp2grule/internal/pkg/logger"
	"github.com/hungpdn/mcp2grule/internal/storage"
	"github.com/spf13/cobra"
)

var (
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Run Server",
		Long:  "Run MCP Server",
		Run:   runServer,
	}
)

// runServer starts the MCP server.
func runServer(_ *cobra.Command, _ []string) {
	ctx := context.Background()

	var store storage.IRulesetStorage
	switch config.App.DatabaseType {
	case config.DatabaseTypeMemory:
		store = storage.NewMemory()
	default:
		logger.Errorf("Unsupported database driver: %s", config.App.DatabaseType)
		os.Exit(exitcode.DatabaseError)
	}

	grule := grule.New(config.App.Grule, store)

	mcpHandler := handler.NewMCPHandler(grule)

	mcpServer := api.NewServer(AppName, Version, mcpHandler)

	if err := mcpServer.Run(ctx); err != nil {
		logger.Errorf("Failed to start MCP server: %v", err)
		os.Exit(exitcode.MCPTransportError)
	}
}
