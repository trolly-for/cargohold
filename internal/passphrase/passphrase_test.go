package passphrase_test

import (
	"os"
	"testing"

	"cargohold/internal/passphrase"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"valid passphrase", "supersecret", nil},
		{"exactly min length", "12345678", nil},
		{"empty string", "", passphrase.ErrEmpty},
		{"only spaces", "       ", passphrase.ErrEmpty},
		{"too short", "abc", passphrase.ErrTooShort},
		{"seven chars", "1234567", passphrase.ErrTooShort},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := passphrase.Validate(tc.input)
			if err != tc.wantErr {
				t.Errorf("Validate(%q) = %v, want %v", tc.input, err, tc.wantErr)
			}
		})
	}
}

func TestReadFromEnvVar(t *testing.T) {
	const testPass = "env-supplied-pass"
	t.Setenv(passphrase.EnvVar, testPass)

	got, err := passphrase.Read("Enter passphrase: ")
	if err != nil {
		t.Fatalf("Read() unexpected error: %v", err)
	}
	if got != testPass {
		t.Errorf("Read() = %q, want %q", got, testPass)
	}
}

func TestReadWithConfirmFromEnvVar(t *testing.T) {
	const testPass = "env-supplied-pass"
	t.Setenv(passphrase.EnvVar, testPass)

	got, err := passphrase.ReadWithConfirm("Enter passphrase: ", "Confirm: ")
	if err != nil {
		t.Fatalf("ReadWithConfirm() unexpected error: %v", err)
	}
	if got != testPass {
		t.Errorf("ReadWithConfirm() = %q, want %q", got, testPass)
	}
}

func TestReadWithConfirmEnvVarTooShort(t *testing.T) {
	t.Setenv(passphrase.EnvVar, "short")

	_, err := passphrase.ReadWithConfirm("Enter passphrase: ", "Confirm: ")
	if err != passphrase.ErrTooShort {
		t.Errorf("ReadWithConfirm() error = %v, want ErrTooShort", err)
	}
}

func TestEnvVarName(t *testing.T) {
	// Ensure the constant is stable so dependent tooling (docs, scripts) can rely on it.
	if passphrase.EnvVar != "CARGOHOLD_PASSPHRASE" {
		t.Errorf("EnvVar = %q, want CARGOHOLD_PASSPHRASE", passphrase.EnvVar)
	}
}

func TestReadEnvVarUnset(t *testing.T) {
	// Make sure env var is cleared; if terminal is not available the call
	// will fail — we just verify it doesn't silently return the env value.
	os.Unsetenv(passphrase.EnvVar)

	// We cannot easily test interactive terminal reads in unit tests,
	// so we only assert the env-var fast-path is skipped (no panic/nil error
	// from env branch). The error from term.ReadPassword on a non-tty is
	// acceptable here.
	_, err := passphrase.Read("Enter passphrase: ")
	// err may be non-nil (no tty in CI), but should not be ErrEmpty or ErrTooShort.
	if err == passphrase.ErrEmpty || err == passphrase.ErrTooShort {
		t.Errorf("Read() returned validation error without env var: %v", err)
	}
}
