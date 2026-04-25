// Package diff provides utilities for comparing two secret bundles
// and reporting keys that were added, removed, or changed between them.
package diff

import "github.com/nicholasgasior/cargohold/internal/bundle"

// Result holds the outcome of comparing two bundles.
type Result struct {
	Added   []string
	Removed []string
	Changed []string
}

// IsEmpty reports whether there are no differences between the bundles.
func (r Result) IsEmpty() bool {
	return len(r.Added) == 0 && len(r.Removed) == 0 && len(r.Changed) == 0
}

// Bundles compares two bundles and returns a Result describing the
// keys that differ. The comparison is key-level only; values are
// compared for equality but are never exposed in the result.
func Bundles(a, b *bundle.Bundle) Result {
	aKeys := keySet(a)
	bKeys := keySet(b)

	var result Result

	for k := range bKeys {
		if _, exists := aKeys[k]; !exists {
			result.Added = append(result.Added, k)
		} else {
			av, _ := a.Get(k)
			bv, _ := b.Get(k)
			if av != bv {
				result.Changed = append(result.Changed, k)
			}
		}
	}

	for k := range aKeys {
		if _, exists := bKeys[k]; !exists {
			result.Removed = append(result.Removed, k)
		}
	}

	sort(result.Added)
	sort(result.Removed)
	sort(result.Changed)

	return result
}

func keySet(b *bundle.Bundle) map[string]struct{} {
	set := make(map[string]struct{})
	for _, k := range b.Keys() {
		set[k] = struct{}{}
	}
	return set
}

func sort(ss []string) {
	for i := 1; i < len(ss); i++ {
		for j := i; j > 0 && ss[j] < ss[j-1]; j-- {
			ss[j], ss[j-1] = ss[j-1], ss[j]
		}
	}
}
