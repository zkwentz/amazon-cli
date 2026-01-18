package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/michaelshimeles/amazon-cli/pkg/models"
)

func TestCartListCommand(t *testing.T) {
	// Create a buffer to capture stdout
	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute the command
	err := cartListCmd.RunE(cartListCmd, []string{})
	if err != nil {
		os.Stdout = oldStdout
		t.Fatalf("cart list command failed: %v", err)
	}

	// Restore stdout and read captured output
	w.Close()
	os.Stdout = oldStdout
	buf.ReadFrom(r)

	// Parse the output as JSON
	output := buf.String()
	if output == "" {
		t.Fatal("cart list command produced no output")
	}

	// Verify it's valid JSON
	var cart models.Cart
	if err := json.Unmarshal([]byte(output), &cart); err != nil {
		t.Fatalf("cart list output is not valid JSON: %v\nOutput: %s", err, output)
	}

	// Verify cart structure
	if cart.Items == nil {
		t.Error("cart items should not be nil")
	}

	// Empty cart should have zero item count
	if cart.ItemCount != 0 {
		t.Errorf("empty cart should have ItemCount = 0, got %d", cart.ItemCount)
	}

	// Empty cart should have zero totals
	if cart.Subtotal != 0 {
		t.Errorf("empty cart should have Subtotal = 0, got %f", cart.Subtotal)
	}
}

func TestCartListCommandJSONFormat(t *testing.T) {
	// Create a buffer to capture stdout
	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute the command
	err := cartListCmd.RunE(cartListCmd, []string{})
	if err != nil {
		os.Stdout = oldStdout
		t.Fatalf("cart list command failed: %v", err)
	}

	// Restore stdout and read captured output
	w.Close()
	os.Stdout = oldStdout
	buf.ReadFrom(r)

	output := buf.String()

	// Verify JSON structure contains expected fields
	expectedFields := []string{
		"\"items\"",
		"\"subtotal\"",
		"\"estimated_tax\"",
		"\"total\"",
		"\"item_count\"",
	}

	for _, field := range expectedFields {
		if !strings.Contains(output, field) {
			t.Errorf("cart list output missing expected field: %s\nOutput: %s", field, output)
		}
	}

	// Verify JSON is properly indented (for readability)
	if !strings.Contains(output, "\n") {
		t.Error("cart list output should be indented JSON")
	}
}

func TestCartCommandExists(t *testing.T) {
	// Verify cart command exists
	if cartCmd == nil {
		t.Fatal("cart command is nil")
	}

	// Verify cart list subcommand exists
	if cartListCmd == nil {
		t.Fatal("cart list command is nil")
	}

	// Verify cart command has correct properties
	if cartCmd.Use != "cart" {
		t.Errorf("cart command Use = %s, want 'cart'", cartCmd.Use)
	}

	// Verify cart list command has correct properties
	if cartListCmd.Use != "list" {
		t.Errorf("cart list command Use = %s, want 'list'", cartListCmd.Use)
	}
}

func TestCartListCommandOutput(t *testing.T) {
	// Create a buffer to capture stdout
	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Execute the command
	err := cartListCmd.RunE(cartListCmd, []string{})
	if err != nil {
		os.Stdout = oldStdout
		t.Fatalf("unexpected error: %v", err)
	}

	// Restore stdout and read captured output
	w.Close()
	os.Stdout = oldStdout
	buf.ReadFrom(r)

	output := buf.String()

	// Parse as JSON to verify structure
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("failed to parse output as JSON: %v", err)
	}

	// Verify required fields exist
	if _, ok := result["items"]; !ok {
		t.Error("output missing 'items' field")
	}
	if _, ok := result["subtotal"]; !ok {
		t.Error("output missing 'subtotal' field")
	}
	if _, ok := result["estimated_tax"]; !ok {
		t.Error("output missing 'estimated_tax' field")
	}
	if _, ok := result["total"]; !ok {
		t.Error("output missing 'total' field")
	}
	if _, ok := result["item_count"]; !ok {
		t.Error("output missing 'item_count' field")
	}
}
