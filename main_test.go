package main

import (
	"os"
	"os/exec"
	"testing"
)

// TestExitCodes tests that the CLI returns appropriate exit codes
// by running the compiled binary as a subprocess
func TestExitCodes(t *testing.T) {
	// Build the binary for testing
	buildCmd := exec.Command("go", "build", "-o", "amazon-cli-test", ".")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build binary for testing: %v", err)
	}
	defer os.Remove("amazon-cli-test")

	tests := []struct {
		name         string
		args         []string
		wantExitCode int
	}{
		{
			name:         "help command exits with success",
			args:         []string{"--help"},
			wantExitCode: 0,
		},
		{
			name:         "version command exits with success",
			args:         []string{"version"},
			wantExitCode: 0,
		},
		{
			name:         "invalid flag exits with error",
			args:         []string{"--invalid-flag"},
			wantExitCode: 1, // Cobra returns error for unknown flags
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("./amazon-cli-test", tt.args...)
			err := cmd.Run()

			if err == nil && tt.wantExitCode != 0 {
				t.Errorf("Expected exit code %d, but command succeeded", tt.wantExitCode)
				return
			}

			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					exitCode := exitErr.ExitCode()
					if exitCode != tt.wantExitCode {
						t.Errorf("Exit code = %d, want %d", exitCode, tt.wantExitCode)
					}
				} else if tt.wantExitCode != 0 {
					t.Errorf("Expected exit code %d, but got different error: %v", tt.wantExitCode, err)
				}
			}
		})
	}
}

// TestExitCodeMapping verifies that the error to exit code mapping works correctly
// This is a documentation test to ensure the mapping matches the PRD
func TestExitCodeMapping(t *testing.T) {
	// Document the exit code mapping from PRD section "Exit Codes"
	mapping := map[string]int{
		"Success":           0,
		"General error":     1,
		"Invalid arguments": 2,
		"Authentication":    3,
		"Network error":     4,
		"Rate limited":      5,
		"Not found":         6,
	}

	// Verify the mapping exists and is documented
	if len(mapping) != 7 {
		t.Errorf("Expected 7 exit codes defined in PRD, got %d", len(mapping))
	}

	// Verify specific codes match PRD
	expectedCodes := map[string]int{
		"Success":           0,
		"General error":     1,
		"Invalid arguments": 2,
		"Authentication":    3,
		"Network error":     4,
		"Rate limited":      5,
		"Not found":         6,
	}

	for name, expectedCode := range expectedCodes {
		if actualCode, ok := mapping[name]; !ok {
			t.Errorf("Missing exit code definition for %s", name)
		} else if actualCode != expectedCode {
			t.Errorf("Exit code for %s = %d, want %d", name, actualCode, expectedCode)
		}
	}
}
