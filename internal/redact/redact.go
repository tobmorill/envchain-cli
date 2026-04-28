// Package redact provides utilities for selectively masking environment
// variable values before display or logging, based on key name patterns.
package redact

import (
	"strings"

	"github.com/envchain/envchain-cli/internal/env"
)

// DefaultPatterns is the list of key substrings that trigger redaction.
var DefaultPatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"PRIVATE",
	"CREDENTIAL",
	"AUTH",
}

// Redactor masks sensitive environment variable values.
type Redactor struct {
	patterns []string
	mask     string
}

// New returns a Redactor using the supplied patterns and mask string.
// If patterns is nil, DefaultPatterns is used. If mask is empty, "***" is used.
func New(patterns []string, mask string) *Redactor {
	if patterns == nil {
		patterns = DefaultPatterns
	}
	if mask == "" {
		mask = "***"
	}
	return &Redactor{patterns: patterns, mask: mask}
}

// IsSensitive reports whether the key matches any redaction pattern.
func (r *Redactor) IsSensitive(key string) bool {
	upper := strings.ToUpper(key)
	for _, p := range r.patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}

// Apply returns a copy of entries where sensitive values are replaced with
// the mask string.
func (r *Redactor) Apply(entries []env.Entry) []env.Entry {
	out := make([]env.Entry, len(entries))
	for i, e := range entries {
		if r.IsSensitive(e.Key) {
			out[i] = env.Entry{Key: e.Key, Value: r.mask}
		} else {
			out[i] = e
		}
	}
	return out
}

// ApplyMap returns a copy of the map where sensitive values are replaced.
func (r *Redactor) ApplyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		if r.IsSensitive(k) {
			out[k] = r.mask
		} else {
			out[k] = v
		}
	}
	return out
}
