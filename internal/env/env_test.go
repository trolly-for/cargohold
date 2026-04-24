package env_test

import (
	"testing"

	"github.com/cargohold/cargohold/internal/env"
)

func TestValidate(t *testing.T) {
	valid := []string{"development", "staging", "production", "Development", " staging "}
	for _, name := range valid {
		if err := env.Validate(name); err != nil {
			t.Errorf("expected %q to be valid, got error: %v", name, err)
		}
	}

	invalid := []string{"", "dev", "prod", "test", "local"}
	for _, name := range invalid {
		if err := env.Validate(name); err == nil {
			t.Errorf("expected %q to be invalid, but got no error", name)
		}
	}
}

func TestNormalize(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"Development", "development"},
		{" staging ", "staging"},
		{"PRODUCTION", "production"},
		{"development", "development"},
	}
	for _, tc := range cases {
		got := env.Normalize(tc.input)
		if got != tc.want {
			t.Errorf("Normalize(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestFromEnvVar(t *testing.T) {
	t.Run("defaults to development", func(t *testing.T) {
		t.Setenv("CARGOHOLD_ENV", "")
		if got := env.FromEnvVar(); got != "development" {
			t.Errorf("expected development, got %q", got)
		}
	})

	t.Run("reads from env var", func(t *testing.T) {
		t.Setenv("CARGOHOLD_ENV", "staging")
		if got := env.FromEnvVar(); got != "staging" {
			t.Errorf("expected staging, got %q", got)
		}
	})

	t.Run("normalizes env var value", func(t *testing.T) {
		t.Setenv("CARGOHOLD_ENV", " PRODUCTION ")
		if got := env.FromEnvVar(); got != "production" {
			t.Errorf("expected production, got %q", got)
		}
	})
}

func TestIsProduction(t *testing.T) {
	if !env.IsProduction("production") {
		t.Error("expected production to be production")
	}
	if !env.IsProduction("PRODUCTION") {
		t.Error("expected PRODUCTION to be production")
	}
	if env.IsProduction("staging") {
		t.Error("expected staging not to be production")
	}
}
