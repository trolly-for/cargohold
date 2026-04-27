// Package redact provides value-masking utilities for the cargohold CLI.
//
// Use Redactor to prevent plaintext secrets from appearing in terminal output,
// audit logs, or error messages. The default mask is "********", but a custom
// mask can be supplied via NewWithMask.
//
// Partial masking is supported for UX scenarios where a hint of the value
// (e.g. first/last two characters) is acceptable.
package redact
