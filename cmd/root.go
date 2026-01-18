package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

var (
	outputFormat string
	quiet        bool
	verbose      bool
	configPath   string
	noColor      bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "amazon-cli",
	Short: "CLI tool for Amazon shopping",
	Long: `amazon-cli is a command-line interface that replaces the Amazon web interface,
enabling programmatic access to core Amazon shopping functionality.

This tool is designed for AI agents and outputs structured JSON for easy parsing.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		PrintError(err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "json", "Output format: json, table, raw")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Suppress non-essential output")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Path to config file (default: ~/.amazon-cli/config.json)")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colored output")
}

// PrintJSON prints data as formatted JSON
func PrintJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// PrintError prints an error in JSON format
func PrintError(err error) {
	var cliErr *models.CLIError
	var ok bool

	if cliErr, ok = err.(*models.CLIError); !ok {
		cliErr = models.NewCLIError(
			models.ErrorCodeAmazonError,
			err.Error(),
			nil,
		)
	}

	errorOutput := map[string]interface{}{
		"error": cliErr,
	}

	encoder := json.NewEncoder(os.Stderr)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(errorOutput)
}

// GetExitCode returns the appropriate exit code for an error
func GetExitCode(err error) int {
	cliErr, ok := err.(*models.CLIError)
	if !ok {
		return 1
	}

	switch cliErr.Code {
	case models.ErrorCodeAuthRequired, models.ErrorCodeAuthExpired:
		return 3
	case models.ErrorCodeNetworkError:
		return 4
	case models.ErrorCodeRateLimited:
		return 5
	case models.ErrorCodeNotFound:
		return 6
	case models.ErrorCodeInvalidInput:
		return 2
	default:
		return 1
	}
}

// LogVerbose prints a message if verbose mode is enabled
func LogVerbose(format string, args ...interface{}) {
	if verbose {
		fmt.Fprintf(os.Stderr, format+"\n", args...)
	}
}
