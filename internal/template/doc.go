// Package template renders a secret bundle into a plain-text format suitable
// for sourcing in a shell session or loading via a dotenv library.
//
// Two formats are supported:
//
//   - "export"  — emits one `export KEY="value"` line per secret, ready to be
//     evaluated by a POSIX shell (e.g. `eval $(cargohold render --format export)`).
//
//   - "dotenv"  — emits one `KEY="value"` line per secret, compatible with
//     .env file conventions used by tools such as docker-compose and godotenv.
//
// Keys are always emitted in lexicographic order so that the output is
// deterministic and diff-friendly.
package template
