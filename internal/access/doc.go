// Package access provides per-project access tracking for envchain-cli.
//
// It records how many times a project chain has been accessed and when it was
// first and last used. This information can be used to surface stale chains,
// build usage reports, or drive cache eviction policies.
//
// Usage:
//
//	st, _ := store.New(dir)
//	m := access.New(st)
//
//	// Record an access event:
//	m.Touch("myproject")
//
//	// Retrieve the record:
//	rec, _ := m.Get("myproject")
//	fmt.Println(rec.Count, rec.LastUsed)
//
//	// Clear the record:
//	m.Reset("myproject")
package access
