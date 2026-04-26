// Package copy provides functionality for duplicating secret bundles
// between environments, with optional key filtering.
package copy

import (
	"errors"
	"fmt"

	"cargohold/internal/bundle"
)

// Options controls the behaviour of a copy operation.
type Options struct {
	// Keys restricts the copy to only the specified keys.
	// If empty, all keys are copied.
	Keys []string
	// Overwrite allows existing keys in dst to be replaced.
	Overwrite bool
}

// Copier copies keys from one bundle into another.
type Copier struct {
	src *bundle.Bundle
}

// New returns a Copier that reads from src.
func New(src *bundle.Bundle) (*Copier, error) {
	if src == nil {
		return nil, errors.New("copy: source bundle must not be nil")
	}
	return &Copier{src: src}, nil
}

// Into copies keys from the source bundle into dst according to opts.
// It returns the number of keys that were written.
func (c *Copier) Into(dst *bundle.Bundle, opts Options) (int, error) {
	if dst == nil {
		return 0, errors.New("copy: destination bundle must not be nil")
	}

	keys := opts.Keys
	if len(keys) == 0 {
		keys = c.src.Keys()
	}

	written := 0
	for _, k := range keys {
		v, ok := c.src.Get(k)
		if !ok {
			return written, fmt.Errorf("copy: key %q not found in source bundle", k)
		}
		_, exists := dst.Get(k)
		if exists && !opts.Overwrite {
			continue
		}
		if err := dst.Set(k, v); err != nil {
			return written, fmt.Errorf("copy: failed to set key %q: %w", k, err)
		}
		written++
	}
	return written, nil
}
