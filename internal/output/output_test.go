package output_test

import (
	"bytes"
	"strings"
	"testing"

	"cargohold/internal/output"
)

func newTestFormatter() (*output.Formatter, *bytes.Buffer, *bytes.Buffer) {
	out := &bytes.Buffer{}
	err := &bytes.Buffer{}
	return output.New(out, err), out, err
}

func TestSuccess(t *testing.T) {
	f, out, _ := newTestFormatter()
	f.Success("bundle created")
	if !strings.Contains(out.String(), "✓") {
		t.Errorf("expected success prefix, got: %q", out.String())
	}
	if !strings.Contains(out.String(), "bundle created") {
		t.Errorf("expected message in output, got: %q", out.String())
	}
}

func TestError(t *testing.T) {
	f, _, errOut := newTestFormatter()
	f.Error("something went wrong")
	if !strings.Contains(errOut.String(), "✗") {
		t.Errorf("expected error prefix, got: %q", errOut.String())
	}
	if !strings.Contains(errOut.String(), "something went wrong") {
		t.Errorf("expected message in error output, got: %q", errOut.String())
	}
}

func TestKeyValue(t *testing.T) {
	f, out, _ := newTestFormatter()
	f.KeyValue("DB_HOST", "localhost")
	expected := "DB_HOST=localhost\n"
	if out.String() != expected {
		t.Errorf("expected %q, got %q", expected, out.String())
	}
}

func TestKeyList(t *testing.T) {
	f, out, _ := newTestFormatter()
	f.KeyList([]string{"ZEBRA", "ALPHA", "MIDDLE"})
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if !strings.Contains(lines[0], "ALPHA") {
		t.Errorf("expected sorted output, first line: %q", lines[0])
	}
	if !strings.Contains(lines[2], "ZEBRA") {
		t.Errorf("expected sorted output, last line: %q", lines[2])
	}
}

func TestBundleHeader(t *testing.T) {
	f, out, _ := newTestFormatter()
	f.BundleHeader("production")
	if !strings.Contains(out.String(), "production") {
		t.Errorf("expected env name in header, got: %q", out.String())
	}
	if strings.Count(out.String(), "\n") < 3 {
		t.Errorf("expected multi-line header, got: %q", out.String())
	}
}

func TestInfo(t *testing.T) {
	f, out, _ := newTestFormatter()
	f.Info("loading secrets")
	if !strings.Contains(out.String(), "loading secrets") {
		t.Errorf("expected info message, got: %q", out.String())
	}
}
