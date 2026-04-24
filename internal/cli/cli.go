// Package cli provides the command-line interface for cargohold.
package cli

import (
	"errors"
	"fmt"
	"os"

	"cargohold/internal/env"
	"cargohold/internal/store"
	"cargohold/internal/vault"
)

// Runner holds the dependencies needed to execute CLI commands.
type Runner struct {
	Store *store.Store
	out   *os.File
}

// New creates a Runner using the default store location.
func New() (*Runner, error) {
	s, err := store.Default()
	if err != nil {
		return nil, fmt.Errorf("cli: init store: %w", err)
	}
	return &Runner{Store: s, out: os.Stdout}, nil
}

// Init creates a new encrypted bundle for the given environment.
func (r *Runner) Init(environment, passphrase string) error {
	environment, err := env.Normalize(environment)
	if err != nil {
		return err
	}
	path, err := r.Store.BundlePath(environment)
	if err != nil {
		return err
	}
	_, err = vault.Init(path, passphrase)
	if err != nil {
		return fmt.Errorf("cli: init %s: %w", environment, err)
	}
	fmt.Fprintf(r.out, "Initialized bundle for environment %q\n", environment)
	return nil
}

// Set stores a key/value pair in the named environment bundle.
func (r *Runner) Set(environment, passphrase, key, value string) error {
	environment, err := env.Normalize(environment)
	if err != nil {
		return err
	}
	v, err := r.openVault(environment, passphrase)
	if err != nil {
		return err
	}
	v.Bundle().Set(key, value)
	if err := v.Save(passphrase); err != nil {
		return fmt.Errorf("cli: set %s/%s: %w", environment, key, err)
	}
	fmt.Fprintf(r.out, "Set %q in %q\n", key, environment)
	return nil
}

// Get retrieves a value by key from the named environment bundle.
func (r *Runner) Get(environment, passphrase, key string) (string, error) {
	environment, err := env.Normalize(environment)
	if err != nil {
		return "", err
	}
	v, err := r.openVault(environment, passphrase)
	if err != nil {
		return "", err
	}
	val, ok := v.Bundle().Get(key)
	if !ok {
		return "", fmt.Errorf("cli: key %q not found in %q", key, environment)
	}
	return val, nil
}

// List prints all keys stored in the named environment bundle.
func (r *Runner) List(environment, passphrase string) ([]string, error) {
	environment, err := env.Normalize(environment)
	if err != nil {
		return nil, err
	}
	v, err := r.openVault(environment, passphrase)
	if err != nil {
		return nil, err
	}
	return v.Bundle().Keys(), nil
}

func (r *Runner) openVault(environment, passphrase string) (*vault.Vault, error) {
	path, err := r.Store.BundlePath(environment)
	if err != nil {
		return nil, err
	}
	if !r.Store.Exists(environment) {
		return nil, errors.New("cli: bundle not found; run init first")
	}
	v, err := vault.Open(path, passphrase)
	if err != nil {
		return nil, fmt.Errorf("cli: open %s: %w", environment, err)
	}
	return v, nil
}
