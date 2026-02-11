package week07

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestRegisterLoginValidate 验证注册登录和 token 校验。
func TestRegisterLoginValidate(t *testing.T) {
	auth := NewAuthService("demo-secret", time.Minute)
	if err := auth.Register("alice", "pass123", "user"); err != nil {
		t.Fatalf("register err: %v", err)
	}
	token, err := auth.Login("alice", "pass123")
	if err != nil {
		t.Fatalf("login err: %v", err)
	}
	claims, err := auth.ValidateToken(token)
	if err != nil {
		t.Fatalf("validate err: %v", err)
	}
	if claims.Rol != "user" {
		t.Fatalf("unexpected role: %s", claims.Rol)
	}
}

// TestAuthMiddlewareMissingToken 验证未登录返回 401。
func TestAuthMiddlewareMissingToken(t *testing.T) {
	auth := NewAuthService("demo-secret", time.Minute)
	h := AuthMiddleware(auth, "user", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	rec := httptest.NewRecorder()
	h(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", rec.Code)
	}
}

// TestAuthMiddlewareInvalidToken 验证非法 token 返回 401。
func TestAuthMiddlewareInvalidToken(t *testing.T) {
	auth := NewAuthService("demo-secret", time.Minute)
	h := AuthMiddleware(auth, "user", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	rec := httptest.NewRecorder()
	h(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", rec.Code)
	}
}

// TestAuthMiddlewareForbidden 验证角色不足返回 403。
func TestAuthMiddlewareForbidden(t *testing.T) {
	auth := NewAuthService("demo-secret", time.Minute)
	_ = auth.Register("bob", "pass123", "user")
	token, _ := auth.Login("bob", "pass123")

	h := AuthMiddleware(auth, "admin", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rec := httptest.NewRecorder()
	h(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Fatalf("expected 403 got %d", rec.Code)
	}
}
