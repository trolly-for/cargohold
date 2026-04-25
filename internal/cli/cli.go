// Package cli wires together the cargohold sub-commands.
package cli

import (
	"errors"
	"fmt"

	"cargohold/internal/bundle"
	"cargohold/internal/env"
	"cargohold/internal/output"
	"cargohold/internal/rotate"
	"cargohold/internal/store"
	"cargohold/internal/vault"
)

// Runner executes cargohold CLI commands.
type Runner struct {
	store     *store.Store
	formatter *output.Formatter
}

// New returns a Runner using the default store location.
func New() (*Runner, error) {
	s, err := store.Default()
	if err != nil {
		return nil, err
	}
	return &Runner{store: s, formatter: output.Default()}, nil
}

// Init creates a new encrypted bundle for the given environment.
func (r *Runner) Init(environment, passphrase string) error {
	e, err := env.Normalize(environment)
	if err != nil {
		return err
	}
	v, err := vault.New(r.store, e)
	if err != nil {
		return err
	}
	b := bundle.New()
	if err := v.Init(b, passphrase); err != nil {
		return err
	}
	r.formatter.Success("initialised bundle for environment: " + e)
	return nil
}

// Set encrypts and stores a key/value pair in the named environment bundle.
func (r *Runner) Set(environment, passphrase, key, value string) error {
	e, err := env.Normalize(environment)
	if err != nil {
		return err
	}
	v, err := vault.New(r.store, e)
	if err != nil {
		return err
	}
	b, err := v.Open(passphrase)
	if err != nil {
		return err
	}
	b.Set(key, value)
	if err := v.Save(b, passphrase); err != nil {
		return err
	}
	r.formatter.Success(fmt.Sprintf("set %s in %s", key, e))
	return nil
}

// Get retrieves a value by key from the named environment bundle.
func (r *Runner) Get(environment, passphrase, key string) error {
	e, err := env.Normalize(environment)
	if err != nil {
		return err
	}
	v, err := vault.New(r.store, e)
	if err != nil {
		return err
	}
	b, err := v.Open(passphrase)
	if err != nil {
		return err
	}
	val, ok := b.Get(key)
	if !ok {
		return fmt.Errorf("key %q not found in %s", key, e)
	}
	r.formatter.KeyValue(key, val)
	return nil
}

// List prints all keys stored in the named environment bundle.
func (r *Runner) List(environment, passphrase string) error {
	e, err := env.Normalize(environment)
	if err != nil {
		return err
	}
	v, err := vault.New(r.store, e)
	if err != nil {
		return err
	}
	b, err := v.Open(passphrase)
	if err != nil {
		return err
	}
	r.formatter.KeyList(b.Keys())
	return nil
}

// Delete removes a key from the named environment bundle.
func (r *Runner) Delete(environment, passphrase, key string) error {
	e, err := env.Normalize(environment)
	if err != nil {
		return err
	}
	v, err := vault.New(r.store, e)
	if err != nil {
		return err
	}
	b, err := v.Open(passphrase)
	if err != nil {
		return err
	}
	if ok := b.Delete(key); !ok {
		return fmt.Errorf("key %q not found in %s", key, e)
	}
	if err := v.Save(b, passphrase); err != nil {
		return err
	}
	r.formatter.Success(fmt.Sprintf("deleted %s from %s", key, e))
	return nil
}

// Rotate re-encrypts the named environment bundle under a new passphrase.
func (r *Runner) Rotate(environment, oldPass, newPass string) error {
	e, err := env.Normalize(environment)
	if err != nil {
		return err
	}
	rot := rotate.New(r.store)
	if err := rot.Rotate(e, oldPass, newPass); err != nil {
		if errors.Is(err, rotate.ErrSamePassphrase) {
			return err
		}
		return fmt.Errorf("rotate failed: %w", err)
	}
	r.formatter.Success(fmt.Sprintf("passphrase rotated for environment: %s", e))
	return nil
}
