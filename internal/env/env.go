package env

import (
	"fmt"
	"os"
	"strings"
)

// ValidEnvironments defines the recognized environment names.
var ValidEnvironments = []string{"development", "staging", "production"}

// Validate checks whether the given environment name is recognized.
func Validate(name string) error {
	normalized := Normalize(name)
	for _, e := range ValidEnvironments {
		if normalized == e {
			return nil
		}
	}
	return fmt.Errorf("unknown environment %q: must be one of %s",
		name, strings.Join(ValidEnvironments, ", "))
}

// Normalize returns the lowercase, trimmed version of an environment name.
func Normalize(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// FromEnvVar reads the current environment from the CARGOHOLD_ENV variable.
// It falls back to "development" when the variable is not set.
func FromEnvVar() string {
	if v := os.Getenv("CARGOHOLD_ENV"); v != "" {
		return Normalize(v)
	}
	return "development"
}

// IsProduction returns true when name resolves to the production environment.
func IsProduction(name string) bool {
	return Normalize(name) == "production"
}
