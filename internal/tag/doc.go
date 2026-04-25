// Package tag provides lightweight key tagging for cargohold secret bundles.
//
// Tags are short lowercase labels (e.g. "db", "api", "infra") that can be
// attached to individual secret keys. They enable selective operations such
// as exporting only the keys relevant to a particular service tier.
//
// A tag.Map is a plain in-memory structure and is intended to be persisted
// alongside the encrypted bundle by the vault layer.
package tag
