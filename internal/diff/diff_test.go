package diff_test

import (
	"strings"
	"testing"

	"github.com/envchain-cli/envchain/internal/diff"
	"github.com/envchain-cli/envchain/internal/env"
)

func entries(pairs ...string) []env.Entry {
	var out []env.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, env.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestCompareAdded(t *testing.T) {
	before := entries()
	after := entries("FOO", "bar")
	changes := diff.Compare(before, after)
	if len(changes) != 1 || changes[0].Kind != diff.Added || changes[0].Key != "FOO" {
		t.Fatalf("expected one Added change, got %+v", changes)
	}
}

func TestCompareRemoved(t *testing.T) {
	before := entries("FOO", "bar")
	after := entries()
	changes := diff.Compare(before, after)
	if len(changes) != 1 || changes[0].Kind != diff.Removed {
		t.Fatalf("expected one Removed change, got %+v", changes)
	}
}

func TestCompareModified(t *testing.T) {
	before := entries("FOO", "old")
	after := entries("FOO", "new")
	changes := diff.Compare(before, after)
	if len(changes) != 1 || changes[0].Kind != diff.Modified {
		t.Fatalf("expected one Modified change, got %+v", changes)
	}
}

func TestCompareUnchanged(t *testing.T) {
	before := entries("FOO", "same")
	after := entries("FOO", "same")
	changes := diff.Compare(before, after)
	if len(changes) != 1 || changes[0].Kind != diff.Unchanged {
		t.Fatalf("expected one Unchanged change, got %+v", changes)
	}
}

func TestCompareSortedOrder(t *testing.T) {
	before := entries("ZZZ", "1", "AAA", "2")
	after := entries("ZZZ", "1", "AAA", "2")
	changes := diff.Compare(before, after)
	if changes[0].Key != "AAA" || changes[1].Key != "ZZZ" {
		t.Fatalf("expected alphabetical order, got %v %v", changes[0].Key, changes[1].Key)
	}
}

func TestSummaryOnlyNonUnchanged(t *testing.T) {
	before := entries("KEEP", "x", "OLD", "y")
	after := entries("KEEP", "x", "NEW", "z")
	changes := diff.Compare(before, after)
	summary := diff.Summary(changes)
	if strings.Contains(summary, "KEEP") {
		t.Errorf("summary should not include unchanged keys, got: %s", summary)
	}
	if !strings.Contains(summary, "- OLD") {
		t.Errorf("summary should include removed OLD, got: %s", summary)
	}
	if !strings.Contains(summary, "+ NEW") {
		t.Errorf("summary should include added NEW, got: %s", summary)
	}
}

func TestMaskingHidesValues(t *testing.T) {
	before := entries()
	after := entries("SECRET", "password123")
	changes := diff.Compare(before, after)
	if strings.Contains(changes[0].New, "password") {
		t.Errorf("masked value should not contain plaintext, got: %s", changes[0].New)
	}
	if !strings.HasPrefix(changes[0].New, "p") {
		t.Errorf("masked value should preserve first char, got: %s", changes[0].New)
	}
}
