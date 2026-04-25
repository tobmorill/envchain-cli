// Package snapshot provides point-in-time captures of environment variable
// chains stored within an envchain project.
//
// A snapshot records the full set of key=value entries belonging to a chain
// at the moment Save is called. Snapshots are stored alongside chain data in
// the same encrypted store under a namespaced key, so no additional storage
// backend is required.
//
// Typical usage:
//
//	mgr := snapshot.New(st, cm)
//
//	// Capture current state before a destructive edit.
//	if err := mgr.Save(project, "before-refactor", passphrase); err != nil {
//	    log.Fatal(err)
//	}
//
//	// Retrieve a snapshot later.
//	snap, err := mgr.Get(project, "before-refactor")
package snapshot
