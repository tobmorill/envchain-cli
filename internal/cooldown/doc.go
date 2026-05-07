// Package cooldown provides per-project cooldown period tracking for
// envchain-cli. A cooldown prevents repeated invocations of sensitive
// operations (e.g. passphrase prompts, secret rotation) within a
// configurable window.
//
// Records are persisted to the shared key-value store and contain the
// start time and configured duration, allowing the caller to determine
// whether the cooldown is still active without relying on in-process
// state.
//
// Usage:
//
//	m := cooldown.New(s)
//	_ = m.Set("myproject", 10*time.Minute)
//	active, _ := m.IsActive("myproject")
package cooldown
