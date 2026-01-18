package main

import (
	"testing"

	"github.com/michaelshimeles/amazon-cli/pkg/models"
)

// TestMain verifies that the main function uses exit codes correctly
// Note: We can't directly test main() due to os.Exit(), but we test
// that cmd.Execute() returns the correct exit codes in cmd/root_test.go

func TestExitCodeConstants(t *testing.T) {
	// Verify the exit code constants are correctly defined
	// This ensures main.go will use the correct values when calling os.Exit()
	tests := []struct {
		name     string
		code     int
		expected int
	}{
		{"ExitSuccess", models.ExitSuccess, 0},
		{"ExitGeneralError", models.ExitGeneralError, 1},
		{"ExitInvalidArgs", models.ExitInvalidArgs, 2},
		{"ExitAuthError", models.ExitAuthError, 3},
		{"ExitNetworkError", models.ExitNetworkError, 4},
		{"ExitRateLimited", models.ExitRateLimited, 5},
		{"ExitNotFound", models.ExitNotFound, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.code != tt.expected {
				t.Errorf("%s = %d, want %d", tt.name, tt.code, tt.expected)
			}
		})
	}
}
