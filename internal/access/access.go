// Package access provides role-based access control for secret bundles,
// restricting which operations are permitted based on a named role.
package access

import "fmt"

// Role represents a named access level.
type Role string

const (
	RoleReader Role = "reader"
	RoleWriter Role = "writer"
	RoleAdmin  Role = "admin"
)

// Op represents an operation that may be gated by a role.
type Op string

const (
	OpGet    Op = "get"
	OpSet    Op = "set"
	OpDelete Op = "delete"
	OpList   Op = "list"
	OpRotate Op = "rotate"
	OpExport Op = "export"
)

// allowed maps each role to the set of permitted operations.
var allowed = map[Role]map[Op]bool{
	RoleReader: {OpGet: true, OpList: true},
	RoleWriter: {OpGet: true, OpSet: true, OpDelete: true, OpList: true},
	RoleAdmin:  {OpGet: true, OpSet: true, OpDelete: true, OpList: true, OpRotate: true, OpExport: true},
}

// Guard enforces access control for a given role.
type Guard struct {
	role Role
}

// New returns a Guard for the given role. Returns an error if the role is unknown.
func New(role Role) (*Guard, error) {
	if _, ok := allowed[role]; !ok {
		return nil, fmt.Errorf("access: unknown role %q", role)
	}
	return &Guard{role: role}, nil
}

// Check returns nil if the role is permitted to perform op, or an error otherwise.
func (g *Guard) Check(op Op) error {
	if allowed[g.role][op] {
		return nil
	}
	return fmt.Errorf("access: role %q is not permitted to perform %q", g.role, op)
}

// Role returns the guard's current role.
func (g *Guard) Role() Role {
	return g.role
}

// ParseRole converts a string to a Role, returning an error for unknown values.
func ParseRole(s string) (Role, error) {
	r := Role(s)
	if _, ok := allowed[r]; !ok {
		return "", fmt.Errorf("access: unknown role %q", s)
	}
	return r, nil
}
