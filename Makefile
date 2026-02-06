# Makefile for Brewls CLI

# Go parameters
GO := go
GO_BUILD_TARGET := ./cmd/brewls
GO_BINARY_NAME := brewls
GO_BINARY_PATH := $(shell $(GO) env GOPATH)/bin/$(GO_BINARY_NAME)

# Default target
.PHONY: all
all: install

# Build the application binary
.PHONY: build
build:
	@echo "Building $(GO_BINARY_NAME)..."
	mkdir -p bin
	$(GO) build -o bin/$(GO_BINARY_NAME) $(GO_BUILD_TARGET)
	@echo "Build complete. Binary: ./bin/$(GO_BINARY_NAME)"

# Install the application using go install
.PHONY: install
install: build
	@echo "Installing $(GO_BINARY_NAME) from ./bin/$(GO_BINARY_NAME) to $(GO_BINARY_PATH)..."
	cp bin/$(GO_BINARY_NAME) $(GO_BINARY_PATH)
	@echo "Installation complete. Ensure $(shell $(GO) env GOPATH)/bin is in your PATH."

# Run all tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GO) test ./... -v

# Run the application
.PHONY: run
run: build
	@echo "Running $(GO_BINARY_NAME)..."
	./$(GO_BINARY_NAME)

# Clean up build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up..."
	$(GO) clean
	@rm -f bin/$(GO_BINARY_NAME)
	@echo "Clean complete."

# Tidy Go module dependencies
.PHONY: tidy
tidy:
	@echo "Tidying Go module dependencies..."
	$(GO) mod tidy
	@echo "Go module dependencies tidied."

# Lint and format check
.PHONY: lint
lint:
	@echo "Running format check..."
	@gofmt_diff=$$(gofmt -l .); \
	if [ -n "$$gofmt_diff" ]; then \
		echo "gofmt found unformatted files:"; \
		echo "$$gofmt_diff"; \
		echo "Run 'make fmt' to format your code."; \
		exit 1; \
	fi
	@echo "Running go vet..."
	$(GO) vet ./...
	@echo "Running golangci-lint..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "golangci-lint not found. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.63.4"; \
		exit 1; \
	fi
	golangci-lint run ./...
	@echo "Lint check complete."

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	gofmt -w .
	@echo "Code formatted."

# Help message
.PHONY: help
help:
	@echo "Usage: make <command>"
	@echo ""
	@echo "Commands:"
	@echo "  all     - Builds, tidies, and installs the application (default)."
	@echo "  build   - Builds the $(GO_BINARY_NAME) executable."
	@echo "  install - Installs the $(GO_BINARY_NAME) executable to GOPATH/bin."
	@echo "  test    - Runs all unit tests."
	@echo "  run     - Builds and runs the application."
	@echo "  clean   - Removes build artifacts."
	@echo "  tidy    - Cleans up go.sum and go.mod files."
	@echo "  lint    - Runs format check, go vet, and golangci-lint."
	@echo "  fmt     - Formats code with gofmt."
	@echo "  help    - Displays this help message."
	@echo ""
	@echo "Prerequisites:"
	@echo "  - Go (version 1.22+ recommended)"
	@echo "  - golangci-lint (for make lint)"
	@echo "  - Homebrew ('brew' command must be in your PATH)"