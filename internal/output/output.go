package output

import (
	"encoding/json"
	"fmt"
	"os"
)

// Format represents the output format type
type Format string

const (
	FormatJSON  Format = "json"
	FormatTable Format = "table"
	FormatRaw   Format = "raw"
)

// Printer handles output formatting
type Printer struct {
	format Format
	quiet  bool
}

// NewPrinter creates a new Printer with the specified format
func NewPrinter(format string, quiet bool) *Printer {
	f := Format(format)
	if f != FormatJSON && f != FormatTable && f != FormatRaw {
		f = FormatJSON
	}
	return &Printer{
		format: f,
		quiet:  quiet,
	}
}

// Print outputs data in the configured format
func (p *Printer) Print(data interface{}) error {
	if p.quiet {
		return nil
	}

	switch p.format {
	case FormatJSON:
		return p.printJSON(data)
	case FormatTable:
		// For now, fall back to JSON
		return p.printJSON(data)
	case FormatRaw:
		fmt.Fprintf(os.Stdout, "%v\n", data)
		return nil
	default:
		return p.printJSON(data)
	}
}

func (p *Printer) printJSON(data interface{}) error {
	output, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, string(output))
	return nil
}

// PrintError outputs an error in the configured format
func (p *Printer) PrintError(err error) error {
	errResponse := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    "GENERAL_ERROR",
			"message": err.Error(),
			"details": map[string]interface{}{},
		},
	}
	return p.Print(errResponse)
}

// JSON is a convenience function to print JSON to stdout
func JSON(data interface{}) error {
	return NewPrinter("json", false).Print(data)
}

// Error is a convenience function to print an error to stderr
func Error(code, message string, details map[string]interface{}) error {
	if details == nil {
		details = map[string]interface{}{}
	}
	errResponse := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
			"details": details,
		},
	}
	output, err := json.MarshalIndent(errResponse, "", "  ")
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, string(output))
	return nil
}
