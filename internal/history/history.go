// Package history tracks a rolling log of key-level mutations within a bundle
// so that users can review what changed, when, and under which environment.
package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Op describes the kind of mutation recorded.
type Op string

const (
	OpSet    Op = "set"
	OpDelete Op = "delete"
)

// Entry is a single history record.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Env       string    `json:"env"`
	Key       string    `json:"key"`
	Op        Op        `json:"op"`
}

// Tracker appends and reads history entries for a given environment bundle.
type Tracker struct {
	path string
}

// New returns a Tracker whose log file lives at dir/<env>.history.jsonl.
func New(dir, env string) (*Tracker, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("history: mkdir %s: %w", dir, err)
	}
	return &Tracker{path: filepath.Join(dir, env+".history.jsonl")}, nil
}

// Record appends a new entry to the log.
func (t *Tracker) Record(env, key string, op Op) error {
	f, err := os.OpenFile(t.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("history: open log: %w", err)
	}
	defer f.Close()

	entry := Entry{
		Timestamp: time.Now().UTC(),
		Env:       env,
		Key:       key,
		Op:        op,
	}
	line, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("history: marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(f, "%s\n", line)
	return err
}

// ReadAll returns all entries from the log in recorded order.
func (t *Tracker) ReadAll() ([]Entry, error) {
	data, err := os.ReadFile(t.path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("history: read log: %w", err)
	}

	var entries []Entry
	for _, line := range splitLines(data) {
		if len(line) == 0 {
			continue
		}
		var e Entry
		if err := json.Unmarshal(line, &e); err != nil {
			return nil, fmt.Errorf("history: parse entry: %w", err)
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
