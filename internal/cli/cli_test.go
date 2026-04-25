package cli_test

import (
	"os"
	"path/filepath"
	"testing"

	clitest "cargohold/internal/cli"
	"cargohold/internal/rotate"
	"cargohold/internal/store"
)

func tempRunner(t *testing.T) *clitest.Runner {
	t.Helper()
	dir, err := os.MkdirTemp("", "cli-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	s := store.New(filepath.Join(dir, "bundles"))
	return clitest.NewWithStore(s)
}

const (
	testEnv  = "staging"
	testPass = "hunter2-passphrase"
)

func TestInitCreatesBundle(t *testing.T) {
	r := tempRunner(t)
	if err := r.Init(testEnv, testPass); err != nil {
		t.Fatalf("Init: %v", err)
	}
}

func TestInitDuplicateErrors(t *testing.T) {
	r := tempRunner(t)
	r.Init(testEnv, testPass) // nolint
	if err := r.Init(testEnv, testPass); err == nil {
		t.Fatal("expected error on duplicate Init")
	}
}

func TestSetAndGet(t *testing.T) {
	r := tempRunner(t)
	r.Init(testEnv, testPass) // nolint

	if err := r.Set(testEnv, testPass, "DB_URL", "postgres://localhost/db"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := r.Get(testEnv, testPass, "DB_URL"); err != nil {
		t.Fatalf("Get: %v", err)
	}
}

func TestGetMissingKey(t *testing.T) {
	r := tempRunner(t)
	r.Init(testEnv, testPass) // nolint

	if err := r.Get(testEnv, testPass, "MISSING"); err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestDeleteKey(t *testing.T) {
	r := tempRunner(t)
	r.Init(testEnv, testPass) // nolint
	r.Set(testEnv, testPass, "TOKEN", "abc123") // nolint

	if err := r.Delete(testEnv, testPass, "TOKEN"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if err := r.Get(testEnv, testPass, "TOKEN"); err == nil {
		t.Fatal("expected key to be gone after Delete")
	}
}

func TestRotateViaRunner(t *testing.T) {
	r := tempRunner(t)
	r.Init(testEnv, testPass) // nolint
	r.Set(testEnv, testPass, "SECRET", "s3cr3t") // nolint

	newPass := "new-passphrase-xyz"
	if err := r.Rotate(testEnv, testPass, newPass); err != nil {
		t.Fatalf("Rotate: %v", err)
	}
	if err := r.Get(testEnv, newPass, "SECRET"); err != nil {
		t.Fatalf("Get after rotate: %v", err)
	}
}

func TestRotateSamePassphraseError(t *testing.T) {
	r := tempRunner(t)
	r.Init(testEnv, testPass) // nolint

	err := r.Rotate(testEnv, testPass, testPass)
	if err == nil {
		t.Fatal("expected error rotating with same passphrase")
	}
	if err != rotate.ErrSamePassphrase {
		t.Fatalf("expected ErrSamePassphrase, got %v", err)
	}
}
