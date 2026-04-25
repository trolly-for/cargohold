// Package policy enforces access and mutation rules on secret bundles
// based on the target environment.
package policy

import (
	"errors"
	"fmt"

	"cargohold/internal/env"
)

// ErrDenied is returned when a policy check blocks an operation.
var ErrDenied = errors.New("policy: operation denied")

// Rule describes a single policy constraint.
type Rule struct {
	// AllowWrite permits set/delete operations on the environment.
	AllowWrite bool
	// RequireConfirmation requires an explicit confirmation flag for destructive ops.
	RequireConfirmation bool
}

// Policy holds per-environment rules.
type Policy struct {
	rules map[string]Rule
}

// New returns a Policy populated with sensible defaults.
// Production environments are write-protected and require confirmation.
func New() *Policy {
	return &Policy{
		rules: map[string]Rule{
			"production": {AllowWrite: false, RequireConfirmation: true},
			"prod":       {AllowWrite: false, RequireConfirmation: true},
			"staging":    {AllowWrite: true, RequireConfirmation: true},
			"development": {AllowWrite: true, RequireConfirmation: false},
			"dev":        {AllowWrite: true, RequireConfirmation: false},
		},
	}
}

// Set overrides the rule for a given environment name.
func (p *Policy) Set(environment string, rule Rule) {
	normalized := env.Normalize(environment)
	p.rules[normalized] = rule
}

// CheckWrite returns an error if write operations are not permitted for
// the given environment. Pass confirmed=true to satisfy confirmation-only guards.
func (p *Policy) CheckWrite(environment string, confirmed bool) error {
	normalized := env.Normalize(environment)
	rule, ok := p.rules[normalized]
	if !ok {
		// Unknown environments are allowed by default.
		return nil
	}
	if !rule.AllowWrite {
		return fmt.Errorf("%w: writes are disabled for environment %q", ErrDenied, normalized)
	}
	if rule.RequireConfirmation && !confirmed {
		return fmt.Errorf("%w: environment %q requires explicit confirmation (--confirm)", ErrDenied, normalized)
	}
	return nil
}

// IsWriteAllowed is a convenience helper that returns false if CheckWrite
// would return an error.
func (p *Policy) IsWriteAllowed(environment string, confirmed bool) bool {
	return p.CheckWrite(environment, confirmed) == nil
}
