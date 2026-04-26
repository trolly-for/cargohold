package importer_test

import (
	"strings"
	"testing"

	"cargohold/internal/bundle"
	importer "cargohold/internal/import"
)

func seedBundle(t *testing.T) *bundle.Bundle {
	t.Helper()
	b, err := bundle.New("test")
	if err != nil {
		t.Fatalf("bundle.New: %v", err)
	}
	return b
}

func TestParseFormatValid(t *testing.T) {
	cases := []struct{ in, want string }{
		{"dotenv", "dotenv"},
		{".env", "dotenv"},
		{"json", "json"},
		{"JSON", "json"},
	}
	for _, c := range cases {
		f, err := importer.ParseFormat(c.in)
		if err != nil {
			t.Errorf("ParseFormat(%q) unexpected error: %v", c.in, err)
		}
		if string(f) != c.want {
			t.Errorf("ParseFormat(%q) = %q, want %q", c.in, f, c.want)
		}
	}
}

func TestParseFormatInvalid(t *testing.T) {
	_, err := importer.ParseFormat("toml")
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestImportDotenv(t *testing.T) {
	src := `# comment
DB_HOST=localhost
DB_PORT=5432
SECRET_KEY="abc123"
OTHER='hello world'
`
	b := seedBundle(t)
	n, err := importer.Import(b, strings.NewReader(src), importer.FormatDotenv)
	if err != nil {
		t.Fatalf("Import: %v", err)
	}
	if n != 4 {
		t.Errorf("imported %d keys, want 4", n)
	}
	if v, _ := b.Get("DB_HOST"); v != "localhost" {
		t.Errorf("DB_HOST = %q, want %q", v, "localhost")
	}
	if v, _ := b.Get("SECRET_KEY"); v != "abc123" {
		t.Errorf("SECRET_KEY = %q, want %q", v, "abc123")
	}
	if v, _ := b.Get("OTHER"); v != "hello world" {
		t.Errorf("OTHER = %q, want %q", v, "hello world")
	}
}

func TestImportDotenvInvalidLine(t *testing.T) {
	src := "NOEQUALS\n"
	b := seedBundle(t)
	_, err := importer.Import(b, strings.NewReader(src), importer.FormatDotenv)
	if err == nil {
		t.Fatal("expected error for line without '='")
	}
}

func TestImportJSON(t *testing.T) {
	src := `{"API_KEY":"secret","REGION":"us-east-1"}`
	b := seedBundle(t)
	n, err := importer.Import(b, strings.NewReader(src), importer.FormatJSON)
	if err != nil {
		t.Fatalf("Import: %v", err)
	}
	if n != 2 {
		t.Errorf("imported %d keys, want 2", n)
	}
	if v, _ := b.Get("API_KEY"); v != "secret" {
		t.Errorf("API_KEY = %q, want %q", v, "secret")
	}
}

func TestImportJSONInvalid(t *testing.T) {
	b := seedBundle(t)
	_, err := importer.Import(b, strings.NewReader(`not json`), importer.FormatJSON)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestImportDotenvEmpty(t *testing.T) {
	b := seedBundle(t)
	n, err := importer.Import(b, strings.NewReader("# only comments\n\n"), importer.FormatDotenv)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 imports, got %d", n)
	}
}
