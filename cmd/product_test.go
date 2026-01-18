package cmd

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestValidateASIN(t *testing.T) {
	tests := []struct {
		name    string
		asin    string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid ASIN",
			asin:    "B08N5WRWNW",
			wantErr: false,
		},
		{
			name:    "valid ASIN with numbers",
			asin:    "B08N5WRW12",
			wantErr: false,
		},
		{
			name:    "too short",
			asin:    "B08N5WRW",
			wantErr: true,
			errMsg:  "must be exactly 10 characters",
		},
		{
			name:    "too long",
			asin:    "B08N5WRWNW1",
			wantErr: true,
			errMsg:  "must be exactly 10 characters",
		},
		{
			name:    "lowercase letters",
			asin:    "b08n5wrwnw",
			wantErr: true,
			errMsg:  "invalid ASIN format",
		},
		{
			name:    "special characters",
			asin:    "B08N5WRW-W",
			wantErr: true,
			errMsg:  "invalid ASIN format",
		},
		{
			name:    "empty string",
			asin:    "",
			wantErr: true,
			errMsg:  "must be exactly 10 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateASIN(tt.asin)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateASIN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err != nil && tt.errMsg != "" {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("validateASIN() error = %v, want error containing %v", err, tt.errMsg)
				}
			}
		})
	}
}

func TestProductGetCmd(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid ASIN",
			args:    []string{"B08N5WRWNW"},
			wantErr: false,
		},
		{
			name:    "no arguments",
			args:    []string{},
			wantErr: true,
			errMsg:  "requires exactly 1 arg(s)",
		},
		{
			name:    "too many arguments",
			args:    []string{"B08N5WRWNW", "extra"},
			wantErr: true,
			errMsg:  "requires exactly 1 arg(s)",
		},
		{
			name:    "invalid ASIN format",
			args:    []string{"invalid"},
			wantErr: true,
			errMsg:  "invalid ASIN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := GetProductCmd()
			getCmd := findSubCommand(cmd, "get")
			if getCmd == nil {
				t.Fatal("get subcommand not found")
			}

			// Capture output
			buf := new(bytes.Buffer)
			getCmd.SetOut(buf)
			getCmd.SetErr(buf)

			// Set args
			getCmd.SetArgs(tt.args)

			// Execute
			err := getCmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("productGetCmd error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && tt.errMsg != "" {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("productGetCmd error = %v, want error containing %v", err, tt.errMsg)
				}
			}

			// If successful, verify JSON output
			if !tt.wantErr {
				output := buf.String()
				var result map[string]interface{}
				if err := json.Unmarshal([]byte(output), &result); err != nil {
					t.Errorf("failed to parse JSON output: %v", err)
				}

				// Check for expected fields
				expectedFields := []string{"asin", "title", "price", "rating", "review_count", "prime", "in_stock"}
				for _, field := range expectedFields {
					if _, ok := result[field]; !ok {
						t.Errorf("expected field %q not found in output", field)
					}
				}

				// Verify ASIN matches
				if asin, ok := result["asin"].(string); ok {
					if asin != tt.args[0] {
						t.Errorf("ASIN in output = %v, want %v", asin, tt.args[0])
					}
				}
			}
		})
	}
}

func TestProductReviewsCmd(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		flags      map[string]string
		wantErr    bool
		errMsg     string
		checkLimit bool
		limit      int
	}{
		{
			name:    "valid ASIN with default limit",
			args:    []string{"B08N5WRWNW"},
			wantErr: false,
		},
		{
			name:       "valid ASIN with custom limit",
			args:       []string{"B08N5WRWNW"},
			flags:      map[string]string{"limit": "5"},
			wantErr:    false,
			checkLimit: true,
			limit:      5,
		},
		{
			name:    "no arguments",
			args:    []string{},
			wantErr: true,
			errMsg:  "requires exactly 1 arg(s)",
		},
		{
			name:    "invalid ASIN",
			args:    []string{"invalid"},
			wantErr: true,
			errMsg:  "invalid ASIN",
		},
		{
			name:    "negative limit",
			args:    []string{"B08N5WRWNW"},
			flags:   map[string]string{"limit": "-1"},
			wantErr: true,
			errMsg:  "limit must be a positive integer",
		},
		{
			name:    "zero limit",
			args:    []string{"B08N5WRWNW"},
			flags:   map[string]string{"limit": "0"},
			wantErr: true,
			errMsg:  "limit must be a positive integer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := GetProductCmd()
			reviewsCmd := findSubCommand(cmd, "reviews")
			if reviewsCmd == nil {
				t.Fatal("reviews subcommand not found")
			}

			// Capture output
			buf := new(bytes.Buffer)
			reviewsCmd.SetOut(buf)
			reviewsCmd.SetErr(buf)

			// Build args with flags
			args := tt.args
			if tt.flags != nil {
				for key, value := range tt.flags {
					args = append(args, "--"+key, value)
				}
			}

			// Set args
			reviewsCmd.SetArgs(args)

			// Execute
			err := reviewsCmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("productReviewsCmd error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil && tt.errMsg != "" {
				if !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("productReviewsCmd error = %v, want error containing %v", err, tt.errMsg)
				}
			}

			// If successful, verify JSON output
			if !tt.wantErr {
				output := buf.String()
				var result map[string]interface{}
				if err := json.Unmarshal([]byte(output), &result); err != nil {
					t.Errorf("failed to parse JSON output: %v", err)
				}

				// Check for expected fields
				expectedFields := []string{"asin", "average_rating", "total_reviews", "reviews"}
				for _, field := range expectedFields {
					if _, ok := result[field]; !ok {
						t.Errorf("expected field %q not found in output", field)
					}
				}

				// Verify ASIN matches
				if asin, ok := result["asin"].(string); ok {
					if asin != tt.args[0] {
						t.Errorf("ASIN in output = %v, want %v", asin, tt.args[0])
					}
				}

				// Verify limit if specified
				if tt.checkLimit {
					if reviews, ok := result["reviews"].([]interface{}); ok {
						if len(reviews) > tt.limit {
							t.Errorf("number of reviews = %d, want at most %d", len(reviews), tt.limit)
						}
					}
				}

				// Verify reviews structure
				if reviews, ok := result["reviews"].([]interface{}); ok {
					for i, review := range reviews {
						reviewMap, ok := review.(map[string]interface{})
						if !ok {
							t.Errorf("review %d is not a map", i)
							continue
						}

						// Check for required review fields
						reviewFields := []string{"rating", "title", "body", "author", "date", "verified"}
						for _, field := range reviewFields {
							if _, ok := reviewMap[field]; !ok {
								t.Errorf("review %d missing field %q", i, field)
							}
						}
					}
				}
			}
		})
	}
}

func TestProductCmdStructure(t *testing.T) {
	cmd := GetProductCmd()

	if cmd == nil {
		t.Fatal("GetProductCmd() returned nil")
	}

	if cmd.Use != "product" {
		t.Errorf("product command Use = %v, want %v", cmd.Use, "product")
	}

	// Check that subcommands exist
	expectedSubcommands := []string{"get", "reviews"}
	for _, subcmd := range expectedSubcommands {
		if findSubCommand(cmd, subcmd) == nil {
			t.Errorf("expected subcommand %q not found", subcmd)
		}
	}
}

// Helper function to find a subcommand by name
func findSubCommand(cmd *cobra.Command, name string) *cobra.Command {
	for _, subcmd := range cmd.Commands() {
		if subcmd.Name() == name {
			return subcmd
		}
	}
	return nil
}
