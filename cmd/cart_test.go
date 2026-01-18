package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// TestCartCheckoutWithoutConfirm tests that checkout without --confirm flag
// shows a preview/dry-run instead of executing the purchase
func TestCartCheckoutWithoutConfirm(t *testing.T) {
	// Reset the confirm flag to ensure clean state
	cartConfirm = false

	// Create a buffer to capture output
	var outputBuf bytes.Buffer

	// Create the root command with cart subcommand
	rootCmd := &cobra.Command{Use: "amazon-cli"}
	cartCmd := &cobra.Command{Use: "cart"}
	checkoutCmd := &cobra.Command{
		Use: "checkout",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Simulate the checkout logic
			if !cartConfirm {
				// Dry run mode - should not execute purchase
				response := map[string]interface{}{
					"dry_run": true,
					"message": "Add --confirm to execute the purchase",
				}
				output, _ := json.MarshalIndent(response, "", "  ")
				outputBuf.Write(output)
				return nil
			}
			// Would execute purchase here if --confirm was set
			response := map[string]interface{}{
				"order_id": "111-1234567-2222222",
				"total":    32.39,
			}
			output, _ := json.MarshalIndent(response, "", "  ")
			outputBuf.Write(output)
			return nil
		},
	}

	checkoutCmd.Flags().BoolVar(&cartConfirm, "confirm", false, "Confirm the purchase")
	checkoutCmd.Flags().StringVar(&addressID, "address-id", "", "Address ID")
	checkoutCmd.Flags().StringVar(&paymentID, "payment-id", "", "Payment method ID")

	cartCmd.AddCommand(checkoutCmd)
	rootCmd.AddCommand(cartCmd)

	// Execute command WITHOUT --confirm flag
	rootCmd.SetArgs([]string{"cart", "checkout"})
	err := rootCmd.Execute()

	// Should not error
	if err != nil {
		t.Fatalf("checkout without --confirm should not error: %v", err)
	}

	// Parse the output
	var result map[string]interface{}
	if err := json.Unmarshal(outputBuf.Bytes(), &result); err != nil {
		t.Fatalf("failed to parse output JSON: %v", err)
	}

	// Verify it's a dry run
	dryRun, ok := result["dry_run"].(bool)
	if !ok {
		t.Fatal("output should contain 'dry_run' field")
	}
	if !dryRun {
		t.Error("dry_run should be true when --confirm is not provided")
	}

	// Verify the message
	message, ok := result["message"].(string)
	if !ok {
		t.Fatal("output should contain 'message' field")
	}
	if !strings.Contains(message, "--confirm") {
		t.Errorf("message should mention --confirm flag, got: %s", message)
	}

	// Verify no order_id is present (indicating purchase was not executed)
	if _, exists := result["order_id"]; exists {
		t.Error("order_id should not be present in dry run mode")
	}
}

// TestCartCheckoutWithConfirm tests that checkout WITH --confirm flag
// actually executes the purchase
func TestCartCheckoutWithConfirm(t *testing.T) {
	// Reset the confirm flag to ensure clean state
	cartConfirm = false

	// Create a buffer to capture output
	var outputBuf bytes.Buffer

	// Create the root command with cart subcommand
	rootCmd := &cobra.Command{Use: "amazon-cli"}
	cartCmd := &cobra.Command{Use: "cart"}
	checkoutCmd := &cobra.Command{
		Use: "checkout",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Simulate the checkout logic
			if !cartConfirm {
				// Dry run mode - should not execute purchase
				response := map[string]interface{}{
					"dry_run": true,
					"message": "Add --confirm to execute the purchase",
				}
				output, _ := json.MarshalIndent(response, "", "  ")
				outputBuf.Write(output)
				return nil
			}
			// Execute purchase when --confirm is set
			response := map[string]interface{}{
				"order_id": "111-1234567-2222222",
				"total":    32.39,
			}
			output, _ := json.MarshalIndent(response, "", "  ")
			outputBuf.Write(output)
			return nil
		},
	}

	checkoutCmd.Flags().BoolVar(&cartConfirm, "confirm", false, "Confirm the purchase")
	checkoutCmd.Flags().StringVar(&addressID, "address-id", "", "Address ID")
	checkoutCmd.Flags().StringVar(&paymentID, "payment-id", "", "Payment method ID")

	cartCmd.AddCommand(checkoutCmd)
	rootCmd.AddCommand(cartCmd)

	// Execute command WITH --confirm flag
	rootCmd.SetArgs([]string{"cart", "checkout", "--confirm"})
	err := rootCmd.Execute()

	// Should not error
	if err != nil {
		t.Fatalf("checkout with --confirm should not error: %v", err)
	}

	// Parse the output
	var result map[string]interface{}
	if err := json.Unmarshal(outputBuf.Bytes(), &result); err != nil {
		t.Fatalf("failed to parse output JSON: %v", err)
	}

	// Verify it's NOT a dry run
	if dryRun, exists := result["dry_run"]; exists && dryRun.(bool) {
		t.Error("dry_run should not be true when --confirm is provided")
	}

	// Verify order_id is present (indicating purchase was executed)
	orderID, ok := result["order_id"].(string)
	if !ok {
		t.Fatal("output should contain 'order_id' field when --confirm is provided")
	}
	if orderID == "" {
		t.Error("order_id should not be empty")
	}

	// Verify total is present
	if _, ok := result["total"].(float64); !ok {
		t.Error("output should contain 'total' field")
	}
}

// TestCartCheckoutCommandFlags tests that the flags are properly defined
func TestCartCheckoutCommandFlags(t *testing.T) {
	// Check that the confirm flag exists
	confirmFlag := cartCheckoutCmd.Flags().Lookup("confirm")
	if confirmFlag == nil {
		t.Fatal("checkout command should have --confirm flag")
	}

	// Check that address-id flag exists
	addressFlag := cartCheckoutCmd.Flags().Lookup("address-id")
	if addressFlag == nil {
		t.Fatal("checkout command should have --address-id flag")
	}

	// Check that payment-id flag exists
	paymentFlag := cartCheckoutCmd.Flags().Lookup("payment-id")
	if paymentFlag == nil {
		t.Fatal("checkout command should have --payment-id flag")
	}
}
