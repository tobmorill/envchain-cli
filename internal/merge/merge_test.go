package merge_test

import (
	"testing"

	"github.com/user/envchain-cli/internal/env"
	"github.com/user/envchain-cli/internal/merge"
)

func makeEntries(pairs ...string) []env.Entry {
	var out []env.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, env.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestMergeNoConflict(t *testing.T) {
	chains := []merge.NamedChain{
		{Name: "base", Entries: makeEntries("FOO", "1", "BAR", "2")},
		{Name: "extra", Entries: makeEntries("BAZ", "3")},
	}
	res, err := merge.Merge(chains, merge.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(res.Entries))
	}
	if res.Origins["FOO"] != "base" || res.Origins["BAZ"] != "extra" {
		t.Errorf("unexpected origins: %v", res.Origins)
	}
}

func TestMergeStrategyFirst(t *testing.T) {
	chains := []merge.NamedChain{
		{Name: "a", Entries: makeEntries("KEY", "first")},
		{Name: "b", Entries: makeEntries("KEY", "second")},
	}
	res, err := merge.Merge(chains, merge.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Entries[0].Value != "first" {
		t.Errorf("expected 'first', got %q", res.Entries[0].Value)
	}
}

func TestMergeStrategyLast(t *testing.T) {
	chains := []merge.NamedChain{
		{Name: "a", Entries: makeEntries("KEY", "first")},
		{Name: "b", Entries: makeEntries("KEY", "second")},
	}
	res, err := merge.Merge(chains, merge.StrategyLast)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Entries[0].Value != "second" {
		t.Errorf("expected 'second', got %q", res.Entries[0].Value)
	}
	if res.Origins["KEY"] != "b" {
		t.Errorf("expected origin 'b', got %q", res.Origins["KEY"])
	}
}

func TestMergeStrategyError(t *testing.T) {
	chains := []merge.NamedChain{
		{Name: "a", Entries: makeEntries("KEY", "v1")},
		{Name: "b", Entries: makeEntries("KEY", "v2")},
	}
	_, err := merge.Merge(chains, merge.StrategyError)
	if err == nil {
		t.Fatal("expected conflict error, got nil")
	}
	var ce *merge.ErrConflict
	if ok := isConflict(err, &ce); !ok {
		t.Fatalf("expected ErrConflict, got %T", err)
	}
	if ce.Key != "KEY" {
		t.Errorf("expected conflict key 'KEY', got %q", ce.Key)
	}
}

func TestMergeEmpty(t *testing.T) {
	res, err := merge.Merge(nil, merge.StrategyFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 0 {
		t.Errorf("expected empty result")
	}
}

func TestNewNamedChainEmptyName(t *testing.T) {
	_, err := merge.NewNamedChain("", nil)
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func isConflict(err error, target **merge.ErrConflict) bool {
	if ce, ok := err.(*merge.ErrConflict); ok {
		*target = ce
		return true
	}
	return false
}
