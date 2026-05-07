// Package budget tracks cumulative value-byte consumption per project and
// enforces a configurable ceiling.
//
// Records are stored as JSON files under a configurable directory, one file
// per project.  The Manager provides Set, Get, Check, and Delete operations.
//
// Check is intentionally lenient: if no record exists for a project, the
// call succeeds so that projects without an explicit budget are unrestricted.
//
// Typical usage:
//
//	m := budget.New(filepath.Join(home, ".envchain", "budgets"))
//	if err := m.Set("myproject", 8192); err != nil {
//		log.Fatal(err)
//	}
//	if err := m.Check("myproject", totalBytes); err != nil {
//		log.Fatal(err)
//	}
package budget
