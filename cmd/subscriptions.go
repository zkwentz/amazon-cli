package cmd

import (
	"encoding/json"
	"fmt"
	"os"

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
	Short: "Manage Amazon Subscribe & Save subscriptions",
	Long: `Manage your Amazon Subscribe & Save subscriptions.

You can list all subscriptions, view details, skip deliveries,
change delivery frequency, cancel subscriptions, and view upcoming deliveries.`,
}

// subscriptionsListCmd lists all subscriptions
var subscriptionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all subscriptions",
	Long:  `List all active and paused Subscribe & Save subscriptions.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := amazon.NewClient()
		response, err := client.GetSubscriptions()
		if err != nil {
			return fmt.Errorf("failed to get subscriptions: %w", err)
		}

		// Output JSON
		output, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// subscriptionsGetCmd gets details for a specific subscription
var subscriptionsGetCmd = &cobra.Command{
	Use:   "get <subscription-id>",
	Short: "Get subscription details",
	Long:  `Get detailed information for a specific subscription.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID := args[0]

		client := amazon.NewClient()
		subscription, err := client.GetSubscription(subscriptionID)
		if err != nil {
			return fmt.Errorf("failed to get subscription: %w", err)
		}

		// Output JSON
		output, err := json.MarshalIndent(subscription, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// subscriptionsSkipCmd skips the next delivery
var subscriptionsSkipCmd = &cobra.Command{
	Use:   "skip <subscription-id>",
	Short: "Skip next delivery",
	Long:  `Skip the next scheduled delivery for a subscription.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID := args[0]

		client := amazon.NewClient()

		// Check if --confirm flag is provided
		if !confirmFlag {
			// Dry run: show what would be skipped
			subscription, err := client.GetSubscription(subscriptionID)
			if err != nil {
				return fmt.Errorf("failed to get subscription: %w", err)
			}

			dryRun := map[string]interface{}{
				"dry_run":       true,
				"subscription":  subscription,
				"message":       "Add --confirm to execute",
				"next_delivery": subscription.NextDelivery,
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

		// Output JSON
		output, err := json.MarshalIndent(subscription, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// subscriptionsFrequencyCmd changes the delivery frequency
var subscriptionsFrequencyCmd = &cobra.Command{
	Use:   "frequency <subscription-id>",
	Short: "Change delivery frequency",
	Long:  `Change the delivery frequency for a subscription (in weeks).`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID := args[0]

		// Validate interval flag is provided
		if intervalFlag == 0 {
			return fmt.Errorf("--interval flag is required")
		}

		client := amazon.NewClient()

		// Check if --confirm flag is provided
		if !confirmFlag {
			// Dry run: show what would change
			subscription, err := client.GetSubscription(subscriptionID)
			if err != nil {
				return fmt.Errorf("failed to get subscription: %w", err)
			}

			dryRun := map[string]interface{}{
				"dry_run":             true,
				"subscription":        subscription,
				"current_frequency":   subscription.FrequencyWeeks,
				"new_frequency":       intervalFlag,
				"message":             "Add --confirm to execute",
			}

			output, err := json.MarshalIndent(dryRun, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to format output: %w", err)
			}

			fmt.Println(string(output))
			return nil
		}

		// Execute frequency change
		subscription, err := client.UpdateFrequency(subscriptionID, intervalFlag)
		if err != nil {
			return fmt.Errorf("failed to update frequency: %w", err)
		}

		// Output JSON
		output, err := json.MarshalIndent(subscription, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// subscriptionsCancelCmd cancels a subscription
var subscriptionsCancelCmd = &cobra.Command{
	Use:   "cancel <subscription-id>",
	Short: "Cancel subscription",
	Long:  `Cancel a Subscribe & Save subscription.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		subscriptionID := args[0]

		client := amazon.NewClient()

		// Check if --confirm flag is provided
		if !confirmFlag {
			// Dry run: show cancellation preview
			subscription, err := client.GetSubscription(subscriptionID)
			if err != nil {
				return fmt.Errorf("failed to get subscription: %w", err)
			}

			dryRun := map[string]interface{}{
				"dry_run":      true,
				"subscription": subscription,
				"message":      "Add --confirm to execute cancellation",
				"warning":      "This action cannot be undone",
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

		// Output JSON
		output, err := json.MarshalIndent(subscription, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// subscriptionsUpcomingCmd lists upcoming deliveries
var subscriptionsUpcomingCmd = &cobra.Command{
	Use:   "upcoming",
	Short: "View upcoming deliveries",
	Long:  `View all upcoming Subscribe & Save deliveries sorted by date.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := amazon.NewClient()
		deliveries, err := client.GetUpcomingDeliveries()
		if err != nil {
			return fmt.Errorf("failed to get upcoming deliveries: %w", err)
		}

		// Output JSON
		output, err := json.MarshalIndent(deliveries, "", "  ")
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

	// Add subcommands to subscriptions
	subscriptionsCmd.AddCommand(subscriptionsListCmd)
	subscriptionsCmd.AddCommand(subscriptionsGetCmd)
	subscriptionsCmd.AddCommand(subscriptionsSkipCmd)
	subscriptionsCmd.AddCommand(subscriptionsFrequencyCmd)
	subscriptionsCmd.AddCommand(subscriptionsCancelCmd)
	subscriptionsCmd.AddCommand(subscriptionsUpcomingCmd)

	// Add flags for skip, frequency, and cancel commands
	subscriptionsSkipCmd.Flags().BoolVar(&confirmFlag, "confirm", false, "Confirm the action")
	subscriptionsFrequencyCmd.Flags().BoolVar(&confirmFlag, "confirm", false, "Confirm the action")
	subscriptionsFrequencyCmd.Flags().IntVar(&intervalFlag, "interval", 0, "Delivery frequency in weeks (1-26)")
	subscriptionsCancelCmd.Flags().BoolVar(&confirmFlag, "confirm", false, "Confirm the cancellation")

	// Set custom error handling
	subscriptionsSkipCmd.SilenceUsage = true
	subscriptionsFrequencyCmd.SilenceUsage = true
	subscriptionsCancelCmd.SilenceUsage = true
	subscriptionsListCmd.SilenceUsage = true
	subscriptionsGetCmd.SilenceUsage = true
	subscriptionsUpcomingCmd.SilenceUsage = true

	// Ensure errors are printed to stderr
	subscriptionsCmd.SetErr(os.Stderr)
}
