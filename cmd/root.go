package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "amazon-cli",
	Short: "CLI for Amazon shopping",
	Long:  `amazon-cli is a command-line interface that provides programmatic access to Amazon shopping functionality.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags will be added here
	rootCmd.PersistentFlags().StringP("output", "o", "json", "Output format: json, table, raw")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "Suppress non-essential output")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().String("config", "", "Path to config file (default: ~/.amazon-cli/config.json)")
	rootCmd.PersistentFlags().Bool("no-color", false, "Disable colored output")
}
