package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/spf13/cobra"
)

// executeCommand is a helper function to execute a command and capture output
func executeCommand(root *cobra.Command, args ...string) (output string, err error) {
	// Redirect stdout to capture JSON output
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	buf := new(bytes.Buffer)
	root.SetOut(buf)
	root.SetErr(buf)
	root.SetArgs(args)

	err = root.Execute()

	// Restore stdout and read captured output
	w.Close()
	os.Stdout = old

	var capturedBuf bytes.Buffer
	capturedBuf.ReadFrom(r)
	output = capturedBuf.String()

	return output, err
}

func TestAuthLoginCmd(t *testing.T) {
	// Create a root command to test with
	rootCmd := &cobra.Command{Use: "amazon-cli"}
	rootCmd.AddCommand(authCmd)

	output, err := executeCommand(rootCmd, "auth", "login")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Parse JSON output
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	// Verify the status field exists
	if status, ok := result["status"]; !ok {
		t.Errorf("Expected 'status' field in output")
	} else if status != "not_implemented" {
		t.Logf("Login status: %v", status)
	}
}

func TestAuthStatusCmd(t *testing.T) {
	rootCmd := &cobra.Command{Use: "amazon-cli"}
	rootCmd.AddCommand(authCmd)

	output, err := executeCommand(rootCmd, "auth", "status")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Parse JSON output
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	// Verify the authenticated field exists
	if _, ok := result["authenticated"]; !ok {
		t.Errorf("Expected 'authenticated' field in output")
	}

	// Since no config exists, authenticated should be false
	if authenticated, ok := result["authenticated"].(bool); ok && authenticated {
		t.Errorf("Expected authenticated to be false when no config exists")
	}
}

func TestAuthLogoutCmd(t *testing.T) {
	rootCmd := &cobra.Command{Use: "amazon-cli"}
	rootCmd.AddCommand(authCmd)

	output, err := executeCommand(rootCmd, "auth", "logout")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Parse JSON output
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	// Verify the status field exists and is "logged_out"
	if status, ok := result["status"]; !ok {
		t.Errorf("Expected 'status' field in output")
	} else if status != "logged_out" {
		t.Errorf("Expected status to be 'logged_out', got %v", status)
	}
}

func TestAuthCmdExists(t *testing.T) {
	rootCmd := &cobra.Command{Use: "amazon-cli"}
	rootCmd.AddCommand(authCmd)

	// Verify auth command exists
	cmd, _, err := rootCmd.Find([]string{"auth"})
	if err != nil {
		t.Fatalf("Failed to find auth command: %v", err)
	}

	if cmd.Use != "auth" {
		t.Errorf("Expected command Use to be 'auth', got %v", cmd.Use)
	}

	// Verify subcommands exist
	subcommands := []string{"login", "status", "logout"}
	for _, subcmd := range subcommands {
		_, _, err := rootCmd.Find([]string{"auth", subcmd})
		if err != nil {
			t.Errorf("Failed to find 'auth %s' subcommand: %v", subcmd, err)
		}
	}
}

func TestOutputJSON(t *testing.T) {
	tests := []struct {
		name string
		data interface{}
		want map[string]interface{}
	}{
		{
			name: "simple map",
			data: map[string]interface{}{"key": "value"},
			want: map[string]interface{}{"key": "value"},
		},
		{
			name: "nested map",
			data: map[string]interface{}{
				"status": "success",
				"data": map[string]interface{}{
					"authenticated": true,
				},
			},
			want: map[string]interface{}{
				"status": "success",
				"data": map[string]interface{}{
					"authenticated": true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			err := outputJSON(tt.data)
			if err != nil {
				t.Fatalf("outputJSON() error = %v", err)
			}

			w.Close()
			os.Stdout = old

			var got map[string]interface{}
			if err := json.NewDecoder(r).Decode(&got); err != nil {
				t.Fatalf("Failed to decode output: %v", err)
			}

			// Basic comparison - just verify the structure
			if len(got) != len(tt.want) {
				t.Errorf("outputJSON() produced different number of fields: got %d, want %d", len(got), len(tt.want))
			}
		})
	}
}
