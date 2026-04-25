// Package diff provides utilities for comparing two sets of environment
// variable entries and producing a human-readable summary of changes.
package diff

import (
	"fmt"
	"sort"
	"strings"

	"github.com/envchain-cli/envchain/internal/env"
)

// ChangeKind describes the type of change for a single key.
type ChangeKind string

const (
	Added    ChangeKind = "added"
	Removed  ChangeKind = "removed"
	Modified ChangeKind = "modified"
	Unchanged ChangeKind = "unchanged"
)

// Change represents a single key-level difference between two env sets.
type Change struct {
	Key  string
	Kind ChangeKind
	Old  string
	New  string
}

// Compare returns the ordered list of changes between two slices of env entries.
// Values are masked so that secrets are not exposed in output.
func Compare(before, after []env.Entry) []Change {
	oldMap := make(map[string]string, len(before))
	for _, e := range before {
		oldMap[e.Key] = e.Value
	}

	newMap := make(map[string]string, len(after))
	for _, e := range after {
		newMap[e.Key] = e.Value
	}

	keys := make(map[string]struct{})
	for k := range oldMap {
		keys[k] = struct{}{}
	}
	for k := range newMap {
		keys[k] = struct{}{}
	}

	sorted := make([]string, 0, len(keys))
	for k := range keys {
		sorted = append(sorted, k)
	}
	sort.Strings(sorted)

	var changes []Change
	for _, k := range sorted {
		ov, inOld := oldMap[k]
		nv, inNew := newMap[k]
		switch {
		case inOld && !inNew:
			changes = append(changes, Change{Key: k, Kind: Removed, Old: mask(ov)})
		case !inOld && inNew:
			changes = append(changes, Change{Key: k, Kind: Added, New: mask(nv)})
		case ov != nv:
			changes = append(changes, Change{Key: k, Kind: Modified, Old: mask(ov), New: mask(nv)})
		default:
			changes = append(changes, Change{Key: k, Kind: Unchanged})
		}
	}
	return changes
}

// Summary returns a compact multi-line string describing only non-unchanged keys.
func Summary(changes []Change) string {
	var sb strings.Builder
	for _, c := range changes {
		switch c.Kind {
		case Added:
			fmt.Fprintf(&sb, "+ %s=%s\n", c.Key, c.New)
		case Removed:
			fmt.Fprintf(&sb, "- %s=%s\n", c.Key, c.Old)
		case Modified:
			fmt.Fprintf(&sb, "~ %s: %s -> %s\n", c.Key, c.Old, c.New)
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

// mask replaces all but the first character of a value with asterisks.
func mask(v string) string {
	if len(v) == 0 {
		return ""
	}
	return string(v[0]) + strings.Repeat("*", len(v)-1)
}
