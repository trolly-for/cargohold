// Package rotate provides passphrase rotation for encrypted bundles.
package rotate

import (
	"errors"
	"fmt"

	"cargohold/internal/bundle"
	"cargohold/internal/store"
	"cargohold/internal/vault"
)

// ErrSamePassphrase is returned when the new passphrase matches the old one.
var ErrSamePassphrase = errors.New("new passphrase must differ from the current passphrase")

// Rotator handles passphrase rotation for a named bundle.
type Rotator struct {
	store *store.Store
}

// New returns a new Rotator backed by the given store.
func New(s *store.Store) *Rotator {
	return &Rotator{store: s}
}

// Rotate decrypts the bundle at env using oldPass, then re-encrypts it with
// newPass and persists the result. The operation is atomic: the store is only
// updated after a successful round-trip.
func (r *Rotator) Rotate(env, oldPass, newPass string) error {
	if oldPass == newPass {
		return ErrSamePassphrase
	}

	// Open the existing vault with the current passphrase.
	v, err := vault.New(r.store, env)
	if err != nil {
		return fmt.Errorf("rotate: open vault: %w", err)
	}

	b, err := v.Open(oldPass)
	if err != nil {
		return fmt.Errorf("rotate: decrypt bundle: %w", err)
	}

	// Re-initialise with the new passphrase using the same bundle data.
	newVault, err := vault.New(r.store, env)
	if err != nil {
		return fmt.Errorf("rotate: create new vault: %w", err)
	}

	if err := newVault.Save(b, newPass); err != nil {
		return fmt.Errorf("rotate: save re-encrypted bundle: %w", err)
	}

	_ = b // bundle is captured inside newVault.Save
	return nil
}

// RotateWithBundle is a lower-level helper that accepts an already-decrypted
// bundle and re-encrypts it under newPass. Useful in tests and advanced flows.
func RotateWithBundle(v *vault.Vault, b *bundle.Bundle, newPass string) error {
	if err := v.Save(b, newPass); err != nil {
		return fmt.Errorf("rotate: save bundle: %w", err)
	}
	return nil
}
