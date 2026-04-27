// Package pin provides key pinning: recording a specific value for a key
// and detecting when it has drifted from the pinned value.
package pin

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

// ErrNotPinned is returned when a key has no pinned value.
var ErrNotPinned = errors.New("key is not pinned")

// ErrDrifted is returned when the current value differs from the pinned value.
var ErrDrifted = errors.New("value has drifted from pinned value")

// Record holds a single pinned key entry.
type Record struct {
	Value   string    `json:"value"`
	PinnedAt time.Time `json:"pinned_at"`
}

// Pinner manages pinned key values for a named bundle.
type Pinner struct {
	mu   sync.RWMutex
	path string
	data map[string]Record
}

// New creates a Pinner backed by the file at path.
// If the file does not exist it starts with an empty set of pins.
func New(path string) (*Pinner, error) {
	p := &Pinner{path: path, data: make(map[string]Record)}
	if err := p.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("pin: load %s: %w", path, err)
	}
	return p, nil
}

// Pin records value as the expected value for key.
func (p *Pinner) Pin(key, value string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.data[key] = Record{Value: value, PinnedAt: time.Now().UTC()}
	return p.save()
}

// Unpin removes the pin for key. Returns ErrNotPinned if key was not pinned.
func (p *Pinner) Unpin(key string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.data[key]; !ok {
		return ErrNotPinned
	}
	delete(p.data, key)
	return p.save()
}

// Check compares current against the pinned value for key.
// Returns ErrNotPinned if no pin exists, ErrDrifted if values differ.
func (p *Pinner) Check(key, current string) error {
	p.mu.RLock()
	defer p.mu.RUnlock()
	rec, ok := p.data[key]
	if !ok {
		return ErrNotPinned
	}
	if rec.Value != current {
		return fmt.Errorf("%w: key %q", ErrDrifted, key)
	}
	return nil
}

// Get returns the Record for key, or ErrNotPinned.
func (p *Pinner) Get(key string) (Record, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	rec, ok := p.data[key]
	if !ok {
		return Record{}, ErrNotPinned
	}
	return rec, nil
}

func (p *Pinner) load() error {
	f, err := os.Open(p.path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&p.data)
}

func (p *Pinner) save() error {
	f, err := os.Create(p.path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(p.data)
}
