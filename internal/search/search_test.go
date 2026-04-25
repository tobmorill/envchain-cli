package search_test

import (
	"os"
	"testing"

	"github.com/user/envchain-cli/internal/chain"
	"github.com/user/envchain-cli/internal/env"
	"github.com/user/envchain-cli/internal/search"
	"github.com/user/envchain-cli/internal/store"
)

const testPassphrase = "hunter2"

func newTempSearch(t *testing.T) (*search.Manager, *chain.Manager) {
	t.Helper()
	dir, err := os.MkdirTemp("", "envchain-search-*")
	if err != nil {
		t.Fatalf("MkdirTemp: %v", err)
	}
	t.Cleanup(func() { os.RemoveAll(dir) })
	st, err := store.New(dir)
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	cm := chain.New(st)
	return search.New(cm), cm
}

func seed(t *testing.T, cm *chain.Manager, project string, keys ...string) {
	t.Helper()
	var entries []env.Entry
	for _, k := range keys {
		entries = append(entries, env.Entry{Key: k, Value: "val"})
	}
	if err := cm.Save(project, testPassphrase, entries); err != nil {
		t.Fatalf("Save(%q): %v", project, err)
	}
}

func TestFindProjectsAll(t *testing.T) {
	sm, _ := newTempSearch(t)
	names := []string{"alpha", "beta", "gamma"}
	results := sm.FindProjects("", names)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
}

func TestFindProjectsFiltered(t *testing.T) {
	sm, _ := newTempSearch(t)
	names := []string{"alpha", "beta", "alphabet"}
	results := sm.FindProjects("alph", names)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, r := range results {
		if r.Project != "alpha" && r.Project != "alphabet" {
			t.Errorf("unexpected project %q", r.Project)
		}
	}
}

func TestFindProjectsCaseInsensitive(t *testing.T) {
	sm, _ := newTempSearch(t)
	names := []string{"MyProject", "other"}
	results := sm.FindProjects("myproject", names)
	if len(results) != 1 || results[0].Project != "MyProject" {
		t.Fatalf("expected MyProject, got %v", results)
	}
}

func TestFindKeysMatchesKey(t *testing.T) {
	sm, cm := newTempSearch(t)
	seed(t, cm, "proj", "DB_HOST", "DB_PORT", "API_KEY")
	results, err := sm.FindKeys("DB_", testPassphrase, []string{"proj"})
	if err != nil {
		t.Fatalf("FindKeys: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestFindKeysWrongPassphraseSkipped(t *testing.T) {
	sm, cm := newTempSearch(t)
	seed(t, cm, "proj", "SECRET")
	results, err := sm.FindKeys("SECRET", "wrongpass", []string{"proj"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results for wrong passphrase, got %d", len(results))
	}
}
