package cli

import (
	"errors"
	"fmt"
	"time"

	"cargohold/internal/expire"
)

// runExpire dispatches the expire sub-commands:
//
//	expire set   <env> <duration>   – set expiry relative to now
//	expire get   <env>              – print current expiry
//	expire clear <env>              – remove expiry
func (r *Runner) runExpire(args []string) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: expire <set|get|clear> <env> [duration]")
	}
	subcmd, env := args[0], args[1]

	bundleDir, err := r.store.BundlePath(env)
	if err != nil {
		return fmt.Errorf("expire: %w", err)
	}

	ex, err := expire.New(bundleDir)
	if err != nil {
		return fmt.Errorf("expire: %w", err)
	}

	switch subcmd {
	case "set":
		if len(args) < 3 {
			return fmt.Errorf("expire set requires a duration argument (e.g. 24h)")
		}
		d, err := time.ParseDuration(args[2])
		if err != nil {
			return fmt.Errorf("expire set: invalid duration %q: %w", args[2], err)
		}
		if d <= 0 {
			return fmt.Errorf("expire set: duration must be positive")
		}
		expiry := time.Now().UTC().Add(d)
		if err := ex.Set(expiry); err != nil {
			return err
		}
		r.out.Success(fmt.Sprintf("expiry set to %s", expiry.Format(time.RFC3339)))

	case "get":
		t, err := ex.Get()
		if errors.Is(err, expire.ErrNotSet) {
			r.out.KeyValue("expiry", "none")
			return nil
		}
		if err != nil {
			return err
		}
		checkErr := ex.Check()
		status := "valid"
		if errors.Is(checkErr, expire.ErrExpired) {
			status = "EXPIRED"
		}
		r.out.KeyValue("expiry", fmt.Sprintf("%s (%s)", t.Format(time.RFC3339), status))

	case "clear":
		if err := ex.Clear(); err != nil {
			return err
		}
		r.out.Success("expiry cleared")

	default:
		return fmt.Errorf("expire: unknown subcommand %q", subcmd)
	}
	return nil
}
