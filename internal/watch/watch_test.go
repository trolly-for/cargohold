package watch_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"cargohold/internal/watch"
)

func TestWatchDetectsChange(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bundle.enc")

	// Create the file before starting the watcher.
	if err := os.WriteFile(path, []byte("v1"), 0o600); err != nil {
		t.Fatal(err)
	}

	w := watch.New(20 * time.Millisecond)
	ch := w.Watch(path, "staging")
	defer w.Stop()

	// Give the watcher one tick to record the initial mod-time.
	time.Sleep(40 * time.Millisecond)

	// Overwrite the file to bump its modification time.
	if err := os.WriteFile(path, []byte("v2"), 0o600); err != nil {
		t.Fatal(err)
	}

	select {
	case ev, ok := <-ch:
		if !ok {
			t.Fatal("channel closed before event")
		}
		if ev.Env != "staging" {
			t.Errorf("env = %q, want %q", ev.Env, "staging")
		}
		if ev.Path != path {
			t.Errorf("path = %q, want %q", ev.Path, path)
		}
		if ev.ModTime.IsZero() {
			t.Error("ModTime should not be zero")
		}
	case <-time.After(500 * time.Millisecond):
		t.Fatal("timed out waiting for change event")
	}
}

func TestWatchNoFalsePositive(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bundle.enc")

	if err := os.WriteFile(path, []byte("unchanged"), 0o600); err != nil {
		t.Fatal(err)
	}

	w := watch.New(20 * time.Millisecond)
	ch := w.Watch(path, "dev")
	defer w.Stop()

	time.Sleep(100 * time.Millisecond)

	select {
	case ev := <-ch:
		t.Errorf("unexpected event for unchanged file: %+v", ev)
	default:
		// correct: no event expected
	}
}

func TestWatchMissingFileNoBlock(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "missing.enc")

	w := watch.New(20 * time.Millisecond)
	ch := w.Watch(path, "prod")

	time.Sleep(80 * time.Millisecond)
	w.Stop()

	// Channel should be closed after Stop with no events.
	select {
	case _, ok := <-ch:
		if ok {
			t.Error("expected channel to be closed")
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for channel close")
	}
}
