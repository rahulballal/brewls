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
	@echo "  help    - Displays this help message."
	@echo ""
	@echo "Prerequisites:"
	@echo "  - Go (version 1.22+ recommended)"
	@echo "  - Homebrew ('brew' command must be in your PATH)"