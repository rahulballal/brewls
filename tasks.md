# Brewls CLI Development Tasks

This document outlines the tasks required to develop the `brewls` CLI, based on the `Gemini.md` project documentation.

## Core Functionality (Based on Project Overview & Key Technologies)

- [ ] **T1: Implement Homebrew Command Execution:**
    - [ ] T1.1: Execute `brew list | xargs brew info --json=v2` to get detailed JSON information for all installed formulae and casks.
- [ ] **T2: Parse Homebrew JSON Output:**
    - [ ] T2.1: Define Go structs to unmarshal the JSON output from `brew info --json=v2`.
    - [ ] T2.2: Handle potential errors during JSON parsing.
- [ ] **T3: Process and Format Data:**
    - [ ] T3.1: Extract Name, Installed Version.
    - [ ] T3.2: **Determine "Installed By" information for each package:**
        - [ ] Build a reverse dependency graph: for each package, list packages that depend on it.
        - [ ] Identify root-level packages (those not depended upon by any other installed package).
        - [ ] Mark root-level packages with an asterisk (*).
    - [ ] T3.3: Generate a clear and readable tabular output for the CLI, including the "Installed By" column.

## Setup and Installation

- [ ] **T4: Local Development Setup:**
    - [ ] T4.1: Ensure Go (1.22+) and Homebrew are installed on the development machine.
    - [ ] T4.2: Initialize Go module (`go mod init brewls`).
    - [ ] T4.3: Set up the recommended directory structure (`cmd/brewls`, `internal`, `pkg`).
    - [ ] T4.4: Create basic `main.go` file.
- [ ] **T5: `go install` Integration:**
    - [ ] T5.1: Ensure the application can be installed via `go install ./cmd/brewls`.

## Testing

- [ ] **T6: Unit Tests:**
    - [ ] T6.1: Write unit tests for Homebrew command execution (mocking external commands if necessary).
    - [ ] T6.2: Write unit tests for JSON parsing logic.
    - [ ] T6.3: Write unit tests for data processing and formatting logic, including the "Installed By" column.
    - [ ] T6.4: Ensure all tests can be run using `go test ./...`.

## Code Quality

- [ ] **T7: Code Formatting:**
    - [ ] T7.1: Ensure all Go code adheres to `gofmt` standards.
- [ ] **T8: Linting:**
    - [ ] T8.1: Integrate `golangci-lint` into the development workflow to catch potential issues.

## Deployment

- [ ] **T9: macOS Binary Build:**
    - [ ] T9.1: Implement a process to build `brewls` binaries specifically for macOS.
- [ ] **T10: Homebrew Formula Creation (Long-term):**
    - [ ] T10.1: Create a Homebrew formula for `brewls`.
    - [ ] T10.2: Ensure compatibility with Homebrew guidelines.

## Project Management & Collaboration

- [ ] **T11: Feature Flag Implementation:**
    - [ ] T11.1: Incorporate feature flags for new features as per trunk-based development.
- [ ] **T12: GitHub Repository Setup (Future):**
    - [ ] T12.1: Create a GitHub repository for the project.
    - [ ] T12.2: Make the project open-source.
    - [ ] T12.3: Establish a pull request review process.
    - [ ] T12.4: Define the process for submitting bug reports or feature requests via GitHub discussions or pull requests.