package access_test

import (
	"testing"

	"cargohold/internal/access"
)

// TestParseRoleAndCheckRoundtrip verifies that parsing a role string and
// constructing a Guard from it produces consistent permission behaviour.
func TestParseRoleAndCheckRoundtrip(t *testing.T) {
	cases := []struct {
		role    string
		op      access.Op
		allowed bool
	}{
		{"reader", access.OpGet, true},
		{"reader", access.OpSet, false},
		{"writer", access.OpSet, true},
		{"writer", access.OpExport, false},
		{"admin", access.OpExport, true},
		{"admin", access.OpRotate, true},
	}

	for _, tc := range cases {
		t.Run(tc.role+"_"+string(tc.op), func(t *testing.T) {
			role, err := access.ParseRole(tc.role)
			if err != nil {
				t.Fatalf("ParseRole: %v", err)
			}
			g, err := access.New(role)
			if err != nil {
				t.Fatalf("New: %v", err)
			}
			err = g.Check(tc.op)
			if tc.allowed && err != nil {
				t.Errorf("expected allowed, got error: %v", err)
			}
			if !tc.allowed && err == nil {
				t.Errorf("expected denied, got nil error")
			}
		})
	}
}
