// Package snapshot provides point-in-time exports of secret bundles.
//
// A snapshot captures all key-value pairs from a bundle, encrypts them
// with a caller-supplied passphrase, and returns an opaque byte slice
// suitable for archiving or transmission.
//
// Usage:
//
//	s := snapshot.New("production")
//	blob, err := s.Export(b, passphrase)   // encrypt bundle → blob
//	snap, err := s.Import(blob, passphrase) // decrypt blob → Snapshot
//
// Blobs can be persisted via FileStore:
//
//	store, err := snapshot.NewFileStore("/var/cargohold/snapshots")
//	path, err  := store.Write("production", blob)
//	blob, err  = store.Read(path)
package snapshot
