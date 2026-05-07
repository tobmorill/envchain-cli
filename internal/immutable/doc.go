// Package immutable tracks which environment variable keys within a project
// are marked as immutable (read-only).
//
// An immutable key may still be read and exported; however, higher-level
// commands (merge, import, rotate) should consult IsImmutable before
// overwriting a value and surface an error or skip the key accordingly.
//
// Data is persisted in the same encrypted store used by other envchain
// subsystems, so the caller must supply the project passphrase on every
// read/write operation.
package immutable
