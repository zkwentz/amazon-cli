package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/internal/output"
)

var (
	cfgFile      string
	outputFormat string
	quiet        bool
	verbose      bool
	noColor      bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "amazon-cli",
	Short: "A CLI tool for Amazon shopping operations",
	Long: `amazon-cli is a command-line interface that replaces the Amazon web interface,
enabling programmatic access to core Amazon shopping functionality including
orders, returns, purchases, and subscriptions.`,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		output.Exit(err)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.amazon-cli/config.json)")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "json", "output format (json, table, raw)")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress non-essential output")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose logging")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable colored output")
}
