// Package tag provides key tagging and filtering for secret bundles.
// Tags are arbitrary string labels attached to individual keys, enabling
// grouping and selective export of secrets by category (e.g. "db", "api").
package tag

import (
	"fmt"
	"regexp"
	"sort"
)

var validTag = regexp.MustCompile(`^[a-z][a-z0-9_-]{0,31}$`)

// Map holds a mapping of key → set of tags.
type Map map[string][]string

// New returns an empty Map.
func New() Map {
	return make(Map)
}

// Validate returns an error if the tag name is not valid.
func Validate(tag string) error {
	if !validTag.MatchString(tag) {
		return fmt.Errorf("tag %q is invalid: must match [a-z][a-z0-9_-]{0,31}", tag)
	}
	return nil
}

// Add attaches tag t to key k. Duplicate tags are ignored.
func (m Map) Add(key, t string) error {
	if err := Validate(t); err != nil {
		return err
	}
	for _, existing := range m[key] {
		if existing == t {
			return nil
		}
	}
	m[key] = append(m[key], t)
	sort.Strings(m[key])
	return nil
}

// Remove detaches tag t from key k. No-op if the tag is not present.
func (m Map) Remove(key, t string) {
	tags := m[key]
	filtered := tags[:0]
	for _, existing := range tags {
		if existing != t {
			filtered = append(filtered, existing)
		}
	}
	if len(filtered) == 0 {
		delete(m, key)
	} else {
		m[key] = filtered
	}
}

// KeysWithTag returns all keys that carry the given tag, sorted.
func (m Map) KeysWithTag(t string) []string {
	var result []string
	for key, tags := range m {
		for _, tag := range tags {
			if tag == t {
				result = append(result, key)
				break
			}
		}
	}
	sort.Strings(result)
	return result
}

// Tags returns the sorted list of tags for key k.
func (m Map) Tags(key string) []string {
	if tags, ok := m[key]; ok {
		out := make([]string, len(tags))
		copy(out, tags)
		return out
	}
	return nil
}
