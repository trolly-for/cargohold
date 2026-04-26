// Package importer provides functionality for importing secrets from
// external formats such as .env files and JSON key-value maps into
// a cargohold bundle.
package importer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"cargohold/internal/bundle"
)

// Format represents a supported import format.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatJSON   Format = "json"
)

// ParseFormat returns a Format from a string, or an error if unrecognised.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(s) {
	case "dotenv", ".env":
		return FormatDotenv, nil
	case "json":
		return FormatJSON, nil
	}
	return "", fmt.Errorf("unsupported import format %q: want dotenv or json", s)
}

// Import reads secrets from r in the given format and sets them on dst.
// Existing keys are overwritten. Returns the number of keys imported.
func Import(dst *bundle.Bundle, r io.Reader, format Format) (int, error) {
	switch format {
	case FormatDotenv:
		return importDotenv(dst, r)
	case FormatJSON:
		return importJSON(dst, r)
	}
	return 0, fmt.Errorf("unknown format: %s", format)
}

func importDotenv(dst *bundle.Bundle, r io.Reader) (int, error) {
	count := 0
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.IndexByte(line, '=')
		if idx < 1 {
			return count, fmt.Errorf("invalid dotenv line: %q", line)
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		val = strings.Trim(val, `"'`)
		if err := dst.Set(key, val); err != nil {
			return count, fmt.Errorf("set %q: %w", key, err)
		}
		count++
	}
	return count, scanner.Err()
}

func importJSON(dst *bundle.Bundle, r io.Reader) (int, error) {
	var m map[string]string
	if err := json.NewDecoder(r).Decode(&m); err != nil {
		return 0, fmt.Errorf("decode json: %w", err)
	}
	count := 0
	for k, v := range m {
		if err := dst.Set(k, v); err != nil {
			return count, fmt.Errorf("set %q: %w", k, err)
		}
		count++
	}
	return count, nil
}
