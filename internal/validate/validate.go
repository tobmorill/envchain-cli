// Package validate provides helpers for checking environment variable
// entries against user-defined rules before they are stored or exported.
package validate

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/envchain-cli/envchain/internal/env"
)

// Rule describes a single validation constraint.
type Rule struct {
	// KeyPattern is an optional regex that the variable name must match.
	KeyPattern string
	// Required lists keys that must be present in the entry set.
	Required []string
	// ForbidEmpty disallows entries whose value is an empty string.
	ForbidEmpty bool
}

// Violation holds a single failed constraint.
type Violation struct {
	Key     string
	Message string
}

func (v Violation) Error() string {
	return fmt.Sprintf("%s: %s", v.Key, v.Message)
}

// Validate checks entries against the provided Rule and returns all
// violations found. A nil/empty slice means the entries are valid.
func Validate(entries []env.Entry, r Rule) []Violation {
	var violations []Violation

	var keyRe *regexp.Regexp
	if r.KeyPattern != "" {
		var err error
		keyRe, err = regexp.Compile(r.KeyPattern)
		if err != nil {
			violations = append(violations, Violation{
				Key:     "<rule>",
				Message: fmt.Sprintf("invalid key_pattern regex: %v", err),
			})
			return violations
		}
	}

	present := make(map[string]bool, len(entries))
	for _, e := range entries {
		present[e.Key] = true

		if keyRe != nil && !keyRe.MatchString(e.Key) {
			violations = append(violations, Violation{
				Key:     e.Key,
				Message: fmt.Sprintf("key does not match pattern %q", r.KeyPattern),
			})
		}

		if r.ForbidEmpty && strings.TrimSpace(e.Value) == "" {
			violations = append(violations, Violation{
				Key:     e.Key,
				Message: "value must not be empty",
			})
		}
	}

	for _, req := range r.Required {
		if !present[req] {
			violations = append(violations, Violation{
				Key:     req,
				Message: "required key is missing",
			})
		}
	}

	return violations
}

// ErrValidation is returned by Check when at least one violation exists.
var ErrValidation = errors.New("validation failed")

// Check is a convenience wrapper that returns a combined error when any
// violations are found, or nil when all entries pass.
func Check(entries []env.Entry, r Rule) error {
	violations := Validate(entries, r)
	if len(violations) == 0 {
		return nil
	}
	msgs := make([]string, len(violations))
	for i, v := range violations {
		msgs[i] = v.Error()
	}
	return fmt.Errorf("%w: %s", ErrValidation, strings.Join(msgs, "; "))
}
