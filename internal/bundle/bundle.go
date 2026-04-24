package bundle

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/cargohold/cargohold/internal/crypto"
)

// Bundle holds a named collection of key-value secret pairs.
type Bundle struct {
	Name    string            `json:"name"`
	Secrets map[string]string `json:"secrets"`
}

// New creates a new empty Bundle with the given name.
func New(name string) *Bundle {
	return &Bundle{
		Name:    name,
		Secrets: make(map[string]string),
	}
}

// Set adds or updates a secret key-value pair in the bundle.
func (b *Bundle) Set(key, value string) {
	b.Secrets[key] = value
}

// Get retrieves a secret value by key. Returns an error if not found.
func (b *Bundle) Get(key string) (string, error) {
	val, ok := b.Secrets[key]
	if !ok {
		return "", errors.New("key not found: " + key)
	}
	return val, nil
}

// Delete removes a secret key from the bundle.
func (b *Bundle) Delete(key string) {
	delete(b.Secrets, key)
}

// Save encrypts and writes the bundle to the given file path using the passphrase.
func (b *Bundle) Save(path, passphrase string) error {
	data, err := json.Marshal(b)
	if err != nil {
		return err
	}

	key := crypto.DeriveKey(passphrase)
	ciphertext, err := crypto.Encrypt(key, data)
	if err != nil {
		return err
	}

	return os.WriteFile(path, ciphertext, 0600)
}

// Load reads and decrypts a bundle from the given file path using the passphrase.
func Load(path, passphrase string) (*Bundle, error) {
	ciphertext, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	key := crypto.DeriveKey(passphrase)
	data, err := crypto.Decrypt(key, ciphertext)
	if err != nil {
		return nil, err
	}

	var b Bundle
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, err
	}

	return &b, nil
}
