package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/michaelshimeles/amazon-cli/internal/amazon"
	"github.com/spf13/cobra"
)

var (
	confirmFlag bool
)

// cartCmd represents the cart command
var cartCmd = &cobra.Command{
	Use:   "cart",
	Short: "Manage shopping cart",
	Long: `Manage your Amazon shopping cart with subcommands for adding, removing,
viewing, and clearing items, as well as completing checkout.`,
}

// clearCmd represents the cart clear command
var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all items from the cart",
	Long: `Remove all items from your Amazon shopping cart.
This command requires the --confirm flag to execute. Without --confirm,
it will display what would be cleared but not actually clear the cart.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := amazon.NewClient()

		// Get current cart to show what will be cleared
		cart, err := client.GetCart()
		if err != nil {
			return fmt.Errorf("failed to get cart: %w", err)
		}

		// If --confirm flag is not provided, show dry run
		if !confirmFlag {
			output := map[string]interface{}{
				"dry_run": true,
				"cart":    cart,
				"message": "Add --confirm to execute cart clear",
			}

			jsonOutput, err := json.MarshalIndent(output, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal output: %w", err)
			}

			fmt.Println(string(jsonOutput))
			return nil
		}

		// Execute cart clear with --confirm flag
		err = client.ClearCart()
		if err != nil {
			return fmt.Errorf("failed to clear cart: %w", err)
		}

		// Get updated cart (should be empty)
		clearedCart, err := client.GetCart()
		if err != nil {
			return fmt.Errorf("failed to get cart after clearing: %w", err)
		}

		output := map[string]interface{}{
			"success": true,
			"message": "Cart cleared successfully",
			"cart":    clearedCart,
		}

		jsonOutput, err := json.MarshalIndent(output, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal output: %w", err)
		}

		fmt.Println(string(jsonOutput))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(cartCmd)
	cartCmd.AddCommand(clearCmd)

	// Add --confirm flag to clear command
	clearCmd.Flags().BoolVar(&confirmFlag, "confirm", false, "Confirm cart clear operation")
}
