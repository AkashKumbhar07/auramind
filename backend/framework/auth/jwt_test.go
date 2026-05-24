package auth

import (
	"testing"
	"time"
)

func TestJWTGenerateAndValidate(t *testing.T) {
	manager := NewJWTManager("test-secret-key", time.Hour)

	token, err := manager.Generate("user-1", "test@example.com", "testuser")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	if token == "" {
		t.Fatal("expected non-empty token")
	}

	claims, err := manager.Validate(token)
	if err != nil {
		t.Fatalf("failed to validate token: %v", err)
	}

	if claims.UserID != "user-1" {
		t.Errorf("expected user-1, got %s", claims.UserID)
	}
	if claims.Email != "test@example.com" {
		t.Errorf("expected test@example.com, got %s", claims.Email)
	}
}

func TestJWTInvalidToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key", time.Hour)

	_, err := manager.Validate("invalid-token")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
}

func TestJWTExpiredToken(t *testing.T) {
	manager := NewJWTManager("test-secret-key", -time.Hour)

	token, err := manager.Generate("user-1", "test@example.com", "testuser")
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	_, err = manager.Validate(token)
	if err == nil {
		t.Fatal("expected error for expired token")
	}
}
