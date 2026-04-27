package expire_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"cargohold/internal/expire"
)

func tempExpirer(t *testing.T) *expire.Expirer {
	t.Helper()
	dir := t.TempDir()
	ex, err := expire.New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return ex
}

func TestSetAndGet(t *testing.T) {
	ex := tempExpirer(t)
	want := time.Now().UTC().Add(24 * time.Hour).Truncate(time.Second)
	if err := ex.Set(want); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, err := ex.Get()
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestGetNotSet(t *testing.T) {
	ex := tempExpirer(t)
	_, err := ex.Get()
	if !errors.Is(err, expire.ErrNotSet) {
		t.Errorf("expected ErrNotSet, got %v", err)
	}
}

func TestCheckValid(t *testing.T) {
	ex := tempExpirer(t)
	if err := ex.Set(time.Now().UTC().Add(time.Hour)); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := ex.Check(); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestCheckExpired(t *testing.T) {
	ex := tempExpirer(t)
	if err := ex.Set(time.Now().UTC().Add(-time.Second)); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := ex.Check(); !errors.Is(err, expire.ErrExpired) {
		t.Errorf("expected ErrExpired, got %v", err)
	}
}

func TestCheckNoExpiry(t *testing.T) {
	ex := tempExpirer(t)
	if err := ex.Check(); !errors.Is(err, expire.ErrNotSet) {
		t.Errorf("expected ErrNotSet, got %v", err)
	}
}

func TestClear(t *testing.T) {
	ex := tempExpirer(t)
	if err := ex.Set(time.Now().UTC().Add(time.Hour)); err != nil {
		t.Fatalf("Set: %v", err)
	}
	if err := ex.Clear(); err != nil {
		t.Fatalf("Clear: %v", err)
	}
	_, err := ex.Get()
	if !errors.Is(err, expire.ErrNotSet) {
		t.Errorf("expected ErrNotSet after clear, got %v", err)
	}
}

func TestClearIdempotent(t *testing.T) {
	ex := tempExpirer(t)
	if err := ex.Clear(); err != nil {
		t.Errorf("Clear on empty should not error: %v", err)
	}
}

func TestNewCreatesDir(t *testing.T) {
	base := t.TempDir()
	dir := base + "/nested/expiry"
	_, err := expire.New(dir)
	if err != nil {
		t.Fatalf("New with nested dir: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("directory not created: %v", err)
	}
}
