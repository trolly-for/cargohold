package redact_test

import (
	"strings"
	"testing"

	"cargohold/internal/redact"
)

func TestValueAlwaysReturnsMask(t *testing.T) {
	r := redact.New()
	if got := r.Value("super-secret"); got != "********" {
		t.Fatalf("expected mask, got %q", got)
	}
	if got := r.Value(""); got != "********" {
		t.Fatalf("expected mask for empty string, got %q", got)
	}
}

func TestNewWithMask(t *testing.T) {
	r := redact.NewWithMask("[REDACTED]")
	if got := r.Value("anything"); got != "[REDACTED]" {
		t.Fatalf("expected custom mask, got %q", got)
	}
}

func TestNewWithEmptyMaskFallsBackToDefault(t *testing.T) {
	r := redact.NewWithMask("")
	if got := r.Mask(); got != "********" {
		t.Fatalf("expected default mask, got %q", got)
	}
}

func TestPartialRevealsEnds(t *testing.T) {
	r := redact.New()
	got := r.Partial("abcdefghij", 2)
	if !strings.HasPrefix(got, "ab") {
		t.Fatalf("expected prefix 'ab', got %q", got)
	}
	if !strings.HasSuffix(got, "ij") {
		t.Fatalf("expected suffix 'ij', got %q", got)
	}
	if strings.Contains(got, "cdefgh") {
		t.Fatalf("middle should be masked, got %q", got)
	}
}

func TestPartialShortValueReturnsMask(t *testing.T) {
	r := redact.New()
	got := r.Partial("ab", 2)
	if got != "********" {
		t.Fatalf("expected full mask for short value, got %q", got)
	}
}

func TestPartialZeroNReturnsMask(t *testing.T) {
	r := redact.New()
	if got := r.Partial("secretvalue", 0); got != "********" {
		t.Fatalf("expected mask for n=0, got %q", got)
	}
}

func TestMapMasksAllValues(t *testing.T) {
	r := redact.New()
	input := map[string]string{
		"DB_PASSWORD": "hunter2",
		"API_KEY":     "abc123",
		"TOKEN":       "xyz",
	}
	out := r.Map(input)
	if len(out) != len(input) {
		t.Fatalf("expected %d keys, got %d", len(input), len(out))
	}
	for k, v := range out {
		if v != "********" {
			t.Errorf("key %q: expected mask, got %q", k, v)
		}
	}
}

func TestMapDoesNotMutateInput(t *testing.T) {
	r := redact.New()
	input := map[string]string{"SECRET": "plaintext"}
	_ = r.Map(input)
	if input["SECRET"] != "plaintext" {
		t.Fatal("original map was mutated")
	}
}
