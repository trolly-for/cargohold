package access_test

import (
	"testing"

	"cargohold/internal/access"
)

func TestParseRoleValid(t *testing.T) {
	for _, r := range []string{"reader", "writer", "admin"} {
		got, err := access.ParseRole(r)
		if err != nil {
			t.Fatalf("ParseRole(%q): unexpected error: %v", r, err)
		}
		if string(got) != r {
			t.Fatalf("ParseRole(%q) = %q, want %q", r, got, r)
		}
	}
}

func TestParseRoleInvalid(t *testing.T) {
	_, err := access.ParseRole("superuser")
	if err == nil {
		t.Fatal("expected error for unknown role, got nil")
	}
}

func TestNewUnknownRoleErrors(t *testing.T) {
	_, err := access.New("ghost")
	if err == nil {
		t.Fatal("expected error for unknown role")
	}
}

func TestReaderPermissions(t *testing.T) {
	g, err := access.New(access.RoleReader)
	if err != nil {
		t.Fatal(err)
	}
	if err := g.Check(access.OpGet); err != nil {
		t.Errorf("reader should be allowed OpGet: %v", err)
	}
	if err := g.Check(access.OpList); err != nil {
		t.Errorf("reader should be allowed OpList: %v", err)
	}
	for _, op := range []access.Op{access.OpSet, access.OpDelete, access.OpRotate, access.OpExport} {
		if err := g.Check(op); err == nil {
			t.Errorf("reader should be denied %q", op)
		}
	}
}

func TestWriterPermissions(t *testing.T) {
	g, _ := access.New(access.RoleWriter)
	for _, op := range []access.Op{access.OpGet, access.OpSet, access.OpDelete, access.OpList} {
		if err := g.Check(op); err != nil {
			t.Errorf("writer should be allowed %q: %v", op, err)
		}
	}
	for _, op := range []access.Op{access.OpRotate, access.OpExport} {
		if err := g.Check(op); err == nil {
			t.Errorf("writer should be denied %q", op)
		}
	}
}

func TestAdminPermissions(t *testing.T) {
	g, _ := access.New(access.RoleAdmin)
	for _, op := range []access.Op{access.OpGet, access.OpSet, access.OpDelete, access.OpList, access.OpRotate, access.OpExport} {
		if err := g.Check(op); err != nil {
			t.Errorf("admin should be allowed %q: %v", op, err)
		}
	}
}

func TestGuardRole(t *testing.T) {
	g, _ := access.New(access.RoleWriter)
	if g.Role() != access.RoleWriter {
		t.Fatalf("Role() = %q, want %q", g.Role(), access.RoleWriter)
	}
}
