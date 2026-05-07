package complexity_test

import (
	"strings"
	"testing"

	"github.com/user/envchain-cli/internal/complexity"
	"github.com/user/envchain-cli/internal/env"
)

func entries(pairs ...string) []env.Entry {
	var out []env.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, env.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestEmptyChainIsLow(t *testing.T) {
	r := complexity.Evaluate(nil)
	if r.Level != complexity.LevelLow {
		t.Fatalf("expected low, got %s", r.Level)
	}
	if len(r.Findings) != 0 {
		t.Fatalf("expected no findings, got %d", len(r.Findings))
	}
}

func TestLongValueAddsFindings(t *testing.T) {
	big := strings.Repeat("x", 300)
	r := complexity.Evaluate(entries("TOKEN", big))
	if r.Level == complexity.LevelLow {
		t.Fatal("expected at least medium for long value")
	}
	if len(r.Findings) == 0 {
		t.Fatal("expected at least one finding")
	}
}

func TestShellInterpolationDetected(t *testing.T) {
	r := complexity.Evaluate(entries("CMD", "$(whoami)"))
	found := false
	for _, f := range r.Findings {
		if f.Key == "CMD" {
			found = true
		}
	}
	if !found {
		t.Fatal("expected finding for CMD")
	}
}

func TestDoubleUnderscoreKey(t *testing.T) {
	r := complexity.Evaluate(entries("APP__DB__HOST", "localhost"))
	if len(r.Findings) == 0 {
		t.Fatal("expected finding for double-underscore key")
	}
}

func TestScoreIsNonNegative(t *testing.T) {
	r := complexity.Evaluate(entries("A", "b", "C", "d"))
	if r.Score < 0 {
		t.Fatalf("score must not be negative, got %f", r.Score)
	}
}

func TestHighLevelForManyProblematicEntries(t *testing.T) {
	var pairs []string
	for i := 0; i < 20; i++ {
		pairs = append(pairs, "KEY__NS", "$(echo "+strings.Repeat("z", 300)+")")
	}
	r := complexity.Evaluate(entries(pairs...))
	if r.Level != complexity.LevelHigh {
		t.Fatalf("expected high, got %s (score %.2f)", r.Level, r.Score)
	}
}
