package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/michaelshimeles/amazon-cli/internal/amazon"
	"github.com/spf13/cobra"
)

var (
	confirmFlag bool
	intervalFlag int
)

// subscriptionsCmd represents the subscriptions command
var subscriptionsCmd = &cobra.Command{
	Use:   "subscriptions",
	Short: "Manage Subscribe & Save subscriptions",
	Long: `Manage Amazon Subscribe & Save subscriptions.

This command provides access to:
- List all subscriptions
- Get subscription details
- Skip next delivery
- Change delivery frequency
- Cancel subscriptions
- View upcoming deliveries`,
}

// subscriptionsListCmd represents the subscriptions list command
var subscriptionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Subscribe & Save subscriptions",
	Long:  `Lists all active and paused Subscribe & Save subscriptions.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := amazon.NewClient()
		response, err := client.GetSubscriptions()
		if err != nil {
			return fmt.Errorf("failed to get subscriptions: %w", err)
		}

		output, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// subscriptionsGetCmd represents the subscriptions get command
var subscriptionsGetCmd = &cobra.Command{
	Use:   "get <subscription-id>",
	Short: "Get details for a specific subscription",
	Long:  `Retrieves detailed information about a specific Subscribe & Save subscription.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID := args[0]
		client := amazon.NewClient()
		subscription, err := client.GetSubscription(subscriptionID)
		if err != nil {
			return fmt.Errorf("failed to get subscription: %w", err)
		}

		output, err := json.MarshalIndent(subscription, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// subscriptionsSkipCmd represents the subscriptions skip command
var subscriptionsSkipCmd = &cobra.Command{
	Use:   "skip <subscription-id>",
	Short: "Skip the next delivery for a subscription",
	Long:  `Skips the next scheduled delivery for a Subscribe & Save subscription.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID := args[0]

		if !confirmFlag {
			dryRun := map[string]interface{}{
				"dry_run": true,
				"message": "Add --confirm to execute",
				"action":  "skip next delivery",
				"subscription_id": subscriptionID,
			}
			output, err := json.MarshalIndent(dryRun, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to format output: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		client := amazon.NewClient()
		subscription, err := client.SkipDelivery(subscriptionID)
		if err != nil {
			return fmt.Errorf("failed to skip delivery: %w", err)
		}

		output, err := json.MarshalIndent(subscription, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// subscriptionsFrequencyCmd represents the subscriptions frequency command
var subscriptionsFrequencyCmd = &cobra.Command{
	Use:   "frequency <subscription-id>",
	Short: "Change the delivery frequency for a subscription",
	Long:  `Changes the delivery frequency (in weeks) for a Subscribe & Save subscription.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID := args[0]

		if intervalFlag <= 0 || intervalFlag > 26 {
			return fmt.Errorf("interval must be between 1 and 26 weeks")
		}

		if !confirmFlag {
			dryRun := map[string]interface{}{
				"dry_run": true,
				"message": "Add --confirm to execute",
				"action":  "change frequency",
				"subscription_id": subscriptionID,
				"new_frequency_weeks": intervalFlag,
			}
			output, err := json.MarshalIndent(dryRun, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to format output: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		client := amazon.NewClient()
		subscription, err := client.UpdateFrequency(subscriptionID, intervalFlag)
		if err != nil {
			return fmt.Errorf("failed to update frequency: %w", err)
		}

		output, err := json.MarshalIndent(subscription, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// subscriptionsCancelCmd represents the subscriptions cancel command
var subscriptionsCancelCmd = &cobra.Command{
	Use:   "cancel <subscription-id>",
	Short: "Cancel a subscription",
	Long:  `Cancels a Subscribe & Save subscription.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID := args[0]

		if !confirmFlag {
			dryRun := map[string]interface{}{
				"dry_run": true,
				"message": "Add --confirm to execute",
				"action":  "cancel subscription",
				"subscription_id": subscriptionID,
			}
			output, err := json.MarshalIndent(dryRun, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to format output: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		client := amazon.NewClient()
		subscription, err := client.CancelSubscription(subscriptionID)
		if err != nil {
			return fmt.Errorf("failed to cancel subscription: %w", err)
		}

		output, err := json.MarshalIndent(subscription, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// subscriptionsUpcomingCmd represents the subscriptions upcoming command
var subscriptionsUpcomingCmd = &cobra.Command{
	Use:   "upcoming",
	Short: "View upcoming deliveries across all subscriptions",
	Long:  `Lists all upcoming deliveries from Subscribe & Save subscriptions, sorted by delivery date.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := amazon.NewClient()
		deliveries, err := client.GetUpcomingDeliveries()
		if err != nil {
			return fmt.Errorf("failed to get upcoming deliveries: %w", err)
		}

		output, err := json.MarshalIndent(deliveries, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(subscriptionsCmd)

	// Add subcommands
	subscriptionsCmd.AddCommand(subscriptionsListCmd)
	subscriptionsCmd.AddCommand(subscriptionsGetCmd)
	subscriptionsCmd.AddCommand(subscriptionsSkipCmd)
	subscriptionsCmd.AddCommand(subscriptionsFrequencyCmd)
	subscriptionsCmd.AddCommand(subscriptionsCancelCmd)
	subscriptionsCmd.AddCommand(subscriptionsUpcomingCmd)

	// Add flags to commands that need confirmation
	subscriptionsSkipCmd.Flags().BoolVar(&confirmFlag, "confirm", false, "Confirm the action")
	subscriptionsCancelCmd.Flags().BoolVar(&confirmFlag, "confirm", false, "Confirm the action")

	// Add flags for frequency command
	subscriptionsFrequencyCmd.Flags().IntVar(&intervalFlag, "interval", 0, "Delivery interval in weeks (1-26)")
	subscriptionsFrequencyCmd.Flags().BoolVar(&confirmFlag, "confirm", false, "Confirm the action")
	subscriptionsFrequencyCmd.MarkFlagRequired("interval")
}
