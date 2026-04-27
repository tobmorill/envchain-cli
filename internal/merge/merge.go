// Package merge provides utilities for merging environment variable sets
// from multiple chains into a single resolved set, with configurable
// conflict resolution strategies.
package merge

import (
	"errors"
	"fmt"

	"github.com/user/envchain-cli/internal/env"
)

// Strategy controls how key conflicts are resolved when merging.
type Strategy int

const (
	// StrategyFirst keeps the value from the first chain that defines the key.
	StrategyFirst Strategy = iota
	// StrategyLast overwrites with the value from the last chain that defines the key.
	StrategyLast
	// StrategyError returns an error if any key appears in more than one chain.
	StrategyError
)

// ErrConflict is returned when StrategyError is used and a duplicate key is found.
type ErrConflict struct {
	Key    string
	Chains []string
}

func (e *ErrConflict) Error() string {
	return fmt.Sprintf("merge conflict: key %q defined in multiple chains: %v", e.Key, e.Chains)
}

// Result holds the merged entries and metadata about the merge.
type Result struct {
	Entries []env.Entry
	// Origins maps each key to the chain name it was sourced from.
	Origins map[string]string
}

// Merge combines multiple named entry slices according to the given strategy.
// The chains parameter is an ordered slice of (name, entries) pairs.
func Merge(chains []NamedChain, strategy Strategy) (*Result, error) {
	if len(chains) == 0 {
		return &Result{Origins: make(map[string]string)}, nil
	}

	seen := make(map[string]string) // key -> chain name
	index := make(map[string]env.Entry)
	order := []string{}

	for _, nc := range chains {
		for _, entry := range nc.Entries {
			if origin, exists := seen[entry.Key]; exists {
				switch strategy {
				case StrategyError:
					return nil, &ErrConflict{Key: entry.Key, Chains: []string{origin, nc.Name}}
				case StrategyFirst:
					continue
				case StrategyLast:
					index[entry.Key] = entry
					seen[entry.Key] = nc.Name
				}
			} else {
				seen[entry.Key] = nc.Name
				index[entry.Key] = entry
				order = append(order, entry.Key)
			}
		}
	}

	result := &Result{Origins: seen}
	for _, k := range order {
		result.Entries = append(result.Entries, index[k])
	}
	return result, nil
}

// NamedChain associates a chain name with its resolved entries.
type NamedChain struct {
	Name    string
	Entries []env.Entry
}

// ErrEmptyChainName is returned when a chain name is blank.
var ErrEmptyChainName = errors.New("chain name must not be empty")

// NewNamedChain constructs a NamedChain, validating that the name is non-empty.
func NewNamedChain(name string, entries []env.Entry) (NamedChain, error) {
	if name == "" {
		return NamedChain{}, ErrEmptyChainName
	}
	return NamedChain{Name: name, Entries: entries}, nil
}
