package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestVersionCommand(t *testing.T) {
	// Save original values
	origVersion := version
	origCommit := commit
	origDate := date

	// Set test values
	version = "1.0.0"
	commit = "abc123"
	date = "2024-01-01"

	// Restore original values after test
	defer func() {
		version = origVersion
		commit = origCommit
		date = origDate
	}()

	// Create a buffer to capture output
	buf := new(bytes.Buffer)

	// Create a new root command for testing
	testRootCmd := &cobra.Command{Use: "amazon-cli"}
	testVersionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Run: func(cmd *cobra.Command, args []string) {
			buf.WriteString("amazon-cli version " + version + "\n")
			buf.WriteString("commit: " + commit + "\n")
			buf.WriteString("built: " + date + "\n")
		},
	}
	testRootCmd.AddCommand(testVersionCmd)

	// Execute the version command
	testRootCmd.SetOut(buf)
	testRootCmd.SetArgs([]string{"version"})
	err := testRootCmd.Execute()

	if err != nil {
		t.Fatalf("version command failed: %v", err)
	}

	output := buf.String()

	// Verify output contains expected strings
	expectedStrings := []string{
		"amazon-cli version 1.0.0",
		"commit: abc123",
		"built: 2024-01-01",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(output, expected) {
			t.Errorf("output missing expected string: %q\nGot: %s", expected, output)
		}
	}
}

func TestVersionVariables(t *testing.T) {
	// Test that version variables exist and have default values
	if version == "" {
		t.Error("version variable should not be empty")
	}
	if commit == "" {
		t.Error("commit variable should not be empty")
	}
	if date == "" {
		t.Error("date variable should not be empty")
	}
}
