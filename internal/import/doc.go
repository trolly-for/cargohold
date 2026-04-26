// Package importer provides utilities for bulk-importing secrets into a
// cargohold bundle from common plaintext formats.
//
// Supported formats:
//
//   - dotenv  — KEY=VALUE lines, with optional quoted values and # comments
//   - json    — flat JSON object mapping string keys to string values
//
// Imported keys overwrite any existing values in the destination bundle.
// The caller is responsible for saving the bundle after a successful import.
package importer
