package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/internal/amazon"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

var (
	returnReason  string
	returnConfirm bool
)

// returnsCmd represents the returns command
var returnsCmd = &cobra.Command{
	Use:   "returns",
	Short: "Manage Amazon returns",
	Long:  `Commands for managing Amazon returns including listing returnable items, checking options, and creating returns.`,
}

// returnsListCmd lists returnable items
var returnsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List returnable items",
	Long:  `List all items that are eligible for return.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := amazon.NewClient()
		items, err := client.GetReturnableItems()
		if err != nil {
			return outputError(err)
		}

		return outputJSON(items)
	},
}

// returnsOptionsCmd gets return options for an item
var returnsOptionsCmd = &cobra.Command{
	Use:   "options <order-id> <item-id>",
	Short: "Get return options for an item",
	Long:  `Get available return methods and options for a specific item.`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		orderID := args[0]
		itemID := args[1]

		client := amazon.NewClient()
		options, err := client.GetReturnOptions(orderID, itemID)
		if err != nil {
			return outputError(err)
		}

		return outputJSON(options)
	},
}

// returnsCreateCmd creates a return
var returnsCreateCmd = &cobra.Command{
	Use:   "create <order-id> <item-id>",
	Short: "Create a return for an item",
	Long: `Initiate a return for a specific item in an order.

Valid return reason codes:
  - defective: Item is defective or doesn't work
  - wrong_item: Received wrong item
  - not_as_described: Item not as described
  - no_longer_needed: No longer needed
  - better_price: Found better price elsewhere
  - other: Other reason

The --confirm flag is required to actually create the return.
Without it, the command will show what would be returned (dry run).`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		orderID := args[0]
		itemID := args[1]

		// Validate reason is provided
		if returnReason == "" {
			return outputError(models.NewCLIError(
				models.ErrCodeInvalidInput,
				"--reason flag is required",
				map[string]interface{}{
					"valid_reasons": models.ValidReturnReasons,
				},
			))
		}

		// Validate reason is valid
		if !models.IsValidReturnReason(returnReason) {
			return outputError(models.NewCLIError(
				models.ErrCodeInvalidInput,
				fmt.Sprintf("Invalid return reason: %s", returnReason),
				map[string]interface{}{
					"valid_reasons": models.ValidReturnReasons,
				},
			))
		}

		client := amazon.NewClient()

		// Dry run mode - show what would be returned
		if !returnConfirm {
			response := models.ReturnCreateResponse{
				DryRun: true,
				WouldReturn: &models.Return{
					OrderID: orderID,
					ItemID:  itemID,
					Reason:  returnReason,
					Status:  "pending_confirmation",
				},
				Message: "Add --confirm flag to execute this return",
			}
			return outputJSON(response)
		}

		// Actually create the return
		returnObj, err := client.CreateReturn(orderID, itemID, returnReason)
		if err != nil {
			return outputError(err)
		}

		response := models.ReturnCreateResponse{
			Return: returnObj,
		}
		return outputJSON(response)
	},
}

// returnsLabelCmd gets a return label
var returnsLabelCmd = &cobra.Command{
	Use:   "label <return-id>",
	Short: "Get return shipping label",
	Long:  `Get the shipping label for an initiated return.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		returnID := args[0]

		client := amazon.NewClient()
		label, err := client.GetReturnLabel(returnID)
		if err != nil {
			return outputError(err)
		}

		return outputJSON(label)
	},
}

// returnsStatusCmd checks return status
var returnsStatusCmd = &cobra.Command{
	Use:   "status <return-id>",
	Short: "Check return status",
	Long:  `Check the current status of a return.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		returnID := args[0]

		client := amazon.NewClient()
		returnObj, err := client.GetReturnStatus(returnID)
		if err != nil {
			return outputError(err)
		}

		return outputJSON(returnObj)
	},
}

func init() {
	rootCmd.AddCommand(returnsCmd)

	// Add subcommands
	returnsCmd.AddCommand(returnsListCmd)
	returnsCmd.AddCommand(returnsOptionsCmd)
	returnsCmd.AddCommand(returnsCreateCmd)
	returnsCmd.AddCommand(returnsLabelCmd)
	returnsCmd.AddCommand(returnsStatusCmd)

	// Flags for create command
	returnsCreateCmd.Flags().StringVarP(&returnReason, "reason", "r", "", "Reason for return (required)")
	returnsCreateCmd.Flags().BoolVar(&returnConfirm, "confirm", false, "Confirm the return creation")
}

// outputJSON outputs data as JSON
func outputJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// outputError outputs an error in JSON format
func outputError(err error) error {
	if cliErr, ok := err.(*models.CLIError); ok {
		response := models.ErrorResponse{
			Error: cliErr,
		}
		encoder := json.NewEncoder(os.Stderr)
		encoder.SetIndent("", "  ")
		encoder.Encode(response)
		return cliErr
	}

	// Convert regular errors to CLIError
	cliErr := models.NewCLIError(
		models.ErrCodeAmazonError,
		err.Error(),
		nil,
	)
	response := models.ErrorResponse{
		Error: cliErr,
	}
	encoder := json.NewEncoder(os.Stderr)
	encoder.SetIndent("", "  ")
	encoder.Encode(response)
	return cliErr
}
