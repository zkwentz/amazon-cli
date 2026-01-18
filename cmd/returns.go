package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/internal/amazon"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

// returnsCmd represents the returns command
var returnsCmd = &cobra.Command{
	Use:   "returns",
	Short: "Manage Amazon returns",
	Long:  `Commands for managing Amazon returns including listing returnable items, creating returns, and checking return status.`,
}

// returnsStatusCmd represents the returns status command
var returnsStatusCmd = &cobra.Command{
	Use:   "status <return-id>",
	Short: "Check the status of a return",
	Long: `Check the current status of an Amazon return by its return ID.

The return status can be one of:
- initiated: Return has been created but item not yet shipped
- shipped: Item has been shipped back to Amazon
- received: Amazon has received the returned item
- refunded: Refund has been processed

Example:
  amazon-cli returns status RET123456789`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		returnID := args[0]

		// Get configuration
		cfg, err := getConfig()
		if err != nil {
			printer := getPrinter()
			handleError(err, printer)
			return
		}

		// Create Amazon client
		client := amazon.NewClient(cfg)

		// Get return status
		returnStatus, err := client.GetReturnStatus(returnID)
		if err != nil {
			printer := getPrinter()
			handleError(err, printer)
			return
		}

		// Output the result
		printer := getPrinter()
		if err := printer.Print(returnStatus); err != nil {
			handleError(err, printer)
			return
		}
	},
}

// returnsListCmd represents the returns list command
var returnsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List returnable items",
	Long:  `List all items that are eligible for return.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := getConfig()
		if err != nil {
			printer := getPrinter()
			handleError(err, printer)
			return
		}

		client := amazon.NewClient(cfg)
		items, err := client.GetReturnableItems()
		if err != nil {
			printer := getPrinter()
			handleError(err, printer)
			return
		}

		printer := getPrinter()
		if err := printer.Print(items); err != nil {
			handleError(err, printer)
			return
		}
	},
}

// returnsOptionsCmd represents the returns options command
var returnsOptionsCmd = &cobra.Command{
	Use:   "options <order-id> <item-id>",
	Short: "Get return options for an item",
	Long:  `Get available return options (shipping methods, drop-off locations) for a specific item.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		orderID := args[0]
		itemID := args[1]

		cfg, err := getConfig()
		if err != nil {
			printer := getPrinter()
			handleError(err, printer)
			return
		}

		client := amazon.NewClient(cfg)
		options, err := client.GetReturnOptions(orderID, itemID)
		if err != nil {
			printer := getPrinter()
			handleError(err, printer)
			return
		}

		printer := getPrinter()
		if err := printer.Print(options); err != nil {
			handleError(err, printer)
			return
		}
	},
}

// returnsCreateCmd represents the returns create command
var returnsCreateCmd = &cobra.Command{
	Use:   "create <order-id> <item-id>",
	Short: "Create a return for an item",
	Long: `Initiate a return for a specific item in an order.

Requires --reason flag with one of:
- defective: Item is defective or doesn't work
- wrong_item: Received wrong item
- not_as_described: Item not as described
- no_longer_needed: No longer needed
- better_price: Found better price elsewhere
- other: Other reason

Requires --confirm flag to actually create the return.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		orderID := args[0]
		itemID := args[1]

		reason, _ := cmd.Flags().GetString("reason")
		confirm, _ := cmd.Flags().GetBool("confirm")

		if reason == "" {
			printer := getPrinter()
			err := models.NewCLIError(
				models.ErrCodeInvalidInput,
				"--reason flag is required",
				nil,
			)
			handleError(err, printer)
			return
		}

		cfg, err := getConfig()
		if err != nil {
			printer := getPrinter()
			handleError(err, printer)
			return
		}

		client := amazon.NewClient(cfg)

		if !confirm {
			// Dry run mode
			printer := getPrinter()
			dryRunResponse := map[string]interface{}{
				"dry_run": true,
				"message": "Add --confirm to execute the return",
				"order_id": orderID,
				"item_id": itemID,
				"reason": reason,
			}
			if err := printer.Print(dryRunResponse); err != nil {
				handleError(err, printer)
			}
			return
		}

		returnResult, err := client.CreateReturn(orderID, itemID, reason)
		if err != nil {
			printer := getPrinter()
			handleError(err, printer)
			return
		}

		printer := getPrinter()
		if err := printer.Print(returnResult); err != nil {
			handleError(err, printer)
			return
		}
	},
}

// returnsLabelCmd represents the returns label command
var returnsLabelCmd = &cobra.Command{
	Use:   "label <return-id>",
	Short: "Get return shipping label",
	Long:  `Get the shipping label for a return.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		returnID := args[0]

		cfg, err := getConfig()
		if err != nil {
			printer := getPrinter()
			handleError(err, printer)
			return
		}

		client := amazon.NewClient(cfg)
		label, err := client.GetReturnLabel(returnID)
		if err != nil {
			printer := getPrinter()
			handleError(err, printer)
			return
		}

		printer := getPrinter()
		if err := printer.Print(label); err != nil {
			handleError(err, printer)
			return
		}
	},
}

func init() {
	// Add returns command to root
	rootCmd.AddCommand(returnsCmd)

	// Add subcommands to returns
	returnsCmd.AddCommand(returnsStatusCmd)
	returnsCmd.AddCommand(returnsListCmd)
	returnsCmd.AddCommand(returnsOptionsCmd)
	returnsCmd.AddCommand(returnsCreateCmd)
	returnsCmd.AddCommand(returnsLabelCmd)

	// Add flags for returns create
	returnsCreateCmd.Flags().String("reason", "", "Reason for return (required)")
	returnsCreateCmd.Flags().Bool("confirm", false, "Confirm the return creation")
}
