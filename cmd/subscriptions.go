package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/michaelshimeles/amazon-cli/internal/amazon"
	"github.com/spf13/cobra"
)

// subscriptionsCmd represents the subscriptions command
var subscriptionsCmd = &cobra.Command{
	Use:   "subscriptions",
	Short: "Manage Subscribe & Save subscriptions",
	Long: `Manage Amazon Subscribe & Save subscriptions including listing,
viewing details, skipping deliveries, changing frequency, and cancelling.`,
}

// subscriptionsListCmd represents the subscriptions list command
var subscriptionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Subscribe & Save subscriptions",
	Long:  `Retrieve and display all active and paused Subscribe & Save subscriptions.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := amazon.NewClient()
		subscriptions, err := client.GetSubscriptions()
		if err != nil {
			return fmt.Errorf("failed to get subscriptions: %w", err)
		}

		output, err := json.MarshalIndent(subscriptions, "", "  ")
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
	Long:  `Retrieve and display detailed information for a specific Subscribe & Save subscription.`,
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
	Long:  `Skip the next scheduled delivery for a Subscribe & Save subscription. Requires --confirm flag.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID := args[0]
		confirm, _ := cmd.Flags().GetBool("confirm")

		client := amazon.NewClient()

		if !confirm {
			// Dry run mode - show what would be skipped
			subscription, err := client.GetSubscription(subscriptionID)
			if err != nil {
				return fmt.Errorf("failed to get subscription: %w", err)
			}

			dryRun := map[string]interface{}{
				"dry_run":           true,
				"subscription":      subscription,
				"next_delivery":     subscription.NextDelivery,
				"message":           "Add --confirm to execute",
			}

			output, err := json.MarshalIndent(dryRun, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to format output: %w", err)
			}

			fmt.Println(string(output))
			return nil
		}

		// Execute skip
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
	Short: "Change delivery frequency for a subscription",
	Long:  `Change the delivery frequency for a Subscribe & Save subscription. Requires --interval and --confirm flags.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID := args[0]
		confirm, _ := cmd.Flags().GetBool("confirm")
		interval, _ := cmd.Flags().GetInt("interval")

		if interval == 0 {
			return fmt.Errorf("--interval flag is required")
		}

		client := amazon.NewClient()

		if !confirm {
			// Dry run mode - show what would change
			subscription, err := client.GetSubscription(subscriptionID)
			if err != nil {
				return fmt.Errorf("failed to get subscription: %w", err)
			}

			dryRun := map[string]interface{}{
				"dry_run":              true,
				"subscription":         subscription,
				"current_frequency":    subscription.FrequencyWeeks,
				"new_frequency":        interval,
				"message":              "Add --confirm to execute",
			}

			output, err := json.MarshalIndent(dryRun, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to format output: %w", err)
			}

			fmt.Println(string(output))
			return nil
		}

		// Execute frequency change
		subscription, err := client.UpdateFrequency(subscriptionID, interval)
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
	Short: "Cancel a Subscribe & Save subscription",
	Long:  `Cancel a Subscribe & Save subscription. Requires --confirm flag to prevent accidental cancellations.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID := args[0]
		confirm, _ := cmd.Flags().GetBool("confirm")

		client := amazon.NewClient()

		if !confirm {
			// Dry run mode - show what would be cancelled
			subscription, err := client.GetSubscription(subscriptionID)
			if err != nil {
				return fmt.Errorf("failed to get subscription: %w", err)
			}

			dryRun := map[string]interface{}{
				"dry_run":      true,
				"subscription": subscription,
				"message":      "Add --confirm to execute cancellation",
			}

			output, err := json.MarshalIndent(dryRun, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to format output: %w", err)
			}

			fmt.Println(string(output))
			return nil
		}

		// Execute cancellation
		subscription, err := client.CancelSubscription(subscriptionID)
		if err != nil {
			return fmt.Errorf("failed to cancel subscription: %w", err)
		}

		result := map[string]interface{}{
			"subscription": subscription,
			"status":       "cancelled",
			"message":      fmt.Sprintf("Subscription %s has been cancelled", subscriptionID),
		}

		output, err := json.MarshalIndent(result, "", "  ")
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
	Short: "View upcoming subscription deliveries",
	Long:  `View all upcoming deliveries across all Subscribe & Save subscriptions, sorted by delivery date.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := amazon.NewClient()
		deliveries, err := client.GetUpcomingDeliveries()
		if err != nil {
			return fmt.Errorf("failed to get upcoming deliveries: %w", err)
		}

		result := map[string]interface{}{
			"upcoming_deliveries": deliveries,
			"count":               len(deliveries),
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
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

	// Add flags
	subscriptionsSkipCmd.Flags().Bool("confirm", false, "Confirm the skip operation")
	subscriptionsFrequencyCmd.Flags().Bool("confirm", false, "Confirm the frequency change")
	subscriptionsFrequencyCmd.Flags().Int("interval", 0, "New delivery interval in weeks (1-26)")
	subscriptionsCancelCmd.Flags().Bool("confirm", false, "Confirm the cancellation")

	// Mark the output directory flag as hidden if needed
	if err := subscriptionsFrequencyCmd.MarkFlagRequired("interval"); err != nil {
		fmt.Fprintf(os.Stderr, "Error marking flag as required: %v\n", err)
	}
}
