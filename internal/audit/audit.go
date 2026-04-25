// Package audit provides a lightweight audit log for recording
// secret bundle operations (init, get, set, delete, rotate).
package audit

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single audit log record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Operation string    `json:"operation"`
	Env       string    `json:"env"`
	Key       string    `json:"key,omitempty"`
	Note      string    `json:"note,omitempty"`
}

// Logger writes audit entries to a newline-delimited JSON file.
type Logger struct {
	path string
}

// New returns a Logger that appends to the file at path.
// The directory must already exist.
func New(path string) *Logger {
	return &Logger{path: path}
}

// Default returns a Logger using ~/.cargohold/audit.log.
func Default() (*Logger, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("audit: resolve home dir: %w", err)
	}
	dir := filepath.Join(home, ".cargohold")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("audit: create dir: %w", err)
	}
	return New(filepath.Join(dir, "audit.log")), nil
}

// Record appends an entry to the audit log.
func (l *Logger) Record(op, env, key, note string) error {
	e := Entry{
		Timestamp: time.Now().UTC(),
		Operation: op,
		Env:       env,
		Key:       key,
		Note:      note,
	}
	data, err := json.Marshal(e)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("audit: open log: %w", err)
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "%s\n", data)
	return err
}

// ReadAll returns all entries from the audit log.
func (l *Logger) ReadAll() ([]Entry, error) {
	data, err := os.ReadFile(l.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("audit: read log: %w", err)
	}
	var entries []Entry
	for _, line := range splitLines(data) {
		if len(line) == 0 {
			continue
		}
		var e Entry
		if err := json.Unmarshal(line, &e); err != nil {
			return nil, fmt.Errorf("audit: parse entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

func splitLines(data []byte) [][]byte {
	var lines [][]byte
	start := 0
	for i, b := range data {
		if b == '\n' {
			lines = append(lines, data[start:i])
			start = i + 1
		}
	}
	if start < len(data) {
		lines = append(lines, data[start:])
	}
	return lines
}
