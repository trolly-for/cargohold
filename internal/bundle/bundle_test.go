package bundle_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/cargohold/internal/bundle"
)

func TestNew(t *testing.T) {
	b := bundle.New("myapp", "staging")
	if b.Name != "myapp" {
		t.Fatalf("expected name %q, got %q", "myapp", b.Name)
	}
	if b.Env != "staging" {
		t.Fatalf("expected env %q, got %q", "staging", b.Env)
	}
	if b.Secrets == nil {
		t.Fatal("expected secrets map to be initialised")
	}
}

func TestSetAndGet(t *testing.T) {
	b := bundle.New("myapp", "production")
	b.Set("DB_PASSWORD", "s3cr3t")

	v, err := b.Get("DB_PASSWORD")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v != "s3cr3t" {
		t.Fatalf("expected %q, got %q", "s3cr3t", v)
	}
}

func TestGetMissingKey(t *testing.T) {
	b := bundle.New("myapp", "production")
	_, err := b.Get("NONEXISTENT")
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestDelete(t *testing.T) {
	b := bundle.New("myapp", "dev")
	b.Set("API_KEY", "abc123")
	b.Delete("API_KEY")

	_, err := b.Get("API_KEY")
	if err == nil {
		t.Fatal("expected error after delete, got nil")
	}
}

func TestSaveAndLoad(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "test-bundle.json")

	b := bundle.New("webapp", "test")
	b.Set("SECRET_KEY", "topsecret")
	b.Set("PORT", "8080")

	if err := b.SaveToFile(path); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	loaded, err := bundle.LoadFromFile(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if loaded.Name != b.Name {
		t.Errorf("name mismatch: want %q got %q", b.Name, loaded.Name)
	}
	if loaded.Env != b.Env {
		t.Errorf("env mismatch: want %q got %q", b.Env, loaded.Env)
	}
	if v, _ := loaded.Get("SECRET_KEY"); v != "topsecret" {
		t.Errorf("secret mismatch: want %q got %q", "topsecret", v)
	}
}

func TestSaveFilePermissions(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "perm-bundle.json")

	b := bundle.New("perm", "ci")
	if err := b.SaveToFile(path); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected file mode 0600, got %v", info.Mode().Perm())
	}
}
