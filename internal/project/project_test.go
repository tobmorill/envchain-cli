package project_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/envchain-cli/envchain/internal/project"
)

func TestResolveUsesDirectoryName(t *testing.T) {
	dir := t.TempDir()
	name, err := project.Resolve(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if want := filepath.Base(dir); name != want {
		t.Errorf("got %q, want %q", name, want)
	}
}

func TestResolveUsesMarkerFile(t *testing.T) {
	dir := t.TempDir()
	if err := project.WriteMarker(dir, "my-project"); err != nil {
		t.Fatalf("WriteMarker: %v", err)
	}
	name, err := project.Resolve(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name != "my-project" {
		t.Errorf("got %q, want %q", name, "my-project")
	}
}

func TestWriteMarkerCreatesFile(t *testing.T) {
	dir := t.TempDir()
	if err := project.WriteMarker(dir, "test-proj"); err != nil {
		t.Fatalf("WriteMarker: %v", err)
	}
	path := filepath.Join(dir, ".envchain")
	if _, err := os.Stat(path); err != nil {
		t.Errorf("marker file not created: %v", err)
	}
}

func TestWriteMarkerOverwrites(t *testing.T) {
	dir := t.TempDir()
	project.WriteMarker(dir, "old-name")
	if err := project.WriteMarker(dir, "new-name"); err != nil {
		t.Fatalf("WriteMarker: %v", err)
	}
	name, err := project.Resolve(dir)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if name != "new-name" {
		t.Errorf("got %q, want %q", name, "new-name")
	}
}

func TestIsValidName(t *testing.T) {
	cases := []struct {
		name  string
		valid bool
	}{
		{"my-project", true},
		{"my_project", true},
		{"MyProject123", true},
		{"", false},
		{"has space", false},
		{"has/slash", false},
		{"has.dot", false},
	}
	for _, tc := range cases {
		t.Run(tc.name+"_valid="+boolStr(tc.valid), func(t *testing.T) {
			if got := project.IsValidName(tc.name); got != tc.valid {
				t.Errorf("IsValidName(%q) = %v, want %v", tc.name, got, tc.valid)
			}
		})
	}
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
