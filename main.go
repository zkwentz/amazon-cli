package main

import (
	"fmt"
	"os"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("amazon-cli %s (commit: %s, built: %s)\n", version, commit, date)
		os.Exit(0)
	}

	fmt.Println("amazon-cli - Amazon Shopping CLI")
	fmt.Println()
	fmt.Println("This is the initial entry point. The CLI will be built with Cobra framework.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  amazon-cli [command]")
	fmt.Println()
	fmt.Println("Available Commands (coming soon):")
	fmt.Println("  auth         Authentication commands (login, status, logout)")
	fmt.Println("  orders       Manage orders (list, get, track, history)")
	fmt.Println("  returns      Manage returns (list, options, create, label, status)")
	fmt.Println("  search       Search for products")
	fmt.Println("  product      Get product details and reviews")
	fmt.Println("  cart         Manage shopping cart")
	fmt.Println("  buy          Quick purchase a product")
	fmt.Println("  subscriptions Manage Subscribe & Save subscriptions")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --help, -h       Show help information")
	fmt.Println("  --version, -v    Show version information")
	fmt.Println()
	fmt.Println("Run 'go mod init github.com/zkwentz/amazon-cli' to initialize the Go module.")
}
