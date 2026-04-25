package cli

import (
	"fmt"

	"cargohold/internal/merge"
)

// runMerge implements the `cargohold merge <src-env> <dst-env>` sub-command.
// It reads both bundles from the vault (prompting for the passphrase once),
// merges src into dst, and persists the updated destination bundle.
func (r *Runner) runMerge(args []string, overwrite bool) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: cargohold merge <src-env> <dst-env>")
	}

	srcEnv := args[0]
	dstEnv := args[1]

	if srcEnv == dstEnv {
		return fmt.Errorf("merge: source and destination environments must differ")
	}

	passphrase, err := r.passphrase()
	if err != nil {
		return fmt.Errorf("merge: reading passphrase: %w", err)
	}

	srcBundle, err := r.vault.Open(srcEnv, passphrase)
	if err != nil {
		return fmt.Errorf("merge: opening source bundle %q: %w", srcEnv, err)
	}

	dstBundle, err := r.vault.Open(dstEnv, passphrase)
	if err != nil {
		return fmt.Errorf("merge: opening destination bundle %q: %w", dstEnv, err)
	}

	opts := merge.Options{Overwrite: overwrite}
	res, err := merge.Bundles(dstBundle, srcBundle, opts)
	if err != nil {
		return fmt.Errorf("merge: %w", err)
	}

	if err := r.vault.Save(dstEnv, dstBundle, passphrase); err != nil {
		return fmt.Errorf("merge: saving destination bundle: %w", err)
	}

	r.out.Success(fmt.Sprintf(
		"merged %q → %q: %d added, %d overwritten, %d skipped",
		srcEnv, dstEnv,
		len(res.Added), len(res.Overwritten), len(res.Skipped),
	))
	return nil
}
