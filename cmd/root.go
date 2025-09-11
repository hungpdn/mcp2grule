package cmd

import (
	"os"

	"github.com/hungpdn/mcp2grule/internal/pkg/exitcode"
	"github.com/hungpdn/mcp2grule/internal/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	AppName = "mcp2grule"
	Version = "v0.0.1"

	rootCmd = &cobra.Command{Use: AppName}
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logger.Errorf("Failed to execute command: %v", err)
		os.Exit(exitcode.GenericError)
	}
}
