package cmd

import (
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zkwentz/amazon-cli/internal/config"
	"github.com/zkwentz/amazon-cli/internal/output"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
	Long:  `Login, logout, and check authentication status with Amazon.`,
}

// authLoginCmd represents the auth login command
var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Amazon",
	Long: `Authenticate with Amazon using browser-based OAuth.
Opens your default browser to Amazon's login page.
After authentication, tokens are stored locally.`,
	Run: func(cmd *cobra.Command, args []string) {
		output.JSON(map[string]interface{}{
			"status":  "login_required",
			"message": "Browser-based login not yet implemented",
		})
	},
}

// authStatusCmd represents the auth status command
var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check authentication status",
	Long:  `Display current authentication status including token expiry.`,
	Run: func(cmd *cobra.Command, args []string) {
		accessToken := viper.GetString("auth.access_token")
		expiresAtStr := viper.GetString("auth.expires_at")

		if accessToken == "" {
			output.JSON(map[string]interface{}{
				"authenticated": false,
				"message":       "Not logged in. Run 'amazon-cli auth login' to authenticate.",
			})
			return
		}

		expiresAt, err := time.Parse(time.RFC3339, expiresAtStr)
		if err != nil {
			output.JSON(map[string]interface{}{
				"authenticated": false,
				"message":       "Invalid token expiry. Please re-authenticate.",
			})
			return
		}

		now := time.Now()
		if now.After(expiresAt) {
			output.JSON(map[string]interface{}{
				"authenticated":      false,
				"expired":            true,
				"expires_at":         expiresAtStr,
				"message":            "Token has expired. Run 'amazon-cli auth login' to re-authenticate.",
			})
			return
		}

		expiresInSeconds := int(expiresAt.Sub(now).Seconds())

		output.JSON(map[string]interface{}{
			"authenticated":      true,
			"expires_at":         expiresAtStr,
			"expires_in_seconds": expiresInSeconds,
		})
	},
}

// authLogoutCmd represents the auth logout command
var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from Amazon",
	Long:  `Clear stored credentials and logout from Amazon.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load config
		cfg, err := config.LoadConfig(config.DefaultConfigPath())
		if err != nil {
			output.Error(models.ErrAmazonError, "Failed to load config: "+err.Error(), nil)
			os.Exit(models.ExitGeneralError)
		}

		// Clear auth tokens
		cfg.ClearAuth()

		// Save config
		err = config.SaveConfig(cfg, config.DefaultConfigPath())
		if err != nil {
			output.Error(models.ErrAmazonError, "Failed to save config: "+err.Error(), nil)
			os.Exit(models.ExitGeneralError)
		}

		// Output JSON
		output.JSON(map[string]interface{}{
			"status": "logged_out",
		})
	},
}

func init() {
	rootCmd.AddCommand(authCmd)

	// Add subcommands
	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authStatusCmd)
	authCmd.AddCommand(authLogoutCmd)
}
