package snapshot

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileStore persists snapshot blobs to a directory on disk.
type FileStore struct {
	dir string
}

// NewFileStore returns a FileStore rooted at dir, creating it if absent.
func NewFileStore(dir string) (*FileStore, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("snapshot store: mkdir %q: %w", dir, err)
	}
	return &FileStore{dir: dir}, nil
}

// Write saves blob under a filename derived from env and the current time.
// Returns the full path of the written file.
func (fs *FileStore) Write(env string, blob []byte) (string, error) {
	ts := time.Now().UTC().Format("20060102T150405Z")
	name := fmt.Sprintf("%s_%s.snap", env, ts)
	path := filepath.Join(fs.dir, name)
	if err := os.WriteFile(path, blob, 0o600); err != nil {
		return "", fmt.Errorf("snapshot store: write: %w", err)
	}
	return path, nil
}

// Read loads a snapshot blob from path.
func (fs *FileStore) Read(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("snapshot store: read: %w", err)
	}
	return data, nil
}

// List returns paths of all snapshots in the store for env.
// If env is empty, all snapshots are returned.
func (fs *FileStore) List(env string) ([]string, error) {
	entries, err := os.ReadDir(fs.dir)
	if err != nil {
		return nil, fmt.Errorf("snapshot store: list: %w", err)
	}
	var paths []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".snap") {
			continue
		}
		if env != "" && !strings.HasPrefix(name, env+"_") {
			continue
		}
		paths = append(paths, filepath.Join(fs.dir, name))
	}
	return paths, nil
}
