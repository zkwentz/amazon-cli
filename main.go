package main

import (
	"os"

	"github.com/michaelshimeles/amazon-cli/cmd"
)

func main() {
	exitCode := cmd.Execute()
	os.Exit(exitCode)
}
