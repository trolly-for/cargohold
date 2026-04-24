// Package env provides helpers for working with named deployment environments
// within cargohold. It defines the set of recognized environment names
// (development, staging, production), normalisation utilities, and a
// convenience function for reading the active environment from the
// CARGOHOLD_ENV environment variable.
//
// Environment names are case-insensitive and whitespace-tolerant so that
// callers can accept user-supplied strings without additional preprocessing.
package env
