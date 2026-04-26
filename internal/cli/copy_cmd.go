package cli

import (
	"fmt"
	"strings"

	copy_ "cargohold/internal/copy"
)

// runCopy implements the `copy <src-env> <dst-env> [key...]` sub-command.
// It opens both bundles, copies the requested keys (or all keys when none
// are specified), then persists the destination bundle.
func (r *Runner) runCopy(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: copy <src-env> <dst-env> [key...]")
	}

	srcEnv := args[0]
	dstEnv := args[1]
	selectedKeys := args[2:]

	if srcEnv == dstEnv {
		return fmt.Errorf("copy: source and destination environments must differ")
	}

	passphrase, err := r.passphrase()
	if err != nil {
		return err
	}

	srcVault, err := r.openVault(srcEnv, passphrase)
	if err != nil {
		return fmt.Errorf("copy: open src %q: %w", srcEnv, err)
	}

	dstVault, err := r.openVault(dstEnv, passphrase)
	if err != nil {
		return fmt.Errorf("copy: open dst %q: %w", dstEnv, err)
	}

	copier, err := copy_.New(srcVault.Bundle())
	if err != nil {
		return err
	}

	opts := copy_.Options{
		Keys:      selectedKeys,
		Overwrite: false,
	}

	n, err := copier.Into(dstVault.Bundle(), opts)
	if err != nil {
		return err
	}

	if err := r.saveVault(dstEnv, dstVault, passphrase); err != nil {
		return fmt.Errorf("copy: save dst %q: %w", dstEnv, err)
	}

	keys := "all keys"
	if len(selectedKeys) > 0 {
		keys = strings.Join(selectedKeys, ", ")
	}
	r.out.Success(fmt.Sprintf("copied %d key(s) (%s) from %s → %s", n, keys, srcEnv, dstEnv))
	return nil
}
