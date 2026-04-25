// Package lock implements lightweight file-based locking for cargohold
// environment bundles.
//
// A lock is represented as a small file on disk inside a dedicated lock
// directory. Acquiring a lock for an environment that is already locked
// returns ErrAlreadyLocked, allowing callers (such as the CLI write
// commands) to abort destructive operations before they begin.
//
// Typical usage:
//
//	locker := lock.New(filepath.Join(storeDir, ".locks"))
//
//	if locker.IsLocked(env) {
//		return fmt.Errorf("environment %q is locked — release it first", env)
//	}
package lock
