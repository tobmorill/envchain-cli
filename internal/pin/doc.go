// Package pin manages per-project pinned environment variable keys.
//
// Pinned keys are stored as a simple JSON list on disk. When a chain is
// loaded into the shell, the caller can consult the pin list to decide
// which variables should always be exported, even when switching between
// multiple chains for the same project.
//
// Usage:
//
//	m := pin.New("/path/to/pin/dir")
//
//	// Pin two keys for a project.
//	m.Set("myproject", []string{"DATABASE_URL", "API_SECRET"})
//
//	// Retrieve pinned keys.
//	keys, err := m.Get("myproject")
//
//	// Remove pins.
//	m.Delete("myproject")
package pin
