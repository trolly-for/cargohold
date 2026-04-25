package tag_test

import (
	"testing"

	"cargohold/internal/tag"
)

func TestValidate(t *testing.T) {
	valid := []string{"db", "api-key", "cache_v2", "x"}
	for _, v := range valid {
		if err := tag.Validate(v); err != nil {
			t.Errorf("expected %q to be valid, got: %v", v, err)
		}
	}

	invalid := []string{"", "DB", "1start", "has space", "toolongtoolongtoolongtoolongtoolong"}
	for _, v := range invalid {
		if err := tag.Validate(v); err == nil {
			t.Errorf("expected %q to be invalid", v)
		}
	}
}

func TestAddAndTags(t *testing.T) {
	m := tag.New()

	if err := m.Add("DB_HOST", "db"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := m.Add("DB_HOST", "infra"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tags := m.Tags("DB_HOST")
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
	if tags[0] != "db" || tags[1] != "infra" {
		t.Errorf("unexpected tags: %v", tags)
	}
}

func TestAddDuplicateIsNoop(t *testing.T) {
	m := tag.New()
	m.Add("KEY", "api")
	m.Add("KEY", "api")

	if len(m.Tags("KEY")) != 1 {
		t.Errorf("expected 1 tag after duplicate add, got %d", len(m.Tags("KEY")))
	}
}

func TestAddInvalidTagErrors(t *testing.T) {
	m := tag.New()
	if err := m.Add("KEY", "Bad Tag"); err == nil {
		t.Error("expected error for invalid tag")
	}
}

func TestRemove(t *testing.T) {
	m := tag.New()
	m.Add("KEY", "db")
	m.Add("KEY", "api")
	m.Remove("KEY", "db")

	tags := m.Tags("KEY")
	if len(tags) != 1 || tags[0] != "api" {
		t.Errorf("expected [api], got %v", tags)
	}
}

func TestRemoveLastTagDeletesEntry(t *testing.T) {
	m := tag.New()
	m.Add("KEY", "db")
	m.Remove("KEY", "db")

	if m.Tags("KEY") != nil {
		t.Error("expected nil tags after removing last tag")
	}
}

func TestKeysWithTag(t *testing.T) {
	m := tag.New()
	m.Add("DB_HOST", "db")
	m.Add("DB_PASS", "db")
	m.Add("API_KEY", "api")

	keys := m.KeysWithTag("db")
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(keys))
	}
	if keys[0] != "DB_HOST" || keys[1] != "DB_PASS" {
		t.Errorf("unexpected keys: %v", keys)
	}

	if len(m.KeysWithTag("missing")) != 0 {
		t.Error("expected empty slice for unknown tag")
	}
}
