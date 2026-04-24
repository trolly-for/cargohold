// Package output provides formatting helpers for CLI output.
package output

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// Formatter writes formatted output to a writer.
type Formatter struct {
	out io.Writer
	err io.Writer
}

// New returns a Formatter writing to the given writers.
func New(out, err io.Writer) *Formatter {
	return &Formatter{out: out, err: err}
}

// Default returns a Formatter writing to stdout and stderr.
func Default() *Formatter {
	return New(os.Stdout, os.Stderr)
}

// Success prints a success message prefixed with "✓".
func (f *Formatter) Success(msg string) {
	fmt.Fprintf(f.out, "✓ %s\n", msg)
}

// Error prints an error message prefixed with "✗" to stderr.
func (f *Formatter) Error(msg string) {
	fmt.Fprintf(f.err, "✗ %s\n", msg)
}

// Info prints an informational message to stdout.
func (f *Formatter) Info(msg string) {
	fmt.Fprintf(f.out, "  %s\n", msg)
}

// KeyValue prints a key=value pair.
func (f *Formatter) KeyValue(key, value string) {
	fmt.Fprintf(f.out, "%s=%s\n", key, value)
}

// KeyList prints a sorted list of keys, one per line.
func (f *Formatter) KeyList(keys []string) {
	sorted := make([]string, len(keys))
	copy(sorted, keys)
	sort.Strings(sorted)
	for _, k := range sorted {
		fmt.Fprintf(f.out, "  - %s\n", k)
	}
}

// BundleHeader prints a section header for a named bundle.
func (f *Formatter) BundleHeader(env string) {
	border := strings.Repeat("-", len(env)+10)
	fmt.Fprintf(f.out, "%s\n  bundle: %s\n%s\n", border, env, border)
}
