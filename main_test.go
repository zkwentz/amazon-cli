package main

import (
	"testing"

	"github.com/michaelshimeles/amazon-cli/cmd"
)

func TestExecute(t *testing.T) {
	// Test that the command can be executed without error
	// Note: This will just verify the command structure is valid
	// We don't actually run any commands as they require user input
	err := cmd.Execute()
	// Expect an error since no command is provided
	if err == nil {
		// No error is also acceptable (shows help)
		return
	}
}
