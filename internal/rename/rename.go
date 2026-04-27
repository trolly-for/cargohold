// Package rename provides key renaming functionality for secret bundles.
// It allows a key to be moved to a new name within the same bundle,
// optionally overwriting an existing key at the destination.
package rename

import (
	"errors"
	"fmt"

	"cargohold/internal/bundle"
)

// ErrSameKey is returned when the source and destination keys are identical.
var ErrSameKey = errors.New("rename: source and destination keys are the same")

// ErrDestExists is returned when the destination key already exists and
// overwrite is not enabled.
var ErrDestExists = errors.New("rename: destination key already exists")

// ErrSrcMissing is returned when the source key does not exist in the bundle.
var ErrSrcMissing = errors.New("rename: source key not found")

// Options controls the behaviour of a rename operation.
type Options struct {
	// Overwrite allows the destination key to be replaced if it already exists.
	Overwrite bool
}

// Key renames oldKey to newKey inside b, applying opts.
// The value associated with oldKey is preserved; oldKey is removed.
func Key(b *bundle.Bundle, oldKey, newKey string, opts Options) error {
	if oldKey == newKey {
		return ErrSameKey
	}

	val, ok := b.Get(oldKey)
	if !ok {
		return fmt.Errorf("%w: %q", ErrSrcMissing, oldKey)
	}

	if _, exists := b.Get(newKey); exists && !opts.Overwrite {
		return fmt.Errorf("%w: %q", ErrDestExists, newKey)
	}

	if err := b.Set(newKey, val); err != nil {
		return fmt.Errorf("rename: set destination: %w", err)
	}

	b.Delete(oldKey)
	return nil
}
