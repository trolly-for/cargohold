package lock_test

import (
	"os"
	"path/filepath"
	"testing"

	"cargohold/internal/lock"
)

func tempLocker(t *testing.T) *lock.Locker {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "locks")
	return lock.New(dir)
}

func TestLockAndIsLocked(t *testing.T) {
	l := tempLocker(t)

	if l.IsLocked("staging") {
		t.Fatal("expected staging to be unlocked initially")
	}

	if err := l.Lock("staging"); err != nil {
		t.Fatalf("Lock: unexpected error: %v", err)
	}

	if !l.IsLocked("staging") {
		t.Fatal("expected staging to be locked after Lock()")
	}
}

func TestLockDuplicateErrors(t *testing.T) {
	l := tempLocker(t)

	if err := l.Lock("production"); err != nil {
		t.Fatalf("first Lock: %v", err)
	}

	if err := l.Lock("production"); err != lock.ErrAlreadyLocked {
		t.Fatalf("expected ErrAlreadyLocked, got %v", err)
	}
}

func TestRelease(t *testing.T) {
	l := tempLocker(t)

	_ = l.Lock("dev")

	if err := l.Release("dev"); err != nil {
		t.Fatalf("Release: unexpected error: %v", err)
	}

	if l.IsLocked("dev") {
		t.Fatal("expected dev to be unlocked after Release()")
	}
}

func TestReleaseNotLockedErrors(t *testing.T) {
	l := tempLocker(t)

	if err := l.Release("staging"); err != lock.ErrNotLocked {
		t.Fatalf("expected ErrNotLocked, got %v", err)
	}
}

func TestLockFileCreatedOnDisk(t *testing.T) {
	dir := t.TempDir()
	l := lock.New(filepath.Join(dir, "locks"))

	if err := l.Lock("production"); err != nil {
		t.Fatalf("Lock: %v", err)
	}

	path := filepath.Join(dir, "locks", "production.lock")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected lock file at %s: %v", path, err)
	}
}

func TestIndependentEnvironments(t *testing.T) {
	l := tempLocker(t)

	_ = l.Lock("production")

	if l.IsLocked("staging") {
		t.Fatal("locking production should not affect staging")
	}
}
