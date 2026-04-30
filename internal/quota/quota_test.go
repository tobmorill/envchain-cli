package quota_test

import (
	"strings"
	"testing"

	"github.com/envchain-cli/envchain/internal/env"
	"github.com/envchain-cli/envchain/internal/quota"
)

func makeEntries(keys []string, valueSize int) []env.Entry {
	entries := make([]env.Entry, len(keys))
	for i, k := range keys {
		entries[i] = env.Entry{Key: k, Value: strings.Repeat("x", valueSize)}
	}
	return entries
}

func TestNoViolationsWithinLimits(t *testing.T) {
	r := quota.Rule{MaxKeys: 5, MaxValueBytes: 1000}
	entries := makeEntries([]string{"A", "B", "C"}, 10)
	if v := quota.Check(entries, r); len(v) != 0 {
		t.Fatalf("expected no violations, got %v", v)
	}
}

func TestKeyCountExceedsLimit(t *testing.T) {
	r := quota.Rule{MaxKeys: 2, MaxValueBytes: 0}
	entries := makeEntries([]string{"A", "B", "C"}, 1)
	v := quota.Check(entries, r)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Field != "keys" {
		t.Errorf("expected field 'keys', got %q", v[0].Field)
	}
	if v[0].Actual != 3 || v[0].Limit != 2 {
		t.Errorf("unexpected limit/actual: %+v", v[0])
	}
}

func TestValueBytesExceedsLimit(t *testing.T) {
	r := quota.Rule{MaxKeys: 0, MaxValueBytes: 10}
	entries := makeEntries([]string{"A", "B"}, 8) // 16 bytes total
	v := quota.Check(entries, r)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Field != "value_bytes" {
		t.Errorf("expected field 'value_bytes', got %q", v[0].Field)
	}
}

func TestBothLimitsBreached(t *testing.T) {
	r := quota.Rule{MaxKeys: 1, MaxValueBytes: 5}
	entries := makeEntries([]string{"A", "B"}, 4)
	v := quota.Check(entries, r)
	if len(v) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(v))
	}
}

func TestDefaultRuleValues(t *testing.T) {
	r := quota.DefaultRule()
	if r.MaxKeys != quota.DefaultMaxKeys {
		t.Errorf("MaxKeys: want %d got %d", quota.DefaultMaxKeys, r.MaxKeys)
	}
	if r.MaxValueBytes != quota.DefaultMaxValueBytes {
		t.Errorf("MaxValueBytes: want %d got %d", quota.DefaultMaxValueBytes, r.MaxValueBytes)
	}
}

func TestZeroLimitDisablesCheck(t *testing.T) {
	r := quota.Rule{MaxKeys: 0, MaxValueBytes: 0}
	entries := makeEntries([]string{"A", "B", "C", "D", "E"}, 512)
	if v := quota.Check(entries, r); len(v) != 0 {
		t.Fatalf("expected no violations when limits are zero, got %v", v)
	}
}
