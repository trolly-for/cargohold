package history_test

import (
	"os"
	"testing"
	"time"

	"cargohold/internal/history"
)

func tempTracker(t *testing.T, env string) *history.Tracker {
	t.Helper()
	dir := t.TempDir()
	tr, err := history.New(dir, env)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return tr
}

func TestReadAllEmptyLog(t *testing.T) {
	tr := tempTracker(t, "dev")
	entries, err := tr.ReadAll()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(entries))
	}
}

func TestRecordAndReadAll(t *testing.T) {
	tr := tempTracker(t, "staging")

	if err := tr.Record("staging", "DB_URL", history.OpSet); err != nil {
		t.Fatalf("Record set: %v", err)
	}
	if err := tr.Record("staging", "OLD_KEY", history.OpDelete); err != nil {
		t.Fatalf("Record delete: %v", err)
	}

	entries, err := tr.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Key != "DB_URL" || entries[0].Op != history.OpSet {
		t.Errorf("entry[0] mismatch: %+v", entries[0])
	}
	if entries[1].Key != "OLD_KEY" || entries[1].Op != history.OpDelete {
		t.Errorf("entry[1] mismatch: %+v", entries[1])
	}
}

func TestEntryTimestampIsUTC(t *testing.T) {
	tr := tempTracker(t, "prod")
	before := time.Now().UTC().Add(-time.Second)

	if err := tr.Record("prod", "SECRET", history.OpSet); err != nil {
		t.Fatalf("Record: %v", err)
	}

	after := time.Now().UTC().Add(time.Second)
	entries, _ := tr.ReadAll()
	ts := entries[0].Timestamp

	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v out of expected range [%v, %v]", ts, before, after)
	}
	if ts.Location() != time.UTC {
		t.Errorf("expected UTC, got %v", ts.Location())
	}
}

func TestRecordAppendsAcrossInstances(t *testing.T) {
	dir := t.TempDir()

	tr1, _ := history.New(dir, "dev")
	_ = tr1.Record("dev", "FIRST", history.OpSet)

	tr2, _ := history.New(dir, "dev")
	_ = tr2.Record("dev", "SECOND", history.OpSet)

	tr3, _ := history.New(dir, "dev")
	entries, err := tr3.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries across instances, got %d", len(entries))
	}
}

func TestNewCreatesDirectory(t *testing.T) {
	base := t.TempDir()
	subdir := base + "/nested/history"

	_, err := history.New(subdir, "dev")
	if err != nil {
		t.Fatalf("New with nested dir: %v", err)
	}
	if _, err := os.Stat(subdir); os.IsNotExist(err) {
		t.Error("expected nested directory to be created")
	}
}
