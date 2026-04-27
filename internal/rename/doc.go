// Package rename provides functionality for renaming keys within a secret bundle.
//
// It supports renaming a key to a new name within the same bundle, with optional
// overwrite behaviour when the destination key already exists.
//
// Example usage:
//
//	err := rename.Key(bundle, "OLD_KEY", "NEW_KEY", false)
//	if err != nil {
//		log.Fatal(err)
//	}
package rename
