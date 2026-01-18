package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/internal/config"
)

var (
	cfgFile      string
	outputFormat string
	quiet        bool
	verbose      bool
	noColor      bool
)

var rootCmd = &cobra.Command{
	Use:   "amazon-cli",
	Short: "CLI tool for managing Amazon orders, returns, purchases, and subscriptions",
	Long: `amazon-cli is a command-line interface that replaces the Amazon web interface,
enabling programmatic access to core Amazon shopping functionality.

Designed for AI agents and power users, it provides full CLI access to:
- Orders management and tracking
- Returns processing
- Product search and purchases
- Subscribe & Save subscriptions

All output is in structured JSON format for easy parsing.`,
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is $HOME/.amazon-cli/config.json)")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "json",
		"output format: json, table, raw")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false,
		"suppress non-essential output")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"enable verbose logging")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false,
		"disable colored output")
}

// getConfigPath returns the config file path
func getConfigPath() string {
	if cfgFile != "" {
		return cfgFile
	}
	return config.GetConfigPath()
}

// loadConfig loads the configuration
func loadConfig() (*config.Config, error) {
	return config.LoadConfig(getConfigPath())
}

// saveConfig saves the configuration
func saveConfig(cfg *config.Config) error {
	return config.SaveConfig(cfg, getConfigPath())
}
