package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/internal/output"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

var (
	cfgFile      string
	outputFormat string
	quiet        bool
	verbose      bool
	noColor      bool
	printer      *output.Printer
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "amazon-cli",
	Short: "CLI for Amazon shopping - orders, returns, purchases, subscriptions",
	Long: `amazon-cli is a command-line interface that replaces the Amazon web interface,
enabling programmatic access to core Amazon shopping functionality.

The CLI outputs structured JSON for seamless AI agent integration and includes
safety rails to prevent accidental purchases.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		printer = output.NewPrinter(outputFormat, quiet)
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		printer := output.NewPrinter("json", false)
		printer.PrintError(err)

		// Map error to exit code
		exitCode := getExitCode(err)
		os.Exit(exitCode)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ~/.amazon-cli/config.json)")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "json", "output format: json, table, raw")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress non-essential output")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose logging")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable colored output")
}

// getExitCode maps errors to exit codes
func getExitCode(err error) int {
	if cliErr, ok := err.(*models.CLIError); ok {
		switch cliErr.Code {
		case models.AuthRequired, models.AuthExpired:
			return 3
		case models.NetworkError:
			return 4
		case models.RateLimited:
			return 5
		case models.NotFound:
			return 6
		case models.InvalidInput:
			return 2
		default:
			return 1
		}
	}
	return 1
}
