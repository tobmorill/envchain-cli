// Package complexity evaluates the structural complexity of an environment
// variable set, producing a numeric score and categorised findings.
package complexity

import (
	"math"
	"strings"

	"github.com/user/envchain-cli/internal/env"
)

// Level describes the overall complexity tier of a chain.
type Level string

const (
	LevelLow    Level = "low"
	LevelMedium Level = "medium"
	LevelHigh   Level = "high"
)

// Finding is a single observation that contributes to the complexity score.
type Finding struct {
	Key     string
	Reason  string
	Penalty float64
}

// Result holds the aggregated complexity assessment.
type Result struct {
	Score    float64
	Level    Level
	Findings []Finding
}

// Evaluate scores the supplied entries and returns a Result.
func Evaluate(entries []env.Entry) Result {
	var findings []Finding
	var score float64

	for _, e := range entries {
		if len(e.Value) > 256 {
			f := Finding{Key: e.Key, Reason: "value exceeds 256 bytes", Penalty: 2.0}
			findings = append(findings, f)
			score += f.Penalty
		}
		if strings.ContainsAny(e.Value, "$`") {
			f := Finding{Key: e.Key, Reason: "value contains shell-interpolation characters", Penalty: 1.5}
			findings = append(findings, f)
			score += f.Penalty
		}
		if strings.Contains(e.Key, "__") {
			f := Finding{Key: e.Key, Reason: "key uses double-underscore namespace separator", Penalty: 0.5}
			findings = append(findings, f)
			score += f.Penalty
		}
	}

	// Base penalty proportional to chain size.
	base := math.Log1p(float64(len(entries))) * 0.8
	score += base

	return Result{
		Score:    math.Round(score*100) / 100,
		Level:    classify(score),
		Findings: findings,
	}
}

func classify(score float64) Level {
	switch {
	case score < 3.0:
		return LevelLow
	case score < 8.0:
		return LevelMedium
	default:
		return LevelHigh
	}
}
