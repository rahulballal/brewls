# Copilot Instructions for brewls

## Project Overview
`brewls` is a Go CLI that extends Homebrew's `brew ls` output. It lists installed formulae and casks with name, installed version, and reverse dependency info ("Installed By"). Root packages (installed directly by the user) are marked with an asterisk (`*`).

## Goals
- Focus on v1 core functionality and stability.
- Keep the project ready for eventual open-source publication.

## Tech Stack
- Language: Go (1.22+)
- Libraries: Go standard library preferred
- Homebrew integration: `os/exec` for `brew list`/`brew info`, parse JSON with `encoding/json`

## Key Behaviors
- Fetch installed packages via `brew list | xargs brew info --json=v2`
- Build reverse dependency graph
- Render clear tabular output

## Repository Structure
- `cmd/` for the CLI entrypoint (`cmd/brewls/main.go`)
- `internal/` for app-specific code
- `pkg/` for reusable library code
- Optional: `docs/`, `scripts/`, `test/`, `build/`

## Dev Workflow
- Run tests with `go test ./...`
- Format code with `gofmt`
- Prefer trunk-based development on `main`
- Use feature flags for safe merging when needed
