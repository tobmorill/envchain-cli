// Package lock manages per-chain unlock sessions for envchain-cli.
//
// A session is a small JSON file written to a runtime directory that records
// the chain name and an expiry timestamp.  When the session file is absent or
// its timestamp has passed, the chain is considered locked and the caller must
// supply a passphrase again.
//
// Typical TTL values range from a few minutes (tight security) to several
// hours (convenience on a personal workstation).  The TTL can be configured
// via the envchain config command.
//
// Usage:
//
//	m := lock.NewManager(runtimeDir)
//	if err := m.Unlock("myproject", 30*time.Minute); err != nil { ... }
//	ok, err := m.IsUnlocked("myproject")
package lock
