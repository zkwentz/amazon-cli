package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/michaelshimeles/amazon-cli/internal/amazon"
	"github.com/spf13/cobra"
)

var (
	quantity  int
	confirm   bool
	addressID string
	paymentID string
)

// cartCmd represents the cart command
var cartCmd = &cobra.Command{
	Use:   "cart",
	Short: "Manage shopping cart",
	Long: `Manage your Amazon shopping cart including adding items, viewing contents,
removing items, clearing the cart, and checking out.

All cart operations output structured JSON for seamless integration.`,
}

// cartAddCmd represents the cart add command
var cartAddCmd = &cobra.Command{
	Use:   "add <asin>",
	Short: "Add an item to the cart",
	Long: `Add an item to your Amazon shopping cart by ASIN.

Example:
  amazon-cli cart add B08N5WRWNW --quantity 2`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		asin := args[0]
		client := amazon.NewClient()

		cart, err := client.AddToCart(asin, quantity)
		if err != nil {
			return fmt.Errorf("failed to add item to cart: %w", err)
		}

		output, err := json.MarshalIndent(cart, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// cartListCmd represents the cart list command
var cartListCmd = &cobra.Command{
	Use:   "list",
	Short: "View cart contents",
	Long: `View all items in your Amazon shopping cart with prices and totals.

Example:
  amazon-cli cart list`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := amazon.NewClient()

		cart, err := client.GetCart()
		if err != nil {
			return fmt.Errorf("failed to get cart: %w", err)
		}

		output, err := json.MarshalIndent(cart, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// cartRemoveCmd represents the cart remove command
var cartRemoveCmd = &cobra.Command{
	Use:   "remove <asin>",
	Short: "Remove an item from the cart",
	Long: `Remove an item from your Amazon shopping cart by ASIN.

Example:
  amazon-cli cart remove B08N5WRWNW`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		asin := args[0]
		client := amazon.NewClient()

		cart, err := client.RemoveFromCart(asin)
		if err != nil {
			return fmt.Errorf("failed to remove item from cart: %w", err)
		}

		output, err := json.MarshalIndent(cart, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// cartClearCmd represents the cart clear command
var cartClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all items from the cart",
	Long: `Remove all items from your Amazon shopping cart.

SAFETY: This command requires the --confirm flag to execute.

Examples:
  amazon-cli cart clear              # Dry run (shows what would happen)
  amazon-cli cart clear --confirm    # Actually clears the cart`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := amazon.NewClient()

		// Get current cart for dry run display
		cart, err := client.GetCart()
		if err != nil {
			return fmt.Errorf("failed to get cart: %w", err)
		}

		if !confirm {
			// Dry run - show what would be cleared
			dryRun := map[string]interface{}{
				"dry_run": true,
				"message": "Add --confirm to execute. This will clear all items from the cart.",
				"cart":    cart,
			}

			output, err := json.MarshalIndent(dryRun, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to format output: %w", err)
			}

			fmt.Println(string(output))
			return nil
		}

		// Actually clear the cart
		err = client.ClearCart()
		if err != nil {
			return fmt.Errorf("failed to clear cart: %w", err)
		}

		result := map[string]interface{}{
			"status":  "success",
			"message": "Cart cleared successfully",
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// cartCheckoutCmd represents the cart checkout command
var cartCheckoutCmd = &cobra.Command{
	Use:   "checkout",
	Short: "Complete checkout and place order",
	Long: `Complete the checkout process and place an order with items in your cart.

SAFETY: This command requires the --confirm flag to execute a purchase.
Without --confirm, it shows a preview of what would be purchased.

Examples:
  # Preview checkout without placing order
  amazon-cli cart checkout --address-id addr123 --payment-id pay456

  # Complete checkout and place order
  amazon-cli cart checkout --address-id addr123 --payment-id pay456 --confirm`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := amazon.NewClient()

		// Validate required parameters
		if addressID == "" {
			return fmt.Errorf("--address-id is required")
		}
		if paymentID == "" {
			return fmt.Errorf("--payment-id is required")
		}

		if !confirm {
			// Dry run - show checkout preview
			preview, err := client.PreviewCheckout(addressID, paymentID)
			if err != nil {
				return fmt.Errorf("failed to preview checkout: %w", err)
			}

			dryRun := map[string]interface{}{
				"dry_run": true,
				"message": "Add --confirm to complete the purchase",
				"preview": preview,
			}

			output, err := json.MarshalIndent(dryRun, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to format output: %w", err)
			}

			fmt.Println(string(output))
			return nil
		}

		// Actually complete the checkout
		confirmation, err := client.CompleteCheckout(addressID, paymentID)
		if err != nil {
			return fmt.Errorf("failed to complete checkout: %w", err)
		}

		output, err := json.MarshalIndent(confirmation, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cartCmd)

	// Add subcommands
	cartCmd.AddCommand(cartAddCmd)
	cartCmd.AddCommand(cartListCmd)
	cartCmd.AddCommand(cartRemoveCmd)
	cartCmd.AddCommand(cartClearCmd)
	cartCmd.AddCommand(cartCheckoutCmd)

	// cart add flags
	cartAddCmd.Flags().IntVarP(&quantity, "quantity", "n", 1, "Quantity of items to add")

	// cart clear flags
	cartClearCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm the action")

	// cart checkout flags
	cartCheckoutCmd.Flags().StringVar(&addressID, "address-id", "", "Address ID for shipping (required)")
	cartCheckoutCmd.Flags().StringVar(&paymentID, "payment-id", "", "Payment method ID (required)")
	cartCheckoutCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm the purchase")
}
