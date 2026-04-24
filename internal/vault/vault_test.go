package vault_test

import (
	"os"
	"testing"

	"cargohold/internal/store"
	"cargohold/internal/vault"
)

func tempVault(t *testing.T, passphrase string) *vault.Vault {
	t.Helper()
	dir, err := os.MkdirTemp("", "vault-test-*")
	if err != nil {
		t.Fatalf("create temp dir: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	s, err := store.New(dir)
	if err != nil {
		t.Fatalf("create store: %v", err)
	}
	return vault.New(s, passphrase)
}

func TestInitAndOpen(t *testing.T) {
	v := tempVault(t, "correct-horse")

	if err := v.Init("prod"); err != nil {
		t.Fatalf("Init: %v", err)
	}

	b, err := v.Open("prod")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil bundle")
	}
}

func TestInitDuplicate(t *testing.T) {
	v := tempVault(t, "pass")
	v.Init("dev")
	if err := v.Init("dev"); err == nil {
		t.Fatal("expected error for duplicate init")
	}
}

func TestOpenMissing(t *testing.T) {
	v := tempVault(t, "pass")
	if _, err := v.Open("ghost"); err == nil {
		t.Fatal("expected error opening non-existent bundle")
	}
}

func TestSaveAndOpenRoundtrip(t *testing.T) {
	v := tempVault(t, "s3cr3t")
	v.Init("staging")

	b, _ := v.Open("staging")
	b.Set("DB_URL", "postgres://localhost/mydb")
	b.Set("API_KEY", "abc123")

	if err := v.Save("staging", b); err != nil {
		t.Fatalf("Save: %v", err)
	}

	b2, err := v.Open("staging")
	if err != nil {
		t.Fatalf("Open after save: %v", err)
	}

	if val, _ := b2.Get("DB_URL"); val != "postgres://localhost/mydb" {
		t.Errorf("DB_URL mismatch: got %q", val)
	}
	if val, _ := b2.Get("API_KEY"); val != "abc123" {
		t.Errorf("API_KEY mismatch: got %q", val)
	}
}

func TestWrongPassphrase(t *testing.T) {
	v1 := tempVault(t, "correct")
	v1.Init("env")

	// Re-open the same store directory with wrong passphrase.
	dir := v1.StorePath()
	s, _ := store.New(dir)
	v2 := vault.New(s, "wrong")

	if _, err := v2.Open("env"); err == nil {
		t.Fatal("expected decryption error with wrong passphrase")
	}
}

func TestListAndDelete(t *testing.T) {
	v := tempVault(t, "pass")
	v.Init("alpha")
	v.Init("beta")

	names, err := v.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 2 {
		t.Fatalf("expected 2 bundles, got %d", len(names))
	}

	v.Delete("alpha")
	names, _ = v.List()
	if len(names) != 1 || names[0] != "beta" {
		t.Errorf("expected only beta after delete, got %v", names)
	}
}
