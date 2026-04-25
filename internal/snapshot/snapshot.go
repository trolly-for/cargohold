// Package snapshot provides point-in-time exports of secret bundles.
// A snapshot captures all key-value pairs from a bundle at a given moment,
// serialised to an encrypted file that can be shared or archived.
package snapshot

import (
	"encoding/json"
	"fmt"
	"time"

	"cargohold/internal/bundle"
	"cargohold/internal/crypto"
)

// Snapshot holds a timestamped copy of a bundle's contents.
type Snapshot struct {
	Environment string            `json:"environment"`
	CreatedAt   time.Time         `json:"created_at"`
	Secrets     map[string]string `json:"secrets"`
}

// Snapshotter creates and restores snapshots for a named environment.
type Snapshotter struct {
	env string
}

// New returns a Snapshotter for the given environment name.
func New(env string) *Snapshotter {
	return &Snapshotter{env: env}
}

// Export serialises all keys from b into an encrypted snapshot blob.
// The returned bytes can be written to disk or transmitted elsewhere.
func (s *Snapshotter) Export(b *bundle.Bundle, passphrase string) ([]byte, error) {
	keys := b.Keys()
	secrets := make(map[string]string, len(keys))
	for _, k := range keys {
		v, err := b.Get(k)
		if err != nil {
			return nil, fmt.Errorf("snapshot export: reading key %q: %w", k, err)
		}
		secrets[k] = v
	}

	snap := Snapshot{
		Environment: s.env,
		CreatedAt:   time.Now().UTC(),
		Secrets:     secrets,
	}

	plain, err := json.Marshal(snap)
	if err != nil {
		return nil, fmt.Errorf("snapshot export: marshal: %w", err)
	}

	key, salt, err := crypto.DeriveKey(passphrase, nil)
	if err != nil {
		return nil, fmt.Errorf("snapshot export: derive key: %w", err)
	}

	ciphertext, err := crypto.Encrypt(key, plain)
	if err != nil {
		return nil, fmt.Errorf("snapshot export: encrypt: %w", err)
	}

	// Prepend salt so Import can re-derive the key.
	return append(salt, ciphertext...), nil
}

// Import decrypts a snapshot blob and returns the Snapshot value.
func (s *Snapshotter) Import(data []byte, passphrase string) (*Snapshot, error) {
	const saltLen = 32
	if len(data) <= saltLen {
		return nil, fmt.Errorf("snapshot import: data too short")
	}

	salt := data[:saltLen]
	ciphertext := data[saltLen:]

	key, _, err := crypto.DeriveKey(passphrase, salt)
	if err != nil {
		return nil, fmt.Errorf("snapshot import: derive key: %w", err)
	}

	plain, err := crypto.Decrypt(key, ciphertext)
	if err != nil {
		return nil, fmt.Errorf("snapshot import: decrypt: %w", err)
	}

	var snap Snapshot
	if err := json.Unmarshal(plain, &snap); err != nil {
		return nil, fmt.Errorf("snapshot import: unmarshal: %w", err)
	}
	return &snap, nil
}
