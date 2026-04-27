// Package redact provides utilities for masking secret values in output,
// preventing accidental exposure of plaintext secrets in logs or terminal output.
package redact

import "strings"

const defaultMask = "********"

// Redactor masks secret values before they are displayed.
type Redactor struct {
	mask string
}

// New returns a Redactor using the default mask string.
func New() *Redactor {
	return &Redactor{mask: defaultMask}
}

// NewWithMask returns a Redactor that replaces secret values with the given mask.
func NewWithMask(mask string) *Redactor {
	if mask == "" {
		mask = defaultMask
	}
	return &Redactor{mask: mask}
}

// Value returns the mask string regardless of the input value.
func (r *Redactor) Value(_ string) string {
	return r.mask
}

// Partial reveals the first n and last n characters of value, masking the rest.
// If the value is too short to reveal both ends, the full mask is returned.
func (r *Redactor) Partial(value string, n int) string {
	if n <= 0 || len(value) <= n*2 {
		return r.mask
	}
	return value[:n] + strings.Repeat("*", len(value)-n*2) + value[len(value)-n:]
}

// Map returns a copy of the input map with all values replaced by the mask.
func (r *Redactor) Map(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k := range secrets {
		out[k] = r.mask
	}
	return out
}

// Mask returns the mask string used by this Redactor.
func (r *Redactor) Mask() string {
	return r.mask
}
