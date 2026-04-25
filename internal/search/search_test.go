package search_test

import (
	"testing"

	"cargohold/internal/bundle"
	"cargohold/internal/search"
)

func seedBundle(t *testing.T) *bundle.Bundle {
	t.Helper()
	b := bundle.New("test")
	pairs := [][2]string{
		{"DB_HOST", "localhost"},
		{"DB_PORT", "5432"},
		{"DB_PASSWORD", "secret"},
		{"APP_SECRET", "topsecret"},
		{"APP_DEBUG", "true"},
		{"LOG_LEVEL", "info"},
	}
	for _, p := range pairs {
		if err := b.Set(p[0], p[1]); err != nil {
			t.Fatalf("seed Set(%q): %v", p[0], err)
		}
	}
	return b
}

func TestByPrefixReturnsMatches(t *testing.T) {
	b := seedBundle(t)
	results := search.ByPrefix(b, "DB_")
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if results[0].Key != "DB_HOST" {
		t.Errorf("unexpected first key: %s", results[0].Key)
	}
}

func TestByPrefixCaseInsensitive(t *testing.T) {
	b := seedBundle(t)
	results := search.ByPrefix(b, "app_")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestByPrefixNoMatch(t *testing.T) {
	b := seedBundle(t)
	results := search.ByPrefix(b, "REDIS_")
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestBySubstringReturnsMatches(t *testing.T) {
	b := seedBundle(t)
	results := search.BySubstring(b, "SECRET")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestBySubstringCaseInsensitive(t *testing.T) {
	b := seedBundle(t)
	results := search.BySubstring(b, "debug")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "APP_DEBUG" {
		t.Errorf("unexpected key: %s", results[0].Key)
	}
}

func TestHasValueReturnsMatches(t *testing.T) {
	b := seedBundle(t)
	results := search.HasValue(b, "secret")
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Key != "DB_PASSWORD" {
		t.Errorf("unexpected key: %s", results[0].Key)
	}
}

func TestHasValueNoMatch(t *testing.T) {
	b := seedBundle(t)
	results := search.HasValue(b, "notpresent")
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got %d", len(results))
	}
}

func TestNilBundleReturnsNil(t *testing.T) {
	if r := search.ByPrefix(nil, "DB"); r != nil {
		t.Error("expected nil for nil bundle")
	}
	if r := search.BySubstring(nil, "x"); r != nil {
		t.Error("expected nil for nil bundle")
	}
	if r := search.HasValue(nil, "x"); r != nil {
		t.Error("expected nil for nil bundle")
	}
}
