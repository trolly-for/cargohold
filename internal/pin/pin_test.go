package pin_test

import (
	"os"
	"path/filepath"
	"testing"

	"cargohold/internal/pin"
)

func tempPinner(t *testing.T) *pin.Pinner {
	t.Helper()
	dir := t.TempDir()
	p, err := pin.New(filepath.Join(dir, "pins.json"))
	if err != nil {
		t.Fatalf("pin.New: %v", err)
	}
	return p
}

func TestPinAndCheck(t *testing.T) {
	p := tempPinner(t)
	if err := p.Pin("DB_URL", "postgres://localhost/dev"); err != nil {
		t.Fatalf("Pin: %v", err)
	}
	if err := p.Check("DB_URL", "postgres://localhost/dev"); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestCheckDrifted(t *testing.T) {
	p := tempPinner(t)
	_ = p.Pin("API_KEY", "secret-original")
	err := p.Check("API_KEY", "secret-changed")
	if err == nil {
		t.Fatal("expected ErrDrifted, got nil")
	}
	if !isError(err, pin.ErrDrifted) {
		t.Errorf("expected ErrDrifted, got %v", err)
	}
}

func TestCheckNotPinned(t *testing.T) {
	p := tempPinner(t)
	err := p.Check("MISSING", "value")
	if !isError(err, pin.ErrNotPinned) {
		t.Errorf("expected ErrNotPinned, got %v", err)
	}
}

func TestUnpin(t *testing.T) {
	p := tempPinner(t)
	_ = p.Pin("TOKEN", "abc123")
	if err := p.Unpin("TOKEN"); err != nil {
		t.Fatalf("Unpin: %v", err)
	}
	if err := p.Check("TOKEN", "abc123"); !isError(err, pin.ErrNotPinned) {
		t.Errorf("expected ErrNotPinned after unpin, got %v", err)
	}
}

func TestUnpinNotPinnedErrors(t *testing.T) {
	p := tempPinner(t)
	if err := p.Unpin("GHOST"); !isError(err, pin.ErrNotPinned) {
		t.Errorf("expected ErrNotPinned, got %v", err)
	}
}

func TestGetRecord(t *testing.T) {
	p := tempPinner(t)
	_ = p.Pin("HOST", "localhost")
	rec, err := p.Get("HOST")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if rec.Value != "localhost" {
		t.Errorf("got value %q, want %q", rec.Value, "localhost")
	}
	if rec.PinnedAt.IsZero() {
		t.Error("PinnedAt should not be zero")
	}
}

func TestPersistsAcrossInstances(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "pins.json")

	p1, _ := pin.New(path)
	_ = p1.Pin("PERSIST", "yes")

	p2, err := pin.New(path)
	if err != nil {
		t.Fatalf("second New: %v", err)
	}
	if err := p2.Check("PERSIST", "yes"); err != nil {
		t.Errorf("expected pinned value to persist, got %v", err)
	}
}

func TestNewMissingFileIsOK(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent.json")
	if _, err := pin.New(path); err != nil {
		t.Errorf("expected no error for missing file, got %v", err)
	}
}

func TestNewCorruptFileErrors(t *testing.T) {
	path := filepath.Join(t.TempDir(), "corrupt.json")
	_ = os.WriteFile(path, []byte("not json{"), 0o600)
	if _, err := pin.New(path); err == nil {
		t.Error("expected error for corrupt file, got nil")
	}
}

// isError reports whether err wraps target.
func isError(err, target error) bool {
	if err == nil {
		return false
	}
	return err == target || containsTarget(err.Error(), target.Error())
}

func containsTarget(msg, target string) bool {
	return len(msg) >= len(target) && (msg == target ||
		len(msg) > 0 && (msg[len(msg)-len(target):] == target ||
			func() bool {
				for i := 0; i+len(target) <= len(msg); i++ {
					if msg[i:i+len(target)] == target {
						return true
					}
				}
				return false
			}()))
}
