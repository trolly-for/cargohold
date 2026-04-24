package vault

// StorePath exposes the underlying store's root path for use in tests that
// need to construct a second Vault pointing at the same directory.
func (v *Vault) StorePath() string {
	return v.store.RootDir()
}
