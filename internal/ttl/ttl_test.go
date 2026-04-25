package ttl_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"cargohold/internal/ttl"
)

func tempTracker(t *testing.T) *ttl.Tracker {
	t.Helper()
	dir := t.TempDir()
	tr, err := ttl.New(dir)
	if err != nil {
		t.Fatalf("ttl.New: %v", err)
	}
	return tr
}

func TestSetAndCheckValid(t *testing.T) {
	tr := tempTracker(t)
	if err := tr.Set("staging", 10*time.Minute); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := tr.Check("staging"); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestCheckExpired(t *testing.T) {
	tr := tempTracker(t)
	if err := tr.Set("staging", -1*time.Second); err != nil {
		t.Fatalf("Set: %v", err)
	}
	err := tr.Check("staging")
	if !errors.Is(err, ttl.ErrExpired) {
		t.Fatalf("expected ErrExpired, got %v", err)
	}
}

func TestCheckNoExpiry(t *testing.T) {
	tr := tempTracker(t)
	err := tr.Check("production")
	if !errors.Is(err, ttl.ErrNoExpiry) {
		t.Fatalf("expected ErrNoExpiry, got %v", err)
	}
}

func TestGetRecord(t *testing.T) {
	tr := tempTracker(t)
	d := 5 * time.Minute
	before := time.Now().UTC().Add(d)
	if err := tr.Set("dev", d); err != nil {
		t.Fatalf("Set: %v", err)
	}
	rec, err := tr.Get("dev")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if rec.Env != "dev" {
		t.Errorf("env = %q, want dev", rec.Env)
	}
	if rec.ExpiresAt.Before(before.Add(-2 * time.Second)) {
		t.Errorf("ExpiresAt %v is too early", rec.ExpiresAt)
	}
}

func TestGetMissingRecord(t *testing.T) {
	tr := tempTracker(t)
	_, err := tr.Get("ghost")
	if !errors.Is(err, ttl.ErrNoExpiry) {
		t.Fatalf("expected ErrNoExpiry, got %v", err)
	}
}

func TestRemove(t *testing.T) {
	tr := tempTracker(t)
	if err := tr.Set("staging", time.Hour); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := tr.Remove("staging"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	err := tr.Check("staging")
	if !errors.Is(err, ttl.ErrNoExpiry) {
		t.Fatalf("expected ErrNoExpiry after remove, got %v", err)
	}
}

func TestRemoveNonExistent(t *testing.T) {
	tr := tempTracker(t)
	if err := tr.Remove("ghost"); err != nil {
		t.Errorf("Remove of non-existent should not error, got %v", err)
	}
}

func TestNewCreatesDir(t *testing.T) {
	base := t.TempDir()
	dir := base + "/nested/ttl"
	_, err := ttl.New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("directory not created: %v", err)
	}
}
