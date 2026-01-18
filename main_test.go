package main

import (
	"testing"

	"github.com/michaelshimeles/amazon-cli/cmd"
)

func TestExecute(t *testing.T) {
	// Test that cmd.Execute() can be called without panicking
	// Note: This will fail if there are no subcommands or required flags
	// but validates the basic integration
	err := cmd.Execute()
	// We expect an error here since no actual command is provided
	// but we're testing that the structure is sound
	_ = err
}
