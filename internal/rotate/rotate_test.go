package rotate_test

import (
	"os"
	"path/filepath"
	"testing"

	"cargohold/internal/bundle"
	"cargohold/internal/rotate"
	"cargohold/internal/store"
	"cargohold/internal/vault"
)

func tempRotator(t *testing.T) (*rotate.Rotator, *store.Store) {
	t.Helper()
	dir, err := os.MkdirTemp("", "rotate-test-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	s := store.New(filepath.Join(dir, "bundles"))
	return rotate.New(s), s
}

func seedBundle(t *testing.T, s *store.Store, env, pass string) {
	t.Helper()
	v, err := vault.New(s, env)
	if err != nil {
		t.Fatalf("vault.New: %v", err)
	}
	b := bundle.New()
	b.Set("KEY", "value")
	if err := v.Init(b, pass); err != nil {
		t.Fatalf("vault.Init: %v", err)
	}
}

func TestRotateSuccess(t *testing.T) {
	r, s := tempRotator(t)
	seedBundle(t, s, "staging", "old-passphrase-1")

	if err := r.Rotate("staging", "old-passphrase-1", "new-passphrase-2"); err != nil {
		t.Fatalf("Rotate: %v", err)
	}

	// Verify the bundle is readable with the new passphrase.
	v, _ := vault.New(s, "staging")
	b, err := v.Open("new-passphrase-2")
	if err != nil {
		t.Fatalf("Open after rotate: %v", err)
	}
	val, ok := b.Get("KEY")
	if !ok || val != "value" {
		t.Errorf("expected KEY=value after rotate, got %q ok=%v", val, ok)
	}
}

func TestRotateOldPassphraseRejected(t *testing.T) {
	r, s := tempRotator(t)
	seedBundle(t, s, "staging", "old-passphrase-1")
	r.Rotate("staging", "old-passphrase-1", "new-passphrase-2") // nolint

	v, _ := vault.New(s, "staging")
	_, err := v.Open("old-passphrase-1")
	if err == nil {
		t.Fatal("expected error opening with old passphrase after rotation")
	}
}

func TestRotateSamePassphraseErrors(t *testing.T) {
	r, s := tempRotator(t)
	seedBundle(t, s, "staging", "passphrase-abc")

	err := r.Rotate("staging", "passphrase-abc", "passphrase-abc")
	if err != rotate.ErrSamePassphrase {
		t.Fatalf("expected ErrSamePassphrase, got %v", err)
	}
}

func TestRotateWrongOldPassphrase(t *testing.T) {
	r, s := tempRotator(t)
	seedBundle(t, s, "staging", "correct-passphrase")

	err := r.Rotate("staging", "wrong-passphrase", "new-passphrase-2")
	if err == nil {
		t.Fatal("expected error with wrong old passphrase")
	}
}

func TestRotateMissingBundle(t *testing.T) {
	r, _ := tempRotator(t)
	err := r.Rotate("nonexistent", "pass1", "pass2")
	if err == nil {
		t.Fatal("expected error rotating missing bundle")
	}
}
