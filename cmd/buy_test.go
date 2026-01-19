package cmd

import (
	"testing"
)

func TestBuyCmd_Configuration(t *testing.T) {
	// Test that the buy command exists and has correct configuration
	if buyCmd.Use != "buy <asin>" {
		t.Errorf("Expected Use='buy <asin>', got '%s'", buyCmd.Use)
	}

	if buyCmd.Short != "Quick purchase an item" {
		t.Errorf("Expected Short='Quick purchase an item', got '%s'", buyCmd.Short)
	}

	expectedLong := `Quickly purchase an item by ASIN without adding to cart first.
Requires --confirm flag to execute the purchase.
Without --confirm, shows a preview of what would be purchased.`

	if buyCmd.Long != expectedLong {
		t.Errorf("Expected Long to mention preview mode and confirm flag")
	}

	if buyCmd.Run == nil {
		t.Error("Expected Run function to be defined")
	}

	// Test that it requires exactly 1 argument
	if buyCmd.Args == nil {
		t.Error("Expected Args validator to be defined")
	}
}

func TestBuyCmd_Flags(t *testing.T) {
	// Test that flags are properly configured
	quantityFlag := buyCmd.Flags().Lookup("quantity")
	if quantityFlag == nil {
		t.Error("Expected --quantity flag to be defined")
	} else {
		if quantityFlag.DefValue != "1" {
			t.Errorf("Expected --quantity default value to be '1', got '%s'", quantityFlag.DefValue)
		}
		if quantityFlag.Shorthand != "n" {
			t.Errorf("Expected --quantity shorthand to be 'n', got '%s'", quantityFlag.Shorthand)
		}
	}

	confirmFlag := buyCmd.Flags().Lookup("confirm")
	if confirmFlag == nil {
		t.Error("Expected --confirm flag to be defined")
	} else {
		if confirmFlag.DefValue != "false" {
			t.Errorf("Expected --confirm default value to be 'false', got '%s'", confirmFlag.DefValue)
		}
	}

	addressIDFlag := buyCmd.Flags().Lookup("address-id")
	if addressIDFlag == nil {
		t.Error("Expected --address-id flag to be defined")
	} else {
		if addressIDFlag.DefValue != "" {
			t.Errorf("Expected --address-id default value to be empty, got '%s'", addressIDFlag.DefValue)
		}
	}

	paymentIDFlag := buyCmd.Flags().Lookup("payment-id")
	if paymentIDFlag == nil {
		t.Error("Expected --payment-id flag to be defined")
	} else {
		if paymentIDFlag.DefValue != "" {
			t.Errorf("Expected --payment-id default value to be empty, got '%s'", paymentIDFlag.DefValue)
		}
	}
}

func TestBuyCmd_ConfirmFlagDefaultIsFalse(t *testing.T) {
	// Test that the confirm flag defaults to false (preview mode)
	confirmFlag := buyCmd.Flags().Lookup("confirm")
	if confirmFlag == nil {
		t.Fatal("Expected --confirm flag to be defined")
	}

	if confirmFlag.DefValue != "false" {
		t.Errorf("Expected --confirm flag to default to false (preview mode), got '%s'", confirmFlag.DefValue)
	}
}

func TestBuyCmd_PreviewModeIsDefault(t *testing.T) {
	// Verify that the command documentation mentions preview as default
	expectedLong := `Quickly purchase an item by ASIN without adding to cart first.
Requires --confirm flag to execute the purchase.
Without --confirm, shows a preview of what would be purchased.`

	if buyCmd.Long != expectedLong {
		t.Errorf("Expected Long description to mention preview mode as default")
	}
}

func TestBuyCmd_HasQuantityFlag(t *testing.T) {
	// Test that the buy command has a quantity flag
	quantityFlag := buyCmd.Flags().Lookup("quantity")
	if quantityFlag == nil {
		t.Error("Expected --quantity/-n flag to be defined")
		return
	}

	if quantityFlag.DefValue != "1" {
		t.Errorf("Expected default quantity to be 1, got '%s'", quantityFlag.DefValue)
	}
}

func TestBuyCmd_HasOptionalAddressAndPaymentFlags(t *testing.T) {
	// Test that the buy command has optional address-id and payment-id flags
	addressIDFlag := buyCmd.Flags().Lookup("address-id")
	if addressIDFlag == nil {
		t.Error("Expected --address-id flag to be defined")
		return
	}

	paymentIDFlag := buyCmd.Flags().Lookup("payment-id")
	if paymentIDFlag == nil {
		t.Error("Expected --payment-id flag to be defined")
		return
	}

	// Both should be optional (empty default)
	if addressIDFlag.DefValue != "" {
		t.Errorf("Expected --address-id to be optional (empty default), got '%s'", addressIDFlag.DefValue)
	}

	if paymentIDFlag.DefValue != "" {
		t.Errorf("Expected --payment-id to be optional (empty default), got '%s'", paymentIDFlag.DefValue)
	}
}

func TestBuyCmd_RequiresExactlyOneArgument(t *testing.T) {
	// Test that the command requires exactly one argument (ASIN)
	if buyCmd.Args == nil {
		t.Error("Expected Args validator to be defined")
	}

	// The use case specifies "<asin>" as required argument
	if buyCmd.Use != "buy <asin>" {
		t.Errorf("Expected Use to show required ASIN argument, got '%s'", buyCmd.Use)
	}
}
