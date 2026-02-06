package brewls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sort"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table" // New import for go-pretty
)

// BrewInfo represents the top-level structure of the JSON output from brew info --json=v2
type BrewInfo struct {
	Formulae []Formula `json:"formulae"`
	Casks    []Cask    `json:"casks"`
}

// Formula represents a Homebrew formula
type Formula struct {
	Name       string      `json:"name"`
	Installed  []Installed `json:"installed"`
	Dependencies []string `json:"dependencies"` // Build dependencies
	InstalledBy  []string // New field: packages that depend on this one
	IsRoot       bool     // New field: true if this is a top-level package (not depended on)
}

// Installed represents an installed version of a formula
type Installed struct {
	Version             string             `json:"version"`
	RuntimeDependencies []RuntimeDependency `json:"runtime_dependencies"`
	InstalledOnRequest bool `json:"installed_on_request"` // This field is crucial for identifying root packages
}

// RuntimeDependency represents a runtime dependency of an installed formula
type RuntimeDependency struct {
	FullName string `json:"full_name"`
	Version  string `json:"version"`
}

// Cask represents a Homebrew cask
type Cask struct {
	Token   string `json:"token"`
	Name    []string `json:"name"` // Display name, if available
	Version string `json:"version"`
	Installed string `json:"installed"` // This seems to represent the installed version for casks
	InstalledBy []string // New field: packages that depend on this one (less common for casks)
	IsRoot      bool     // New field: true if this is a top-level package
	// Homebrew cask info often just lists depends_on for macOS versions or other casks/formulae,
	// direct dependencies are not as clear-cut as for formulae.
	// For simplicity, we'll primarily rely on formulae dependencies for now.
	// If cask-to-cask or cask-to-formula dependencies are needed,
	// the JSON structure for casks needs to be re-evaluated.
}

// execCommand is a global variable to allow mocking os/exec.Command in tests.
var ExecCommand = exec.Command // Exported for testing

// lookPath is a global variable to allow mocking os/exec.LookPath in tests.
var LookPath = exec.LookPath // Exported for testing

// ExecuteBrewInfoCommand runs the brew command and returns its JSON output as a string.
func ExecuteBrewInfoCommand() (string, error) {
	// Check if "brew" command is available
	_, err := LookPath("brew")
	if err != nil {
		return "", fmt.Errorf("Homebrew 'brew' command not found in PATH: %w. Please ensure Homebrew is installed and configured correctly.", err)
	}

	cmdString := "brew list | xargs brew info --json=v2"
	cmd := ExecCommand("bash", "-c", cmdString) // Use exported ExecCommand

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("command finished with error: %w; Stderr: %s", err, stderrBuf.String())
	}

	return stdoutBuf.String(), nil
}

// ParseBrewInfoJSON unmarshals the JSON string into a BrewInfo struct.
func ParseBrewInfoJSON(jsonInput string) (*BrewInfo, error) {
	var brewInfo BrewInfo
	err := json.Unmarshal([]byte(jsonInput), &brewInfo)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}
	return &brewInfo, nil
}

// BuildReverseDependencyGraph processes BrewInfo to determine which packages are installed by others.
// It populates the InstalledBy field for each Formula and Cask and identifies root packages.
func BuildReverseDependencyGraph(info *BrewInfo) {
	// Map to store which packages install a given package
	installedByMap := make(map[string][]string)
	
	// Collect all installed package names for quick lookup
	allInstalledPackages := make(map[string]struct{})
	for _, f := range info.Formulae {
		allInstalledPackages[f.Name] = struct{}{}
	}
	for _, c := range info.Casks {
		allInstalledPackages[c.Token] = struct{}{} // Use token for cask names
	}

	// Process Formulae dependencies
	for _, f := range info.Formulae {
		// Combine build and runtime dependencies
		var dependencies []string
		dependencies = append(dependencies, f.Dependencies...)
		if len(f.Installed) > 0 {
			for _, rd := range f.Installed[len(f.Installed)-1].RuntimeDependencies {
				dependencies = append(dependencies, rd.FullName)
			}
		}

		for _, dep := range UniqueAndSortStrings(dependencies) { // Ensure unique and sorted dependencies
			// Only consider dependencies that are actually installed
			if _, ok := allInstalledPackages[dep]; ok {
				installedByMap[dep] = append(installedByMap[dep], f.Name)
			}
		}
	}

	// Casks dependencies are not as straightforward in brew info --json=v2,
	// so for simplicity, we'll primarily rely on formulae dependencies for now.
	// If cask-to-cask or cask-to-formula dependencies are needed,
	// the JSON structure for casks needs to be re-evaluated.

	// Populate InstalledBy for Formulae and determine IsRoot
	for i := range info.Formulae {
		info.Formulae[i].InstalledBy = UniqueAndSortStrings(installedByMap[info.Formulae[i].Name])
		// A formula is a root if it was installed on request AND nothing else depends on it
		if len(info.Formulae[i].Installed) > 0 && info.Formulae[i].Installed[len(info.Formulae[i].Installed)-1].InstalledOnRequest && len(info.Formulae[i].InstalledBy) == 0 {
			info.Formulae[i].IsRoot = true
		} else {
			info.Formulae[i].IsRoot = false
		}
	}

	// Populate InstalledBy for Casks and determine IsRoot
	for i := range info.Casks {
		info.Casks[i].InstalledBy = UniqueAndSortStrings(installedByMap[info.Casks[i].Token]) // Use token for lookup
		// A cask is a root if nothing else depends on it (and it's installed, which is implied by being in the list).
		// For casks, 'installed_on_request' equivalent is not readily available in the mock JSON,
		// so we'll treat it as root if nothing depends on it.
		if len(info.Casks[i].InstalledBy) == 0 {
			info.Casks[i].IsRoot = true
		} else {
			info.Casks[i].IsRoot = false
		}
	}
}


// FormatBrewOutput generates the formatted tabular output for formulae and casks.
// It now accepts an io.Writer interface, making it more testable.
func FormatBrewOutput(brewInfo *BrewInfo, writer io.Writer) {
	// --- Process and Format Formulae ---
	fmt.Fprintln(writer, "\n--- Homebrew Formulae ---")
	
	// Create a new go-pretty table writer
	formulaeTable := table.NewWriter()
	formulaeTable.SetOutputMirror(writer) // Set the output writer
	// Changed header from "Dependencies" to "Installed By"
	formulaeTable.AppendHeader(table.Row{"Name", "Version", "Installed By"})

	for _, formula := range brewInfo.Formulae {
		installedVersion := "N/A"

		if len(formula.Installed) > 0 {
			installedVersion = formula.Installed[len(formula.Installed)-1].Version
		}
		
		// Determine display name for formulae
		displayName := formula.Name
		if formula.IsRoot {
			displayName += " *"
		}

		formulaeTable.AppendRow(table.Row{
			displayName,
			installedVersion,
			strings.Join(formula.InstalledBy, ", "), // Use the new InstalledBy field
		})
	}
	formulaeTable.Render()

	// --- Process and Format Casks ---
	fmt.Fprintln(writer, "\n--- Homebrew Casks ---")
	casksTable := table.NewWriter()
	casksTable.SetOutputMirror(writer) // Set the output writer
	// Changed header from "Dependencies" to "Installed By"
	casksTable.AppendHeader(table.Row{"Name", "Version", "Installed By"})

	for _, cask := range brewInfo.Casks {
		displayName := cask.Token
		if len(cask.Name) > 0 {
			displayName = cask.Name[0]
		}
		if cask.IsRoot {
			displayName += " *"
		}
		
		casksTable.AppendRow(table.Row{
			displayName,
			cask.Installed,
			strings.Join(cask.InstalledBy, ", "), // Use the new InstalledBy field
		})
	}
	casksTable.Render()
}

// UniqueAndSortStrings is a helper function to remove duplicates and sort strings.
// It now returns an empty slice instead of nil for empty input.
func UniqueAndSortStrings(s []string) []string { // Exported
    if len(s) == 0 {
        return []string{} // Return empty slice instead of nil
    }
    seen := make(map[string]struct{})
    var result []string
    for _, val := range s {
        if _, ok := seen[val]; !ok {
            seen[val] = struct{}{}
            result = append(result, val)
        }
    }
    sort.Strings(result)
    return result
}
