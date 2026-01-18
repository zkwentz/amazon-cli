package main

import (
	"os"

	"github.com/michaelshimeles/amazon-cli/cmd"
	"github.com/michaelshimeles/amazon-cli/pkg/logger"
)

func main() {
	logger.Debug("Starting amazon-cli application")
	if err := cmd.Execute(); err != nil {
		logger.Error("Command execution failed", "error", err)
		os.Exit(1)
	}
	logger.Debug("Application exiting normally")
}
