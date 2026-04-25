// Package ttl provides time-to-live expiry tracking for secret bundles.
// Each bundle can be assigned an expiry time; expired bundles are flagged
// so callers can warn or refuse access accordingly.
package ttl

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// ErrExpired is returned when a bundle's TTL has elapsed.
var ErrExpired = errors.New("bundle has expired")

// ErrNoExpiry is returned when no expiry record exists for a bundle.
var ErrNoExpiry = errors.New("no expiry set for bundle")

// Record holds the expiry metadata persisted alongside a bundle.
type Record struct {
	Env       string    `json:"env"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Tracker manages TTL records stored in a directory.
type Tracker struct {
	dir string
}

// New returns a Tracker that stores expiry files under dir.
func New(dir string) (*Tracker, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	return &Tracker{dir: dir}, nil
}

func (t *Tracker) recordPath(env string) string {
	return filepath.Join(t.dir, env+".ttl.json")
}

// Set stores an expiry for env that is valid for duration d from now.
func (t *Tracker) Set(env string, d time.Duration) error {
	rec := Record{
		Env:       env,
		ExpiresAt: time.Now().UTC().Add(d),
	}
	data, err := json.Marshal(rec)
	if err != nil {
		return err
	}
	return os.WriteFile(t.recordPath(env), data, 0o600)
}

// Check returns ErrNoExpiry if no record exists, ErrExpired if the TTL has
// elapsed, or nil if the bundle is still valid.
func (t *Tracker) Check(env string) error {
	data, err := os.ReadFile(t.recordPath(env))
	if errors.Is(err, os.ErrNotExist) {
		return ErrNoExpiry
	}
	if err != nil {
		return err
	}
	var rec Record
	if err := json.Unmarshal(data, &rec); err != nil {
		return err
	}
	if time.Now().UTC().After(rec.ExpiresAt) {
		return ErrExpired
	}
	return nil
}

// Remove deletes the expiry record for env. It is not an error if none exists.
func (t *Tracker) Remove(env string) error {
	err := os.Remove(t.recordPath(env))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

// Get returns the Record for env or an error.
func (t *Tracker) Get(env string) (Record, error) {
	data, err := os.ReadFile(t.recordPath(env))
	if errors.Is(err, os.ErrNotExist) {
		return Record{}, ErrNoExpiry
	}
	if err != nil {
		return Record{}, err
	}
	var rec Record
	if err := json.Unmarshal(data, &rec); err != nil {
		return Record{}, err
	}
	return rec, nil
}
