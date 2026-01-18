package output

import (
	"encoding/json"
	"fmt"
	"os"
)

// Printer handles output formatting
type Printer struct {
	format string
	quiet  bool
}

// NewPrinter creates a new Printer
func NewPrinter(format string, quiet bool) *Printer {
	return &Printer{
		format: format,
		quiet:  quiet,
	}
}

// Print outputs data in the configured format
func (p *Printer) Print(data interface{}) error {
	if p.quiet {
		return nil
	}

	switch p.format {
	case "json":
		return p.printJSON(data)
	default:
		return p.printJSON(data)
	}
}

func (p *Printer) printJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// PrintError outputs an error in JSON format
func (p *Printer) PrintError(code, message string, details interface{}) error {
	errorData := map[string]interface{}{
		"error": map[string]interface{}{
			"code":    code,
			"message": message,
			"details": details,
		},
	}

	encoder := json.NewEncoder(os.Stderr)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(errorData); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", message)
		return err
	}
	return nil
}
