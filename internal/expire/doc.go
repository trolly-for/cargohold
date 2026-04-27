// Package expire provides lightweight expiry enforcement for cargohold bundles.
//
// An Expirer persists a single RFC3339 timestamp alongside a bundle's data
// directory. Callers should invoke Check before any read or write operation;
// if the bundle has passed its expiry time, Check returns ErrExpired and the
// operation should be aborted.
//
// Expiry can be extended by calling Set with a new future time, or removed
// entirely with Clear.
package expire
