// Package rotate implements passphrase rotation for cargohold secret bundles.
//
// Rotation decrypts an existing bundle using the current passphrase and
// immediately re-encrypts it with a new passphrase. The underlying secret
// data is never written to disk in plaintext — the re-encryption happens
// entirely in memory before the new ciphertext is persisted.
//
// Basic usage:
//
//	r := rotate.New(store)
//	if err := r.Rotate("production", oldPass, newPass); err != nil {
//		log.Fatal(err)
//	}
package rotate
