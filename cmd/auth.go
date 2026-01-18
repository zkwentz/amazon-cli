package cmd

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/internal/config"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
	Long:  `Manage authentication for Amazon CLI`,
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check current authentication status",
	Long:  `Check whether you are currently authenticated and when your session expires`,
	RunE:  runAuthStatus,
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authStatusCmd)
}

func runAuthStatus(cmd *cobra.Command, args []string) error {
	printer := GetPrinter()

	// Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		printer.PrintError("CONFIG_ERROR", "Failed to load configuration", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	// Check if tokens exist
	if cfg.Auth.AccessToken == "" || cfg.Auth.RefreshToken == "" {
		response := map[string]interface{}{
			"authenticated": false,
		}
		return printer.Print(response)
	}

	// Calculate time until expiration
	now := time.Now()
	expiresIn := cfg.Auth.ExpiresAt.Sub(now)
	expiresInSeconds := int(expiresIn.Seconds())

	// Check if token is expired
	isExpired := expiresInSeconds <= 0

	response := map[string]interface{}{
		"authenticated":      true,
		"expires_at":         cfg.Auth.ExpiresAt.Format(time.RFC3339),
		"expires_in_seconds": expiresInSeconds,
		"expired":            isExpired,
	}

	return printer.Print(response)
}
