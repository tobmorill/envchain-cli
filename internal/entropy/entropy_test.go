package entropy_test

import (
	"testing"

	"github.com/envchain/envchain-cli/internal/entropy"
	"github.com/envchain/envchain-cli/internal/env"
)

func entries(pairs ...string) []env.Entry {
	var out []env.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, env.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestScoreEmptyString(t *testing.T) {
	if got := entropy.Score(""); got != 0 {
		t.Fatalf("expected 0 for empty string, got %f", got)
	}
}

func TestScoreSingleChar(t *testing.T) {
	if got := entropy.Score("aaaa"); got != 0 {
		t.Fatalf("expected 0 for uniform string, got %f", got)
	}
}

func TestScoreHighEntropyString(t *testing.T) {
	// A random-looking base64 token should score well above 3.5
	score := entropy.Score("aB3xQ9mZpL2wRvTn")
	if score < entropy.DefaultHighThreshold {
		t.Fatalf("expected high entropy score, got %f", score)
	}
}

func TestScoreLowEntropyString(t *testing.T) {
	score := entropy.Score("hello")
	if score >= entropy.DefaultHighThreshold {
		t.Fatalf("expected low entropy score for 'hello', got %f", score)
	}
}

func TestAnalyzeReturnsResultPerEntry(t *testing.T) {
	a := entropy.New()
	in := entries("TOKEN", "aB3xQ9mZpL2wRvTn", "ENV", "production")
	results := a.Analyze(in)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestAnalyzeIsHighFlag(t *testing.T) {
	a := entropy.New()
	in := entries("SECRET", "aB3xQ9mZpL2wRvTn", "NAME", "alice")
	results := a.Analyze(in)

	if !results[0].IsHigh {
		t.Errorf("expected SECRET to be flagged as high entropy")
	}
	if results[1].IsHigh {
		t.Errorf("expected NAME not to be flagged as high entropy")
	}
}

func TestHighEntropyFilter(t *testing.T) {
	a := entropy.New()
	in := entries(
		"SECRET", "aB3xQ9mZpL2wRvTn",
		"ENV", "production",
		"TOKEN", "xK8!mN2@pQ5#rS7$",
	)
	all := a.Analyze(in)
	high := entropy.HighEntropy(all)

	if len(high) != 2 {
		t.Fatalf("expected 2 high-entropy results, got %d", len(high))
	}
}

func TestAnalyzeEmptyInput(t *testing.T) {
	a := entropy.New()
	results := a.Analyze(nil)
	if len(results) != 0 {
		t.Fatalf("expected empty results for nil input, got %d", len(results))
	}
}

func TestCustomThreshold(t *testing.T) {
	a := &entropy.Analyzer{HighThreshold: 1.0}
	in := entries("WORD", "hello")
	results := a.Analyze(in)
	if !results[0].IsHigh {
		t.Errorf("expected 'hello' to be high entropy with threshold 1.0")
	}
}
