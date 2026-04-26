// Package copy provides a Copier that duplicates secret keys from one
// bundle into another. It supports full or selective key copying and
// can optionally overwrite existing values in the destination bundle.
//
// Typical usage:
//
//	c, err := copy.New(srcBundle)
//	if err != nil { ... }
//	n, err := c.Into(dstBundle, copy.Options{Overwrite: false})
//
// When Keys is left empty all keys from the source are considered.
// Keys that already exist in the destination are silently skipped
// unless Overwrite is set to true.
package copy
