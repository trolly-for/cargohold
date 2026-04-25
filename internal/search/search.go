// Package search provides key lookup and filtering utilities for secret bundles.
package search

import (
	"strings"

	"cargohold/internal/bundle"
)

// Result holds a matched key and its value.
type Result struct {
	Key   string
	Value string
}

// ByPrefix returns all entries whose keys start with the given prefix.
// The prefix match is case-insensitive.
func ByPrefix(b *bundle.Bundle, prefix string) []Result {
	if b == nil {
		return nil
	}
	lower := strings.ToLower(prefix)
	var results []Result
	for _, k := range b.Keys() {
		if strings.HasPrefix(strings.ToLower(k), lower) {
			v, _ := b.Get(k)
			results = append(results, Result{Key: k, Value: v})
		}
	}
	sortResults(results)
	return results
}

// BySubstring returns all entries whose keys contain the given substring.
// The match is case-insensitive.
func BySubstring(b *bundle.Bundle, substr string) []Result {
	if b == nil {
		return nil
	}
	lower := strings.ToLower(substr)
	var results []Result
	for _, k := range b.Keys() {
		if strings.Contains(strings.ToLower(k), lower) {
			v, _ := b.Get(k)
			results = append(results, Result{Key: k, Value: v})
		}
	}
	sortResults(results)
	return results
}

// HasValue returns all entries whose values equal the given string (exact, case-sensitive).
func HasValue(b *bundle.Bundle, value string) []Result {
	if b == nil {
		return nil
	}
	var results []Result
	for _, k := range b.Keys() {
		v, _ := b.Get(k)
		if v == value {
			results = append(results, Result{Key: k, Value: v})
		}
	}
	sortResults(results)
	return results
}

// sortResults sorts a result slice by key lexicographically.
func sortResults(r []Result) {
	for i := 1; i < len(r); i++ {
		for j := i; j > 0 && r[j].Key < r[j-1].Key; j-- {
			r[j], r[j-1] = r[j-1], r[j]
		}
	}
}
