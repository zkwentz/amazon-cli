package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/michaelshimeles/amazon-cli/internal/amazon"
	"github.com/spf13/cobra"
)

var (
	quantity  int
	confirm   bool
	addressID string
	paymentID string
)

// cartCmd represents the cart parent command
var cartCmd = &cobra.Command{
	Use:   "cart",
	Short: "Manage shopping cart",
	Long: `Manage your Amazon shopping cart including adding items, viewing contents,
removing items, clearing the cart, and checking out.

Examples:
  amazon-cli cart add B08N5WRWNW --quantity 2
  amazon-cli cart list
  amazon-cli cart remove B08N5WRWNW
  amazon-cli cart clear --confirm
  amazon-cli cart checkout --confirm`,
}

// cartAddCmd adds an item to the cart
var cartAddCmd = &cobra.Command{
	Use:   "add <asin>",
	Short: "Add an item to the cart",
	Long: `Add an item to the cart by ASIN.

Example:
  amazon-cli cart add B08N5WRWNW --quantity 2`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		asin := args[0]

		// Validate ASIN format
		if !isValidASIN(asin) {
			return fmt.Errorf("invalid ASIN format: %s (must be 10 alphanumeric characters)", asin)
		}

		// Validate quantity
		if quantity <= 0 {
			return fmt.Errorf("quantity must be positive")
		}

		// Create Amazon client
		client := amazon.NewClient()

		// Add to cart
		cart, err := client.AddToCart(asin, quantity)
		if err != nil {
			return fmt.Errorf("failed to add to cart: %w", err)
		}

		// Output result as JSON
		return outputJSON(cart)
	},
}

// cartListCmd displays cart contents
var cartListCmd = &cobra.Command{
	Use:   "list",
	Short: "View cart contents",
	Long: `View all items in your shopping cart with prices and totals.

Example:
  amazon-cli cart list`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create Amazon client
		client := amazon.NewClient()

		// Get cart
		cart, err := client.GetCart()
		if err != nil {
			return fmt.Errorf("failed to get cart: %w", err)
		}

		// Output result as JSON
		return outputJSON(cart)
	},
}

// cartRemoveCmd removes an item from the cart
var cartRemoveCmd = &cobra.Command{
	Use:   "remove <asin>",
	Short: "Remove an item from the cart",
	Long: `Remove an item from the cart by ASIN.

Example:
  amazon-cli cart remove B08N5WRWNW`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		asin := args[0]

		// Validate ASIN format
		if !isValidASIN(asin) {
			return fmt.Errorf("invalid ASIN format: %s (must be 10 alphanumeric characters)", asin)
		}

		// Create Amazon client
		client := amazon.NewClient()

		// Remove from cart
		cart, err := client.RemoveFromCart(asin)
		if err != nil {
			return fmt.Errorf("failed to remove from cart: %w", err)
		}

		// Output result as JSON
		return outputJSON(cart)
	},
}

// cartClearCmd clears all items from the cart
var cartClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all items from the cart",
	Long: `Clear all items from the cart. Requires --confirm flag to execute.

Without --confirm, shows a preview of what would be cleared.

Examples:
  amazon-cli cart clear              # Preview mode
  amazon-cli cart clear --confirm    # Actually clear the cart`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create Amazon client
		client := amazon.NewClient()

		// Get current cart to show what would be cleared
		cart, err := client.GetCart()
		if err != nil {
			return fmt.Errorf("failed to get cart: %w", err)
		}

		// If --confirm not provided, show dry run
		if !confirm {
			result := map[string]interface{}{
				"dry_run":      true,
				"cart":         cart,
				"message":      "Add --confirm to execute",
				"items_count":  cart.ItemCount,
				"total_amount": cart.Total,
			}
			return outputJSON(result)
		}

		// Actually clear the cart
		err = client.ClearCart()
		if err != nil {
			return fmt.Errorf("failed to clear cart: %w", err)
		}

		// Output success message
		result := map[string]interface{}{
			"status":  "success",
			"message": "Cart cleared successfully",
		}
		return outputJSON(result)
	},
}

// cartCheckoutCmd completes checkout
var cartCheckoutCmd = &cobra.Command{
	Use:   "checkout",
	Short: "Complete checkout and place order",
	Long: `Complete checkout and place an order. Requires --confirm flag to execute.

Without --confirm, shows a preview of what would be purchased.

Examples:
  amazon-cli cart checkout                                    # Preview mode
  amazon-cli cart checkout --confirm                          # Use default address/payment
  amazon-cli cart checkout --confirm --address-id addr_123    # Specify address
  amazon-cli cart checkout --confirm --payment-id pay_456     # Specify payment method`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create Amazon client
		client := amazon.NewClient()

		// Use defaults if not provided
		if addressID == "" {
			addressID = "addr_default"
		}
		if paymentID == "" {
			paymentID = "pay_default"
		}

		// If --confirm not provided, show preview
		if !confirm {
			preview, err := client.PreviewCheckout(addressID, paymentID)
			if err != nil {
				return fmt.Errorf("failed to preview checkout: %w", err)
			}

			result := map[string]interface{}{
				"dry_run": true,
				"preview": preview,
				"message": "Add --confirm to execute",
			}
			return outputJSON(result)
		}

		// Actually complete checkout
		confirmation, err := client.CompleteCheckout(addressID, paymentID)
		if err != nil {
			return fmt.Errorf("failed to complete checkout: %w", err)
		}

		// Output order confirmation
		return outputJSON(confirmation)
	},
}

// isValidASIN validates ASIN format (10 alphanumeric characters)
func isValidASIN(asin string) bool {
	match, _ := regexp.MatchString(`^[A-Z0-9]{10}$`, asin)
	return match
}

// outputJSON outputs data as JSON
func outputJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func init() {
	// Register cart command with root
	rootCmd.AddCommand(cartCmd)

	// Register subcommands
	cartCmd.AddCommand(cartAddCmd)
	cartCmd.AddCommand(cartListCmd)
	cartCmd.AddCommand(cartRemoveCmd)
	cartCmd.AddCommand(cartClearCmd)
	cartCmd.AddCommand(cartCheckoutCmd)

	// Add flags
	cartAddCmd.Flags().IntVarP(&quantity, "quantity", "n", 1, "Quantity to add")
	cartClearCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm the action")
	cartCheckoutCmd.Flags().BoolVar(&confirm, "confirm", false, "Confirm the purchase")
	cartCheckoutCmd.Flags().StringVar(&addressID, "address-id", "", "Address ID to use")
	cartCheckoutCmd.Flags().StringVar(&paymentID, "payment-id", "", "Payment method ID to use")
}
