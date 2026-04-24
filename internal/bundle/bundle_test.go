package bundle_test

import (
	"os"
	"testing"

	"github.com/cargohold/cargohold/internal/bundle"
)

func TestNew(t *testing.T) {
	b := bundle.New("production")
	if b.Name != "production" {
		t.Fatalf("expected name 'production', got %q", b.Name)
	}
	if len(b.Secrets) != 0 {
		t.Fatal("expected empty secrets map")
	}
}

func TestSetAndGet(t *testing.T) {
	b := bundle.New("test")
	b.Set("DB_URL", "postgres://localhost/mydb")

	val, err := b.Get("DB_URL")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "postgres://localhost/mydb" {
		t.Fatalf("expected 'postgres://localhost/mydb', got %q", val)
	}
}

func TestGetMissingKey(t *testing.T) {
	b := bundle.New("test")
	_, err := b.Get("MISSING_KEY")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestDelete(t *testing.T) {
	b := bundle.New("test")
	b.Set("API_KEY", "abc123")
	b.Delete("API_KEY")

	_, err := b.Get("API_KEY")
	if err == nil {
		t.Fatal("expected error after deleting key")
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "bundle-*.enc")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	b := bundle.New("staging")
	b.Set("SECRET_KEY", "supersecret")
	b.Set("API_TOKEN", "token-xyz")

	passphrase := "mypassphrase"
	if err := b.Save(tmpFile.Name(), passphrase); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := bundle.Load(tmpFile.Name(), passphrase)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Name != "staging" {
		t.Fatalf("expected name 'staging', got %q", loaded.Name)
	}

	val, err := loaded.Get("SECRET_KEY")
	if err != nil || val != "supersecret" {
		t.Fatalf("expected 'supersecret', got %q (err: %v)", val, err)
	}
}

func TestLoadWrongPassphrase(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "bundle-*.enc")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	b := bundle.New("test")
	b.Set("KEY", "value")

	if err := b.Save(tmpFile.Name(), "correctpass"); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	_, err = bundle.Load(tmpFile.Name(), "wrongpass")
	if err == nil {
		t.Fatal("expected error when loading with wrong passphrase")
	}
}
