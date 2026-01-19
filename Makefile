# amazon-cli Makefile

BINARY_NAME=amazon-cli
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=$(GOCMD) fmt

# Platforms
PLATFORMS=darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64

.PHONY: all build build-all test lint clean install fmt tidy help

all: build

## build: Build for current platform
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) .

## build-all: Build for all platforms
build-all:
	@mkdir -p dist
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*} GOARCH=$${platform#*/} \
		$(GOBUILD) $(LDFLAGS) -o dist/$(BINARY_NAME)-$${platform%/*}-$${platform#*/}$$([ "$${platform%/*}" = "windows" ] && echo ".exe") . ; \
		echo "Built: $(BINARY_NAME)-$${platform%/*}-$${platform#*/}" ; \
	done

## test: Run all tests
test:
	$(GOTEST) -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

## lint: Run linter
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

## clean: Remove build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -rf dist/
	rm -f coverage.out coverage.html

## install: Install to GOPATH/bin
install:
	$(GOBUILD) $(LDFLAGS) -o $(GOPATH)/bin/$(BINARY_NAME) .

## fmt: Format code
fmt:
	$(GOFMT) ./...

## tidy: Tidy dependencies
tidy:
	$(GOMOD) tidy

## version: Print version info
version:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Date: $(DATE)"

## help: Show this help
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed 's/^/ /'
