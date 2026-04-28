// Package scope provides utilities for resolving and managing environment
// variable scopes, allowing chains to inherit or override values from a
// parent (global) scope before being applied to the shell.
package scope

import (
	"errors"

	"github.com/user/envchain-cli/internal/env"
)

// ErrCircularInheritance is returned when a scope inherits from itself.
var ErrCircularInheritance = errors.New("scope: circular inheritance detected")

// Scope represents a named layer of environment variables.
type Scope struct {
	Name   string
	Parent string // empty means no parent
	Entries []env.Entry
}

// Resolve merges parent entries into child entries, with child values taking
// precedence over parent values for duplicate keys.
// If parent and child share the same Name, ErrCircularInheritance is returned.
func Resolve(parent, child Scope) ([]env.Entry, error) {
	if parent.Name != "" && parent.Name == child.Name {
		return nil, ErrCircularInheritance
	}

	merged := make(map[string]env.Entry, len(parent.Entries)+len(child.Entries))

	for _, e := range parent.Entries {
		merged[e.Key] = e
	}
	for _, e := range child.Entries {
		merged[e.Key] = e
	}

	result := make([]env.Entry, 0, len(merged))
	for _, e := range merged {
		result = append(result, e)
	}

	sortEntries(result)
	return result, nil
}

// ResolveChain resolves a linear chain of scopes from outermost (index 0) to
// innermost (last index). Each scope overrides keys from the previous ones.
// Returns ErrCircularInheritance if any two adjacent scopes share the same name.
func ResolveChain(scopes []Scope) ([]env.Entry, error) {
	if len(scopes) == 0 {
		return nil, nil
	}

	current := scopes[0]
	for i := 1; i < len(scopes); i++ {
		resolved, err := Resolve(current, scopes[i])
		if err != nil {
			return nil, err
		}
		current = Scope{
			Name:    scopes[i].Name,
			Entries: resolved,
		}
	}
	return current.Entries, nil
}

// sortEntries sorts entries by key for deterministic output.
func sortEntries(entries []env.Entry) {
	for i := 1; i < len(entries); i++ {
		for j := i; j > 0 && entries[j].Key < entries[j-1].Key; j-- {
			entries[j], entries[j-1] = entries[j-1], entries[j]
		}
	}
}
