package cmd

import (
	"testing"

	"github.com/zkwentz/amazon-cli/pkg/models"
)

func TestReturnsCommandStructure(t *testing.T) {
	cmd := GetReturnsCmd()

	if cmd == nil {
		t.Fatal("GetReturnsCmd() returned nil")
	}

	if cmd.Use != "returns" {
		t.Errorf("Expected Use to be 'returns', got '%s'", cmd.Use)
	}

	// Check that all subcommands are registered
	expectedSubcommands := []string{"list", "options", "create", "label", "status"}
	commands := cmd.Commands()

	if len(commands) != len(expectedSubcommands) {
		t.Errorf("Expected %d subcommands, got %d", len(expectedSubcommands), len(commands))
	}

	for _, expectedName := range expectedSubcommands {
		found := false
		for _, cmd := range commands {
			if cmd.Name() == expectedName {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected subcommand '%s' not found", expectedName)
		}
	}
}

func TestListReturnsCommand(t *testing.T) {
	cmd := listReturnsCmd

	if cmd.Use != "list" {
		t.Errorf("Expected Use to be 'list', got '%s'", cmd.Use)
	}

	if cmd.Args != nil {
		// list should not require arguments
		err := cmd.Args(cmd, []string{})
		if err != nil {
			t.Errorf("list command should not require arguments")
		}
	}
}

func TestOptionsReturnsCommand(t *testing.T) {
	cmd := optionsReturnsCmd

	if cmd.Use != "options <order-id> <item-id>" {
		t.Errorf("Expected Use to be 'options <order-id> <item-id>', got '%s'", cmd.Use)
	}

	// Test that it requires exactly 2 arguments
	err := cmd.Args(cmd, []string{"order1"})
	if err == nil {
		t.Error("options command should require 2 arguments, but accepted 1")
	}

	err = cmd.Args(cmd, []string{"order1", "item1"})
	if err != nil {
		t.Errorf("options command should accept 2 arguments, got error: %v", err)
	}

	err = cmd.Args(cmd, []string{"order1", "item1", "extra"})
	if err == nil {
		t.Error("options command should require exactly 2 arguments, but accepted 3")
	}
}

func TestCreateReturnsCommand(t *testing.T) {
	cmd := createReturnsCmd

	if cmd.Use != "create <order-id> <item-id>" {
		t.Errorf("Expected Use to be 'create <order-id> <item-id>', got '%s'", cmd.Use)
	}

	// Test that it requires exactly 2 arguments
	err := cmd.Args(cmd, []string{"order1"})
	if err == nil {
		t.Error("create command should require 2 arguments, but accepted 1")
	}

	err = cmd.Args(cmd, []string{"order1", "item1"})
	if err != nil {
		t.Errorf("create command should accept 2 arguments, got error: %v", err)
	}

	// Check that flags are registered
	reasonFlag := cmd.Flags().Lookup("reason")
	if reasonFlag == nil {
		t.Error("create command should have --reason flag")
	}

	confirmFlag := cmd.Flags().Lookup("confirm")
	if confirmFlag == nil {
		t.Error("create command should have --confirm flag")
	}
}

func TestLabelReturnsCommand(t *testing.T) {
	cmd := labelReturnsCmd

	if cmd.Use != "label <return-id>" {
		t.Errorf("Expected Use to be 'label <return-id>', got '%s'", cmd.Use)
	}

	// Test that it requires exactly 1 argument
	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("label command should require 1 argument, but accepted 0")
	}

	err = cmd.Args(cmd, []string{"return1"})
	if err != nil {
		t.Errorf("label command should accept 1 argument, got error: %v", err)
	}

	err = cmd.Args(cmd, []string{"return1", "extra"})
	if err == nil {
		t.Error("label command should require exactly 1 argument, but accepted 2")
	}
}

func TestStatusReturnsCommand(t *testing.T) {
	cmd := statusReturnsCmd

	if cmd.Use != "status <return-id>" {
		t.Errorf("Expected Use to be 'status <return-id>', got '%s'", cmd.Use)
	}

	// Test that it requires exactly 1 argument
	err := cmd.Args(cmd, []string{})
	if err == nil {
		t.Error("status command should require 1 argument, but accepted 0")
	}

	err = cmd.Args(cmd, []string{"return1"})
	if err != nil {
		t.Errorf("status command should accept 1 argument, got error: %v", err)
	}

	err = cmd.Args(cmd, []string{"return1", "extra"})
	if err == nil {
		t.Error("status command should require exactly 1 argument, but accepted 2")
	}
}

func TestReturnReasonValidation(t *testing.T) {
	validReasons := []string{
		"defective",
		"wrong_item",
		"not_as_described",
		"no_longer_needed",
		"better_price",
		"other",
	}

	for _, reason := range validReasons {
		if !models.IsValidReturnReason(reason) {
			t.Errorf("Expected '%s' to be a valid return reason", reason)
		}
	}

	invalidReasons := []string{
		"invalid",
		"",
		"DefEcTiVe", // case sensitive
		"wrong-item", // underscore vs dash
	}

	for _, reason := range invalidReasons {
		if models.IsValidReturnReason(reason) {
			t.Errorf("Expected '%s' to be an invalid return reason", reason)
		}
	}
}
