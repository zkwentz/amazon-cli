package output

import (
	"encoding/json"
	"os"
)

// Printer handles output formatting
type Printer struct {
	format string
	quiet  bool
}

// NewPrinter creates a new output printer
func NewPrinter(format string, quiet bool) *Printer {
	return &Printer{
		format: format,
		quiet:  quiet,
	}
}

// Print outputs data in the specified format
func (p *Printer) Print(data interface{}) error {
	if p.quiet {
		return nil
	}

	switch p.format {
	case "json", "":
		return p.printJSON(data)
	default:
		return p.printJSON(data)
	}
}

// printJSON outputs data as JSON
func (p *Printer) printJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// PrintError outputs an error
func (p *Printer) PrintError(err error) error {
	errorOutput := map[string]interface{}{
		"error": map[string]interface{}{
			"message": err.Error(),
		},
	}
	return p.printJSON(errorOutput)
}
