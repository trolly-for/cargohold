package cli_test

import (
	"os"
	"path/filepath"
	"testing"

	"cargohold/internal/cli"
	"cargohold/internal/store"
)

func tempRunner(t *testing.T) *cli.Runner {
	t.Helper()
	dir := t.TempDir()
	s, err := store.New(filepath.Join(dir, "bundles"))
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return &cli.Runner{Store: s}
}

func TestInitCreatesBundle(t *testing.T) {
	r := tempRunner(t)
	if err := r.Init("staging", "secret"); err != nil {
		t.Fatalf("Init: %v", err)
	}
	envs, err := r.Store.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(envs) != 1 || envs[0] != "staging" {
		t.Errorf("expected [staging], got %v", envs)
	}
}

func TestInitDuplicateErrors(t *testing.T) {
	r := tempRunner(t)
	if err := r.Init("dev", "pass"); err != nil {
		t.Fatalf("first Init: %v", err)
	}
	if err := r.Init("dev", "pass"); err == nil {
		t.Fatal("expected error on duplicate init")
	}
}

func TestSetAndGet(t *testing.T) {
	r := tempRunner(t)
	const pass = "hunter2"
	if err := r.Init("dev", pass); err != nil {
		t.Fatalf("Init: %v", err)
	}
	if err := r.Set("dev", pass, "DB_URL", "postgres://localhost/dev"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	val, err := r.Get("dev", pass, "DB_URL")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if val != "postgres://localhost/dev" {
		t.Errorf("got %q, want %q", val, "postgres://localhost/dev")
	}
}

func TestGetMissingKey(t *testing.T) {
	r := tempRunner(t)
	const pass = "pass"
	if err := r.Init("dev", pass); err != nil {
		t.Fatalf("Init: %v", err)
	}
	if _, err := r.Get("dev", pass, "MISSING"); err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestSetWrongPassphraseErrors(t *testing.T) {
	r := tempRunner(t)
	if err := r.Init("prod", "correct"); err != nil {
		t.Fatalf("Init: %v", err)
	}
	if err := r.Set("prod", "wrong", "KEY", "val"); err == nil {
		t.Fatal("expected error with wrong passphrase")
	}
}

func TestListKeys(t *testing.T) {
	r := tempRunner(t)
	const pass = "pass"
	if err := r.Init("dev", pass); err != nil {
		t.Fatalf("Init: %v", err)
	}
	for _, kv := range [][2]string{{"A", "1"}, {"B", "2"}, {"C", "3"}} {
		if err := r.Set("dev", pass, kv[0], kv[1]); err != nil {
			t.Fatalf("Set %s: %v", kv[0], err)
		}
	}
	keys, err := r.List("dev", pass)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d: %v", len(keys), keys)
	}
}

func TestOpenMissingBundleErrors(t *testing.T) {
	r := tempRunner(t)
	if _, err := r.Get("ghost", "pass", "KEY"); err == nil {
		t.Fatal("expected error for missing bundle")
	}
}

// ensure test file compiles even when os/filepath are only used in helpers
var _ = os.DevNull
var _ = filepath.Separator
