package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestSetVersion(t *testing.T) {
	testVersion := "1.2.3"
	SetVersion(testVersion)

	if rootCmd.Version != testVersion {
		t.Errorf("Expected version %s, got %s", testVersion, rootCmd.Version)
	}
}

func TestVersionFlag(t *testing.T) {
	// Set a test version
	testVersion := "test-version"
	SetVersion(testVersion)

	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	// Set args to trigger version flag
	rootCmd.SetArgs([]string{"--version"})

	// Execute command
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Check output contains version
	output := buf.String()
	if !strings.Contains(output, testVersion) {
		t.Errorf("Expected output to contain version %s, got: %s", testVersion, output)
	}
}
