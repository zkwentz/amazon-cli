package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/zkwentz/amazon-cli/internal/output"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

var (
	returnsReason  string
	returnsConfirm bool
)

// returnsCmd represents the returns command
var returnsCmd = &cobra.Command{
	Use:   "returns",
	Short: "Manage returns",
	Long:  `List returnable items, get return options, and create returns.`,
}

// returnsCreateCmd represents the returns create command
var returnsCreateCmd = &cobra.Command{
	Use:   "create <order-id> <item-id>",
	Short: "Create a return",
	Long: `Create a return request for an order item.
Requires --reason flag with one of: defective, wrong_item, not_as_described, no_longer_needed, better_price, other.
Without --confirm, shows a preview of the return. With --confirm, submits the return.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		orderID := args[0]
		itemID := args[1]

		// Validate orderID is not empty
		if orderID == "" {
			output.Error(models.ErrInvalidInput, "order ID cannot be empty", nil)
			os.Exit(models.ExitInvalidArgs)
		}

		// Validate itemID is not empty
		if itemID == "" {
			output.Error(models.ErrInvalidInput, "item ID cannot be empty", nil)
			os.Exit(models.ExitInvalidArgs)
		}

		// Validate reason is provided
		if returnsReason == "" {
			output.Error(models.ErrInvalidInput, "reason is required (use --reason flag with: defective, wrong_item, not_as_described, no_longer_needed, better_price, other)", nil)
			os.Exit(models.ExitInvalidArgs)
		}

		c := getClient()

		if !returnsConfirm {
			// Dry run - show preview
			output.JSON(map[string]interface{}{
				"dry_run":  true,
				"order_id": orderID,
				"item_id":  itemID,
				"reason":   returnsReason,
				"message":  "Add --confirm to submit the return",
			})
			return
		}

		// Execute return creation
		ret, err := c.CreateReturn(orderID, itemID, returnsReason)
		if err != nil {
			output.Error(models.ErrInvalidInput, err.Error(), nil)
			os.Exit(models.ExitInvalidArgs)
		}

		output.JSON(ret)
	},
}

// returnsLabelCmd represents the returns label command
var returnsLabelCmd = &cobra.Command{
	Use:   "label <return-id>",
	Short: "Get return label",
	Long:  `Retrieve the shipping label for a return.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		returnID := args[0]

		// Validate returnID is not empty
		if returnID == "" {
			output.Error(models.ErrInvalidInput, "return ID cannot be empty", nil)
			os.Exit(models.ExitInvalidArgs)
		}

		c := getClient()

		// Get return label
		label, err := c.GetReturnLabel(returnID)
		if err != nil {
			output.Error(models.ErrInvalidInput, err.Error(), nil)
			os.Exit(models.ExitInvalidArgs)
		}

		output.JSON(label)
	},
}

// returnsStatusCmd represents the returns status command
var returnsStatusCmd = &cobra.Command{
	Use:   "status <return-id>",
	Short: "Get return status",
	Long:  `Retrieve the current status of a return.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		returnID := args[0]

		// Validate returnID is not empty
		if returnID == "" {
			output.Error(models.ErrInvalidInput, "return ID cannot be empty", nil)
			os.Exit(models.ExitInvalidArgs)
		}

		c := getClient()

		// Get return status
		ret, err := c.GetReturnStatus(returnID)
		if err != nil {
			output.Error(models.ErrInvalidInput, err.Error(), nil)
			os.Exit(models.ExitInvalidArgs)
		}

		output.JSON(ret)
	},
}

func init() {
	rootCmd.AddCommand(returnsCmd)

	// Add subcommands
	returnsCmd.AddCommand(returnsCreateCmd)
	returnsCmd.AddCommand(returnsLabelCmd)
	returnsCmd.AddCommand(returnsStatusCmd)

	// Flags for returns create
	returnsCreateCmd.Flags().StringVar(&returnsReason, "reason", "", "Return reason (required): defective, wrong_item, not_as_described, no_longer_needed, better_price, other")
	returnsCreateCmd.Flags().BoolVar(&returnsConfirm, "confirm", false, "Confirm the return creation")
	returnsCreateCmd.MarkFlagRequired("reason")
}
