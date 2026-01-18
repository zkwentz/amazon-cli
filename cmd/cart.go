package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/michaelshimeles/amazon-cli/internal/amazon"
	"github.com/spf13/cobra"
)

var (
	quantity int
)

// cartCmd represents the cart command
var cartCmd = &cobra.Command{
	Use:   "cart",
	Short: "Manage shopping cart",
	Long: `Cart commands allow you to manage your Amazon shopping cart.
You can add items, remove items, view cart contents, clear the cart, and checkout.`,
}

// cartAddCmd represents the cart add command
var cartAddCmd = &cobra.Command{
	Use:   "add <asin>",
	Short: "Add an item to the cart",
	Long:  `Add an item to the cart by specifying its ASIN and optionally the quantity.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		asin := args[0]

		client := amazon.NewClient()
		cart, err := client.AddToCart(asin, quantity)
		if err != nil {
			return fmt.Errorf("failed to add item to cart: %w", err)
		}

		output, err := json.MarshalIndent(cart, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal cart: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// cartListCmd represents the cart list command
var cartListCmd = &cobra.Command{
	Use:   "list",
	Short: "View cart contents",
	Long:  `Display all items currently in your shopping cart with prices and totals.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client := amazon.NewClient()
		cart, err := client.GetCart()
		if err != nil {
			return fmt.Errorf("failed to get cart: %w", err)
		}

		output, err := json.MarshalIndent(cart, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal cart: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// cartRemoveCmd represents the cart remove command
var cartRemoveCmd = &cobra.Command{
	Use:   "remove <asin>",
	Short: "Remove an item from the cart",
	Long:  `Remove an item from the cart by specifying its ASIN.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		asin := args[0]

		client := amazon.NewClient()
		cart, err := client.RemoveFromCart(asin)
		if err != nil {
			return fmt.Errorf("failed to remove item from cart: %w", err)
		}

		output, err := json.MarshalIndent(cart, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal cart: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

// cartClearCmd represents the cart clear command
var cartClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all items from the cart",
	Long:  `Remove all items from your shopping cart. Requires --confirm flag to execute.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		confirm, _ := cmd.Flags().GetBool("confirm")

		if !confirm {
			dryRun := map[string]interface{}{
				"dry_run": true,
				"message": "Add --confirm to execute this command",
				"action":  "clear cart",
			}
			output, err := json.MarshalIndent(dryRun, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal response: %w", err)
			}
			fmt.Println(string(output))
			return nil
		}

		client := amazon.NewClient()
		err := client.ClearCart()
		if err != nil {
			return fmt.Errorf("failed to clear cart: %w", err)
		}

		result := map[string]interface{}{
			"status":  "success",
			"message": "Cart cleared successfully",
		}
		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal response: %w", err)
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

	// Flags for cart add
	cartAddCmd.Flags().IntVar(&quantity, "quantity", 1, "Quantity to add")

	// Flags for cart clear
	cartClearCmd.Flags().Bool("confirm", false, "Confirm clearing the cart")
}
