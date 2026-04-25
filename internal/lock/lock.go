// Package lock provides environment-level write-locking for bundles.
// A lock prevents accidental modifications to a bundle until explicitly
// released, acting as a safety mechanism for sensitive environments.
package lock

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ErrAlreadyLocked is returned when attempting to lock an already-locked bundle.
var ErrAlreadyLocked = errors.New("bundle is already locked")

// ErrNotLocked is returned when attempting to release a lock that does not exist.
var ErrNotLocked = errors.New("bundle is not locked")

// Locker manages lock files for environment bundles.
type Locker struct {
	dir string
}

// New returns a Locker that stores lock files under dir.
func New(dir string) *Locker {
	return &Locker{dir: dir}
}

// Lock creates a lock file for the given environment.
// Returns ErrAlreadyLocked if a lock already exists.
func (l *Locker) Lock(env string) error {
	path := l.lockPath(env)
	if err := os.MkdirAll(l.dir, 0o700); err != nil {
		return fmt.Errorf("lock: create dir: %w", err)
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		if os.IsExist(err) {
			return ErrAlreadyLocked
		}
		return fmt.Errorf("lock: create file: %w", err)
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "%d\n", time.Now().Unix())
	return err
}

// Release removes the lock file for the given environment.
// Returns ErrNotLocked if no lock exists.
func (l *Locker) Release(env string) error {
	path := l.lockPath(env)
	err := os.Remove(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrNotLocked
		}
		return fmt.Errorf("lock: remove: %w", err)
	}
	return nil
}

// IsLocked reports whether the given environment currently has a lock.
func (l *Locker) IsLocked(env string) bool {
	_, err := os.Stat(l.lockPath(env))
	return err == nil
}

func (l *Locker) lockPath(env string) string {
	return filepath.Join(l.dir, env+".lock")
}
