package logging

import (
	"testing"
)

func TestNewDevelopment(t *testing.T) {
	logger := NewDevelopment()
	if logger == nil {
		t.Fatal("expected non-nil logger")
	}
	logger.Info("test log message")
}

func TestNewProduction(t *testing.T) {
	logger := NewProduction()
	if logger == nil {
		t.Fatal("expected non-nil logger")
	}
}
