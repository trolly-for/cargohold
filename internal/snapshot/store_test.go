package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"

	"cargohold/internal/snapshot"
)

func tempFileStore(t *testing.T) *snapshot.FileStore {
	t.Helper()
	dir := t.TempDir()
	fs, err := snapshot.NewFileStore(dir)
	if err != nil {
		t.Fatalf("NewFileStore: %v", err)
	}
	return fs
}

func TestWriteAndRead(t *testing.T) {
	fs := tempFileStore(t)
	blob := []byte("encryptedpayload")

	path, err := fs.Write("staging", blob)
	if err != nil {
		t.Fatalf("Write: %v", err)
	}

	got, err := fs.Read(path)
	if err != nil {
		t.Fatalf("Read: %v", err)
	}
	if string(got) != string(blob) {
		t.Errorf("blob mismatch: got %q, want %q", got, blob)
	}
}

func TestListFiltersByEnv(t *testing.T) {
	fs := tempFileStore(t)

	if _, err := fs.Write("staging", []byte("a")); err != nil {
		t.Fatalf("Write staging: %v", err)
	}
	if _, err := fs.Write("prod", []byte("b")); err != nil {
		t.Fatalf("Write prod: %v", err)
	}

	paths, err := fs.List("staging")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(paths) != 1 {
		t.Fatalf("expected 1 staging snapshot, got %d", len(paths))
	}
	if filepath.Base(paths[0])[:7] != "staging" {
		t.Errorf("unexpected path %q", paths[0])
	}
}

func TestListAll(t *testing.T) {
	fs := tempFileStore(t)
	for _, env := range []string{"dev", "staging", "prod"} {
		if _, err := fs.Write(env, []byte(env)); err != nil {
			t.Fatalf("Write %s: %v", env, err)
		}
	}

	paths, err := fs.List("")
	if err != nil {
		t.Fatalf("List all: %v", err)
	}
	if len(paths) != 3 {
		t.Errorf("expected 3 snapshots, got %d", len(paths))
	}
}

func TestReadMissing(t *testing.T) {
	fs := tempFileStore(t)
	_, err := fs.Read(filepath.Join(os.TempDir(), "nonexistent.snap"))
	if err == nil {
		t.Fatal("expected error reading missing file")
	}
}
