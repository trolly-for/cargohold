package cli

import (
	"fmt"

	"cargohold/internal/redact"
)

// runRedact handles the `cargohold redact <env> <key> [--partial]` command.
// It prints the masked value for a key, optionally showing a partial hint.
func (r *Runner) runRedact(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: cargohold redact <env> <key> [--partial]")
	}

	env := args[0]
	key := args[1]
	partial := len(args) >= 3 && args[2] == "--partial"

	passphrase, err := r.passphrase(env)
	if err != nil {
		return err
	}

	v, err := r.openVault(env, passphrase)
	if err != nil {
		return err
	}

	value, ok := v.Get(key)
	if !ok {
		return fmt.Errorf("key %q not found in environment %q", key, env)
	}

	red := redact.New()
	var display string
	if partial {
		display = red.Partial(value, 2)
	} else {
		display = red.Value(value)
	}

	r.out.KeyValue(key, display)
	return nil
}
