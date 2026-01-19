package main

import (
	"os"

	"github.com/zkwentz/amazon-cli/cmd"
)

var version = "dev"

func main() {
	cmd.SetVersion(version)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
