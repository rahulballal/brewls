package brewls

import (
	"os"
	"reflect"
	"testing"
)

func TestParseFeatureFlags(t *testing.T) {
	tests := []struct {
		name     string
		raw      string
		expected map[string]struct{}
	}{
		{
			name:     "empty",
			raw:      "",
			expected: map[string]struct{}{},
		},
		{
			name: "trims and lowercases",
			raw:  " Foo,bar , BAZ ",
			expected: map[string]struct{}{
				"foo": {},
				"bar": {},
				"baz": {},
			},
		},
		{
			name: "drops blanks and duplicates",
			raw:  "alpha,,ALPHA, ,beta",
			expected: map[string]struct{}{
				"alpha": {},
				"beta":  {},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseFeatureFlags(tt.raw)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("ParseFeatureFlags(%q) = %v, want %v", tt.raw, got, tt.expected)
			}
		})
	}
}

func TestIsFeatureEnabled(t *testing.T) {
	original := os.Getenv(FeatureFlagsEnv)
	defer func() {
		_ = os.Setenv(FeatureFlagsEnv, original)
	}()

	_ = os.Setenv(FeatureFlagsEnv, "alpha, beta")

	if !IsFeatureEnabled("alpha") {
		t.Fatalf("expected alpha to be enabled")
	}
	if !IsFeatureEnabled(" BETA ") {
		t.Fatalf("expected beta to be enabled")
	}
	if IsFeatureEnabled("gamma") {
		t.Fatalf("expected gamma to be disabled")
	}
}
