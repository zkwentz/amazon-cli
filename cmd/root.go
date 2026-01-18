package cmd

import (
	"fmt"
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

var rootCmd = &cobra.Command{
	Use:   "amazon-cli",
	Short: "CLI for Amazon shopping",
	Long: `amazon-cli is a command-line interface for Amazon shopping.
It provides programmatic access to orders, returns, purchases, and subscriptions.`,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "json", "Output format: json, table, raw")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Suppress non-essential output")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().StringVar(&configPath, "config", "", "Path to config file (default: ~/.amazon-cli/config.json)")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colored output")
}

// GetPrinter returns an output printer based on flags
func GetPrinter() *output.Printer {
	return output.NewPrinter(outputFormat, quiet)
}
