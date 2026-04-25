// Package merge provides functionality for merging two secret bundles.
// Keys from the source bundle are applied on top of the destination bundle,
// with optional overwrite control for conflicting keys.
package merge

import (
	"fmt"

	"cargohold/internal/bundle"
)

// Options controls the behaviour of a merge operation.
type Options struct {
	// Overwrite allows source keys to replace existing keys in the destination.
	Overwrite bool
}

// Result summarises what changed after a merge.
type Result struct {
	Added    []string
	Skipped  []string
	Overwritten []string
}

// Bundles merges src into dst according to opts.
// It returns a Result describing every key that was touched.
func Bundles(dst, src *bundle.Bundle, opts Options) (Result, error) {
	if dst == nil {
		return Result{}, fmt.Errorf("merge: destination bundle must not be nil")
	}
	if src == nil {
		return Result{}, fmt.Errorf("merge: source bundle must not be nil")
	}

	var res Result

	for _, key := range src.Keys() {
		val, err := src.Get(key)
		if err != nil {
			return Result{}, fmt.Errorf("merge: reading source key %q: %w", key, err)
		}

		_, exists := dst.Get(key) //nolint:errcheck
		if exists == nil {
			// key already present in destination
			if !opts.Overwrite {
				res.Skipped = append(res.Skipped, key)
				continue
			}
			res.Overwritten = append(res.Overwritten, key)
		} else {
			res.Added = append(res.Added, key)
		}

		if err := dst.Set(key, val); err != nil {
			return Result{}, fmt.Errorf("merge: writing key %q to destination: %w", key, err)
		}
	}

	return res, nil
}
