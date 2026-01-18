package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/internal/amazon"
	"github.com/zkwentz/amazon-cli/internal/output"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

var (
	returnReason  string
	confirmReturn bool
)

// returnsCmd represents the returns command
var returnsCmd = &cobra.Command{
	Use:   "returns",
	Short: "Manage Amazon returns",
	Long: `Manage Amazon returns including listing returnable items, getting return options,
creating returns, and tracking return status.`,
}

// listReturnsCmd represents the returns list command
var listReturnsCmd = &cobra.Command{
	Use:   "list",
	Short: "List returnable items",
	Long:  `List all items eligible for return from recent orders.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := amazon.NewClient()
		printer := output.NewPrinter("json", false)

		items, err := client.GetReturnableItems()
		if err != nil {
			printer.PrintError(err)
			return err
		}

		return printer.Print(map[string]interface{}{
			"returnable_items": items,
			"count":            len(items),
		})
	},
}

// optionsReturnsCmd represents the returns options command
var optionsReturnsCmd = &cobra.Command{
	Use:   "options <order-id> <item-id>",
	Short: "Get return options for an item",
	Long:  `Get available return options (methods, dropoff locations, fees) for a specific item.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		orderID := args[0]
		itemID := args[1]

		client := amazon.NewClient()
		printer := output.NewPrinter("json", false)

		options, err := client.GetReturnOptions(orderID, itemID)
		if err != nil {
			printer.PrintError(err)
			return err
		}

		return printer.Print(map[string]interface{}{
			"order_id": orderID,
			"item_id":  itemID,
			"options":  options,
		})
	},
}

// createReturnsCmd represents the returns create command
var createReturnsCmd = &cobra.Command{
	Use:   "create <order-id> <item-id>",
	Short: "Initiate a return",
	Long: `Initiate a return for a specific item. Requires --reason and --confirm flags.
Without --confirm, this command will show what would be returned (dry run).

Valid reason codes:
  - defective: Item is defective or doesn't work
  - wrong_item: Received wrong item
  - not_as_described: Item not as described
  - no_longer_needed: No longer needed
  - better_price: Found better price elsewhere
  - other: Other reason`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		orderID := args[0]
		itemID := args[1]

		client := amazon.NewClient()
		printer := output.NewPrinter("json", false)

		// Validate reason code
		if returnReason == "" {
			return fmt.Errorf("--reason flag is required")
		}

		if !models.IsValidReturnReason(returnReason) {
			return fmt.Errorf("invalid return reason: %s. Valid reasons: defective, wrong_item, not_as_described, no_longer_needed, better_price, other", returnReason)
		}

		// Dry run mode (without --confirm)
		if !confirmReturn {
			return printer.Print(map[string]interface{}{
				"dry_run": true,
				"would_return": map[string]string{
					"order_id": orderID,
					"item_id":  itemID,
					"reason":   returnReason,
				},
				"message": "Add --confirm to execute this return",
			})
		}

		// Execute the return
		returnInfo, err := client.CreateReturn(orderID, itemID, returnReason)
		if err != nil {
			printer.PrintError(err)
			return err
		}

		return printer.Print(returnInfo)
	},
}

// labelReturnsCmd represents the returns label command
var labelReturnsCmd = &cobra.Command{
	Use:   "label <return-id>",
	Short: "Get return label",
	Long:  `Get the return shipping label for an initiated return.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		returnID := args[0]

		client := amazon.NewClient()
		printer := output.NewPrinter("json", false)

		label, err := client.GetReturnLabel(returnID)
		if err != nil {
			printer.PrintError(err)
			return err
		}

		return printer.Print(map[string]interface{}{
			"return_id": returnID,
			"label":     label,
		})
	},
}

// statusReturnsCmd represents the returns status command
var statusReturnsCmd = &cobra.Command{
	Use:   "status <return-id>",
	Short: "Check return status",
	Long:  `Check the current status of a return (initiated, shipped, received, refunded).`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		returnID := args[0]

		client := amazon.NewClient()
		printer := output.NewPrinter("json", false)

		status, err := client.GetReturnStatus(returnID)
		if err != nil {
			printer.PrintError(err)
			return err
		}

		return printer.Print(status)
	},
}

func init() {
	// Add flags to create command
	createReturnsCmd.Flags().StringVarP(&returnReason, "reason", "r", "", "Return reason code (required)")
	createReturnsCmd.Flags().BoolVar(&confirmReturn, "confirm", false, "Confirm the return (required to execute)")

	// Add subcommands to returns
	returnsCmd.AddCommand(listReturnsCmd)
	returnsCmd.AddCommand(optionsReturnsCmd)
	returnsCmd.AddCommand(createReturnsCmd)
	returnsCmd.AddCommand(labelReturnsCmd)
	returnsCmd.AddCommand(statusReturnsCmd)
}

// GetReturnsCmd returns the returns command for registration with root
func GetReturnsCmd() *cobra.Command {
	return returnsCmd
}
