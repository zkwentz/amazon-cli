package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

// OutputFormat represents the output format type
type OutputFormat string

const (
	// FormatJSON outputs data in JSON format
	FormatJSON OutputFormat = "json"
	// FormatTable outputs data in a human-readable table format
	FormatTable OutputFormat = "table"
	// FormatRaw outputs data as raw string representation
	FormatRaw OutputFormat = "raw"
)

// Printer handles output formatting and printing
type Printer struct {
	format OutputFormat
	quiet  bool
}

// NewPrinter creates a new Printer with the specified format and quiet mode
func NewPrinter(format string, quiet bool) *Printer {
	outputFormat := FormatJSON
	switch strings.ToLower(format) {
	case "json":
		outputFormat = FormatJSON
	case "table":
		outputFormat = FormatTable
	case "raw":
		outputFormat = FormatRaw
	default:
		outputFormat = FormatJSON
	}

	return &Printer{
		format: outputFormat,
		quiet:  quiet,
	}
}

// Print outputs the data according to the configured format
func (p *Printer) Print(data interface{}) error {
	if p.quiet {
		return nil
	}

	switch p.format {
	case FormatJSON:
		return p.printJSON(data)
	case FormatTable:
		return p.printTable(data)
	case FormatRaw:
		return p.printRaw(data)
	default:
		return p.printJSON(data)
	}
}

// printJSON marshals data to JSON and prints to stdout
func (p *Printer) printJSON(data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}
	fmt.Println(string(jsonData))
	return nil
}

// printTable formats data as a human-readable table
// This is a basic implementation that can be extended for specific data types
func (p *Printer) printTable(data interface{}) error {
	// For now, convert to JSON and display in a basic format
	// This can be extended to handle specific structs with proper table formatting
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data for table: %w", err)
	}

	// Create a simple table with key-value pairs using the new API
	table := tablewriter.NewWriter(os.Stdout)
	table.Header("Field", "Value")

	// Parse as map to display key-value pairs
	var dataMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &dataMap); err == nil {
		for key, value := range dataMap {
			valueStr := fmt.Sprintf("%v", value)
			// Truncate long values
			if len(valueStr) > 100 {
				valueStr = valueStr[:97] + "..."
			}
			if err := table.Append(key, valueStr); err != nil {
				return fmt.Errorf("failed to append table row: %w", err)
			}
		}
		if err := table.Render(); err != nil {
			return fmt.Errorf("failed to render table: %w", err)
		}
		return nil
	}

	// Fallback to raw JSON display
	fmt.Println(string(jsonData))
	return nil
}

// printRaw prints the raw string representation of data
func (p *Printer) printRaw(data interface{}) error {
	fmt.Printf("%v\n", data)
	return nil
}

// PrintError formats and outputs an error according to the configured format
func (p *Printer) PrintError(err error) error {
	if err == nil {
		return nil
	}

	// Check if it's already a CLIError
	if cliErr, ok := err.(*models.CLIError); ok {
		return p.printCLIError(cliErr)
	}

	// Convert standard error to CLIError
	cliErr := models.NewCLIError("ERROR", err.Error())
	return p.printCLIError(cliErr)
}

// printCLIError formats and prints a CLIError
func (p *Printer) printCLIError(err *models.CLIError) error {
	errorOutput := map[string]interface{}{
		"error": err,
	}

	switch p.format {
	case FormatJSON:
		jsonData, marshalErr := json.MarshalIndent(errorOutput, "", "  ")
		if marshalErr != nil {
			return fmt.Errorf("failed to marshal error JSON: %w", marshalErr)
		}
		fmt.Fprintln(os.Stderr, string(jsonData))
	case FormatTable, FormatRaw:
		// For table and raw formats, print a simplified error message
		fmt.Fprintf(os.Stderr, "Error [%s]: %s\n", err.Code, err.Message)
		if len(err.Details) > 0 {
			fmt.Fprintf(os.Stderr, "Details: %v\n", err.Details)
		}
	default:
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	return nil
}

// PrintSuccess is a convenience method to print a success message
func (p *Printer) PrintSuccess(message string) error {
	if p.quiet {
		return nil
	}

	data := map[string]string{
		"status":  "success",
		"message": message,
	}

	return p.Print(data)
}

// IsQuiet returns whether the printer is in quiet mode
func (p *Printer) IsQuiet() bool {
	return p.quiet
}

// GetFormat returns the current output format
func (p *Printer) GetFormat() OutputFormat {
	return p.format
}
