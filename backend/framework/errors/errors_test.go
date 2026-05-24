package errors

import (
	"testing"
)

func TestNew(t *testing.T) {
	err := New(KindNotFound, "market not found")
	if err.Kind != KindNotFound {
		t.Errorf("expected KindNotFound, got %v", err.Kind)
	}
	if err.Message != "market not found" {
		t.Errorf("expected 'market not found', got %s", err.Message)
	}
}

func TestWrap(t *testing.T) {
	original := New(KindBadRequest, "bad request")
	err := Wrap(KindInternal, "wrapped", original)
	if err.Kind != KindInternal {
		t.Errorf("expected KindInternal, got %v", err.Kind)
	}
	if err.Err != original {
		t.Errorf("expected wrapped error")
	}
}

func TestIsNotFound(t *testing.T) {
	err := NotFound("x")
	if !IsNotFound(err) {
		t.Errorf("expected IsNotFound to be true")
	}
}

func TestKindOf(t *testing.T) {
	err := Unauthorized("x")
	if KindOf(err) != KindUnauthorized {
		t.Errorf("expected KindUnauthorized")
	}
}
