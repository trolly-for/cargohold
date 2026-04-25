// Package template provides functionality for rendering secret bundles
// into environment-variable export scripts or dotenv-formatted files.
package template

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"cargohold/internal/bundle"
)

// Format represents the output format for rendering a bundle.
type Format string

const (
	// FormatExport renders keys as shell export statements.
	FormatExport Format = "export"
	// FormatDotenv renders keys in .env file format.
	FormatDotenv Format = "dotenv"
)

// ErrUnknownFormat is returned when an unrecognised format is requested.
type ErrUnknownFormat struct {
	Given string
}

func (e ErrUnknownFormat) Error() string {
	return fmt.Sprintf("unknown template format %q: must be \"export\" or \"dotenv\"", e.Given)
}

// Render writes the contents of b to w using the specified format.
// Keys are emitted in lexicographic order for deterministic output.
func Render(w io.Writer, b *bundle.Bundle, f Format) error {
	switch f {
	case FormatExport:
		return renderExport(w, b)
	case FormatDotenv:
		return renderDotenv(w, b)
	default:
		return ErrUnknownFormat{Given: string(f)}
	}
}

// ParseFormat converts a raw string into a Format, returning an error for
// unrecognised values.
func ParseFormat(s string) (Format, error) {
	f := Format(strings.ToLower(strings.TrimSpace(s)))
	switch f {
	case FormatExport, FormatDotenv:
		return f, nil
	default:
		return "", ErrUnknownFormat{Given: s}
	}
}

func sortedKeys(b *bundle.Bundle) []string {
	keys := b.Keys()
	sort.Strings(keys)
	return keys
}

func renderExport(w io.Writer, b *bundle.Bundle) error {
	for _, k := range sortedKeys(b) {
		v, _ := b.Get(k)
		if _, err := fmt.Fprintf(w, "export %s=%q\n", k, v); err != nil {
			return err
		}
	}
	return nil
}

func renderDotenv(w io.Writer, b *bundle.Bundle) error {
	for _, k := range sortedKeys(b) {
		v, _ := b.Get(k)
		if _, err := fmt.Fprintf(w, "%s=%q\n", k, v); err != nil {
			return err
		}
	}
	return nil
}
