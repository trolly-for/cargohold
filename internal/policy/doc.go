// Package policy provides environment-aware write guards for cargohold.
//
// A Policy holds a set of Rules keyed by environment name. Each Rule
// controls whether write operations (set, delete, merge) are permitted
// and whether the caller must supply an explicit confirmation flag.
//
// Default rules:
//
//	"production" / "prod"  — writes disabled entirely
//	"staging"              — writes allowed with --confirm
//	"development" / "dev"  — writes allowed without confirmation
//
// Unknown environments are permitted without restriction so that
// custom environment names work out of the box.
package policy
