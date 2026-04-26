package cli_test

import (
	"testing"
)

func TestCopyAllKeys(t *testing.T) {
	r, cleanup := tempRunner(t)
	defer cleanup()

	if err := r.Run([]string{"init", "staging"}); err != nil {
		t.Fatalf("init staging: %v", err)
	}
	if err := r.Run([]string{"init", "production"}); err != nil {
		t.Fatalf("init production: %v", err)
	}
	if err := r.Run([]string{"set", "staging", "DB_URL", "postgres://staging"}); err != nil {
		t.Fatalf("set: %v", err)
	}
	if err := r.Run([]string{"set", "staging", "API_KEY", "abc123"}); err != nil {
		t.Fatalf("set: %v", err)
	}

	if err := r.Run([]string{"copy", "staging", "production"}); err != nil {
		t.Fatalf("copy: %v", err)
	}

	for _, key := range []string{"DB_URL", "API_KEY"} {
		if err := r.Run([]string{"get", "production", key}); err != nil {
			t.Errorf("get production %s after copy: %v", key, err)
		}
	}
}

func TestCopySelectedKey(t *testing.T) {
	r, cleanup := tempRunner(t)
	defer cleanup()

	r.Run([]string{"init", "dev"})
	r.Run([]string{"init", "staging"})
	r.Run([]string{"set", "dev", "SECRET", "s3cr3t"})
	r.Run([]string{"set", "dev", "OTHER", "other"})

	if err := r.Run([]string{"copy", "dev", "staging", "SECRET"}); err != nil {
		t.Fatalf("copy: %v", err)
	}

	if err := r.Run([]string{"get", "staging", "SECRET"}); err != nil {
		t.Errorf("expected SECRET in staging: %v", err)
	}
	if err := r.Run([]string{"get", "staging", "OTHER"}); err == nil {
		t.Error("OTHER should not have been copied")
	}
}

func TestCopySameEnvErrors(t *testing.T) {
	r, cleanup := tempRunner(t)
	defer cleanup()

	r.Run([]string{"init", "dev"})
	if err := r.Run([]string{"copy", "dev", "dev"}); err == nil {
		t.Fatal("expected error when src == dst")
	}
}

func TestCopyTooFewArgsErrors(t *testing.T) {
	r, cleanup := tempRunner(t)
	defer cleanup()

	if err := r.Run([]string{"copy", "dev"}); err == nil {
		t.Fatal("expected error for too few args")
	}
}
