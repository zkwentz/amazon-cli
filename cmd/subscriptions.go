package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/michaelshimeles/amazon-cli/internal/amazon"
	"github.com/spf13/cobra"
)

var (
	confirmFlag  bool
	intervalFlag int
)

// subscriptionsCmd represents the subscriptions parent command
var subscriptionsCmd = &cobra.Command{
	Use:   "subscriptions",
	Short: "Manage Subscribe & Save subscriptions",
	Long: `Manage Subscribe & Save subscriptions.

View, modify, and cancel your Subscribe & Save subscriptions, adjust delivery
frequencies, skip deliveries, and view upcoming deliveries.`,
}

// subscriptionsListCmd represents the subscriptions list command
var subscriptionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Subscribe & Save subscriptions",
	Long: `List all Subscribe & Save subscriptions with details including:
- Subscription ID
- Product ASIN and title
- Price and discount percentage
- Delivery frequency
- Next delivery date
- Status (active, paused, cancelled)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := amazon.NewClient()

		subscriptions, err := client.GetSubscriptions()
		if err != nil {
			return fmt.Errorf("failed to get subscriptions: %w", err)
		}

		// Output as JSON
		output, err := json.MarshalIndent(subscriptions, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal response: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// subscriptionsGetCmd represents the subscriptions get command
var subscriptionsGetCmd = &cobra.Command{
	Use:   "get <subscription-id>",
	Short: "Get details for a specific subscription",
	Long:  `Get detailed information about a specific Subscribe & Save subscription.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID := args[0]
		client := amazon.NewClient()

		subscription, err := client.GetSubscription(subscriptionID)
		if err != nil {
			return fmt.Errorf("failed to get subscription: %w", err)
		}

		// Output as JSON
		output, err := json.MarshalIndent(subscription, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal response: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// subscriptionsSkipCmd represents the subscriptions skip command
var subscriptionsSkipCmd = &cobra.Command{
	Use:   "skip <subscription-id>",
	Short: "Skip next delivery for a subscription",
	Long:  `Skip the next delivery for a Subscribe & Save subscription.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID := args[0]
		client := amazon.NewClient()

		// Check if --confirm flag is provided
		if !confirmFlag {
			subscription, err := client.GetSubscription(subscriptionID)
			if err != nil {
				return fmt.Errorf("failed to get subscription: %w", err)
			}

			dryRun := map[string]interface{}{
				"dry_run": true,
				"would_skip": map[string]interface{}{
					"subscription_id": subscription.SubscriptionID,
					"title":           subscription.Title,
					"next_delivery":   subscription.NextDelivery,
				},
				"message": "Add --confirm to execute",
			}

			output, err := json.MarshalIndent(dryRun, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal response: %w", err)
			}

			fmt.Println(string(output))
			return nil
		}

		// Execute skip delivery
		subscription, err := client.SkipDelivery(subscriptionID)
		if err != nil {
			return fmt.Errorf("failed to skip delivery: %w", err)
		}

		// Output as JSON
		output, err := json.MarshalIndent(subscription, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal response: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// subscriptionsFrequencyCmd represents the subscriptions frequency command
var subscriptionsFrequencyCmd = &cobra.Command{
	Use:   "frequency <subscription-id>",
	Short: "Change delivery frequency for a subscription",
	Long:  `Change the delivery frequency for a Subscribe & Save subscription.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID := args[0]

		// Validate interval flag is provided
		if intervalFlag == 0 {
			return fmt.Errorf("--interval flag is required (delivery frequency in weeks)")
		}

		// Validate interval is reasonable (1-26 weeks)
		if intervalFlag < 1 || intervalFlag > 26 {
			return fmt.Errorf("interval must be between 1 and 26 weeks")
		}

		client := amazon.NewClient()

		// Check if --confirm flag is provided
		if !confirmFlag {
			subscription, err := client.GetSubscription(subscriptionID)
			if err != nil {
				return fmt.Errorf("failed to get subscription: %w", err)
			}

			dryRun := map[string]interface{}{
				"dry_run": true,
				"would_change": map[string]interface{}{
					"subscription_id":     subscription.SubscriptionID,
					"title":               subscription.Title,
					"current_frequency":   subscription.FrequencyWeeks,
					"new_frequency":       intervalFlag,
				},
				"message": "Add --confirm to execute",
			}

			output, err := json.MarshalIndent(dryRun, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal response: %w", err)
			}

			fmt.Println(string(output))
			return nil
		}

		// Execute frequency update
		subscription, err := client.UpdateFrequency(subscriptionID, intervalFlag)
		if err != nil {
			return fmt.Errorf("failed to update frequency: %w", err)
		}

		// Output as JSON
		output, err := json.MarshalIndent(subscription, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal response: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// subscriptionsCancelCmd represents the subscriptions cancel command
var subscriptionsCancelCmd = &cobra.Command{
	Use:   "cancel <subscription-id>",
	Short: "Cancel a subscription",
	Long:  `Cancel a Subscribe & Save subscription.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID := args[0]
		client := amazon.NewClient()

		// Check if --confirm flag is provided
		if !confirmFlag {
			subscription, err := client.GetSubscription(subscriptionID)
			if err != nil {
				return fmt.Errorf("failed to get subscription: %w", err)
			}

			dryRun := map[string]interface{}{
				"dry_run": true,
				"would_cancel": map[string]interface{}{
					"subscription_id": subscription.SubscriptionID,
					"title":           subscription.Title,
					"price":           subscription.Price,
				},
				"message": "Add --confirm to execute",
			}

			output, err := json.MarshalIndent(dryRun, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal response: %w", err)
			}

			fmt.Println(string(output))
			return nil
		}

		// Execute cancellation
		subscription, err := client.CancelSubscription(subscriptionID)
		if err != nil {
			return fmt.Errorf("failed to cancel subscription: %w", err)
		}

		// Output as JSON
		output, err := json.MarshalIndent(subscription, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal response: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// subscriptionsUpcomingCmd represents the subscriptions upcoming command
var subscriptionsUpcomingCmd = &cobra.Command{
	Use:   "upcoming",
	Short: "View upcoming subscription deliveries",
	Long:  `View all upcoming Subscribe & Save deliveries sorted by delivery date.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := amazon.NewClient()

		deliveries, err := client.GetUpcomingDeliveries()
		if err != nil {
			return fmt.Errorf("failed to get upcoming deliveries: %w", err)
		}

		// Output as JSON
		output, err := json.MarshalIndent(deliveries, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal response: %w", err)
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

	// Add flags for skip, frequency, and cancel commands
	subscriptionsSkipCmd.Flags().BoolVar(&confirmFlag, "confirm", false, "Confirm skip delivery action")
	subscriptionsFrequencyCmd.Flags().BoolVar(&confirmFlag, "confirm", false, "Confirm frequency change action")
	subscriptionsFrequencyCmd.Flags().IntVar(&intervalFlag, "interval", 0, "Delivery frequency in weeks (1-26)")
	subscriptionsCancelCmd.Flags().BoolVar(&confirmFlag, "confirm", false, "Confirm cancellation action")
}
