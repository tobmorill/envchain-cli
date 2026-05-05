// Package notes provides per-project free-text annotation storage for
// envchain-cli.
//
// Notes are arbitrary text blobs attached to a named project. They are
// encrypted at rest using the same passphrase as the associated chain,
// ensuring that annotations (which may contain sensitive context such as
// rotation schedules or credential owners) are never stored in plaintext.
//
// Basic usage:
//
//	m := notes.New(s)           // s is a *store.Store
//	err := m.Set("myproject", "rotate quarterly", passphrase)
//	n, err := m.Get("myproject", passphrase)
//	fmt.Println(n.Body)
package notes
