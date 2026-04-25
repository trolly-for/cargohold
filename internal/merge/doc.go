// Package merge implements bundle-to-bundle merging for cargohold.
//
// It provides a single entry point, Bundles, which copies keys from a source
// bundle into a destination bundle.  The caller controls whether conflicting
// keys should be overwritten or silently skipped via Options.
//
// A Result is returned describing every key that was added, skipped, or
// overwritten so that callers (e.g. the CLI) can display a human-readable
// summary of what changed.
package merge
