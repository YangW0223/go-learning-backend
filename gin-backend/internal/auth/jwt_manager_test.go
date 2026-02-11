package auth

import (
	"testing"
	"time"
)

func TestJWTManager_GenerateAndParse(t *testing.T) {
	manager := NewJWTManager("test-secret-123", "test-issuer", 10*time.Minute)
	token, err := manager.Generate("u1", "u1@example.com", "user")
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	claims, err := manager.Parse(token)
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if claims.UserID != "u1" || claims.Email != "u1@example.com" || claims.Role != "user" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

func TestJWTManager_ParseIssuerMismatch(t *testing.T) {
	managerA := NewJWTManager("test-secret-123", "issuer-a", 10*time.Minute)
	managerB := NewJWTManager("test-secret-123", "issuer-b", 10*time.Minute)
	token, err := managerA.Generate("u1", "u1@example.com", "user")
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}

	if _, err := managerB.Parse(token); err == nil {
		t.Fatal("expected issuer mismatch error")
	}
}
