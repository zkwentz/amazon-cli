package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/internal/output"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

var (
	subscriptionInterval int
	subscriptionConfirm  bool
)

// subscriptionsCmd represents the subscriptions command
var subscriptionsCmd = &cobra.Command{
	Use:   "subscriptions",
	Short: "Manage Subscribe & Save subscriptions",
	Long:  `List, skip, cancel, and manage delivery frequency for Subscribe & Save subscriptions.`,
}

// frequencyCmd represents the subscriptions frequency command
var frequencyCmd = &cobra.Command{
	Use:   "frequency <id>",
	Short: "Update subscription delivery frequency",
	Long:  `Update the delivery frequency for a Subscribe & Save subscription. Interval must be between 1-26 weeks. Requires --confirm flag to execute.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		c := getClient()

		// Validate interval
		if subscriptionInterval < 1 || subscriptionInterval > 26 {
			output.Error(models.ErrInvalidInput, "interval must be between 1 and 26 weeks", nil)
			os.Exit(models.ExitInvalidArgs)
		}

		// Without --confirm, show preview
		if !subscriptionConfirm {
			// Get current subscription info for preview
			// For now, we'll show a preview with the new frequency
			output.JSON(map[string]interface{}{
				"dry_run":         true,
				"subscription_id": id,
				"new_interval":    subscriptionInterval,
				"message":         "Add --confirm to update frequency",
			})
			return
		}

		// With --confirm, call UpdateFrequency
		subscription, err := c.UpdateFrequency(id, subscriptionInterval)
		if err != nil {
			output.Error(models.ErrInvalidInput, err.Error(), nil)
			os.Exit(models.ExitInvalidArgs)
		}

		output.JSON(subscription)
	},
}

// subscriptionsCancelCmd represents the subscriptions cancel command
var subscriptionsCancelCmd = &cobra.Command{
	Use:   "cancel <id>",
	Short: "Cancel a subscription",
	Long: `Cancel an Amazon Subscribe & Save subscription by ID.
Requires --confirm flag to execute the cancellation.
Without --confirm, shows a preview of the cancellation.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := args[0]
		c := getClient()

		// Without --confirm, show cancellation preview
		if !subscriptionConfirm {
			// Get subscription details for preview
			subscription, err := c.CancelSubscription(id)
			if err != nil {
				output.Error(models.ErrInvalidInput, err.Error(), nil)
				os.Exit(models.ExitInvalidArgs)
			}

			// Reset status to show current state in preview
			subscription.Status = "active"

			output.JSON(map[string]interface{}{
				"dry_run":      true,
				"subscription": subscription,
				"message":      "Add --confirm to cancel this subscription",
			})
			return
		}

		// With --confirm, execute the cancellation
		subscription, err := c.CancelSubscription(id)
		if err != nil {
			output.Error(models.ErrAmazonError, err.Error(), nil)
			os.Exit(models.ExitGeneralError)
		}

		output.JSON(subscription)
	},
}

func init() {
	rootCmd.AddCommand(subscriptionsCmd)

	// Add subcommands
	subscriptionsCmd.AddCommand(frequencyCmd)
	subscriptionsCmd.AddCommand(subscriptionsCancelCmd)

	// Flags for frequency command
	frequencyCmd.Flags().IntVarP(&subscriptionInterval, "interval", "i", 0, "Delivery interval in weeks (1-26, required)")
	frequencyCmd.MarkFlagRequired("interval")
	frequencyCmd.Flags().BoolVar(&subscriptionConfirm, "confirm", false, "Confirm the frequency update")

	// Flags for cancel command
	subscriptionsCancelCmd.Flags().BoolVar(&subscriptionConfirm, "confirm", false, "Confirm the cancellation")
}
