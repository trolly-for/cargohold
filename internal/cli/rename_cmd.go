package cli

import (
	"errors"
	"fmt"

	"cargohold/internal/rename"
)

// runRename handles the "rename" subcommand.
// Usage: cargohold rename <env> <old-key> <new-key> [--overwrite]
func (r *Runner) runRename(args []string) error {
	if len(args) < 3 {
		return errors.New("usage: cargohold rename <env> <old-key> <new-key> [--overwrite]")
	}

	env := args[0]
	oldKey := args[1]
	newKey := args[2]

	overwrite := false
	for _, flag := range args[3:] {
		if flag == "--overwrite" {
			overwrite = true
		}
	}

	passphrase, err := r.passphrase(env)
	if err != nil {
		return fmt.Errorf("rename: read passphrase: %w", err)
	}

	v, err := r.vault.Open(env, passphrase)
	if err != nil {
		return fmt.Errorf("rename: open vault: %w", err)
	}

	if err := rename.Key(v, oldKey, newKey, overwrite); err != nil {
		return fmt.Errorf("rename: %w", err)
	}

	if err := r.vault.Save(env, passphrase, v); err != nil {
		return fmt.Errorf("rename: save vault: %w", err)
	}

	r.out.Success(fmt.Sprintf("renamed %q to %q in %s", oldKey, newKey, env))
	return nil
}
