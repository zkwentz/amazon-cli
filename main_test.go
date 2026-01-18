package main

import (
	"testing"

	"github.com/michaelshimeles/amazon-cli/cmd"
)

func TestMain(t *testing.T) {
	// Test that Execute doesn't panic
	// Note: This will print help text since no args are provided
	err := cmd.Execute()
	// It's okay if Execute returns an error for invalid args
	// We just want to ensure it doesn't panic
	_ = err
}
