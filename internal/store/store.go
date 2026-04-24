package store

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	DefaultDir      = ".cargohold"
	BundleExtension = ".bundle"
)

// Store manages the filesystem location where encrypted bundles are persisted.
type Store struct {
	BaseDir string
}

// New creates a Store rooted at baseDir, creating the directory if needed.
func New(baseDir string) (*Store, error) {
	if err := os.MkdirAll(baseDir, 0700); err != nil {
		return nil, fmt.Errorf("store: failed to create base directory %q: %w", baseDir, err)
	}
	return &Store{BaseDir: baseDir}, nil
}

// Default returns a Store in the user's home directory under DefaultDir.
func Default() (*Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("store: could not determine home directory: %w", err)
	}
	return New(filepath.Join(home, DefaultDir))
}

// BundlePath returns the full path for a named bundle.
func (s *Store) BundlePath(name string) string {
	return filepath.Join(s.BaseDir, name+BundleExtension)
}

// Exists reports whether a bundle with the given name exists on disk.
func (s *Store) Exists(name string) bool {
	_, err := os.Stat(s.BundlePath(name))
	return err == nil
}

// List returns the names of all bundles currently stored.
func (s *Store) List() ([]string, error) {
	entries, err := os.ReadDir(s.BaseDir)
	if err != nil {
		return nil, fmt.Errorf("store: failed to read directory: %w", err)
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if ext := filepath.Ext(e.Name()); ext == BundleExtension {
			names = append(names, e.Name()[:len(e.Name())-len(BundleExtension)])
		}
	}
	return names, nil
}

// Delete removes a bundle from the store.
func (s *Store) Delete(name string) error {
	path := s.BundlePath(name)
	if err := os.Remove(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("store: bundle %q does not exist", name)
		}
		return fmt.Errorf("store: failed to delete bundle %q: %w", name, err)
	}
	return nil
}
