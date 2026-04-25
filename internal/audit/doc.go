// Package audit provides a lightweight, append-only audit log for
// cargohold operations.
//
// Each operation (init, get, set, delete, rotate) is recorded as a
// newline-delimited JSON entry containing a UTC timestamp, the
// operation name, the target environment, an optional key name, and
// an optional free-form note.
//
// Entries are written to ~/.cargohold/audit.log by default, with
// file permissions restricted to the owning user (0600).
package audit
