// Package diff compares two [bundle.Bundle] instances and reports
// which keys were added, removed, or had their values changed.
//
// It is intentionally value-safe: changed values are detected but
// never included in the returned [Result], so secret material is
// not inadvertently exposed through diff output.
//
// Typical usage:
//
//	r := diff.Bundles(bundleA, bundleB)
//	if !r.IsEmpty() {
//		fmt.Println("added:", r.Added)
//		fmt.Println("removed:", r.Removed)
//		fmt.Println("changed:", r.Changed)
//	}
package diff
