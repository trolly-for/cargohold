package cli_test

import (
	"os"
	"path/filepath"
	"testing"

	"cargohold/internal/cli"
)

func writeTempFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return p
}

func TestImportDotenvIntoBundle(t *testing.T) {
	dir := t.TempDir()
	r := tempRunner(t, dir)

	t.Setenv("CARGOHOLD_PASSPHRASE", "supersecret123")

	if err := r.Run([]string{"init", "dev"}); err != nil {
		t.Fatalf("init: %v", err)
	}

	dotenv := writeTempFile(t, dir, "secrets.env", "FOO=bar\nBAZ=qux\n")

	if err := r.Run([]string{"import", "dev", dotenv, "--format", "dotenv"}); err != nil {
		t.Fatalf("import: %v", err)
	}

	out, err := r.Run([]string{"get", "dev", "FOO"}); _ = out
	if err != nil {
		t.Fatalf("get FOO after import: %v", err)
	}
}

func TestImportJSONIntoBundle(t *testing.T) {
	dir := t.TempDir()
	r := tempRunner(t, dir)

	t.Setenv("CARGOHOLD_PASSPHRASE", "supersecret123")

	if err := r.Run([]string{"init", "staging"}); err != nil {
		t.Fatalf("init: %v", err)
	}

	jsonFile := writeTempFile(t, dir, "secrets.json", `{"API_KEY":"tok123","REGION":"eu-west-1"}`)

	if err := r.Run([]string{"import", "staging", jsonFile, "--format", "json"}); err != nil {
		t.Fatalf("import json: %v", err)
	}
}

func TestImportMissingFileErrors(t *testing.T) {
	dir := t.TempDir()
	r := tempRunner(t, dir)

	t.Setenv("CARGOHOLD_PASSPHRASE", "supersecret123")

	if err := r.Run([]string{"init", "dev"}); err != nil {
		t.Fatalf("init: %v", err)
	}

	err := r.Run([]string{"import", "dev", filepath.Join(dir, "nonexistent.env")})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestImportTooFewArgsErrors(t *testing.T) {
	dir := t.TempDir()
	r := tempRunner(t, dir)
	err := r.Run([]string{"import", "dev"})
	if err == nil {
		t.Fatal("expected usage error")
	}
}

var _ = cli.NewWithStore // ensure export_test.go symbol is referenced
