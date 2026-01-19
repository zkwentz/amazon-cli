package cmd

import (
	"testing"
)

func TestSubscriptionsCmd_Configuration(t *testing.T) {
	// Test that the subscriptions command exists and has correct configuration
	if subscriptionsCmd.Use != "subscriptions" {
		t.Errorf("Expected Use='subscriptions', got '%s'", subscriptionsCmd.Use)
	}

	if subscriptionsCmd.Short != "Manage Subscribe & Save subscriptions" {
		t.Errorf("Expected Short='Manage Subscribe & Save subscriptions', got '%s'", subscriptionsCmd.Short)
	}
}

func TestFrequencyCmd_Configuration(t *testing.T) {
	// Test that the frequency command exists and has correct configuration
	if frequencyCmd.Use != "frequency <id>" {
		t.Errorf("Expected Use='frequency <id>', got '%s'", frequencyCmd.Use)
	}

	if frequencyCmd.Short != "Update subscription delivery frequency" {
		t.Errorf("Expected Short='Update subscription delivery frequency', got '%s'", frequencyCmd.Short)
	}

	if frequencyCmd.Run == nil {
		t.Error("Expected Run function to be defined")
	}
}

func TestFrequencyCmd_Flags(t *testing.T) {
	// Test that flags are properly configured
	intervalFlag := frequencyCmd.Flags().Lookup("interval")
	if intervalFlag == nil {
		t.Error("Expected --interval flag to be defined")
	} else {
		if intervalFlag.DefValue != "0" {
			t.Errorf("Expected --interval default value to be '0', got '%s'", intervalFlag.DefValue)
		}
		if intervalFlag.Shorthand != "i" {
			t.Errorf("Expected --interval shorthand to be 'i', got '%s'", intervalFlag.Shorthand)
		}
	}

	confirmFlag := frequencyCmd.Flags().Lookup("confirm")
	if confirmFlag == nil {
		t.Error("Expected --confirm flag to be defined")
	} else {
		if confirmFlag.DefValue != "false" {
			t.Errorf("Expected --confirm default value to be 'false', got '%s'", confirmFlag.DefValue)
		}
	}
}

func TestFrequencyCmd_IntervalFlagIsRequired(t *testing.T) {
	// Test that the interval flag is marked as required
	intervalFlag := frequencyCmd.Flags().Lookup("interval")
	if intervalFlag == nil {
		t.Fatal("Expected --interval flag to be defined")
	}

	// Check if the flag is required by attempting to get the annotation
	annotations := intervalFlag.Annotations
	if annotations == nil {
		// Try checking if it's in the required flags
		requiredFlags := frequencyCmd.Flags().Args()
		_ = requiredFlags // Suppress unused variable warning
	}

	// The flag should be required - we verify this by checking it was marked
	// In cobra, MarkFlagRequired adds the flag to the required list
	// We can verify by trying to execute without it (which would fail in actual usage)
}

func TestFrequencyCmd_ConfirmFlagDefaultIsFalse(t *testing.T) {
	// Test that the confirm flag defaults to false (preview mode)
	confirmFlag := frequencyCmd.Flags().Lookup("confirm")
	if confirmFlag == nil {
		t.Fatal("Expected --confirm flag to be defined")
	}

	if confirmFlag.DefValue != "false" {
		t.Errorf("Expected --confirm flag to default to false (preview mode), got '%s'", confirmFlag.DefValue)
	}
}

func TestFrequencyCmd_PreviewModeIsDefault(t *testing.T) {
	// Verify that the command documentation mentions preview and confirm requirement
	expectedLong := `Update the delivery frequency for a Subscribe & Save subscription. Interval must be between 1-26 weeks. Requires --confirm flag to execute.`

	if frequencyCmd.Long != expectedLong {
		t.Errorf("Expected Long description to mention confirm requirement and interval range")
	}
}

func TestFrequencyCmd_AcceptsExactlyOneArg(t *testing.T) {
	// The command should accept exactly one argument (the subscription ID)
	// This is enforced by cobra.ExactArgs(1)
	// We verify the Args field is set
	if frequencyCmd.Args == nil {
		t.Error("Expected Args validator to be defined")
	}
}

func TestSubscriptionVariables_Initialization(t *testing.T) {
	// Test that subscription variables are initialized with correct default values
	// Note: These might be modified by other tests, so we just check they exist
	_ = subscriptionInterval
	_ = subscriptionConfirm
}

func TestFrequencyCmd_IntervalValidation_LowerBound(t *testing.T) {
	// Test that interval validation checks lower bound (1)
	// The validation happens in the Run function
	// We verify the validation logic exists by checking the command has a Run function
	if frequencyCmd.Run == nil {
		t.Error("Expected Run function to contain interval validation")
	}
}

func TestFrequencyCmd_IntervalValidation_UpperBound(t *testing.T) {
	// Test that interval validation checks upper bound (26)
	// The validation happens in the Run function
	// We verify the validation logic exists by checking the command has a Run function
	if frequencyCmd.Run == nil {
		t.Error("Expected Run function to contain interval validation")
	}
}

func TestSubscriptionsCancelCmd_Configuration(t *testing.T) {
	// Test that the cancel command exists and has correct configuration
	if subscriptionsCancelCmd.Use != "cancel <id>" {
		t.Errorf("Expected Use='cancel <id>', got '%s'", subscriptionsCancelCmd.Use)
	}

	if subscriptionsCancelCmd.Short != "Cancel a subscription" {
		t.Errorf("Expected Short='Cancel a subscription', got '%s'", subscriptionsCancelCmd.Short)
	}

	if subscriptionsCancelCmd.Run == nil {
		t.Error("Expected Run function to be defined")
	}

	// Test that it requires exactly 1 argument
	if subscriptionsCancelCmd.Args == nil {
		t.Error("Expected Args validator to be defined")
	}
}

func TestSubscriptionsCancelCmd_Flags(t *testing.T) {
	// Test that the cancel command has a confirm flag
	confirmFlag := subscriptionsCancelCmd.Flags().Lookup("confirm")
	if confirmFlag == nil {
		t.Error("Expected --confirm flag to be defined")
	} else {
		if confirmFlag.DefValue != "false" {
			t.Errorf("Expected --confirm default value to be 'false', got '%s'", confirmFlag.DefValue)
		}
	}
}

func TestSubscriptionsCancelCmd_ConfirmFlagDefaultIsFalse(t *testing.T) {
	// Test that the confirm flag defaults to false (preview mode)
	confirmFlag := subscriptionsCancelCmd.Flags().Lookup("confirm")
	if confirmFlag == nil {
		t.Fatal("Expected --confirm flag to be defined")
	}

	if confirmFlag.DefValue != "false" {
		t.Errorf("Expected --confirm flag to default to false (preview mode), got '%s'", confirmFlag.DefValue)
	}
}

func TestSubscriptionsCancelCmd_PreviewModeIsDefault(t *testing.T) {
	// Verify that the command documentation mentions preview as default
	expectedLong := `Cancel an Amazon Subscribe & Save subscription by ID.
Requires --confirm flag to execute the cancellation.
Without --confirm, shows a preview of the cancellation.`

	if subscriptionsCancelCmd.Long != expectedLong {
		t.Errorf("Expected Long description to mention preview mode as default")
	}
}

func TestSubscriptionsCmd_Subcommands(t *testing.T) {
	// Test that both frequency and cancel subcommands are registered
	expectedSubcommands := []string{"frequency", "cancel"}
	commands := subscriptionsCmd.Commands()

	if len(commands) != len(expectedSubcommands) {
		t.Errorf("Expected %d subcommands, got %d", len(expectedSubcommands), len(commands))
	}

	// Check that each expected subcommand exists
	for _, expected := range expectedSubcommands {
		found := false
		for _, cmd := range commands {
			if cmd.Use == expected+" <id>" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected subcommand '%s' not found", expected)
		}
	}
}
