package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/michaelshimeles/amazon-cli/internal/amazon"
	"github.com/spf13/cobra"
)

// cartCmd represents the cart parent command
var cartCmd = &cobra.Command{
	Use:   "cart",
	Short: "Manage shopping cart",
	Long: `Manage your Amazon shopping cart including:
- View cart contents and totals
- Add items to cart
- Remove items from cart
- Clear cart
- Checkout`,
}

// cartListCmd represents the cart list command
var cartListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all items in the cart",
	Long: `Retrieve and display all items currently in your Amazon shopping cart.

This command returns the complete cart contents including:
- All items with ASIN, title, price, and quantity
- Item subtotals
- Cart subtotal
- Estimated tax
- Total amount
- Item count

The output is in JSON format for easy parsing by AI agents.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Create Amazon client
		client := amazon.NewClient()

		// Get cart contents
		cart, err := client.GetCart()
		if err != nil {
			return fmt.Errorf("failed to get cart: %w", err)
		}

		// Output cart as JSON
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(cart); err != nil {
			return fmt.Errorf("failed to encode cart to JSON: %w", err)
		}

		return nil
	},
}

func init() {
	// Add cart command to root
	rootCmd.AddCommand(cartCmd)

	// Add cart list subcommand to cart
	cartCmd.AddCommand(cartListCmd)
}
