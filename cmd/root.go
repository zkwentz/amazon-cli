package cmd

import (
	"fmt"

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

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "amazon-cli",
	Short: "CLI for Amazon shopping",
	Long: `amazon-cli is a command-line interface that replaces the Amazon web interface,
enabling programmatic access to core Amazon shopping functionality including
orders, returns, purchases, and subscriptions.`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Persistent flags available to all commands
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.amazon-cli/config.json)")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "json", "output format: json, table, raw")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress non-essential output")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose logging")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable colored output")
}

// getConfig loads the configuration from the config file
func getConfig() (*config.Config, error) {
	configPath := cfgFile
	if configPath == "" {
		configPath = config.GetConfigPath()
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return cfg, nil
}

// getPrinter returns an output printer based on global flags
func getPrinter() *output.Printer {
	return output.NewPrinter(outputFormat, quiet)
}
