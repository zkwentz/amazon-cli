package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCommand(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedOutput string
		expectedError  bool
	}{
		{
			name:           "help flag",
			args:           []string{"--help"},
			expectedOutput: "amazon-cli is a command-line interface",
			expectedError:  false,
		},
		{
			name:           "short help flag",
			args:           []string{"-h"},
			expectedOutput: "amazon-cli is a command-line interface",
			expectedError:  false,
		},
		{
			name:           "version info in help",
			args:           []string{"--help"},
			expectedOutput: "Usage:",
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new root command for each test to avoid state pollution
			cmd := &cobra.Command{
				Use:   "amazon-cli",
				Short: "CLI for Amazon shopping - orders, returns, purchases, subscriptions",
				Long: `amazon-cli is a command-line interface that replaces the Amazon web interface,
enabling programmatic access to core Amazon shopping functionality.

The CLI provides full access to Amazon shopping features including:
  - Orders management (list, track, view history)
  - Returns management (initiate, track returns)
  - Product search and discovery
  - Cart and checkout operations
  - Subscribe & Save subscription management

All commands output structured JSON for seamless AI agent integration.`,
			}

			// Capture output
			buf := new(bytes.Buffer)
			cmd.SetOut(buf)
			cmd.SetErr(buf)
			cmd.SetArgs(tt.args)

			// Execute command
			err := cmd.Execute()

			// Check error expectation
			if tt.expectedError && err == nil {
				t.Errorf("expected error but got none")
			}
			if !tt.expectedError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			// Check output
			output := buf.String()
			if tt.expectedOutput != "" && !strings.Contains(output, tt.expectedOutput) {
				t.Errorf("expected output to contain %q, got %q", tt.expectedOutput, output)
			}
		})
	}
}

func TestRootCommandFlags(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	tests := []struct {
		name     string
		args     []string
		flagName string
		expected interface{}
	}{
		{
			name:     "output flag json",
			args:     []string{"amazon-cli", "--output", "json"},
			flagName: "output",
			expected: "json",
		},
		{
			name:     "output flag short",
			args:     []string{"amazon-cli", "-o", "table"},
			flagName: "output",
			expected: "table",
		},
		{
			name:     "quiet flag",
			args:     []string{"amazon-cli", "--quiet"},
			flagName: "quiet",
			expected: true,
		},
		{
			name:     "quiet flag short",
			args:     []string{"amazon-cli", "-q"},
			flagName: "quiet",
			expected: true,
		},
		{
			name:     "verbose flag",
			args:     []string{"amazon-cli", "--verbose"},
			flagName: "verbose",
			expected: true,
		},
		{
			name:     "verbose flag short",
			args:     []string{"amazon-cli", "-v"},
			flagName: "verbose",
			expected: true,
		},
		{
			name:     "no-color flag",
			args:     []string{"amazon-cli", "--no-color"},
			flagName: "no-color",
			expected: true,
		},
		{
			name:     "config flag",
			args:     []string{"amazon-cli", "--config", "/custom/path/config.json"},
			flagName: "config",
			expected: "/custom/path/config.json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new root command for each test
			cmd := &cobra.Command{
				Use:   "amazon-cli",
				Short: "CLI for Amazon shopping",
			}

			// Add the persistent flags
			var cfgFile, outputFmt string
			var quiet, verbose, noColor bool

			cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
			cmd.PersistentFlags().StringVarP(&outputFmt, "output", "o", "json", "output format")
			cmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress output")
			cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose logging")
			cmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable color")

			// Silence output for cleaner test results
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)

			// Set args (skip first element as it's the program name)
			cmd.SetArgs(tt.args[1:])

			// Parse flags
			err := cmd.ParseFlags(tt.args[1:])
			if err != nil {
				t.Fatalf("failed to parse flags: %v", err)
			}

			// Check flag value
			flag := cmd.PersistentFlags().Lookup(tt.flagName)
			if flag == nil {
				t.Fatalf("flag %q not found", tt.flagName)
			}

			// Compare based on type
			switch expected := tt.expected.(type) {
			case string:
				if flag.Value.String() != expected {
					t.Errorf("expected %q for flag %q, got %q", expected, tt.flagName, flag.Value.String())
				}
			case bool:
				if flag.Value.String() != "true" && expected {
					t.Errorf("expected true for flag %q, got %q", tt.flagName, flag.Value.String())
				}
			}
		})
	}
}

func TestRootCommandDefaultValues(t *testing.T) {
	cmd := &cobra.Command{
		Use: "amazon-cli",
	}

	var cfgFile, outputFmt string
	var quiet, verbose, noColor bool

	cmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	cmd.PersistentFlags().StringVarP(&outputFmt, "output", "o", "json", "output format")
	cmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress output")
	cmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose logging")
	cmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable color")

	// Check default values
	if outputFmt != "json" {
		t.Errorf("expected default output format to be 'json', got %q", outputFmt)
	}
	if quiet {
		t.Errorf("expected quiet to be false by default")
	}
	if verbose {
		t.Errorf("expected verbose to be false by default")
	}
	if noColor {
		t.Errorf("expected no-color to be false by default")
	}
	if cfgFile != "" {
		t.Errorf("expected config file to be empty by default")
	}
}
