package main

import (
	"os"

	"github.com/michaelshimeles/amazon-cli/cmd"
	"github.com/michaelshimeles/amazon-cli/internal/output"
)

func main() {
	// Recover from any panics and output as JSON errors
	defer output.WrapPanic()

	// Execute the command and handle errors appropriately
	if err := cmd.Execute(); err != nil {
		exitCode := output.HandleError(err)
		os.Exit(exitCode)
	}
}
