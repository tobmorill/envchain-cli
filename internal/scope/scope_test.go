package scope_test

import (
	"testing"

	"github.com/user/envchain-cli/internal/env"
	"github.com/user/envchain-cli/internal/scope"
)

func makeEntries(pairs ...string) []env.Entry {
	if len(pairs)%2 != 0 {
		panic("makeEntries: pairs must be even")
	}
	out := make([]env.Entry, 0, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		out = append(out, env.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestResolveChildOverridesParent(t *testing.T) {
	parent := scope.Scope{Name: "global", Entries: makeEntries("FOO", "parent", "BAR", "shared")}
	child := scope.Scope{Name: "project", Entries: makeEntries("BAR", "child", "BAZ", "new")}

	result, err := scope.Resolve(parent, child)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := make(map[string]string, len(result))
	for _, e := range result {
		got[e.Key] = e.Value
	}

	if got["FOO"] != "parent" {
		t.Errorf("FOO: want parent, got %q", got["FOO"])
	}
	if got["BAR"] != "child" {
		t.Errorf("BAR: want child, got %q", got["BAR"])
	}
	if got["BAZ"] != "new" {
		t.Errorf("BAZ: want new, got %q", got["BAZ"])
	}
}

func TestResolveCircularInheritance(t *testing.T) {
	s := scope.Scope{Name: "loop", Entries: makeEntries("X", "1")}
	_, err := scope.Resolve(s, s)
	if err != scope.ErrCircularInheritance {
		t.Fatalf("want ErrCircularInheritance, got %v", err)
	}
}

func TestResolveEmptyParent(t *testing.T) {
	parent := scope.Scope{}
	child := scope.Scope{Name: "project", Entries: makeEntries("KEY", "val")}

	result, err := scope.Resolve(parent, child)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 1 || result[0].Key != "KEY" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestResolveChainLinear(t *testing.T) {
	scopes := []scope.Scope{
		{Name: "base", Entries: makeEntries("A", "1", "B", "2")},
		{Name: "mid", Entries: makeEntries("B", "overridden", "C", "3")},
		{Name: "top", Entries: makeEntries("C", "final")},
	}

	result, err := scope.ResolveChain(scopes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := make(map[string]string)
	for _, e := range result {
		got[e.Key] = e.Value
	}

	if got["A"] != "1" {
		t.Errorf("A: want 1, got %q", got["A"])
	}
	if got["B"] != "overridden" {
		t.Errorf("B: want overridden, got %q", got["B"])
	}
	if got["C"] != "final" {
		t.Errorf("C: want final, got %q", got["C"])
	}
}

func TestResolveChainEmpty(t *testing.T) {
	result, err := scope.ResolveChain(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Errorf("expected nil result, got %v", result)
	}
}
