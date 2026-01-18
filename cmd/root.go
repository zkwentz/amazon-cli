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

var rootCmd = &cobra.Command{
	Use:   "amazon-cli",
	Short: "CLI for Amazon shopping - orders, returns, purchases, subscriptions",
	Long: `amazon-cli is a command-line interface that replaces the Amazon web interface,
enabling programmatic access to core Amazon shopping functionality.

The primary target users are AI agents, with the tool designed for publication
on ClawdHub alongside a companion skills.md file.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		printer := output.NewPrinter(outputFormat, quiet)
		printer.PrintError(err)
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
