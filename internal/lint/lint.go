// Package lint provides heuristic checks for environment variable chains,
// warning about common mistakes such as keys with whitespace, values that look
// like unexpanded shell variables, or duplicate keys across merged chains.
package lint

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/user/envchain-cli/internal/env"
)

// Severity indicates how serious a lint finding is.
type Severity string

const (
	Warn  Severity = "warn"
	Error Severity = "error"
)

// Finding describes a single lint issue.
type Finding struct {
	Key      string
	Message  string
	Severity Severity
}

func (f Finding) String() string {
	return fmt.Sprintf("[%s] %s: %s", f.Severity, f.Key, f.Message)
}

// Check runs all built-in lint rules against entries and returns any findings.
func Check(entries []env.Entry) []Finding {
	var findings []Finding
	seen := make(map[string]bool)

	for _, e := range entries {
		// Duplicate key detection.
		norm := strings.ToUpper(e.Key)
		if seen[norm] {
			findings = append(findings, Finding{
				Key:      e.Key,
				Message:  "duplicate key",
				Severity: Error,
			})
		}
		seen[norm] = true

		// Key contains whitespace.
		if strings.IndexFunc(e.Key, unicode.IsSpace) >= 0 {
			findings = append(findings, Finding{
				Key:      e.Key,
				Message:  "key contains whitespace",
				Severity: Error,
			})
		}

		// Value looks like an unexpanded shell variable.
		if strings.Contains(e.Value, "${") || (strings.HasPrefix(e.Value, "$") && len(e.Value) > 1 && e.Value[1] != '(') {
			findings = append(findings, Finding{
				Key:      e.Key,
				Message:  "value appears to contain an unexpanded shell variable",
				Severity: Warn,
			})
		}

		// Empty value warning.
		if strings.TrimSpace(e.Value) == "" {
			findings = append(findings, Finding{
				Key:      e.Key,
				Message:  "value is empty",
				Severity: Warn,
			})
		}
	}

	return findings
}
