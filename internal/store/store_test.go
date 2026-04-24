package store_test

import (
	"os"
	"path/filepath"
	"testing"

	"cargohold/internal/store"
)

func tempStore(t *testing.T) *store.Store {
	t.Helper()
	dir := t.TempDir()
	s, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return s
}

func TestNew(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "dir")
	s, err := store.New(dir)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if s.BaseDir != dir {
		t.Errorf("expected BaseDir %q, got %q", dir, s.BaseDir)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}

func TestBundlePath(t *testing.T) {
	s := tempStore(t)
	got := s.BundlePath("production")
	want := filepath.Join(s.BaseDir, "production"+store.BundleExtension)
	if got != want {
		t.Errorf("BundlePath: got %q, want %q", got, want)
	}
}

func TestExistsAndList(t *testing.T) {
	s := tempStore(t)

	if s.Exists("dev") {
		t.Error("expected bundle to not exist yet")
	}

	// Create a fake bundle file.
	if err := os.WriteFile(s.BundlePath("dev"), []byte("data"), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := os.WriteFile(s.BundlePath("prod"), []byte("data"), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	if !s.Exists("dev") {
		t.Error("expected bundle 'dev' to exist")
	}

	names, err := s.List()
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(names) != 2 {
		t.Errorf("expected 2 bundles, got %d: %v", len(names), names)
	}
}

func TestDelete(t *testing.T) {
	s := tempStore(t)

	if err := os.WriteFile(s.BundlePath("staging"), []byte("x"), 0600); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	if err := s.Delete("staging"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if s.Exists("staging") {
		t.Error("expected bundle to be deleted")
	}
}

func TestDeleteMissing(t *testing.T) {
	s := tempStore(t)
	if err := s.Delete("nonexistent"); err == nil {
		t.Error("expected error when deleting nonexistent bundle")
	}
}
