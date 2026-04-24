// Package passphrase provides utilities for reading and validating
// passphrases used to encrypt and decrypt secret bundles.
//
// Passphrases can be supplied interactively via terminal prompt or
// non-interactively via the CARGOHOLD_PASSPHRASE environment variable,
// which is useful in CI/CD pipelines and automated workflows.
//
// A valid passphrase must be at least 8 characters long.
package passphrase
