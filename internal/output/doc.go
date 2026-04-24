// Package output provides lightweight formatting utilities for cargohold's
// CLI layer. It abstracts stdout/stderr writes behind a Formatter type,
// making it easy to inject writers in tests and keep presentation logic
// out of command handlers.
//
// Usage:
//
//	f := output.Default()        // writes to os.Stdout / os.Stderr
//	f.Success("bundle created")
//	f.KeyValue("API_KEY", "***")
package output
