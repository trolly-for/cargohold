package merge_test

import (
	"testing"

	"cargohold/internal/bundle"
	"cargohold/internal/merge"
)

func seedBundle(t *testing.T, pairs map[string]string) *bundle.Bundle {
	t.Helper()
	b := bundle.New("test")
	for k, v := range pairs {
		if err := b.Set(k, v); err != nil {
			t.Fatalf("seed: Set(%q): %v", k, err)
		}
	}
	return b
}

func TestMergeAddsNewKeys(t *testing.T) {
	dst := seedBundle(t, map[string]string{"A": "1"})
	src := seedBundle(t, map[string]string{"B": "2", "C": "3"})

	res, err := merge.Bundles(dst, src, merge.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Added) != 2 {
		t.Errorf("expected 2 added, got %d", len(res.Added))
	}
	if len(res.Skipped) != 0 || len(res.Overwritten) != 0 {
		t.Errorf("expected no skipped/overwritten")
	}
}

func TestMergeSkipsConflictsByDefault(t *testing.T) {
	dst := seedBundle(t, map[string]string{"A": "original"})
	src := seedBundle(t, map[string]string{"A": "new", "B": "2"})

	res, err := merge.Bundles(dst, src, merge.Options{Overwrite: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "A" {
		t.Errorf("expected A to be skipped, got %v", res.Skipped)
	}
	val, _ := dst.Get("A")
	if val != "original" {
		t.Errorf("expected original value preserved, got %q", val)
	}
}

func TestMergeOverwritesWhenEnabled(t *testing.T) {
	dst := seedBundle(t, map[string]string{"A": "original"})
	src := seedBundle(t, map[string]string{"A": "new"})

	res, err := merge.Bundles(dst, src, merge.Options{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Overwritten) != 1 {
		t.Errorf("expected 1 overwritten, got %d", len(res.Overwritten))
	}
	val, _ := dst.Get("A")
	if val != "new" {
		t.Errorf("expected value updated to 'new', got %q", val)
	}
}

func TestMergeNilDstErrors(t *testing.T) {
	src := seedBundle(t, map[string]string{"A": "1"})
	_, err := merge.Bundles(nil, src, merge.Options{})
	if err == nil {
		t.Fatal("expected error for nil dst")
	}
}

func TestMergeNilSrcErrors(t *testing.T) {
	dst := seedBundle(t, map[string]string{"A": "1"})
	_, err := merge.Bundles(dst, nil, merge.Options{})
	if err == nil {
		t.Fatal("expected error for nil src")
	}
}

func TestMergeEmptySourceIsNoop(t *testing.T) {
	dst := seedBundle(t, map[string]string{"A": "1"})
	src := bundle.New("empty")

	res, err := merge.Bundles(dst, src, merge.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Added)+len(res.Skipped)+len(res.Overwritten) != 0 {
		t.Errorf("expected empty result for empty source")
	}
}

func TestMergeOverwritePreservesNonConflictingKeys(t *testing.T) {
	dst := seedBundle(t, map[string]string{"A": "original", "B": "keep"})
	src := seedBundle(t, map[string]string{"A": "new", "C": "added"})

	res, err := merge.Bundles(dst, src, merge.Options{Overwrite: true})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Overwritten) != 1 || res.Overwritten[0] != "A" {
		t.Errorf("expected only A overwritten, got %v", res.Overwritten)
	}
	if len(res.Added) != 1 || res.Added[0] != "C" {
		t.Errorf("expected only C added, got %v", res.Added)
	}
	val, _ := dst.Get("B")
	if val != "keep" {
		t.Errorf("expected B to be unchanged, got %q", val)
	}
}
