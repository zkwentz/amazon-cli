package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/internal/amazon"
)

var (
	useBrowser bool
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
	Long:  "Authenticate with Amazon and manage your session",
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to Amazon",
	Long: `Login to Amazon using browser-based authentication.

This will open a browser window where you can login to Amazon.
Once logged in, your session cookies will be captured and stored
for use with subsequent commands.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		authManager, err := amazon.NewAuthManager(cfg)
		if err != nil {
			return fmt.Errorf("failed to create auth manager: %w", err)
		}

		if err := authManager.Login(useBrowser); err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		if err := saveConfig(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		output, err := json.MarshalIndent(map[string]interface{}{
			"status":        "authenticated",
			"auth_method":   cfg.Auth.AuthMethod,
			"cookies_count": len(cfg.Auth.Cookies),
		}, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(output))
		return nil
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check authentication status",
	Long:  "Display current authentication status and session information",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		authManager, err := amazon.NewAuthManager(cfg)
		if err != nil {
			return fmt.Errorf("failed to create auth manager: %w", err)
		}

		status := authManager.GetAuthStatus()
		output, err := json.MarshalIndent(status, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(output))
		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout and clear credentials",
	Long:  "Clear all stored authentication credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		authManager, err := amazon.NewAuthManager(cfg)
		if err != nil {
			return fmt.Errorf("failed to create auth manager: %w", err)
		}

		if err := authManager.Logout(); err != nil {
			return fmt.Errorf("logout failed: %w", err)
		}

		if err := saveConfig(cfg); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}

		output, err := json.MarshalIndent(map[string]string{
			"status": "logged_out",
		}, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(output))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(statusCmd)
	authCmd.AddCommand(logoutCmd)

	loginCmd.Flags().BoolVar(&useBrowser, "browser", true,
		"use browser-based authentication (cookie capture)")
}
