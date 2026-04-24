package vault

import (
	"errors"

	"cargohold/internal/bundle"
	"cargohold/internal/crypto"
	"cargohold/internal/store"
)

// Vault provides high-level operations for loading, mutating, and persisting
// encrypted secret bundles.
type Vault struct {
	store      *store.Store
	passphrase string
}

// New creates a Vault backed by the given store and passphrase.
func New(s *store.Store, passphrase string) *Vault {
	return &Vault{store: s, passphrase: passphrase}
}

// Open decrypts and loads an existing bundle from the store.
// Returns an error if the bundle does not exist or the passphrase is wrong.
func (v *Vault) Open(name string) (*bundle.Bundle, error) {
	if !v.store.Exists(name) {
		return nil, errors.New("bundle not found: " + name)
	}

	ciphertext, err := v.store.Load(name)
	if err != nil {
		return nil, err
	}

	key, err := crypto.DeriveKey(v.passphrase, nil)
	if err != nil {
		return nil, err
	}

	plaintext, err := crypto.Decrypt(key, ciphertext)
	if err != nil {
		return nil, errors.New("failed to decrypt bundle (wrong passphrase?)")
	}

	return bundle.Unmarshal(plaintext)
}

// Save encrypts and persists a bundle to the store.
func (v *Vault) Save(name string, b *bundle.Bundle) error {
	plaintext, err := bundle.Marshal(b)
	if err != nil {
		return err
	}

	key, err := crypto.DeriveKey(v.passphrase, nil)
	if err != nil {
		return err
	}

	ciphertext, err := crypto.Encrypt(key, plaintext)
	if err != nil {
		return err
	}

	return v.store.Save(name, ciphertext)
}

// Init creates a new empty bundle and saves it under name.
// Returns an error if a bundle with that name already exists.
func (v *Vault) Init(name string) error {
	if v.store.Exists(name) {
		return errors.New("bundle already exists: " + name)
	}
	return v.Save(name, bundle.New())
}

// Delete removes a bundle from the store.
func (v *Vault) Delete(name string) error {
	return v.store.Delete(name)
}

// List returns all bundle names known to the store.
func (v *Vault) List() ([]string, error) {
	return v.store.List()
}
