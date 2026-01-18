package cmd

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/michaelshimeles/amazon-cli/internal/amazon"
	"github.com/spf13/cobra"
)

var (
	quantity int
)

// cartCmd represents the cart parent command
var cartCmd = &cobra.Command{
	Use:   "cart",
	Short: "Manage shopping cart",
	Long: `Manage your Amazon shopping cart including adding items,
viewing cart contents, removing items, and checkout operations.`,
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

		// Validate ASIN format (10 alphanumeric characters)
		if !isValidASIN(asin) {
			return fmt.Errorf("invalid ASIN format: %s (must be 10 alphanumeric characters)", asin)
		}

		// Validate quantity
		if quantity <= 0 {
			return fmt.Errorf("quantity must be positive, got: %d", quantity)
		}

		// Create Amazon client
		client := amazon.NewClient()

		// Add item to cart
		cart, err := client.AddToCart(asin, quantity)
		if err != nil {
			return fmt.Errorf("failed to add item to cart: %w", err)
		}

		// Output cart as JSON
		output, err := json.MarshalIndent(cart, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to format output: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cartCmd)
	cartCmd.AddCommand(cartAddCmd)

	// Add flags to cart add command
	cartAddCmd.Flags().IntVarP(&quantity, "quantity", "n", 1, "Quantity of items to add")
}

// isValidASIN validates that an ASIN is in the correct format
// ASINs are 10 alphanumeric characters
func isValidASIN(asin string) bool {
	// ASIN must be exactly 10 characters and alphanumeric
	match, _ := regexp.MatchString(`^[A-Z0-9]{10}$`, asin)
	return match
}
