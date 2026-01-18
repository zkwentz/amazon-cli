package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
	Short: "Amazon CLI - Command-line interface for Amazon shopping",
	Long: `amazon-cli is a command-line interface that replaces the Amazon web interface,
enabling programmatic access to core Amazon shopping functionality.

This tool outputs structured JSON for seamless AI agent integration and includes
safety rails to prevent accidental purchases.`,
}

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
