// export_test.go exposes internal fields for white-box testing.
package cli

import (
	"os"

	"cargohold/internal/store"
)

// Allow tests to inject a pre-built Store and redirect output.
func (r *Runner) SetOut(f *os.File) {
	r.out = f
}

// NewWithStore creates a Runner with a caller-supplied store, used in tests.
func NewWithStore(s *store.Store) *Runner {
	return &Runner{Store: s, out: os.Stdout}
}
