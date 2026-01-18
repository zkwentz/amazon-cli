package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version information - set via ldflags during build
var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `Print the version, commit hash, and build date of amazon-cli.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("amazon-cli version %s\n", version)
		fmt.Printf("commit: %s\n", commit)
		fmt.Printf("built: %s\n", date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
