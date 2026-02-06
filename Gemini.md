# Brewls Project

## 1. Project Overview
This is a CLI that mimics the 'ls' command of bash but presents the following details for installed Homebrew packages: Name of the package, installed version, and its dependencies. It also includes packages installed as casks.

## 2. Goals
- Focus on the v1 iteration of the project, ensuring core functionality is robust and stable.
- Evaluate the possibility of publishing the project to GitHub as an open-source initiative.

## 3. Key Technologies
- **Language:** Go
- **Frameworks/Libraries:** Primarily Go's standard library, particularly `os/exec` for command execution (e.g., `brew list`, `brew info`) and `encoding/json` for processing Homebrew's JSON output.
- **Database:** None (CLI application)
- **Other Tools:** Homebrew (as an external dependency for information retrieval)

## 4. Directory Structure
- `cmd/`: Main applications for the project. Each application should have its own directory (e.g., `cmd/brewls/main.go`).
- `internal/`: Private application code that you don't want other applications or projects to import. This code is typically specific to this project.
- `pkg/`: Library code that's safe to use by external applications. If others can import your code, put it here.
- `api/`: (Optional) API definitions (e.g., OpenAPI specs, Protocol Buffers).
- `build/`: (Optional) Packaging and CI/CD related assets.
- `docs/`: (Optional) Project documentation.
- `scripts/`: (Optional) Various helper scripts.
- `test/`: (Optional) External test apps and test data.

## 5. Setup/Installation
Instructions on how to set up and install the project.

1.  **Prerequisites:**
    - Go (version 1.22 or higher recommended)
    - Homebrew (https://brew.sh/)

2.  **Clone the repository:**
    \`\`\`bash
    # The repository URL will be provided once the project is hosted on GitHub.
    # For local development, ensure you are in the project root.
    \`\`\`

3.  **Install the application:**
    \`\`\`bash
    go install ./cmd/brewls
    \`\`\`
    This will install the `brewls` executable to your Go bin directory, making it available in your PATH.

4.  **Running the application:**
    \`\`\`bash
    brewls
    \`\`\`

## 6. Running Tests
Explain how to run tests for the project.

\`\`\`bash
go test ./...
\`\`\`

## 7. Deployment
Instructions on how to deploy the application.

Long term, the goal is to make `brewls` available via Homebrew. This implies:
- Building and packaging binaries specifically for macOS.
- Creating and maintaining a Homebrew formula.
- Ensuring compatibility and adherence to Homebrew's guidelines.

## 8. Contribution Guidelines
Contributions to `brewls` are welcome! Please follow these guidelines:

1.  **Branching Strategy:** We utilize a Trunk-Based Development approach. All development occurs directly on the `main` branch.
2.  **Pull Request Process:**
    - Features should be developed behind feature flags to allow safe merging into `main`.
    - Pull requests should be descriptive and clearly outline the changes.
    - Code will be reviewed before merging.
3.  **Coding Style Guidelines:**
    - All Go code must be formatted using `gofmt`.
    - We recommend running `golangci-lint` to catch potential issues.
4.  **Issue Tracking:** We currently do not use a formal issue tracker. Please submit bug reports or feature requests via GitHub discussions or pull requests.
