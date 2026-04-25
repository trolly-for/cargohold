package diff_test

import (
	"testing"

	"github.com/nicholasgasior/cargohold/internal/bundle"
	"github.com/nicholasgasior/cargohold/internal/diff"
)

func seedBundle(t *testing.T, pairs ...string) *bundle.Bundle {
	t.Helper()
	b := bundle.New()
	for i := 0; i+1 < len(pairs); i += 2 {
		if err := b.Set(pairs[i], pairs[i+1]); err != nil {
			t.Fatalf("seed: set %s: %v", pairs[i], err)
		}
	}
	return b
}

func TestIdenticalBundlesProduceEmptyResult(t *testing.T) {
	a := seedBundle(t, "KEY_A", "val1", "KEY_B", "val2")
	b := seedBundle(t, "KEY_A", "val1", "KEY_B", "val2")
	r := diff.Bundles(a, b)
	if !r.IsEmpty() {
		t.Errorf("expected empty diff, got %+v", r)
	}
}

func TestAddedKeys(t *testing.T) {
	a := seedBundle(t, "KEY_A", "val1")
	b := seedBundle(t, "KEY_A", "val1", "KEY_B", "val2")
	r := diff.Bundles(a, b)
	if len(r.Added) != 1 || r.Added[0] != "KEY_B" {
		t.Errorf("expected Added=[KEY_B], got %v", r.Added)
	}
	if len(r.Removed) != 0 || len(r.Changed) != 0 {
		t.Errorf("unexpected removed/changed: %+v", r)
	}
}

func TestRemovedKeys(t *testing.T) {
	a := seedBundle(t, "KEY_A", "val1", "KEY_B", "val2")
	b := seedBundle(t, "KEY_A", "val1")
	r := diff.Bundles(a, b)
	if len(r.Removed) != 1 || r.Removed[0] != "KEY_B" {
		t.Errorf("expected Removed=[KEY_B], got %v", r.Removed)
	}
}

func TestChangedKeys(t *testing.T) {
	a := seedBundle(t, "KEY_A", "old")
	b := seedBundle(t, "KEY_A", "new")
	r := diff.Bundles(a, b)
	if len(r.Changed) != 1 || r.Changed[0] != "KEY_A" {
		t.Errorf("expected Changed=[KEY_A], got %v", r.Changed)
	}
	if len(r.Added) != 0 || len(r.Removed) != 0 {
		t.Errorf("unexpected added/removed: %+v", r)
	}
}

func TestEmptyBundles(t *testing.T) {
	a := bundle.New()
	b := bundle.New()
	r := diff.Bundles(a, b)
	if !r.IsEmpty() {
		t.Errorf("expected empty diff for two empty bundles, got %+v", r)
	}
}

func TestResultsAreSorted(t *testing.T) {
	a := seedBundle(t, "Z_KEY", "v", "A_KEY", "v")
	b := seedBundle(t, "M_KEY", "v", "B_KEY", "v")
	r := diff.Bundles(a, b)
	for i := 1; i < len(r.Added); i++ {
		if r.Added[i] < r.Added[i-1] {
			t.Errorf("Added not sorted: %v", r.Added)
		}
	}
	for i := 1; i < len(r.Removed); i++ {
		if r.Removed[i] < r.Removed[i-1] {
			t.Errorf("Removed not sorted: %v", r.Removed)
		}
	}
}
