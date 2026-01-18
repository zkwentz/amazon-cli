package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Display configuration information",
	Long:  `Display the current configuration path and settings.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath := config.GetConfigPath()

		output := map[string]string{
			"config_path": configPath,
		}

		jsonOutput, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(jsonOutput))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}
