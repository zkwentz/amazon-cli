package cmd

import (
	"fmt"
	"os"

	"github.com/michaelshimeles/amazon-cli/internal/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	output  string
	quiet   bool
	verbose bool
	noColor bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "amazon-cli",
	Short: "CLI for Amazon shopping - orders, returns, purchases, subscriptions",
	Long: `amazon-cli is a command-line interface that replaces the Amazon web interface,
enabling programmatic access to core Amazon shopping functionality.

The CLI provides full access to Amazon shopping features including:
- Orders management (list, get details, track shipments)
- Returns management (list returnable items, initiate returns)
- Search and product information
- Cart and checkout operations
- Subscribe & Save management

All commands output structured JSON for seamless AI agent integration.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.amazon-cli/config.json)")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "json", "Output format: json, table, raw")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Suppress non-essential output")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "Disable colored output")

	// Bind flags to viper
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("quiet", rootCmd.PersistentFlags().Lookup("quiet"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("no-color", rootCmd.PersistentFlags().Lookup("no-color"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Initialize logger based on verbose flag
	logger.InitLogger(verbose)

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			logger.Error("failed to get home directory", "error", err)
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Search config in home directory with name ".amazon-cli/config" (without extension).
		configPath := home + "/.amazon-cli"
		viper.AddConfigPath(configPath)
		viper.SetConfigType("json")
		viper.SetConfigName("config")

		logger.Debug("config path set", "path", configPath)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logger.Debug("config file loaded", "file", viper.ConfigFileUsed())
	} else {
		logger.Debug("no config file found", "error", err)
	}
}
