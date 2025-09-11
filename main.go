package main

import (
	"os"

	"github.com/hungpdn/mcp2grule/cmd"
	"github.com/hungpdn/mcp2grule/internal/pkg/exitcode"
)

func main() {
	cmd.Execute()
	os.Exit(exitcode.Success)
}
