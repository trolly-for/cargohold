// Package access implements role-based access control for cargohold operations.
//
// Three built-in roles are provided:
//
//   - reader  — may get and list secrets
//   - writer  — may get, set, delete, and list secrets
//   - admin   — full access including rotate and export
//
// Use New to obtain a Guard for a role, then call Check before performing
// any sensitive operation:
//
//	guard, err := access.New(access.RoleWriter)
//	if err := guard.Check(access.OpSet); err != nil {
//	    return err
//	}
package access
