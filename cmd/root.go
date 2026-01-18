package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/internal/config"
	"github.com/zkwentz/amazon-cli/internal/output"
)

var (
	cfgFile      string
	outputFormat string
	quiet        bool
	verbose      bool
	noColor      bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "amazon-cli",
	Short: "CLI for Amazon shopping - orders, returns, purchases, subscriptions",
	Long: `amazon-cli is a command-line interface for Amazon shopping functionality.
It provides programmatic access to orders, returns, purchases, and subscriptions,
with structured JSON output designed for AI agent integration.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.amazon-cli/config.json)")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "json", "output format: json, table, raw")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress non-essential output")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose logging")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable colored output")
}

// getConfig loads the configuration
func getConfig() (*config.Config, error) {
	configPath := cfgFile
	if configPath == "" {
		var err error
		configPath, err = config.GetConfigPath()
		if err != nil {
			return nil, err
		}
	}

	return config.LoadConfig(configPath)
}

// getPrinter returns an output printer based on global flags
func getPrinter() *output.Printer {
	return output.NewPrinter(outputFormat, quiet)
}

// handleError prints an error and exits with the appropriate code
func handleError(err error, printer *output.Printer) {
	if printer != nil {
		printer.PrintError(err)
	} else {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	}

	// Exit with appropriate code based on error type
	// For now, exit with 1 for all errors
	// TODO: Map error codes to exit codes as per PRD
	os.Exit(1)
}
