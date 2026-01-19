package main

import (
	"os"

	"github.com/zkwentz/amazon-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
