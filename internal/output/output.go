package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/zkwentz/amazon-cli/pkg/models"
)

// OutputFormat represents the output format type
type OutputFormat string

const (
	// JSON output format - structured JSON output
	JSON OutputFormat = "json"
	// Table output format - human-readable table
	Table OutputFormat = "table"
	// Raw output format - raw string representation
	Raw OutputFormat = "raw"
)

// Printer handles formatted output to stdout
type Printer struct {
	format OutputFormat
	quiet  bool
}

// NewPrinter creates a new Printer with the specified format and quiet mode
func NewPrinter(format string, quiet bool) *Printer {
	outputFormat := JSON // default to JSON
	switch format {
	case "json":
		outputFormat = JSON
	case "table":
		outputFormat = Table
	case "raw":
		outputFormat = Raw
	}

	return &Printer{
		format: outputFormat,
		quiet:  quiet,
	}
}

// Print outputs data in the configured format
// For JSON: marshals with indentation and prints to stdout
// For Table: uses tablewriter package for human-readable output
// For Raw: prints raw string representation
func (p *Printer) Print(data interface{}) error {
	if p.quiet {
		return nil
	}

	switch p.format {
	case JSON:
		return p.printJSON(data)
	case Table:
		return p.printTable(data)
	case Raw:
		return p.printRaw(data)
	default:
		return p.printJSON(data)
	}
}

// printJSON marshals data to JSON with indentation and prints to stdout
func (p *Printer) printJSON(data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	fmt.Println(string(jsonData))
	return nil
}

// printTable outputs data in a human-readable table format
func (p *Printer) printTable(data interface{}) error {
	// Convert data to JSON first to handle it uniformly
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data for table: %w", err)
	}

	// Parse JSON into a map or slice
	var parsed interface{}
	if err := json.Unmarshal(jsonData, &parsed); err != nil {
		return fmt.Errorf("failed to parse data for table: %w", err)
	}

	table := tablewriter.NewWriter(os.Stdout)

	// Handle different data types
	switch v := parsed.(type) {
	case map[string]interface{}:
		// Single object - display as key-value pairs
		for key, value := range v {
			if err := table.Append([]string{key, fmt.Sprintf("%v", value)}); err != nil {
				return fmt.Errorf("failed to append table row: %w", err)
			}
		}
	case []interface{}:
		// Array of objects - display as rows
		if len(v) > 0 {
			// Use first item to determine headers
			if firstItem, ok := v[0].(map[string]interface{}); ok {
				headers := make([]string, 0, len(firstItem))
				for key := range firstItem {
					headers = append(headers, key)
				}
				// Convert []string to []any for Header method
				headerAny := make([]any, len(headers))
				for i, h := range headers {
					headerAny[i] = h
				}
				table.Header(headerAny...)

				// Add rows
				for _, item := range v {
					if itemMap, ok := item.(map[string]interface{}); ok {
						row := make([]string, len(headers))
						for i, header := range headers {
							row[i] = fmt.Sprintf("%v", itemMap[header])
						}
						if err := table.Append(row); err != nil {
							return fmt.Errorf("failed to append table row: %w", err)
						}
					}
				}
			}
		}
	default:
		// For other types, just print as string
		fmt.Println(fmt.Sprintf("%v", v))
		return nil
	}

	if err := table.Render(); err != nil {
		return fmt.Errorf("failed to render table: %w", err)
	}
	return nil
}

// printRaw outputs the raw string representation of data
func (p *Printer) printRaw(data interface{}) error {
	fmt.Printf("%v\n", data)
	return nil
}

// PrintError formats and outputs errors as JSON
// Format: {"error": {"code": "...", "message": "...", "details": {}}}
func (p *Printer) PrintError(err error) error {
	if err == nil {
		return nil
	}

	var cliErr *models.CLIError

	// Check if it's already a CLIError
	if e, ok := err.(*models.CLIError); ok {
		cliErr = e
	} else {
		// Wrap generic errors as AMAZON_ERROR
		cliErr = models.NewCLIError(models.AmazonError, err.Error())
	}

	// Create error response wrapper
	errResponse := models.ErrorResponse{
		Error: cliErr,
	}

	// Always output errors as JSON regardless of format setting
	jsonData, err := json.MarshalIndent(errResponse, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal error JSON: %w", err)
	}

	fmt.Println(string(jsonData))
	return nil
}
