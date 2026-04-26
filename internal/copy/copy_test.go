package copy_test

import (
	"testing"

	"cargohold/internal/bundle"
	copy_ "cargohold/internal/copy"
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

func TestCopyAllKeys(t *testing.T) {
	src := seedBundle(t, map[string]string{"A": "1", "B": "2"})
	dst := bundle.New("dst")
	c, err := copy_.New(src)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	n, err := c.Into(dst, copy_.Options{})
	if err != nil {
		t.Fatalf("Into: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 written, got %d", n)
	}
	for _, k := range []string{"A", "B"} {
		if _, ok := dst.Get(k); !ok {
			t.Errorf("key %q missing from dst", k)
		}
	}
}

func TestCopySelectedKeys(t *testing.T) {
	src := seedBundle(t, map[string]string{"X": "x", "Y": "y", "Z": "z"})
	dst := bundle.New("dst")
	c, _ := copy_.New(src)
	n, err := c.Into(dst, copy_.Options{Keys: []string{"X", "Z"}})
	if err != nil {
		t.Fatalf("Into: %v", err)
	}
	if n != 2 {
		t.Errorf("expected 2 written, got %d", n)
	}
	if _, ok := dst.Get("Y"); ok {
		t.Error("key Y should not have been copied")
	}
}

func TestCopySkipsExistingWithoutOverwrite(t *testing.T) {
	src := seedBundle(t, map[string]string{"K": "new"})
	dst := seedBundle(t, map[string]string{"K": "old"})
	c, _ := copy_.New(src)
	n, err := c.Into(dst, copy_.Options{Overwrite: false})
	if err != nil {
		t.Fatalf("Into: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 written, got %d", n)
	}
	if v, _ := dst.Get("K"); v != "old" {
		t.Errorf("expected old value preserved, got %q", v)
	}
}

func TestCopyOverwritesWhenEnabled(t *testing.T) {
	src := seedBundle(t, map[string]string{"K": "new"})
	dst := seedBundle(t, map[string]string{"K": "old"})
	c, _ := copy_.New(src)
	_, err := c.Into(dst, copy_.Options{Overwrite: true})
	if err != nil {
		t.Fatalf("Into: %v", err)
	}
	if v, _ := dst.Get("K"); v != "new" {
		t.Errorf("expected new value, got %q", v)
	}
}

func TestCopyMissingKeyErrors(t *testing.T) {
	src := seedBundle(t, map[string]string{"A": "1"})
	dst := bundle.New("dst")
	c, _ := copy_.New(src)
	_, err := c.Into(dst, copy_.Options{Keys: []string{"MISSING"}})
	if err == nil {
		t.Fatal("expected error for missing key, got nil")
	}
}

func TestCopyNilSrcErrors(t *testing.T) {
	_, err := copy_.New(nil)
	if err == nil {
		t.Fatal("expected error for nil src")
	}
}

func TestCopyNilDstErrors(t *testing.T) {
	src := seedBundle(t, map[string]string{"A": "1"})
	c, _ := copy_.New(src)
	_, err := c.Into(nil, copy_.Options{})
	if err == nil {
		t.Fatal("expected error for nil dst")
	}
}
