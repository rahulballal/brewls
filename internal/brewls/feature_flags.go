package brewls

import (
	"os"
	"strings"
)

// FeatureFlagsEnv is the env var used to enable feature flags.
const FeatureFlagsEnv = "BREWLS_FEATURE_FLAGS"

// IsFeatureEnabled reports whether a feature flag is enabled via env var.
func IsFeatureEnabled(flag string) bool {
	flag = normalizeFeatureFlag(flag)
	if flag == "" {
		return false
	}
	flags := ParseFeatureFlags(os.Getenv(FeatureFlagsEnv))
	_, ok := flags[flag]
	return ok
}

// ParseFeatureFlags parses a comma-separated list of feature flags.
func ParseFeatureFlags(raw string) map[string]struct{} {
	result := make(map[string]struct{})
	for _, item := range strings.Split(raw, ",") {
		name := normalizeFeatureFlag(item)
		if name == "" {
			continue
		}
		result[name] = struct{}{}
	}
	return result
}

func normalizeFeatureFlag(flag string) string {
	return strings.ToLower(strings.TrimSpace(flag))
}
