package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/internal/output"
)

var (
	outputFormat string
	quiet        bool
	verbose      bool
	configPath   string
	noColor      bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "amazon-cli",
	Short: "CLI for Amazon shopping operations",
	Long: `amazon-cli is a command-line interface that provides programmatic access
to Amazon shopping functionality including orders, returns, purchases, and subscriptions.

Designed for AI agents and power users to interact with Amazon programmatically.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		printer := output.NewPrinter(outputFormat, quiet)
		printer.PrintError(err)
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

// GetPrinter returns a configured output printer
func GetPrinter() *output.Printer {
	return output.NewPrinter(outputFormat, quiet)
}
