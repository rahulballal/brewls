package brewls_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"

	"brewls/internal/brewls"
)

// Mock for exec.Command
func mockExecCommand(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

// TestHelperProcess is not a real test. It's a helper for mocking.
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)

	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, `No command`+"\n")
		os.Exit(2)
	}

	cmd := args[0]
	switch cmd {
	case "bash":
		// Expecting "bash", "-c", "brew list | xargs brew info --json=v2"
		if len(args) < 3 || args[1] != "-c" {
			fmt.Fprintf(os.Stderr, `Expected -c argument`+"\n")
			os.Exit(2)
		}
		shellCommand := args[2]
		if shellCommand == "brew list | xargs brew info --json=v2" {
			mockOutput := os.Getenv("MOCK_BREW_OUTPUT")
			mockError := os.Getenv("MOCK_BREW_ERROR")
			if mockError != "" {
				fmt.Fprint(os.Stderr, mockError) // Changed to stderr as stderrBuf is used in real code
				os.Exit(1)
			}
			fmt.Fprint(os.Stdout, mockOutput)
			return
		} else {
			fmt.Fprintf(os.Stderr, `Unexpected shell command: %s`+"\n", shellCommand)
			os.Exit(2)
		}
	default:
		fmt.Fprintf(os.Stderr, `Unknown command: %s`+"\n", cmd)
		os.Exit(2)
	}
}

func TestExecuteBrewInfoCommand(t *testing.T) {
	// Save original ExecCommand and LookPath and restore after test
	oldExecCommand := brewls.ExecCommand
	oldLookPath := brewls.LookPath
	defer func() {
		brewls.ExecCommand = oldExecCommand
		brewls.LookPath = oldLookPath
	}()

	tests := []struct {
		name              string
		mockLookPathError error // Error to return from LookPath("brew")
		mockOutput        string
		mockError         string // Error to return from command execution
		expectedErr       bool
		expectedOut       string
		expectedErrMsg    string // Expected substring in error message
	}{
		{
			name:              "successful execution",
			mockLookPathError: nil,
			mockOutput:        `{"formulae": [], "casks": []}`,
			mockError:         "",
			expectedErr:       false,
			expectedOut:       `{"formulae": [], "casks": []}`,
			expectedErrMsg:    "",
		},
		{
			name:              "brew command not found",
			mockLookPathError: errors.New("not found in path"),
			mockOutput:        "",
			mockError:         "",
			expectedErr:       true,
			expectedOut:       "",
			expectedErrMsg:    "Homebrew 'brew' command not found in PATH",
		},
		{
			name:              "command returns error",
			mockLookPathError: nil,
			mockOutput:        "",
			mockError:         "brew command failed",
			expectedErr:       true,
			expectedOut:       "",
			expectedErrMsg:    "command finished with error",
		},
		{
			name:              "empty output",
			mockLookPathError: nil,
			mockOutput:        "",
			mockError:         "",
			expectedErr:       false,
			expectedOut:       "",
			expectedErrMsg:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock LookPath
			brewls.LookPath = func(file string) (string, error) {
				if file == "brew" {
					return "", tt.mockLookPathError
				}
				return oldLookPath(file) // Fallback to original for other LookPath calls if any
			}

			// Mock ExecCommand only if LookPath succeeds
			brewls.ExecCommand = func(command string, args ...string) *exec.Cmd {
				cmd := mockExecCommand(command, args...)
				cmd.Env = append(cmd.Env, "MOCK_BREW_OUTPUT="+tt.mockOutput)
				cmd.Env = append(cmd.Env, "MOCK_BREW_ERROR="+tt.mockError)
				return cmd
			}

			output, err := brewls.ExecuteBrewInfoCommand()

			if tt.expectedErr {
				if err == nil {
					t.Errorf("Expected an error but got none")
				} else if !strings.Contains(err.Error(), tt.expectedErrMsg) {
					t.Errorf("Expected error message to contain %q, but got %q", tt.expectedErrMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error but got: %v", err)
				}
				if output != tt.expectedOut {
					t.Errorf("Expected output %q, but got %q", tt.expectedOut, output)
				}
			}
		})
	}
}

func TestParseBrewInfoJSON(t *testing.T) {
	tests := []struct {
		name        string
		jsonInput   string
		expected    *brewls.BrewInfo
		expectedErr bool
	}{
		{
			name:        "valid empty JSON",
			jsonInput:   `{"formulae": [], "casks": []}`,
			expected:    &brewls.BrewInfo{Formulae: []brewls.Formula{}, Casks: []brewls.Cask{}},
			expectedErr: false,
		},
		{
			name: "valid JSON with formula and cask",
			jsonInput: `{
				"formulae": [
					{
						"name": "test-formula",
						"installed": [{"version": "1.0.0", "runtime_dependencies": [{"full_name": "dep1"}]}],
						"dependencies": ["build-dep1"]
					}
				],
				"casks": [
					{
						"token": "test-cask",
						"name": ["Test Cask App"],
						"version": "2.0.0",
						"installed": "2.0.0"
					}
				]
			}`,
			expected: &brewls.BrewInfo{
				Formulae: []brewls.Formula{
					{
						Name: "test-formula",
						Installed: []brewls.Installed{
							{
								Version: "1.0.0",
								RuntimeDependencies: []brewls.RuntimeDependency{
									{FullName: "dep1"},
								},
							},
						},
						Dependencies: []string{"build-dep1"},
					},
				},
				Casks: []brewls.Cask{
					{
						Token:     "test-cask",
						Name:      []string{"Test Cask App"},
						Version:   "2.0.0",
						Installed: "2.0.0",
					},
				},
			},
			expectedErr: false,
		},
		{
			name:        "invalid JSON",
			jsonInput:   `{"formulae": [{ "name": "invalid", ]}`,
			expected:    nil,
			expectedErr: true,
		},
		{
			name: "partially valid JSON (missing fields)",
			jsonInput: `{
				"formulae": [
					{ "name": "partial-formula" }
				],
				"casks": []
			}`,
			expected: &brewls.BrewInfo{
				Formulae: []brewls.Formula{
					{
						Name:         "partial-formula",
						Installed:    nil, // Expected nil for missing slice from Unmarshal
						Dependencies: nil, // Expected nil for missing slice from Unmarshal
					},
				},
				Casks: []brewls.Cask{},
			},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := brewls.ParseBrewInfoJSON(tt.jsonInput)

			if tt.expectedErr {
				if err == nil {
					t.Errorf("Expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect an error but got: %v", err)
				}
				if !reflect.DeepEqual(got, tt.expected) {
					t.Errorf("Expected BrewInfo %v, but got %v", tt.expected, got)
				}
			}
		})
	}
}

func TestFormatBrewOutput(t *testing.T) {
	originalFlags := os.Getenv(brewls.FeatureFlagsEnv)
	defer func() {
		_ = os.Setenv(brewls.FeatureFlagsEnv, originalFlags)
	}()
	_ = os.Setenv(brewls.FeatureFlagsEnv, "")

	tests := []struct {
		name           string
		brewInfo       *brewls.BrewInfo
		expectedOutput string
	}{
		{
			name:     "empty brewInfo",
			brewInfo: &brewls.BrewInfo{Formulae: []brewls.Formula{}, Casks: []brewls.Cask{}},
			expectedOutput: `
--- Homebrew Formulae ---
+------+---------+--------------+
| NAME | VERSION | INSTALLED BY |
+------+---------+--------------+
+------+---------+--------------+

--- Homebrew Casks ---
+------+---------+--------------+
| NAME | VERSION | INSTALLED BY |
+------+---------+--------------+
+------+---------+--------------+
`,
		},
		{
			name: "single formula and cask",
			brewInfo: &brewls.BrewInfo{
				Formulae: []brewls.Formula{
					{
						Name: "test-formula",
						Installed: []brewls.Installed{
							{
								Version: "1.0.0",
								// Note: InstalledOnRequest is false by default, so test-formula will not be a root here
								RuntimeDependencies: []brewls.RuntimeDependency{
									{FullName: "dep1"},
									{FullName: "dep2"},
								},
							},
						},
						Dependencies: []string{"build-dep1"},
					},
				},
				Casks: []brewls.Cask{
					{
						Token:     "test-cask",
						Name:      []string{"Test Cask App"},
						Installed: "2.0.0",
					},
				},
			},
			expectedOutput: `
--- Homebrew Formulae ---
+--------------+---------+--------------+
| NAME         | VERSION | INSTALLED BY |
+--------------+---------+--------------+
| test-formula | 1.0.0   |              |
+--------------+---------+--------------+

--- Homebrew Casks ---
+-----------------+---------+--------------+
| NAME            | VERSION | INSTALLED BY |
+-----------------+---------+--------------+
| Test Cask App * | 2.0.0   |              |
+-----------------+---------+--------------+
`,
		},
		{
			name: "formula with no installed version",
			brewInfo: &brewls.BrewInfo{
				Formulae: []brewls.Formula{
					{
						Name:         "no-install-formula",
						Installed:    nil,
						Dependencies: []string{"depA"},
					},
				},
				Casks: []brewls.Cask{},
			},
			expectedOutput: `
--- Homebrew Formulae ---
+--------------------+---------+--------------+
| NAME               | VERSION | INSTALLED BY |
+--------------------+---------+--------------+
| no-install-formula | N/A     |              |
+--------------------+---------+--------------+

--- Homebrew Casks ---
+------+---------+--------------+
| NAME | VERSION | INSTALLED BY |
+------+---------+--------------+
+------+---------+--------------+
`,
		},
		{
			name: "formula with multiple installed versions (should use last)",
			brewInfo: &brewls.BrewInfo{
				Formulae: []brewls.Formula{
					{
						Name: "multi-version",
						Installed: []brewls.Installed{
							{Version: "1.0.0"},
							{Version: "1.1.0"},
						},
						Dependencies: nil,
					},
				},
				Casks: []brewls.Cask{},
			},
			expectedOutput: `
--- Homebrew Formulae ---
+---------------+---------+--------------+
| NAME          | VERSION | INSTALLED BY |
+---------------+---------+--------------+
| multi-version | 1.1.0   |              |
+---------------+---------+--------------+

--- Homebrew Casks ---
+------+---------+--------------+
| NAME | VERSION | INSTALLED BY |
+------+---------+--------------+
+------+---------+--------------+
`,
		},
		{
			name: "complex mock brewInfo with reverse dependencies and root markers",
			brewInfo: &brewls.BrewInfo{
				Formulae: []brewls.Formula{
					{ // PackageA: Root, depends on B
						Name: "packageA",
						Installed: []brewls.Installed{
							{Version: "1.0.0", InstalledOnRequest: true, RuntimeDependencies: []brewls.RuntimeDependency{{FullName: "packageB"}}},
						},
						Dependencies: []string{}, // Build deps
					},
					{ // PackageB: Not root, depended on by A and D, depends on C
						Name: "packageB",
						Installed: []brewls.Installed{
							{Version: "1.1.0", InstalledOnRequest: false, RuntimeDependencies: []brewls.RuntimeDependency{{FullName: "packageC"}}},
						},
						Dependencies: []string{},
					},
					{ // PackageC: Not root, depended on by B
						Name: "packageC",
						Installed: []brewls.Installed{
							{Version: "1.2.0", InstalledOnRequest: false},
						},
						Dependencies: []string{},
					},
					{ // PackageD: Root, depends on B
						Name: "packageD",
						Installed: []brewls.Installed{
							{Version: "2.0.0", InstalledOnRequest: true, RuntimeDependencies: []brewls.RuntimeDependency{{FullName: "packageB"}}},
						},
						Dependencies: []string{},
					},
					{ // PackageF: Root, no deps
						Name: "packageF",
						Installed: []brewls.Installed{
							{Version: "3.0.0", InstalledOnRequest: true},
						},
						Dependencies: []string{},
					},
				},
				Casks: []brewls.Cask{
					{ // CaskE: Root, no explicit dependencies in the JSON mock for simplicity
						Token:     "caskE",
						Name:      []string{"Cask E App"},
						Installed: "1.0.0",
					},
				},
			},
			expectedOutput: `
--- Homebrew Formulae ---
+------------+---------+--------------------+
| NAME       | VERSION | INSTALLED BY       |
+------------+---------+--------------------+
| packageA * | 1.0.0   |                    |
| packageB   | 1.1.0   | packageA, packageD |
| packageC   | 1.2.0   | packageB           |
| packageD * | 2.0.0   |                    |
| packageF * | 3.0.0   |                    |
+------------+---------+--------------------+

--- Homebrew Casks ---
+--------------+---------+--------------+
| NAME         | VERSION | INSTALLED BY |
+--------------+---------+--------------+
| Cask E App * | 1.0.0   |              |
+--------------+---------+--------------+
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Need to call BuildReverseDependencyGraph before formatting to populate InstalledBy and IsRoot
			// A copy is made so that the original tt.brewInfo is not modified,
			// which can cause unexpected side effects in subsequent tests.
			brewInfoCopy := *tt.brewInfo
			brewls.BuildReverseDependencyGraph(&brewInfoCopy)

			var buf bytes.Buffer
			brewls.FormatBrewOutput(&brewInfoCopy, &buf) // Pass bytes.Buffer directly

			normalizedExpected := strings.TrimSpace(strings.ReplaceAll(tt.expectedOutput, "\r\n", "\n"))
			normalizedActual := strings.TrimSpace(strings.ReplaceAll(buf.String(), "\r\n", "\n"))

			if normalizedActual != normalizedExpected {
				t.Errorf("Expected output:\n%q\nGot:\n%q", normalizedExpected, normalizedActual)
			}
		})
	}
}

func TestUniqueAndSortStrings(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "nil slice",
			input:    nil,
			expected: []string{}, // Expect empty slice now
		},
		{
			name:     "no duplicates, unsorted",
			input:    []string{"c", "a", "b"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "with duplicates",
			input:    []string{"a", "b", "a", "c", "b"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "already sorted and unique",
			input:    []string{"alpha", "beta", "gamma"},
			expected: []string{"alpha", "beta", "gamma"},
		},
		{
			name:     "mixed case duplicates (should treat as different)",
			input:    []string{"Apple", "apple", "Orange"},
			expected: []string{"Apple", "Orange", "apple"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := brewls.UniqueAndSortStrings(tt.input)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("UniqueAndSortStrings(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}
