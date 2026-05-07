// Package visibility manages per-project key visibility settings for
// envchain-cli. Each environment variable key within a project can be
// individually marked as either visible (displayed in plain text during
// listing and inspection commands) or hidden (redacted to a placeholder
// such as "***" regardless of the active redact policy).
//
// Visibility settings are stored as JSON records in the shared key-value
// store and are keyed by project name. The zero value for any key that
// has not been explicitly configured is LevelVisible, preserving
// backward-compatible behaviour for existing projects.
package visibility
