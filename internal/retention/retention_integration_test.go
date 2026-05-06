package retention_test

import (
	"testing"
	"time"

	"github.com/envchain-cli/internal/retention"
)

// TestMultipleProjectsAreIsolated ensures that setting a policy for one
// project does not affect another.
func TestMultipleProjectsAreIsolated(t *testing.T) {
	m := newTempManager(t)

	alpha := retention.Policy{Project: "alpha", MaxAge: time.Hour, MaxVersions: 5}
	beta := retention.Policy{Project: "beta", MaxAge: 48 * time.Hour, MaxVersions: 20}

	if err := m.Set(alpha); err != nil {
		t.Fatalf("Set alpha: %v", err)
	}
	if err := m.Set(beta); err != nil {
		t.Fatalf("Set beta: %v", err)
	}

	gotAlpha, err := m.Get("alpha")
	if err != nil {
		t.Fatalf("Get alpha: %v", err)
	}
	gotBeta, err := m.Get("beta")
	if err != nil {
		t.Fatalf("Get beta: %v", err)
	}

	if gotAlpha.MaxAge != time.Hour {
		t.Errorf("alpha MaxAge: got %v, want %v", gotAlpha.MaxAge, time.Hour)
	}
	if gotBeta.MaxVersions != 20 {
		t.Errorf("beta MaxVersions: got %d, want 20", gotBeta.MaxVersions)
	}
}

// TestOverwriteUpdatesPolicy verifies that calling Set twice replaces the
// previous policy rather than merging.
func TestOverwriteUpdatesPolicy(t *testing.T) {
	m := newTempManager(t)

	first := retention.Policy{Project: "proj", MaxAge: time.Hour, MaxVersions: 3}
	if err := m.Set(first); err != nil {
		t.Fatalf("Set first: %v", err)
	}

	second := retention.Policy{Project: "proj", MaxAge: 72 * time.Hour, MaxVersions: 50}
	if err := m.Set(second); err != nil {
		t.Fatalf("Set second: %v", err)
	}

	got, err := m.Get("proj")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.MaxAge != 72*time.Hour {
		t.Errorf("MaxAge: got %v, want %v", got.MaxAge, 72*time.Hour)
	}
	if got.MaxVersions != 50 {
		t.Errorf("MaxVersions: got %d, want 50", got.MaxVersions)
	}
}
