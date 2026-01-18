package cmd

import (
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

	rootConfig *config.Config
	printer    *output.Printer
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "amazon-cli",
	Short: "CLI for Amazon shopping operations",
	Long: `amazon-cli is a command-line interface for Amazon shopping,
enabling programmatic access to orders, returns, purchases, and subscriptions.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Load config
		cfg, err := config.LoadConfig(cfgFile)
		if err != nil {
			return err
		}
		rootConfig = cfg

		// Override output format if specified
		if outputFormat != "" {
			rootConfig.Defaults.OutputFormat = outputFormat
		}

		// Initialize printer
		printer = output.NewPrinter(rootConfig.Defaults.OutputFormat, quiet)

		return nil
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Persistent flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.amazon-cli/config.json)")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "output format: json, table, raw (default: json)")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress non-essential output")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose logging")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable colored output")
}
