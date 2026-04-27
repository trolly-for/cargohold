package cli_test

import (
	"strings"
	"testing"
)

func TestRedactMasksValue(t *testing.T) {
	run, buf := tempRunner(t)

	if err := run("init", "dev"); err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := run("set", "dev", "API_KEY", "supersecret"); err != nil {
		t.Fatalf("set: %v", err)
	}

	buf.Reset()
	if err := run("redact", "dev", "API_KEY"); err != nil {
		t.Fatalf("redact: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, "supersecret") {
		t.Fatalf("plaintext value leaked in output: %q", out)
	}
	if !strings.Contains(out, "********") {
		t.Fatalf("expected mask in output, got: %q", out)
	}
}

func TestRedactPartialShowsHint(t *testing.T) {
	run, buf := tempRunner(t)

	if err := run("init", "dev"); err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := run("set", "dev", "TOKEN", "abcdefghij"); err != nil {
		t.Fatalf("set: %v", err)
	}

	buf.Reset()
	if err := run("redact", "dev", "TOKEN", "--partial"); err != nil {
		t.Fatalf("redact --partial: %v", err)
	}

	out := buf.String()
	if strings.Contains(out, "abcdefghij") {
		t.Fatalf("full plaintext should not appear: %q", out)
	}
	if !strings.HasPrefix(strings.TrimSpace(strings.SplitN(out, "=", 2)[1]), "ab") {
		t.Fatalf("expected partial hint starting with 'ab', got: %q", out)
	}
}

func TestRedactMissingKeyErrors(t *testing.T) {
	run, _ := tempRunner(t)

	if err := run("init", "dev"); err != nil {
		t.Fatalf("init: %v", err)
	}

	err := run("redact", "dev", "MISSING_KEY")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRedactTooFewArgsErrors(t *testing.T) {
	run, _ := tempRunner(t)

	err := run("redact", "dev")
	if err == nil {
		t.Fatal("expected error for too few args")
	}
	if !strings.Contains(err.Error(), "usage") {
		t.Fatalf("unexpected error: %v", err)
	}
}
