package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/michaelshimeles/amazon-cli/internal/amazon"
	"github.com/spf13/cobra"
)

var (
	subscriptionsConfirm bool
	subscriptionsInterval int
)

// subscriptionsCmd represents the subscriptions parent command
var subscriptionsCmd = &cobra.Command{
	Use:   "subscriptions",
	Short: "Manage Subscribe & Save subscriptions",
	Long: `Manage your Subscribe & Save subscriptions including:
- List all subscriptions
- View subscription details
- Skip next delivery
- Change delivery frequency
- Cancel subscriptions
- View upcoming deliveries`,
}

// subscriptionsListCmd lists all subscriptions
var subscriptionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Subscribe & Save subscriptions",
	Long:  "Retrieves and displays all active and paused Subscribe & Save subscriptions.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := amazon.NewClient()
		subscriptions, err := client.GetSubscriptions()
		if err != nil {
			return fmt.Errorf("failed to get subscriptions: %w", err)
		}

		output, err := json.MarshalIndent(subscriptions, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal subscriptions: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// subscriptionsGetCmd gets details for a specific subscription
var subscriptionsGetCmd = &cobra.Command{
	Use:   "get <subscription-id>",
	Short: "Get details for a specific subscription",
	Long:  "Retrieves detailed information about a specific Subscribe & Save subscription.",
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
			return fmt.Errorf("failed to marshal subscription: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// subscriptionsSkipCmd skips the next delivery for a subscription
var subscriptionsSkipCmd = &cobra.Command{
	Use:   "skip <subscription-id>",
	Short: "Skip the next delivery for a subscription",
	Long: `Skips the next scheduled delivery for a Subscribe & Save subscription.
The delivery will be rescheduled to the following period based on your frequency setting.

REQUIRES --confirm flag to execute. Without it, shows a preview of what would be skipped.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID := args[0]
		client := amazon.NewClient()

		// Get current subscription details
		subscription, err := client.GetSubscription(subscriptionID)
		if err != nil {
			return fmt.Errorf("failed to get subscription: %w", err)
		}

		// If --confirm not provided, show dry run
		if !subscriptionsConfirm {
			dryRun := map[string]interface{}{
				"dry_run":     true,
				"subscription": subscription,
				"message":     "Add --confirm to skip this delivery",
			}
			output, err := json.MarshalIndent(dryRun, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal dry run: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		// Execute the skip
		updatedSubscription, err := client.SkipDelivery(subscriptionID)
		if err != nil {
			return fmt.Errorf("failed to skip delivery: %w", err)
		}

		output, err := json.MarshalIndent(updatedSubscription, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal updated subscription: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// subscriptionsFrequencyCmd changes the delivery frequency
var subscriptionsFrequencyCmd = &cobra.Command{
	Use:   "frequency <subscription-id>",
	Short: "Change the delivery frequency for a subscription",
	Long: `Changes how often you receive Subscribe & Save deliveries.

REQUIRES --confirm flag to execute. Without it, shows a preview of what would change.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID := args[0]
		client := amazon.NewClient()

		// Validate interval flag was provided
		if subscriptionsInterval <= 0 {
			return fmt.Errorf("--interval flag is required and must be positive")
		}

		// Get current subscription details
		subscription, err := client.GetSubscription(subscriptionID)
		if err != nil {
			return fmt.Errorf("failed to get subscription: %w", err)
		}

		// If --confirm not provided, show dry run
		if !subscriptionsConfirm {
			dryRun := map[string]interface{}{
				"dry_run":         true,
				"subscription":    subscription,
				"new_frequency":   subscriptionsInterval,
				"message":         "Add --confirm to change the frequency",
			}
			output, err := json.MarshalIndent(dryRun, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal dry run: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		// Execute the frequency change
		updatedSubscription, err := client.UpdateFrequency(subscriptionID, subscriptionsInterval)
		if err != nil {
			return fmt.Errorf("failed to update frequency: %w", err)
		}

		output, err := json.MarshalIndent(updatedSubscription, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal updated subscription: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// subscriptionsCancelCmd cancels a subscription
var subscriptionsCancelCmd = &cobra.Command{
	Use:   "cancel <subscription-id>",
	Short: "Cancel a subscription",
	Long: `Cancels a Subscribe & Save subscription. This will stop all future deliveries.

REQUIRES --confirm flag to execute. Without it, shows a preview of what would be cancelled.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID := args[0]
		client := amazon.NewClient()

		// Get current subscription details
		subscription, err := client.GetSubscription(subscriptionID)
		if err != nil {
			return fmt.Errorf("failed to get subscription: %w", err)
		}

		// If --confirm not provided, show dry run
		if !subscriptionsConfirm {
			dryRun := map[string]interface{}{
				"dry_run":      true,
				"subscription": subscription,
				"message":      "Add --confirm to cancel this subscription",
			}
			output, err := json.MarshalIndent(dryRun, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal dry run: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		// Execute the cancellation
		cancelledSubscription, err := client.CancelSubscription(subscriptionID)
		if err != nil {
			return fmt.Errorf("failed to cancel subscription: %w", err)
		}

		output, err := json.MarshalIndent(cancelledSubscription, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal cancelled subscription: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// subscriptionsUpcomingCmd lists upcoming deliveries
var subscriptionsUpcomingCmd = &cobra.Command{
	Use:   "upcoming",
	Short: "View upcoming subscription deliveries",
	Long:  "Lists all upcoming Subscribe & Save deliveries across all active subscriptions, sorted by delivery date.",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := amazon.NewClient()
		deliveries, err := client.GetUpcomingDeliveries()
		if err != nil {
			return fmt.Errorf("failed to get upcoming deliveries: %w", err)
		}

		output, err := json.MarshalIndent(deliveries, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal deliveries: %w", err)
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

	// Add flags for skip command
	subscriptionsSkipCmd.Flags().BoolVar(&subscriptionsConfirm, "confirm", false, "Confirm the skip action")

	// Add flags for frequency command
	subscriptionsFrequencyCmd.Flags().IntVar(&subscriptionsInterval, "interval", 0, "Delivery frequency in weeks")
	subscriptionsFrequencyCmd.Flags().BoolVar(&subscriptionsConfirm, "confirm", false, "Confirm the frequency change")
	subscriptionsFrequencyCmd.MarkFlagRequired("interval")

	// Add flags for cancel command
	subscriptionsCancelCmd.Flags().BoolVar(&subscriptionsConfirm, "confirm", false, "Confirm the cancellation")

	// Redirect stderr to stdout to avoid usage message on errors
	cobra.OnInitialize(func() {
		subscriptionsSkipCmd.SilenceUsage = true
		subscriptionsFrequencyCmd.SilenceUsage = true
		subscriptionsCancelCmd.SilenceUsage = true
	})
}

func setSubscriptionsStderr(w *os.File) {
	subscriptionsCmd.SetErr(w)
	subscriptionsListCmd.SetErr(w)
	subscriptionsGetCmd.SetErr(w)
	subscriptionsSkipCmd.SetErr(w)
	subscriptionsFrequencyCmd.SetErr(w)
	subscriptionsCancelCmd.SetErr(w)
	subscriptionsUpcomingCmd.SetErr(w)
}
