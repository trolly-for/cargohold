// Package passphrase handles reading and validating passphrases
// from the terminal or environment variables for use with vault encryption.
package passphrase

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

const (
	// EnvVar is the environment variable used to supply a passphrase non-interactively.
	EnvVar = "CARGOHOLD_PASSPHRASE"

	minLength = 8
)

// ErrTooShort is returned when a passphrase does not meet the minimum length.
var ErrTooShort = fmt.Errorf("passphrase must be at least %d characters", minLength)

// ErrEmpty is returned when an empty passphrase is provided.
var ErrEmpty = errors.New("passphrase must not be empty")

// ErrMismatch is returned when confirmation passphrase does not match.
var ErrMismatch = errors.New("passphrases do not match")

// Read reads a passphrase from the environment variable if set,
// otherwise prompts the user interactively via the terminal.
func Read(prompt string) (string, error) {
	if val := os.Getenv(EnvVar); val != "" {
		return val, nil
	}
	return readFromTerminal(prompt)
}

// ReadWithConfirm reads a passphrase and asks the user to confirm it.
// Intended for use during vault initialisation.
func ReadWithConfirm(prompt, confirmPrompt string) (string, error) {
	if val := os.Getenv(EnvVar); val != "" {
		if err := Validate(val); err != nil {
			return "", err
		}
		return val, nil
	}

	first, err := readFromTerminal(prompt)
	if err != nil {
		return "", err
	}
	if err := Validate(first); err != nil {
		return "", err
	}

	second, err := readFromTerminal(confirmPrompt)
	if err != nil {
		return "", err
	}

	if first != second {
		return "", ErrMismatch
	}
	return first, nil
}

// Validate checks that a passphrase meets minimum requirements.
func Validate(p string) error {
	if strings.TrimSpace(p) == "" {
		return ErrEmpty
	}
	if len(p) < minLength {
		return ErrTooShort
	}
	return nil
}

func readFromTerminal(prompt string) (string, error) {
	fmt.Fprint(os.Stderr, prompt)
	bytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr) // newline after hidden input
	if err != nil {
		return "", fmt.Errorf("reading passphrase: %w", err)
	}
	return string(bytes), nil
}
