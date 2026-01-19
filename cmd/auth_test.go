package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TestAuthLoginCmd(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Create a new command instance for testing
	cmd := &cobra.Command{
		Use: "login",
		Run: authLoginCmd.Run,
	}

	// Execute the command
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Command execution failed: %v", err)
	}

	// Restore stdout and read captured output
	w.Close()
	os.Stdout = oldStdout
	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Parse the JSON output
	var result map[string]interface{}
	err = json.Unmarshal(buf.Bytes(), &result)
	if err != nil {
		t.Fatalf("Failed to parse JSON output: %v\nOutput: %s", err, buf.String())
	}

	// Verify the status field
	status, ok := result["status"].(string)
	if !ok {
		t.Fatal("status field is missing or not a string")
	}
	if status != "login_required" {
		t.Errorf("Expected status 'login_required', got '%s'", status)
	}

	// Verify the message field
	message, ok := result["message"].(string)
	if !ok {
		t.Fatal("message field is missing or not a string")
	}
	expectedMessage := "Browser-based login not yet implemented"
	if message != expectedMessage {
		t.Errorf("Expected message '%s', got '%s'", expectedMessage, message)
	}
}
