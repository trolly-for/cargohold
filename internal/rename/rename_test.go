package rename_test

import (
	"testing"

	"cargohold/internal/bundle"
	"cargohold/internal/rename"
)

func seedBundle(t *testing.T, pairs map[string]string) *bundle.Bundle {
	t.Helper()
	b := bundle.New("test")
	for k, v := range pairs {
		if err := b.Set(k, v); err != nil {
			t.Fatalf("seed: set %q: %v", k, err)
		}
	}
	return b
}

func TestRenameMovesValue(t *testing.T) {
	b := seedBundle(t, map[string]string{"OLD_KEY": "secret"})

	if err := rename.Key(b, "OLD_KEY", "NEW_KEY", rename.Options{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, ok := b.Get("OLD_KEY"); ok {
		t.Error("old key should have been removed")
	}
	val, ok := b.Get("NEW_KEY")
	if !ok {
		t.Fatal("new key should exist")
	}
	if val != "secret" {
		t.Errorf("expected %q, got %q", "secret", val)
	}
}

func TestRenameSameKeyErrors(t *testing.T) {
	b := seedBundle(t, map[string]string{"KEY": "val"})

	err := rename.Key(b, "KEY", "KEY", rename.Options{})
	if err == nil {
		t.Fatal("expected error for same key")
	}
	if !errors.Is(err, rename.ErrSameKey) {
		t.Errorf("expected ErrSameKey, got %v", err)
	}
}

func TestRenameMissingSourceErrors(t *testing.T) {
	b := bundle.New("test")

	err := rename.Key(b, "MISSING", "DEST", rename.Options{})
	if err == nil {
		t.Fatal("expected error for missing source")
	}
	if !errors.Is(err, rename.ErrSrcMissing) {
		t.Errorf("expected ErrSrcMissing, got %v", err)
	}
}

func TestRenameDestExistsWithoutOverwriteErrors(t *testing.T) {
	b := seedBundle(t, map[string]string{"SRC": "a", "DST": "b"})

	err := rename.Key(b, "SRC", "DST", rename.Options{Overwrite: false})
	if err == nil {
		t.Fatal("expected error when destination exists")
	}
	if !errors.Is(err, rename.ErrDestExists) {
		t.Errorf("expected ErrDestExists, got %v", err)
	}
}

func TestRenameDestExistsWithOverwrite(t *testing.T) {
	b := seedBundle(t, map[string]string{"SRC": "new-val", "DST": "old-val"})

	if err := rename.Key(b, "SRC", "DST", rename.Options{Overwrite: true}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	val, ok := b.Get("DST")
	if !ok {
		t.Fatal("DST should exist after overwrite")
	}
	if val != "new-val" {
		t.Errorf("expected %q, got %q", "new-val", val)
	}
	if _, ok := b.Get("SRC"); ok {
		t.Error("SRC should have been removed")
	}
}
