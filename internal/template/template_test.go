package template_test

import (
	"strings"
	"testing"

	"cargohold/internal/bundle"
	"cargohold/internal/template"
)

func seedBundle(t *testing.T) *bundle.Bundle {
	t.Helper()
	b := bundle.New("test")
	b.Set("DB_HOST", "localhost")
	b.Set("DB_PASS", "s3cr3t")
	b.Set("API_KEY", "abc123")
	return b
}

func TestParseFormatValid(t *testing.T) {
	for _, tc := range []struct {
		input string
		want  template.Format
	}{
		{"export", template.FormatExport},
		{"EXPORT", template.FormatExport},
		{"dotenv", template.FormatDotenv},
		{"Dotenv", template.FormatDotenv},
	} {
		got, err := template.ParseFormat(tc.input)
		if err != nil {
			t.Fatalf("ParseFormat(%q): unexpected error: %v", tc.input, err)
		}
		if got != tc.want {
			t.Errorf("ParseFormat(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestParseFormatInvalid(t *testing.T) {
	_, err := template.ParseFormat("yaml")
	if err == nil {
		t.Fatal("expected error for unknown format, got nil")
	}
	if !strings.Contains(err.Error(), "yaml") {
		t.Errorf("error message should mention the bad format, got: %v", err)
	}
}

func TestRenderExport(t *testing.T) {
	b := seedBundle(t)
	var sb strings.Builder
	if err := template.Render(&sb, b, template.FormatExport); err != nil {
		t.Fatalf("Render: %v", err)
	}
	out := sb.String()
	for _, want := range []string{
		`export API_KEY="abc123"`,
		`export DB_HOST="localhost"`,
		`export DB_PASS="s3cr3t"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\ngot:\n%s", want, out)
		}
	}
}

func TestRenderDotenv(t *testing.T) {
	b := seedBundle(t)
	var sb strings.Builder
	if err := template.Render(&sb, b, template.FormatDotenv); err != nil {
		t.Fatalf("Render: %v", err)
	}
	out := sb.String()
	for _, want := range []string{
		`API_KEY="abc123"`,
		`DB_HOST="localhost"`,
		`DB_PASS="s3cr3t"`,
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\ngot:\n%s", want, out)
		}
		if strings.Contains(out, "export "+want) {
			t.Errorf("dotenv output must not contain 'export' prefix")
		}
	}
}

func TestRenderOutputIsSorted(t *testing.T) {
	b := seedBundle(t)
	var sb strings.Builder
	_ = template.Render(&sb, b, template.FormatDotenv)
	lines := strings.Split(strings.TrimSpace(sb.String()), "\n")
	for i := 1; i < len(lines); i++ {
		if lines[i] < lines[i-1] {
			t.Errorf("output not sorted: %q comes after %q", lines[i], lines[i-1])
		}
	}
}

func TestRenderUnknownFormat(t *testing.T) {
	b := seedBundle(t)
	var sb strings.Builder
	err := template.Render(&sb, b, template.Format("toml"))
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestRenderEmptyBundle(t *testing.T) {
	b := bundle.New("empty")
	var sb strings.Builder
	if err := template.Render(&sb, b, template.FormatExport); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sb.Len() != 0 {
		t.Errorf("expected empty output for empty bundle, got: %q", sb.String())
	}
}
