// Package ttl provides a lightweight in-memory key/value cache with
// per-entry time-to-live support.
//
// It is used by envchain to hold decrypted passphrases in memory for a
// configurable duration, reducing the number of times the user must
// re-enter their passphrase within a single shell session.
//
// Entries are evicted lazily on Get or explicitly via Purge. The cache
// is safe for concurrent use.
package ttl
