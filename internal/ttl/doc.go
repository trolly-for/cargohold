// Package ttl provides lightweight time-to-live (TTL) expiry tracking for
// cargohold secret bundles.
//
// A Tracker persists a JSON record alongside each bundle that records when
// the bundle's secrets should be considered stale. Callers can use Check to
// gate access and Set / Remove to manage the lifecycle of expiry records.
//
// Typical usage:
//
//	tr, err := ttl.New(filepath.Join(storeDir, ".ttl"))
//	if err != nil { ... }
//
//	// Assign a 24-hour TTL when a bundle is written.
//	tr.Set("production", 24*time.Hour)
//
//	// Verify the bundle has not expired before reading secrets.
//	if err := tr.Check("production"); errors.Is(err, ttl.ErrExpired) {
//	    // warn or refuse access
//	}
package ttl
