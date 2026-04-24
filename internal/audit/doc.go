// Package audit implements an append-only audit log for envchain-cli.
//
// Each operation that reads, writes, or deletes a chain is recorded as a
// JSON-Lines entry containing a timestamp, the event kind, the project name,
// and an optional human-readable message.
//
// The log file is stored under the envchain data directory (e.g.
// ~/.local/share/envchain/audit.jsonl on Linux) and is only readable by the
// current user (mode 0600).
//
// Usage:
//
//	logger, err := audit.NewLogger(path)
//	if err != nil { ... }
//
//	logger.Record(audit.EventSave, "myproject", "")
//
//	events, err := logger.ReadAll()
package audit
