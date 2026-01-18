package main

import (
	"os"

	"github.com/michaelshimeles/amazon-cli/cmd"
	"github.com/michaelshimeles/amazon-cli/pkg/models"
)

func main() {
	if err := cmd.Execute(); err != nil {
		exitCode := models.GetExitCodeFromError(err)
		os.Exit(exitCode)
	}
}
