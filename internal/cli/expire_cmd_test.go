package cli_test

import (
	"strings"
	"testing"
)

func TestExpireSetAndGet(t *testing.T) {
	r, _ := tempRunner(t)

	if err := r.Run([]string{"init", "dev"}); err != nil {
		t.Fatalf("init: %v", err)
	}

	if err := r.Run([]string{"expire", "set", "dev", "48h"}); err != nil {
		t.Fatalf("expire set: %v", err)
	}

	out := captureOutput(t, r, []string{"expire", "get", "dev"})
	if !strings.Contains(out, "valid") {
		t.Errorf("expected 'valid' in output, got: %s", out)
	}
}

func TestExpireGetNotSet(t *testing.T) {
	r, _ := tempRunner(t)
	if err := r.Run([]string{"init", "dev"}); err != nil {
		t.Fatalf("init: %v", err)
	}

	out := captureOutput(t, r, []string{"expire", "get", "dev"})
	if !strings.Contains(out, "none") {
		t.Errorf("expected 'none' in output, got: %s", out)
	}
}

func TestExpireClear(t *testing.T) {
	r, _ := tempRunner(t)
	if err := r.Run([]string{"init", "dev"}); err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := r.Run([]string{"expire", "set", "dev", "1h"}); err != nil {
		t.Fatalf("expire set: %v", err)
	}
	if err := r.Run([]string{"expire", "clear", "dev"}); err != nil {
		t.Fatalf("expire clear: %v", err)
	}

	out := captureOutput(t, r, []string{"expire", "get", "dev"})
	if !strings.Contains(out, "none") {
		t.Errorf("expected 'none' after clear, got: %s", out)
	}
}

func TestExpireSetNegativeDurationErrors(t *testing.T) {
	r, _ := tempRunner(t)
	if err := r.Run([]string{"init", "dev"}); err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := r.Run([]string{"expire", "set", "dev", "-1h"}); err == nil {
		t.Error("expected error for negative duration")
	}
}

func TestExpireUnknownSubcommandErrors(t *testing.T) {
	r, _ := tempRunner(t)
	if err := r.Run([]string{"init", "dev"}); err != nil {
		t.Fatalf("init: %v", err)
	}
	if err := r.Run([]string{"expire", "purge", "dev"}); err == nil {
		t.Error("expected error for unknown subcommand")
	}
}

func TestExpireTooFewArgsErrors(t *testing.T) {
	r, _ := tempRunner(t)
	if err := r.Run([]string{"expire", "set"}); err == nil {
		t.Error("expected error for too few args")
	}
}
