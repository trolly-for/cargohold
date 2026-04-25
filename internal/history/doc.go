// Package history provides a lightweight append-only mutation log for
// cargohold bundles.
//
// Each environment bundle maintains its own JSONL log file recording every
// set and delete operation along with the affected key and a UTC timestamp.
// Entries are never removed; the log grows monotonically and can be read
// back in full for auditing or debugging purposes.
//
// Usage:
//
//	tr, err := history.New("/var/cargohold/history", "production")
//	if err != nil { ... }
//
//	// record a mutation
//	_ = tr.Record("production", "DB_PASSWORD", history.OpSet)
//
//	// read the full log
//	entries, _ := tr.ReadAll()
package history
