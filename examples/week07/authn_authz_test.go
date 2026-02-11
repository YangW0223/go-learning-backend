// 详细注释: package week07
package week07

// 详细注释: import (
import (
	// 详细注释: "net/http"
	"net/http"
	// 详细注释: "net/http/httptest"
	"net/http/httptest"
	// 详细注释: "testing"
	"testing"
	// 详细注释: "time"
	"time"
	// 详细注释: )
)

// TestRegisterLoginValidate 验证注册登录和 token 校验。
// 详细注释: func TestRegisterLoginValidate(t *testing.T) {
func TestRegisterLoginValidate(t *testing.T) {
	// 详细注释: auth := NewAuthService("demo-secret", time.Minute)
	auth := NewAuthService("demo-secret", time.Minute)
	// 详细注释: if err := auth.Register("alice", "pass123", "user"); err != nil {
	if err := auth.Register("alice", "pass123", "user"); err != nil {
		// 详细注释: t.Fatalf("register err: %v", err)
		t.Fatalf("register err: %v", err)
		// 详细注释: }
	}
	// 详细注释: token, err := auth.Login("alice", "pass123")
	token, err := auth.Login("alice", "pass123")
	// 详细注释: if err != nil {
	if err != nil {
		// 详细注释: t.Fatalf("login err: %v", err)
		t.Fatalf("login err: %v", err)
		// 详细注释: }
	}
	// 详细注释: claims, err := auth.ValidateToken(token)
	claims, err := auth.ValidateToken(token)
	// 详细注释: if err != nil {
	if err != nil {
		// 详细注释: t.Fatalf("validate err: %v", err)
		t.Fatalf("validate err: %v", err)
		// 详细注释: }
	}
	// 详细注释: if claims.Rol != "user" {
	if claims.Rol != "user" {
		// 详细注释: t.Fatalf("unexpected role: %s", claims.Rol)
		t.Fatalf("unexpected role: %s", claims.Rol)
		// 详细注释: }
	}
	// 详细注释: }
}

// TestAuthMiddlewareMissingToken 验证未登录返回 401。
// 详细注释: func TestAuthMiddlewareMissingToken(t *testing.T) {
func TestAuthMiddlewareMissingToken(t *testing.T) {
	// 详细注释: auth := NewAuthService("demo-secret", time.Minute)
	auth := NewAuthService("demo-secret", time.Minute)
	// 详细注释: h := AuthMiddleware(auth, "user", func(w http.ResponseWriter, _ *http.Request) {
	h := AuthMiddleware(auth, "user", func(w http.ResponseWriter, _ *http.Request) {
		// 详细注释: w.WriteHeader(http.StatusOK)
		w.WriteHeader(http.StatusOK)
		// 详细注释: })
	})

	// 详细注释: req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	// 详细注释: rec := httptest.NewRecorder()
	rec := httptest.NewRecorder()
	// 详细注释: h(rec, req)
	h(rec, req)

	// 详细注释: if rec.Code != http.StatusUnauthorized {
	if rec.Code != http.StatusUnauthorized {
		// 详细注释: t.Fatalf("expected 401 got %d", rec.Code)
		t.Fatalf("expected 401 got %d", rec.Code)
		// 详细注释: }
	}
	// 详细注释: }
}

// TestAuthMiddlewareInvalidToken 验证非法 token 返回 401。
// 详细注释: func TestAuthMiddlewareInvalidToken(t *testing.T) {
func TestAuthMiddlewareInvalidToken(t *testing.T) {
	// 详细注释: auth := NewAuthService("demo-secret", time.Minute)
	auth := NewAuthService("demo-secret", time.Minute)
	// 详细注释: h := AuthMiddleware(auth, "user", func(w http.ResponseWriter, _ *http.Request) {
	h := AuthMiddleware(auth, "user", func(w http.ResponseWriter, _ *http.Request) {
		// 详细注释: w.WriteHeader(http.StatusOK)
		w.WriteHeader(http.StatusOK)
		// 详细注释: })
	})

	// 详细注释: req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	// 详细注释: req.Header.Set("Authorization", "Bearer invalid-token")
	req.Header.Set("Authorization", "Bearer invalid-token")
	// 详细注释: rec := httptest.NewRecorder()
	rec := httptest.NewRecorder()
	// 详细注释: h(rec, req)
	h(rec, req)

	// 详细注释: if rec.Code != http.StatusUnauthorized {
	if rec.Code != http.StatusUnauthorized {
		// 详细注释: t.Fatalf("expected 401 got %d", rec.Code)
		t.Fatalf("expected 401 got %d", rec.Code)
		// 详细注释: }
	}
	// 详细注释: }
}

// TestAuthMiddlewareForbidden 验证角色不足返回 403。
// 详细注释: func TestAuthMiddlewareForbidden(t *testing.T) {
func TestAuthMiddlewareForbidden(t *testing.T) {
	// 详细注释: auth := NewAuthService("demo-secret", time.Minute)
	auth := NewAuthService("demo-secret", time.Minute)
	// 详细注释: _ = auth.Register("bob", "pass123", "user")
	_ = auth.Register("bob", "pass123", "user")
	// 详细注释: token, _ := auth.Login("bob", "pass123")
	token, _ := auth.Login("bob", "pass123")

	// 详细注释: h := AuthMiddleware(auth, "admin", func(w http.ResponseWriter, _ *http.Request) {
	h := AuthMiddleware(auth, "admin", func(w http.ResponseWriter, _ *http.Request) {
		// 详细注释: w.WriteHeader(http.StatusOK)
		w.WriteHeader(http.StatusOK)
		// 详细注释: })
	})

	// 详细注释: req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	req := httptest.NewRequest(http.MethodGet, "/secure", nil)
	// 详细注释: req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Authorization", "Bearer "+token)
	// 详细注释: rec := httptest.NewRecorder()
	rec := httptest.NewRecorder()
	// 详细注释: h(rec, req)
	h(rec, req)

	// 详细注释: if rec.Code != http.StatusForbidden {
	if rec.Code != http.StatusForbidden {
		// 详细注释: t.Fatalf("expected 403 got %d", rec.Code)
		t.Fatalf("expected 403 got %d", rec.Code)
		// 详细注释: }
	}
	// 详细注释: }
}
