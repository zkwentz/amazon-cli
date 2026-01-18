package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// Helper function to execute a command and capture output
func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	err = root.Execute()
	return buf.String(), err
}

func TestOrdersListCommand(t *testing.T) {
	// Test basic list command
	output, err := executeCommand(ordersCmd, "list")
	if err != nil {
		t.Fatalf("ordersListCmd failed: %v", err)
	}

	if !strings.Contains(output, "orders") {
		t.Errorf("Expected output to contain 'orders', got: %s", output)
	}
}

func TestOrdersListWithLimit(t *testing.T) {
	// Test list command with limit flag
	output, err := executeCommand(ordersCmd, "list", "--limit", "5")
	if err != nil {
		t.Fatalf("ordersListCmd with limit failed: %v", err)
	}

	if !strings.Contains(output, "Limit: 5") {
		t.Errorf("Expected output to contain 'Limit: 5', got: %s", output)
	}
}

func TestOrdersListWithStatus(t *testing.T) {
	// Test list command with status flag
	output, err := executeCommand(ordersCmd, "list", "--status", "delivered")
	if err != nil {
		t.Fatalf("ordersListCmd with status failed: %v", err)
	}

	if !strings.Contains(output, "Status: delivered") {
		t.Errorf("Expected output to contain 'Status: delivered', got: %s", output)
	}
}

func TestOrdersGetCommand(t *testing.T) {
	// Test get command with order ID
	orderID := "123-4567890-1234567"
	output, err := executeCommand(ordersCmd, "get", orderID)
	if err != nil {
		t.Fatalf("ordersGetCmd failed: %v", err)
	}

	if !strings.Contains(output, orderID) {
		t.Errorf("Expected output to contain order ID '%s', got: %s", orderID, output)
	}
}

func TestOrdersGetCommandNoArgs(t *testing.T) {
	// Test get command without order ID should fail
	_, err := executeCommand(ordersCmd, "get")
	if err == nil {
		t.Error("Expected error when calling get without order ID, got nil")
	}
}

func TestOrdersTrackCommand(t *testing.T) {
	// Test track command with order ID
	orderID := "123-4567890-1234567"
	output, err := executeCommand(ordersCmd, "track", orderID)
	if err != nil {
		t.Fatalf("ordersTrackCmd failed: %v", err)
	}

	if !strings.Contains(output, orderID) {
		t.Errorf("Expected output to contain order ID '%s', got: %s", orderID, output)
	}

	if !strings.Contains(output, "tracking") {
		t.Errorf("Expected output to contain 'tracking', got: %s", output)
	}
}

func TestOrdersTrackCommandNoArgs(t *testing.T) {
	// Test track command without order ID should fail
	_, err := executeCommand(ordersCmd, "track")
	if err == nil {
		t.Error("Expected error when calling track without order ID, got nil")
	}
}

func TestOrdersHistoryCommand(t *testing.T) {
	// Test history command
	output, err := executeCommand(ordersCmd, "history")
	if err != nil {
		t.Fatalf("ordersHistoryCmd failed: %v", err)
	}

	if !strings.Contains(output, "orders") {
		t.Errorf("Expected output to contain 'orders', got: %s", output)
	}
}

func TestOrdersHistoryWithYear(t *testing.T) {
	// Test history command with year flag
	output, err := executeCommand(ordersCmd, "history", "--year", "2024")
	if err != nil {
		t.Fatalf("ordersHistoryCmd with year failed: %v", err)
	}

	if !strings.Contains(output, "2024") {
		t.Errorf("Expected output to contain '2024', got: %s", output)
	}
}

func TestOrdersHistoryWithFormat(t *testing.T) {
	// Test history command with format flag
	output, err := executeCommand(ordersCmd, "history", "--format", "json")
	if err != nil {
		t.Fatalf("ordersHistoryCmd with format failed: %v", err)
	}

	if !strings.Contains(output, "json") {
		t.Errorf("Expected output to contain 'json', got: %s", output)
	}
}

func TestOrdersCommandStructure(t *testing.T) {
	// Test that all subcommands are registered
	subcommands := []string{"list", "get", "track", "history"}

	for _, subcmd := range subcommands {
		found := false
		for _, cmd := range ordersCmd.Commands() {
			if cmd.Name() == subcmd {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected subcommand '%s' to be registered", subcmd)
		}
	}
}

func TestOrdersListFlags(t *testing.T) {
	// Test that list command has required flags
	limitFlag := ordersListCmd.Flags().Lookup("limit")
	if limitFlag == nil {
		t.Error("Expected 'limit' flag to exist on list command")
	}

	statusFlag := ordersListCmd.Flags().Lookup("status")
	if statusFlag == nil {
		t.Error("Expected 'status' flag to exist on list command")
	}
}

func TestOrdersHistoryFlags(t *testing.T) {
	// Test that history command has required flags
	yearFlag := ordersHistoryCmd.Flags().Lookup("year")
	if yearFlag == nil {
		t.Error("Expected 'year' flag to exist on history command")
	}

	formatFlag := ordersHistoryCmd.Flags().Lookup("format")
	if formatFlag == nil {
		t.Error("Expected 'format' flag to exist on history command")
	}
}
