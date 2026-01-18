package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/internal/amazon"
	"github.com/zkwentz/amazon-cli/internal/config"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication with Amazon",
	Long:  `Commands for logging in, checking authentication status, and logging out of Amazon.`,
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Amazon account",
	Long: `Opens your browser to authenticate with Amazon using OAuth.
After successful authentication, access tokens are stored locally.`,
	RunE: runLogin,
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check authentication status",
	Long:  `Check if you are currently authenticated and when your tokens expire.`,
	RunE:  runStatus,
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from Amazon account",
	Long:  `Clear stored authentication tokens from local config.`,
	RunE:  runLogout,
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(statusCmd)
	authCmd.AddCommand(logoutCmd)
}

func runLogin(cmd *cobra.Command, args []string) error {
	// Load existing config
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		return outputError("INVALID_CONFIG", fmt.Sprintf("Failed to load config: %v", err))
	}

	// Check for client ID and secret
	authClient := amazon.NewAuthClient(cfg)
	if authClient.ClientID == "" || authClient.ClientSecret == "" {
		return outputError("AUTH_REQUIRED", "Amazon OAuth credentials not configured. Please set AMAZON_CLIENT_ID and AMAZON_CLIENT_SECRET environment variables. Register your app at https://developer.amazon.com/")
	}

	// Perform login
	tokens, err := authClient.Login()
	if err != nil {
		return outputError("AUTH_FAILED", fmt.Sprintf("Authentication failed: %v", err))
	}

	// Save tokens to config
	cfg.Auth.AccessToken = tokens.AccessToken
	cfg.Auth.RefreshToken = tokens.RefreshToken
	cfg.Auth.ExpiresAt = time.Now().Add(time.Duration(tokens.ExpiresIn) * time.Second)

	if err := config.SaveConfig(cfg, cfgFile); err != nil {
		return outputError("CONFIG_SAVE_FAILED", fmt.Sprintf("Failed to save config: %v", err))
	}

	// Output success
	result := map[string]interface{}{
		"status":     "authenticated",
		"expires_at": cfg.Auth.ExpiresAt.Format(time.RFC3339),
	}

	return outputJSON(result)
}

func runStatus(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		return outputError("INVALID_CONFIG", fmt.Sprintf("Failed to load config: %v", err))
	}

	if !cfg.IsAuthenticated() {
		result := map[string]interface{}{
			"authenticated": false,
		}
		return outputJSON(result)
	}

	expiresIn := int(time.Until(cfg.Auth.ExpiresAt).Seconds())
	if expiresIn < 0 {
		expiresIn = 0
	}

	result := map[string]interface{}{
		"authenticated":       true,
		"expires_at":          cfg.Auth.ExpiresAt.Format(time.RFC3339),
		"expires_in_seconds":  expiresIn,
		"token_expired":       cfg.IsTokenExpired(),
	}

	return outputJSON(result)
}

func runLogout(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig(cfgFile)
	if err != nil {
		return outputError("INVALID_CONFIG", fmt.Sprintf("Failed to load config: %v", err))
	}

	// Clear auth tokens
	cfg.Auth.AccessToken = ""
	cfg.Auth.RefreshToken = ""
	cfg.Auth.ExpiresAt = time.Time{}

	if err := config.SaveConfig(cfg, cfgFile); err != nil {
		return outputError("CONFIG_SAVE_FAILED", fmt.Sprintf("Failed to save config: %v", err))
	}

	result := map[string]interface{}{
		"status": "logged_out",
	}

	return outputJSON(result)
}

func outputJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func outputError(code, message string) error {
	errResponse := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
			"details": map[string]interface{}{},
		},
	}

	encoder := json.NewEncoder(os.Stderr)
	encoder.SetIndent("", "  ")
	encoder.Encode(errResponse)

	return fmt.Errorf("%s: %s", code, message)
}
