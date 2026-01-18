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
	Long: `Manage Amazon Subscribe & Save subscriptions.

This command allows you to:
- List all subscriptions
- Get details for a specific subscription
- Skip next delivery
- Change delivery frequency
- Cancel subscriptions
- View upcoming deliveries`,
}

// subscriptionsGetCmd represents the subscriptions get command
var subscriptionsGetCmd = &cobra.Command{
	Use:   "get <subscription-id>",
	Short: "Get details for a specific subscription",
	Long: `Get detailed information about a specific Subscribe & Save subscription.

Example:
  amazon-cli subscriptions get S01-1234567-8901234

Output includes:
- Subscription ID
- Product ASIN and title
- Price and discount percentage
- Delivery frequency
- Next delivery date
- Status (active, paused, cancelled)
- Quantity`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		subscriptionID := args[0]

		// Create Amazon client
		client := amazon.NewClient()

		// Get subscription details
		subscription, err := client.GetSubscription(subscriptionID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Output JSON response
		jsonOutput, err := json.MarshalIndent(subscription, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(string(jsonOutput))
	},
}

func init() {
	rootCmd.AddCommand(subscriptionsCmd)
	subscriptionsCmd.AddCommand(subscriptionsGetCmd)
}
