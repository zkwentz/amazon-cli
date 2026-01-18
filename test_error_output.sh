#!/bin/bash
# Test script to verify error outputs are valid JSON

echo "Testing error output is valid JSON..."

# Build the CLI
go build -o amazon-cli . || exit 1

# Test 1: Invalid command should output JSON error
echo "Test 1: Invalid command"
output=$(./amazon-cli invalidcommand 2>&1)
if echo "$output" | jq . >/dev/null 2>&1; then
    echo "✓ Invalid command outputs valid JSON"
else
    echo "✗ Invalid command does not output valid JSON"
    echo "Output: $output"
    exit 1
fi

echo ""
echo "All error output tests passed!"
exit 0
