package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/michaelshimeles/amazon-cli/internal/amazon"
	"github.com/spf13/cobra"
)

var (
	subscriptionConfirm bool
	subscriptionInterval int
)

// subscriptionsCmd represents the subscriptions parent command
var subscriptionsCmd = &cobra.Command{
	Use:   "subscriptions",
	Short: "Manage Subscribe & Save subscriptions",
	Long: `Manage your Amazon Subscribe & Save subscriptions.

Commands allow you to list, view, skip, change frequency, and cancel subscriptions.`,
}

// subscriptionsListCmd lists all subscriptions
var subscriptionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Subscribe & Save subscriptions",
	Long:  `Retrieves and displays all active and paused Subscribe & Save subscriptions.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := amazon.NewClient()
		response, err := client.GetSubscriptions()
		if err != nil {
			return fmt.Errorf("failed to get subscriptions: %w", err)
		}

		output, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal response: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// subscriptionsGetCmd gets details for a specific subscription
var subscriptionsGetCmd = &cobra.Command{
	Use:   "get <subscription-id>",
	Short: "Get subscription details",
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
			return fmt.Errorf("failed to marshal response: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// subscriptionsSkipCmd skips the next delivery
var subscriptionsSkipCmd = &cobra.Command{
	Use:   "skip <subscription-id>",
	Short: "Skip next delivery",
	Long:  `Skips the next scheduled delivery for a subscription.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID := args[0]

		client := amazon.NewClient()

		// Without --confirm, show preview
		if !subscriptionConfirm {
			subscription, err := client.GetSubscription(subscriptionID)
			if err != nil {
				return fmt.Errorf("failed to get subscription: %w", err)
			}

			preview := map[string]interface{}{
				"dry_run": true,
				"subscription": subscription,
				"message": "Add --confirm to skip the next delivery",
			}

			output, err := json.MarshalIndent(preview, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal preview: %w", err)
			}

			fmt.Println(string(output))
			return nil
		}

		// With --confirm, execute skip
		subscription, err := client.SkipDelivery(subscriptionID)
		if err != nil {
			return fmt.Errorf("failed to skip delivery: %w", err)
		}

		output, err := json.MarshalIndent(subscription, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal response: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// subscriptionsFrequencyCmd changes the delivery frequency
var subscriptionsFrequencyCmd = &cobra.Command{
	Use:   "frequency <subscription-id>",
	Short: "Change delivery frequency",
	Long: `Changes the delivery frequency for a subscription.

Frequency is specified in weeks (1-26).`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID := args[0]

		// Validate interval flag is provided
		if subscriptionInterval == 0 {
			return fmt.Errorf("--interval flag is required")
		}

		// Validate interval range
		if subscriptionInterval < 1 || subscriptionInterval > 26 {
			return fmt.Errorf("interval must be between 1 and 26 weeks")
		}

		client := amazon.NewClient()

		// Without --confirm, show preview
		if !subscriptionConfirm {
			subscription, err := client.GetSubscription(subscriptionID)
			if err != nil {
				return fmt.Errorf("failed to get subscription: %w", err)
			}

			preview := map[string]interface{}{
				"dry_run": true,
				"subscription_id": subscription.SubscriptionID,
				"current_frequency_weeks": subscription.FrequencyWeeks,
				"new_frequency_weeks": subscriptionInterval,
				"message": "Add --confirm to update the frequency",
			}

			output, err := json.MarshalIndent(preview, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal preview: %w", err)
			}

			fmt.Println(string(output))
			return nil
		}

		// With --confirm, execute frequency update
		subscription, err := client.UpdateFrequency(subscriptionID, subscriptionInterval)
		if err != nil {
			return fmt.Errorf("failed to update frequency: %w", err)
		}

		output, err := json.MarshalIndent(subscription, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal response: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// subscriptionsCancelCmd cancels a subscription
var subscriptionsCancelCmd = &cobra.Command{
	Use:   "cancel <subscription-id>",
	Short: "Cancel subscription",
	Long:  `Cancels a Subscribe & Save subscription.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID := args[0]

		client := amazon.NewClient()

		// Without --confirm, show preview
		if !subscriptionConfirm {
			subscription, err := client.GetSubscription(subscriptionID)
			if err != nil {
				return fmt.Errorf("failed to get subscription: %w", err)
			}

			preview := map[string]interface{}{
				"dry_run": true,
				"subscription": subscription,
				"message": "Add --confirm to cancel this subscription",
			}

			output, err := json.MarshalIndent(preview, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal preview: %w", err)
			}

			fmt.Println(string(output))
			return nil
		}

		// With --confirm, execute cancellation
		subscription, err := client.CancelSubscription(subscriptionID)
		if err != nil {
			return fmt.Errorf("failed to cancel subscription: %w", err)
		}

		output, err := json.MarshalIndent(subscription, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal response: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// subscriptionsUpcomingCmd shows upcoming deliveries
var subscriptionsUpcomingCmd = &cobra.Command{
	Use:   "upcoming",
	Short: "View upcoming deliveries",
	Long:  `Displays all upcoming Subscribe & Save deliveries sorted by date.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := amazon.NewClient()
		deliveries, err := client.GetUpcomingDeliveries()
		if err != nil {
			return fmt.Errorf("failed to get upcoming deliveries: %w", err)
		}

		response := map[string]interface{}{
			"upcoming_deliveries": deliveries,
		}

		output, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal response: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

func init() {
	// Add subscriptions command to root
	rootCmd.AddCommand(subscriptionsCmd)

	// Add subcommands
	subscriptionsCmd.AddCommand(subscriptionsListCmd)
	subscriptionsCmd.AddCommand(subscriptionsGetCmd)
	subscriptionsCmd.AddCommand(subscriptionsSkipCmd)
	subscriptionsCmd.AddCommand(subscriptionsFrequencyCmd)
	subscriptionsCmd.AddCommand(subscriptionsCancelCmd)
	subscriptionsCmd.AddCommand(subscriptionsUpcomingCmd)

	// Add flags to skip command
	subscriptionsSkipCmd.Flags().BoolVar(&subscriptionConfirm, "confirm", false, "Confirm the action")

	// Add flags to frequency command
	subscriptionsFrequencyCmd.Flags().IntVar(&subscriptionInterval, "interval", 0, "Delivery frequency in weeks (1-26)")
	subscriptionsFrequencyCmd.Flags().BoolVar(&subscriptionConfirm, "confirm", false, "Confirm the action")

	// Add flags to cancel command
	subscriptionsCancelCmd.Flags().BoolVar(&subscriptionConfirm, "confirm", false, "Confirm the action")
}
