// Package quota enforces per-project limits on the number of stored
// environment variable entries and the total size of their values.
package quota

import (
	"errors"
	"fmt"

	"github.com/envchain-cli/envchain/internal/env"
)

// DefaultMaxKeys is the default maximum number of keys allowed per project.
const DefaultMaxKeys = 100

// DefaultMaxValueBytes is the default maximum total byte size of all values.
const DefaultMaxValueBytes = 64 * 1024 // 64 KiB

// ErrQuotaExceeded is returned when a quota limit would be breached.
var ErrQuotaExceeded = errors.New("quota exceeded")

// Rule describes the limits applied to a project's chain.
type Rule struct {
	MaxKeys       int
	MaxValueBytes int
}

// DefaultRule returns a Rule populated with the package-level defaults.
func DefaultRule() Rule {
	return Rule{
		MaxKeys:       DefaultMaxKeys,
		MaxValueBytes: DefaultMaxValueBytes,
	}
}

// Violation records a single quota breach.
type Violation struct {
	Field   string
	Limit   int
	Actual  int
	Message string
}

func (v Violation) Error() string { return v.Message }

// Check evaluates entries against r and returns any violations found.
// An empty slice means no limits were exceeded.
func Check(entries []env.Entry, r Rule) []Violation {
	var violations []Violation

	if r.MaxKeys > 0 && len(entries) > r.MaxKeys {
		violations = append(violations, Violation{
			Field:  "keys",
			Limit:  r.MaxKeys,
			Actual: len(entries),
			Message: fmt.Sprintf(
				"key count %d exceeds limit of %d", len(entries), r.MaxKeys,
			),
		})
	}

	if r.MaxValueBytes > 0 {
		total := 0
		for _, e := range entries {
			total += len(e.Value)
		}
		if total > r.MaxValueBytes {
			violations = append(violations, Violation{
				Field:  "value_bytes",
				Limit:  r.MaxValueBytes,
				Actual: total,
				Message: fmt.Sprintf(
					"total value size %d bytes exceeds limit of %d bytes", total, r.MaxValueBytes,
				),
			})
		}
	}

	return violations
}
