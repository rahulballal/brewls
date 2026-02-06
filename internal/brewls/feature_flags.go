package brewls

import (
"os"
"strings"
)

const featureSortOutput = "sort-output"

// FeatureEnabled reports whether a feature flag is enabled via BREWLS_FEATURES.
func FeatureEnabled(name string) bool {
name = strings.TrimSpace(name)
if name == "" {
return false
}

raw := os.Getenv("BREWLS_FEATURES")
if raw == "" {
return false
}

for _, feature := range strings.Split(raw, ",") {
if strings.TrimSpace(feature) == name {
return true
}
}

return false
}
