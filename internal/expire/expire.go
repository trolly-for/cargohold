// Package expire provides bundle-level expiry enforcement for cargohold.
// A bundle can be marked with an expiry time; once expired, read and write
// operations are blocked until the expiry is cleared or extended.
package expire

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ErrExpired is returned when an operation is attempted on an expired bundle.
var ErrExpired = errors.New("bundle has expired")

// ErrNotSet is returned when no expiry has been configured for the bundle.
var ErrNotSet = errors.New("no expiry set")

// Expirer manages expiry metadata for a single bundle.
type Expirer struct {
	path string
}

// New returns an Expirer whose state is persisted under dir.
func New(dir string) (*Expirer, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("expire: mkdir: %w", err)
	}
	return &Expirer{path: filepath.Join(dir, "expiry")}, nil
}

// Set writes an expiry time for the bundle.
func (e *Expirer) Set(t time.Time) error {
	data := []byte(t.UTC().Format(time.RFC3339))
	if err := os.WriteFile(e.path, data, 0o600); err != nil {
		return fmt.Errorf("expire: set: %w", err)
	}
	return nil
}

// Clear removes any expiry that has been set.
func (e *Expirer) Clear() error {
	err := os.Remove(e.path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("expire: clear: %w", err)
	}
	return nil
}

// Get returns the expiry time. Returns ErrNotSet if no expiry is configured.
func (e *Expirer) Get() (time.Time, error) {
	data, err := os.ReadFile(e.path)
	if errors.Is(err, os.ErrNotExist) {
		return time.Time{}, ErrNotSet
	}
	if err != nil {
		return time.Time{}, fmt.Errorf("expire: get: %w", err)
	}
	t, err := time.Parse(time.RFC3339, string(data))
	if err != nil {
		return time.Time{}, fmt.Errorf("expire: parse: %w", err)
	}
	return t, nil
}

// Check returns ErrExpired if the bundle has passed its expiry time,
// ErrNotSet if no expiry is configured, or nil if still valid.
func (e *Expirer) Check() error {
	t, err := e.Get()
	if err != nil {
		return err
	}
	if time.Now().UTC().After(t) {
		return ErrExpired
	}
	return nil
}
