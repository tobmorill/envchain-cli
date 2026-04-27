package validate_test

import (
	"errors"
	"testing"

	"github.com/envchain-cli/envchain/internal/env"
	"github.com/envchain-cli/envchain/internal/validate"
)

func entries(pairs ...string) []env.Entry {
	var out []env.Entry
	for i := 0; i+1 < len(pairs); i += 2 {
		out = append(out, env.Entry{Key: pairs[i], Value: pairs[i+1]})
	}
	return out
}

func TestNoViolationsOnEmptyRule(t *testing.T) {
	e := entries("FOO", "bar", "BAZ", "qux")
	violations := validate.Validate(e, validate.Rule{})
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestRequiredMissingKey(t *testing.T) {
	e := entries("FOO", "bar")
	r := validate.Rule{Required: []string{"FOO", "SECRET"}}
	violations := validate.Validate(e, r)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "SECRET" {
		t.Errorf("expected violation for SECRET, got %s", violations[0].Key)
	}
}

func TestForbidEmptyValue(t *testing.T) {
	e := entries("TOKEN", "", "HOST", "localhost")
	r := validate.Rule{ForbidEmpty: true}
	violations := validate.Validate(e, r)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "TOKEN" {
		t.Errorf("expected violation for TOKEN, got %s", violations[0].Key)
	}
}

func TestKeyPatternMismatch(t *testing.T) {
	e := entries("FOO_BAR", "ok", "lowercase", "bad")
	r := validate.Rule{KeyPattern: `^[A-Z][A-Z0-9_]*$`}
	violations := validate.Validate(e, r)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "lowercase" {
		t.Errorf("unexpected key %s", violations[0].Key)
	}
}

func TestInvalidKeyPatternReturnsViolation(t *testing.T) {
	e := entries("FOO", "bar")
	r := validate.Rule{KeyPattern: `[invalid(`}
	violations := validate.Validate(e, r)
	if len(violations) != 1 || violations[0].Key != "<rule>" {
		t.Fatalf("expected rule violation, got %v", violations)
	}
}

func TestCheckReturnsNilOnValid(t *testing.T) {
	e := entries("DB_URL", "postgres://localhost")
	err := validate.Check(e, validate.Rule{Required: []string{"DB_URL"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckWrapsErrValidation(t *testing.T) {
	e := entries("FOO", "bar")
	err := validate.Check(e, validate.Rule{Required: []string{"MISSING"}})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, validate.ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestMultipleViolationsCombined(t *testing.T) {
	e := entries("bad_key", "")
	r := validate.Rule{
		KeyPattern:  `^[A-Z_]+$`,
		ForbidEmpty: true,
		Required:    []string{"MUST_EXIST"},
	}
	violations := validate.Validate(e, r)
	if len(violations) != 3 {
		t.Fatalf("expected 3 violations, got %d: %v", len(violations), violations)
	}
}
