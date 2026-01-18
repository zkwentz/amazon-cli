package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/internal/config"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
	Long:  `Manage authentication for Amazon CLI. Includes login, logout, and status commands.`,
}

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Clear stored credentials",
	Long:  `Clear stored authentication credentials from the configuration file.`,
	RunE:  runLogout,
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(logoutCmd)
}

func runLogout(cmd *cobra.Command, args []string) error {
	// Get config path
	configPath := cfgFile
	if configPath == "" {
		path, err := config.GetConfigPath()
		if err != nil {
			return models.NewCLIError(
				models.NetworkError,
				"Failed to determine config path",
				map[string]interface{}{"error": err.Error()},
			)
		}
		configPath = path
	}

	// Load config
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return models.NewCLIError(
			models.NetworkError,
			"Failed to load configuration",
			map[string]interface{}{"error": err.Error()},
		)
	}

	// Clear auth section
	cfg.Auth = config.AuthConfig{
		AccessToken:  "",
		RefreshToken: "",
		ExpiresAt:    time.Time{},
	}

	// Save config
	if err := config.SaveConfig(cfg, configPath); err != nil {
		return models.NewCLIError(
			models.NetworkError,
			"Failed to save configuration",
			map[string]interface{}{"error": err.Error()},
		)
	}

	// Output success
	result := map[string]string{
		"status": "logged_out",
	}

	return printer.Print(result)
}
