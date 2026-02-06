# Brewls CLI Development Tasks

This document outlines the tasks required to develop the `brewls` CLI, based on the `.github/copilot-instructions.md` project documentation.

## Core Functionality (Based on Project Overview & Key Technologies)

- [x] **T1: Implement Homebrew Command Execution:**
    - [x] T1.1: Execute `brew list | xargs brew info --json=v2` to get detailed JSON information for all installed formulae and casks.
- [x] **T2: Parse Homebrew JSON Output:**
    - [x] T2.1: Define Go structs to unmarshal the JSON output from `brew info --json=v2`.
    - [x] T2.2: Handle potential errors during JSON parsing.
- [x] **T3: Process and Format Data:**
    - [x] T3.1: Extract Name, Installed Version.
    - [x] T3.2: **Determine "Installed By" information for each package:**
        - [x] Build a reverse dependency graph: for each package, list packages that depend on it.
        - [x] Identify root-level packages (those not depended upon by any other installed package).
        - [x] Mark root-level packages with an asterisk (*).
    - [x] T3.3: Generate a clear and readable tabular output for the CLI, including the "Installed By" column.

## Setup and Installation

- [x] **T4: Local Development Setup:**
    - [x] T4.1: Ensure Go (1.22+) and Homebrew are installed on the development machine.
    - [x] T4.2: Initialize Go module (`go mod init brewls`).
    - [x] T4.3: Set up the recommended directory structure (`cmd/brewls`, `internal`, `pkg`).
    - [x] T4.4: Create basic `main.go` file.
- [x] **T5: `go install` Integration:**
    - [x] T5.1: Ensure the application can be installed via `go install ./cmd/brewls`.

## Testing

- [x] **T6: Unit Tests:**
    - [x] T6.1: Write unit tests for Homebrew command execution (mocking external commands if necessary).
    - [x] T6.2: Write unit tests for JSON parsing logic.
    - [x] T6.3: Write unit tests for data processing and formatting logic, including the "Installed By" column.
    - [x] T6.4: Ensure all tests can be run using `go test ./...`.

## Code Quality

- [x] **T7: Code Formatting:**
    - [x] T7.1: Ensure all Go code adheres to `gofmt` standards.
- [x] **T8: Linting:**
    - [x] T8.1: Integrate `golangci-lint` into the development workflow to catch potential issues.

## Deployment

- [x] **T9: macOS Binary Build:**
    - [x] T9.1: Implement a process to build `brewls` binaries specifically for macOS.
- [x] **T10: Homebrew Formula Creation (Long-term):**
    - [x] T10.1: Create a Homebrew formula for `brewls`.
    - [x] T10.2: Ensure compatibility with Homebrew guidelines.

## Project Management & Collaboration

- [ ] **T11: Feature Flag Implementation:**
    - [x] T11.1: Incorporate feature flags for new features as per trunk-based development.
- [x] **T12: GitHub Repository Setup (Future):**
    - [x] T12.1: Create a GitHub repository for the project.
    - [x] T12.2: Make the project open-source.
    - [x] T12.3: Establish a pull request review process.