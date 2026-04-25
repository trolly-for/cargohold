package snapshot_test

import (
	"testing"

	"cargohold/internal/bundle"
	"cargohold/internal/snapshot"
)

const testPass = "correct-horse-battery-staple"

func seedBundle(t *testing.T) *bundle.Bundle {
	t.Helper()
	b := bundle.New()
	if err := b.Set("DB_HOST", "localhost"); err != nil {
		t.Fatalf("seed: %v", err)
	}
	if err := b.Set("API_KEY", "supersecret"); err != nil {
		t.Fatalf("seed: %v", err)
	}
	return b
}

func TestExportImportRoundtrip(t *testing.T) {
	s := snapshot.New("staging")
	b := seedBundle(t)

	blob, err := s.Export(b, testPass)
	if err != nil {
		t.Fatalf("Export: %v", err)
	}

	snap, err := s.Import(blob, testPass)
	if err != nil {
		t.Fatalf("Import: %v", err)
	}

	if snap.Environment != "staging" {
		t.Errorf("environment: got %q, want %q", snap.Environment, "staging")
	}
	if snap.Secrets["DB_HOST"] != "localhost" {
		t.Errorf("DB_HOST: got %q", snap.Secrets["DB_HOST"])
	}
	if snap.Secrets["API_KEY"] != "supersecret" {
		t.Errorf("API_KEY: got %q", snap.Secrets["API_KEY"])
	}
	if snap.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestImportWrongPassphrase(t *testing.T) {
	s := snapshot.New("prod")
	b := seedBundle(t)

	blob, err := s.Export(b, testPass)
	if err != nil {
		t.Fatalf("Export: %v", err)
	}

	_, err = s.Import(blob, "wrong-passphrase")
	if err == nil {
		t.Fatal("expected error with wrong passphrase, got nil")
	}
}

func TestImportTruncatedData(t *testing.T) {
	s := snapshot.New("dev")
	_, err := s.Import([]byte("short"), testPass)
	if err == nil {
		t.Fatal("expected error for truncated data")
	}
}

func TestExportEmptyBundle(t *testing.T) {
	s := snapshot.New("dev")
	b := bundle.New()

	blob, err := s.Export(b, testPass)
	if err != nil {
		t.Fatalf("Export empty bundle: %v", err)
	}

	snap, err := s.Import(blob, testPass)
	if err != nil {
		t.Fatalf("Import empty bundle: %v", err)
	}
	if len(snap.Secrets) != 0 {
		t.Errorf("expected 0 secrets, got %d", len(snap.Secrets))
	}
}
