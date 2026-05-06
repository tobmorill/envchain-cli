package label_test

import (
	"errors"
	"testing"

	"github.com/envchain/envchain-cli/internal/label"
)

func newTempManager(t *testing.T) *label.Manager {
	t.Helper()
	m, err := label.New(t.TempDir())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return m
}

func TestSetAndGet(t *testing.T) {
	m := newTempManager(t)
	input := label.Labels{"env": "production", "team": "platform"}
	if err := m.Set("myproject", input); err != nil {
		t.Fatalf("Set: %v", err)
	}
	got, err := m.Get("myproject")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	for k, v := range input {
		if got[k] != v {
			t.Errorf("label[%q]: want %q, got %q", k, v, got[k])
		}
	}
}

func TestGetNotFound(t *testing.T) {
	m := newTempManager(t)
	_, err := m.Get("ghost")
	if !errors.Is(err, label.ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}

func TestSetEmptyProjectReturnsError(t *testing.T) {
	m := newTempManager(t)
	if err := m.Set("", label.Labels{"k": "v"}); err == nil {
		t.Fatal("expected error for empty project name")
	}
}

func TestSetOverwritesPrevious(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj", label.Labels{"a": "1"})
	_ = m.Set("proj", label.Labels{"b": "2"})
	got, err := m.Get("proj")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if _, ok := got["a"]; ok {
		t.Error("old label 'a' should have been replaced")
	}
	if got["b"] != "2" {
		t.Errorf("want b=2, got %q", got["b"])
	}
}

func TestDelete(t *testing.T) {
	m := newTempManager(t)
	_ = m.Set("proj", label.Labels{"x": "y"})
	if err := m.Delete("proj"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := m.Get("proj")
	if !errors.Is(err, label.ErrNotFound) {
		t.Fatalf("want ErrNotFound after delete, got %v", err)
	}
}

func TestDeleteNoop(t *testing.T) {
	m := newTempManager(t)
	if err := m.Delete("nonexistent"); err != nil {
		t.Fatalf("Delete on missing project should not error: %v", err)
	}
}
