package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLI_Help_Works(t *testing.T) {
	// Get the project root directory (two levels up from test/e2e)
	projectRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("Failed to get project root: %v", err)
	}

	// Build the binary
	binaryPath := filepath.Join(projectRoot, "amazon-cli")
	buildCmd := exec.Command("go", "build", "-o", binaryPath)
	buildCmd.Dir = projectRoot
	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build binary: %v\nOutput: %s", err, buildOutput)
	}

	// Clean up the binary after the test
	defer func() {
		if err := os.Remove(binaryPath); err != nil && !os.IsNotExist(err) {
			t.Logf("Warning: Failed to clean up binary: %v", err)
		}
	}()

	// Run the binary with --help flag
	helpCmd := exec.Command(binaryPath, "--help")
	helpCmd.Dir = projectRoot
	output, err := helpCmd.CombinedOutput()

	// Verify exit code is 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			t.Fatalf("Command exited with non-zero status: %d\nOutput: %s", exitErr.ExitCode(), output)
		}
		t.Fatalf("Failed to run command: %v\nOutput: %s", err, output)
	}

	// Verify output contains "amazon-cli"
	outputStr := string(output)
	if !strings.Contains(outputStr, "amazon-cli") {
		t.Errorf("Expected output to contain 'amazon-cli', but it didn't.\nOutput: %s", outputStr)
	}
}
