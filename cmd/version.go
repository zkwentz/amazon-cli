package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of amazon-cli",
	Long:  `Print the version number, commit hash, and build date of amazon-cli.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("amazon-cli version %s (commit: %s, built: %s)\n", version, commit, date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
