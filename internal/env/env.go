package env

import (
	"errors"
	"fmt"
	"strings"
)

// ErrInvalidEntry is returned when an environment variable entry is malformed.
var ErrInvalidEntry = errors.New("invalid environment variable entry")

// Entry represents a single environment variable key-value pair.
type Entry struct {
	Key   string
	Value string
}

// Parse parses a string in the form KEY=VALUE into an Entry.
// The key must be non-empty and must not contain '='.
func Parse(s string) (Entry, error) {
	parts := strings.SplitN(s, "=", 2)
	if len(parts) != 2 {
		return Entry{}, fmt.Errorf("%w: %q (expected KEY=VALUE)", ErrInvalidEntry, s)
	}
	key := parts[0]
	if key == "" {
		return Entry{}, fmt.Errorf("%w: key must not be empty", ErrInvalidEntry)
	}
	return Entry{Key: key, Value: parts[1]}, nil
}

// String returns the entry formatted as KEY=VALUE.
func (e Entry) String() string {
	return e.Key + "=" + e.Value
}

// ParseAll parses a slice of KEY=VALUE strings into a map.
// Returns an error if any entry is malformed or if a key appears more than once.
func ParseAll(lines []string) (map[string]string, error) {
	result := make(map[string]string, len(lines))
	for _, line := range lines {
		entry, err := Parse(line)
		if err != nil {
			return nil, err
		}
		if _, exists := result[entry.Key]; exists {
			return nil, fmt.Errorf("%w: duplicate key %q", ErrInvalidEntry, entry.Key)
		}
		result[entry.Key] = entry.Value
	}
	return result, nil
}

// ToLines converts a map of environment variables to a sorted slice of KEY=VALUE strings.
func ToLines(vars map[string]string) []string {
	lines := make([]string, 0, len(vars))
	for k, v := range vars {
		lines = append(lines, k+"="+v)
	}
	return lines
}

// ExportScript returns a shell script snippet that exports all variables.
func ExportScript(vars map[string]string) string {
	var sb strings.Builder
	for k, v := range vars {
		fmt.Fprintf(&sb, "export %s=%q\n", k, v)
	}
	return sb.String()
}
