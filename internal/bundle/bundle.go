package bundle

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Bundle represents a named collection of secrets for a specific environment.
type Bundle struct {
	Name      string            `json:"name"`
	Env       string            `json:"env"`
	Secrets   map[string]string `json:"secrets"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// New creates a new empty Bundle with the given name and environment.
func New(name, env string) *Bundle {
	now := time.Now().UTC()
	return &Bundle{
		Name:      name,
		Env:       env,
		Secrets:   make(map[string]string),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// Set adds or updates a secret key-value pair in the bundle.
func (b *Bundle) Set(key, value string) {
	b.Secrets[key] = value
	b.UpdatedAt = time.Now().UTC()
}

// Get retrieves a secret value by key. Returns an error if the key does not exist.
func (b *Bundle) Get(key string) (string, error) {
	v, ok := b.Secrets[key]
	if !ok {
		return "", fmt.Errorf("key %q not found in bundle %q", key, b.Name)
	}
	return v, nil
}

// Delete removes a secret key from the bundle.
func (b *Bundle) Delete(key string) {
	delete(b.Secrets, key)
	b.UpdatedAt = time.Now().UTC()
}

// SaveToFile serialises the bundle as JSON and writes it to path.
func (b *Bundle) SaveToFile(path string) error {
	data, err := json.MarshalIndent(b, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal bundle: %w", err)
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("write bundle file: %w", err)
	}
	return nil
}

// LoadFromFile reads a JSON bundle file from path and returns a Bundle.
func LoadFromFile(path string) (*Bundle, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read bundle file: %w", err)
	}
	var b Bundle
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("unmarshal bundle: %w", err)
	}
	return &b, nil
}
