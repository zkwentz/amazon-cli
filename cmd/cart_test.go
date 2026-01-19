package cmd

import (
	"testing"
)

func TestCartCheckoutCmd_Configuration(t *testing.T) {
	// Test that the checkout command exists and has correct configuration
	if cartCheckoutCmd.Use != "checkout" {
		t.Errorf("Expected Use='checkout', got '%s'", cartCheckoutCmd.Use)
	}

	if cartCheckoutCmd.Short != "Checkout cart" {
		t.Errorf("Expected Short='Checkout cart', got '%s'", cartCheckoutCmd.Short)
	}

	if cartCheckoutCmd.Run == nil {
		t.Error("Expected Run function to be defined")
	}
}

func TestCartCheckoutCmd_Flags(t *testing.T) {
	// Test that flags are properly configured
	confirmFlag := cartCheckoutCmd.Flags().Lookup("confirm")
	if confirmFlag == nil {
		t.Error("Expected --confirm flag to be defined")
	} else {
		if confirmFlag.DefValue != "false" {
			t.Errorf("Expected --confirm default value to be 'false', got '%s'", confirmFlag.DefValue)
		}
	}

	addressIDFlag := cartCheckoutCmd.Flags().Lookup("address-id")
	if addressIDFlag == nil {
		t.Error("Expected --address-id flag to be defined")
	} else {
		if addressIDFlag.DefValue != "" {
			t.Errorf("Expected --address-id default value to be empty, got '%s'", addressIDFlag.DefValue)
		}
	}

	paymentIDFlag := cartCheckoutCmd.Flags().Lookup("payment-id")
	if paymentIDFlag == nil {
		t.Error("Expected --payment-id flag to be defined")
	} else {
		if paymentIDFlag.DefValue != "" {
			t.Errorf("Expected --payment-id default value to be empty, got '%s'", paymentIDFlag.DefValue)
		}
	}
}

func TestCartCheckoutCmd_ConfirmFlagDefaultIsFalse(t *testing.T) {
	// Test that the confirm flag defaults to false (preview mode)
	confirmFlag := cartCheckoutCmd.Flags().Lookup("confirm")
	if confirmFlag == nil {
		t.Fatal("Expected --confirm flag to be defined")
	}

	if confirmFlag.DefValue != "false" {
		t.Errorf("Expected --confirm flag to default to false (preview mode), got '%s'", confirmFlag.DefValue)
	}
}

func TestCartCheckoutCmd_PreviewModeIsDefault(t *testing.T) {
	// Verify that the command documentation mentions preview as default
	expectedLong := `Complete purchase of items in cart.
Requires --confirm flag to execute the purchase.
Without --confirm, shows a preview of the order.`

	if cartCheckoutCmd.Long != expectedLong {
		t.Errorf("Expected Long description to mention preview mode as default")
	}
}

func TestCartClearCmd_Configuration(t *testing.T) {
	// Test that the clear command exists and has correct configuration
	if cartClearCmd.Use != "clear" {
		t.Errorf("Expected Use='clear', got '%s'", cartClearCmd.Use)
	}

	if cartClearCmd.Short != "Clear all items from cart" {
		t.Errorf("Expected Short='Clear all items from cart', got '%s'", cartClearCmd.Short)
	}

	if cartClearCmd.Run == nil {
		t.Error("Expected Run function to be defined")
	}
}

func TestCartClearCmd_Flags(t *testing.T) {
	// Test that the clear command has a confirm flag
	confirmFlag := cartClearCmd.Flags().Lookup("confirm")
	if confirmFlag == nil {
		t.Error("Expected --confirm flag to be defined")
	} else {
		if confirmFlag.DefValue != "false" {
			t.Errorf("Expected --confirm default value to be 'false', got '%s'", confirmFlag.DefValue)
		}
	}
}

func TestCartAddCmd_Configuration(t *testing.T) {
	// Test that the add command exists and has correct configuration
	if cartAddCmd.Use != "add <asin>" {
		t.Errorf("Expected Use='add <asin>', got '%s'", cartAddCmd.Use)
	}

	if cartAddCmd.Short != "Add item to cart" {
		t.Errorf("Expected Short='Add item to cart', got '%s'", cartAddCmd.Short)
	}

	if cartAddCmd.Run == nil {
		t.Error("Expected Run function to be defined")
	}

	// Test that it requires exactly 1 argument
	if cartAddCmd.Args == nil {
		t.Error("Expected Args validator to be defined")
	}
}

func TestCartAddCmd_Flags(t *testing.T) {
	// Test that the add command has a quantity flag
	quantityFlag := cartAddCmd.Flags().Lookup("quantity")
	if quantityFlag == nil {
		t.Error("Expected --quantity flag to be defined")
	} else {
		if quantityFlag.DefValue != "1" {
			t.Errorf("Expected --quantity default value to be '1', got '%s'", quantityFlag.DefValue)
		}
	}
}

func TestCartListCmd_Configuration(t *testing.T) {
	// Test that the list command exists and has correct configuration
	if cartListCmd.Use != "list" {
		t.Errorf("Expected Use='list', got '%s'", cartListCmd.Use)
	}

	if cartListCmd.Short != "View cart contents" {
		t.Errorf("Expected Short='View cart contents', got '%s'", cartListCmd.Short)
	}

	if cartListCmd.Run == nil {
		t.Error("Expected Run function to be defined")
	}
}

func TestCartRemoveCmd_Configuration(t *testing.T) {
	// Test that the remove command exists and has correct configuration
	if cartRemoveCmd.Use != "remove <asin>" {
		t.Errorf("Expected Use='remove <asin>', got '%s'", cartRemoveCmd.Use)
	}

	if cartRemoveCmd.Short != "Remove item from cart" {
		t.Errorf("Expected Short='Remove item from cart', got '%s'", cartRemoveCmd.Short)
	}

	if cartRemoveCmd.Run == nil {
		t.Error("Expected Run function to be defined")
	}

	// Test that it requires exactly 1 argument
	if cartRemoveCmd.Args == nil {
		t.Error("Expected Args validator to be defined")
	}
}

func TestCartCmd_Subcommands(t *testing.T) {
	// Test that all subcommands are registered
	expectedSubcommands := []string{"add", "list", "remove", "clear", "checkout"}
	commands := cartCmd.Commands()

	if len(commands) != len(expectedSubcommands) {
		t.Errorf("Expected %d subcommands, got %d", len(expectedSubcommands), len(commands))
	}

	// Check that each expected subcommand exists
	for _, expected := range expectedSubcommands {
		found := false
		for _, cmd := range commands {
			if cmd.Use == expected || cmd.Use == expected+" <asin>" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected subcommand '%s' not found", expected)
		}
	}
}

func TestCartCmd_Configuration(t *testing.T) {
	// Test the main cart command configuration
	if cartCmd.Use != "cart" {
		t.Errorf("Expected Use='cart', got '%s'", cartCmd.Use)
	}

	if cartCmd.Short != "Manage shopping cart" {
		t.Errorf("Expected Short='Manage shopping cart', got '%s'", cartCmd.Short)
	}

	expectedLong := `Add, remove, view, and checkout items in your Amazon shopping cart.`
	if cartCmd.Long != expectedLong {
		t.Errorf("Expected Long='%s', got '%s'", expectedLong, cartCmd.Long)
	}
}

func TestGetClient_ReturnsSameInstance(t *testing.T) {
	// Test that getClient returns the same instance across calls
	c1 := getClient()
	if c1 == nil {
		t.Error("Expected getClient() to return non-nil client")
	}

	c2 := getClient()
	if c2 == nil {
		t.Error("Expected getClient() to return non-nil client")
	}

	if c1 != c2 {
		t.Error("Expected getClient() to return the same client instance")
	}
}
