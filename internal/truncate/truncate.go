// Package truncate provides helpers for safely shortening environment
// variable values for display purposes without exposing sensitive data.
package truncate

import "strings"

const (
	// DefaultMaxLen is the default maximum number of visible characters.
	DefaultMaxLen = 32
	// DefaultMask is the suffix appended when a value is truncated.
	DefaultMask = "…"
)

// Options controls how a value is truncated.
type Options struct {
	// MaxLen is the maximum number of runes to keep before appending Mask.
	MaxLen int
	// Mask is the string appended to indicate truncation.
	Mask string
	// RedactAll replaces the entire value with asterisks regardless of length.
	RedactAll bool
}

// Value shortens s according to opts. If opts is nil, defaults are used.
func Value(s string, opts *Options) string {
	o := defaults(opts)
	if o.RedactAll {
		return strings.Repeat("*", min(len(s), 8))
	}
	runes := []rune(s)
	if len(runes) <= o.MaxLen {
		return s
	}
	return string(runes[:o.MaxLen]) + o.Mask
}

// Preview returns a short preview of s suitable for terminal display.
// It always uses DefaultMaxLen and DefaultMask.
func Preview(s string) string {
	return Value(s, nil)
}

// Redact replaces every character in s with '*', capped at 8 characters,
// so the presence of a value is visible without leaking its content.
func Redact(s string) string {
	if s == "" {
		return ""
	}
	return Value(s, &Options{RedactAll: true})
}

func defaults(o *Options) Options {
	if o == nil {
		return Options{MaxLen: DefaultMaxLen, Mask: DefaultMask}
	}
	out := *o
	if out.MaxLen <= 0 {
		out.MaxLen = DefaultMaxLen
	}
	if out.Mask == "" {
		out.Mask = DefaultMask
	}
	return out
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
