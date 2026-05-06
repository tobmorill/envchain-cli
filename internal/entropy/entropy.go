// Package entropy measures the randomness/complexity of environment variable values
// and flags entries that may be secrets based on Shannon entropy scoring.
package entropy

import (
	"math"

	"github.com/envchain/envchain-cli/internal/env"
)

// Result holds the entropy analysis for a single entry.
type Result struct {
	Key     string
	Value   string
	Score   float64 // Shannon entropy in bits
	IsHigh  bool    // true when Score >= HighThreshold
}

// DefaultHighThreshold is the entropy score above which a value is considered
// high-entropy (likely a secret or random token).
const DefaultHighThreshold = 3.5

// Analyzer computes entropy scores for environment variable entries.
type Analyzer struct {
	HighThreshold float64
}

// New returns an Analyzer with the default high-entropy threshold.
func New() *Analyzer {
	return &Analyzer{HighThreshold: DefaultHighThreshold}
}

// Score computes the Shannon entropy of a string in bits per character.
func Score(s string) float64 {
	if len(s) == 0 {
		return 0
	}
	freq := make(map[rune]int)
	for _, ch := range s {
		freq[ch]++
	}
	n := float64(len([]rune(s)))
	var h float64
	for _, count := range freq {
		p := float64(count) / n
		h -= p * math.Log2(p)
	}
	return h
}

// Analyze returns entropy Results for each entry in the provided slice.
func (a *Analyzer) Analyze(entries []env.Entry) []Result {
	results := make([]Result, 0, len(entries))
	for _, e := range entries {
		s := Score(e.Value)
		results = append(results, Result{
			Key:    e.Key,
			Value:  e.Value,
			Score:  s,
			IsHigh: s >= a.HighThreshold,
		})
	}
	return results
}

// HighEntropy filters the results and returns only those with IsHigh == true.
func HighEntropy(results []Result) []Result {
	var out []Result
	for _, r := range results {
		if r.IsHigh {
			out = append(out, r)
		}
	}
	return out
}
