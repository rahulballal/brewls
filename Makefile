# Makefile for Brewls CLI

# Go parameters
GO := go
GO_BUILD_TARGET := ./cmd/brewls
GO_BINARY_NAME := brewls
GO_BINARY_PATH := $(shell $(GO) env GOPATH)/bin/$(GO_BINARY_NAME)

# Default target
.PHONY: all
all: check-brew install

# Check if brew command is available
.PHONY: check-brew
check-brew:
	@if ! command -v brew &> /dev/null; then \
		echo "Error: Homebrew 'brew' command not found."; \
		echo "Please install Homebrew (https://brew.sh/) before proceeding."; \
		exit 1; \
	fi
	@echo "Homebrew 'brew' command found."

# Build the application binary
.PHONY: build
build: check-brew
	@echo "Building $(GO_BINARY_NAME)..."
	$(GO) build -o $(GO_BINARY_NAME) $(GO_BUILD_TARGET)
	@echo "Build complete. Binary: ./$(GO_BINARY_NAME)"

# Install the application using go install
.PHONY: install
install: check-brew
	@echo "Installing $(GO_BINARY_NAME) to $(GO_BINARY_PATH)..."
	$(GO) install $(GO_BUILD_TARGET)
	@echo "Installation complete. Ensure $(shell $(GO) env GOPATH)/bin is in your PATH."

# Run all tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GO) test ./... -v

# Run the application
.PHONY: run
run: check-brew build
	@echo "Running $(GO_BINARY_NAME)..."
	./$(GO_BINARY_NAME)

# Clean up build artifacts
.PHONY: clean
clean:
	@echo "Cleaning up..."
	$(GO) clean
	@rm -f $(GO_BINARY_NAME)
	@echo "Clean complete."

# Help message
.PHONY: help
help:
	@echo "Usage: make <command>"
	@echo ""
	@echo "Commands:"
	@echo "  all     - Builds and installs the application (default)."
	@echo "  build   - Builds the $(GO_BINARY_NAME) executable."
	@echo "  install - Installs the $(GO_BINARY_NAME) executable to GOPATH/bin."
	@echo "  test    - Runs all unit tests."
	@echo "  run     - Builds and runs the application."
	@echo "  clean   - Removes build artifacts."
	@echo "  help    - Displays this help message."
	@echo ""
	@echo "Prerequisites:"
	@echo "  - Go (version 1.22+ recommended)"
	@echo "  - Homebrew ('brew' command must be in your PATH)"