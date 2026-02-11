// 详细注释: package week11
package week11

// 详细注释: import (
import (
	// 详细注释: "errors"
	"errors"
	// 详细注释: "net/http"
	"net/http"
	// 详细注释: "net/http/httptest"
	"net/http/httptest"
	// 详细注释: "testing"
	"testing"
	// 详细注释: )
)

// TestLoadConfigMissing 验证缺少关键配置时报错。
// 详细注释: func TestLoadConfigMissing(t *testing.T) {
func TestLoadConfigMissing(t *testing.T) {
	// 详细注释: _, err := LoadConfigFromEnv(func(_ string) string { return "" })
	_, err := LoadConfigFromEnv(func(_ string) string { return "" })
	// 详细注释: if !errors.Is(err, ErrMissingConfig) {
	if !errors.Is(err, ErrMissingConfig) {
		// 详细注释: t.Fatalf("expected ErrMissingConfig, got %v", err)
		t.Fatalf("expected ErrMissingConfig, got %v", err)
		// 详细注释: }
	}
	// 详细注释: }
}

// TestLoadConfigSuccess 验证完整配置可正常加载。
// 详细注释: func TestLoadConfigSuccess(t *testing.T) {
func TestLoadConfigSuccess(t *testing.T) {
	// 详细注释: values := map[string]string{
	values := map[string]string{
		// 详细注释: "APP_PORT":    "18080",
		"APP_PORT": "18080",
		// 详细注释: "DB_DSN":      "postgres://demo",
		"DB_DSN": "postgres://demo",
		// 详细注释: "JWT_SECRET":  "secret",
		"JWT_SECRET": "secret",
		// 详细注释: "APP_ENV":     "test",
		"APP_ENV": "test",
		// 详细注释: "APP_VERSION": "v1.0.0",
		"APP_VERSION": "v1.0.0",
		// 详细注释: }
	}
	// 详细注释: cfg, err := LoadConfigFromEnv(func(k string) string { return values[k] })
	cfg, err := LoadConfigFromEnv(func(k string) string { return values[k] })
	// 详细注释: if err != nil {
	if err != nil {
		// 详细注释: t.Fatalf("unexpected err: %v", err)
		t.Fatalf("unexpected err: %v", err)
		// 详细注释: }
	}
	// 详细注释: if cfg.Port != "18080" || cfg.Version != "v1.0.0" {
	if cfg.Port != "18080" || cfg.Version != "v1.0.0" {
		// 详细注释: t.Fatalf("unexpected cfg: %+v", cfg)
		t.Fatalf("unexpected cfg: %+v", cfg)
		// 详细注释: }
	}
	// 详细注释: }
}

// TestServerReadiness 验证 healthz 和 readyz 行为。
// 详细注释: func TestServerReadiness(t *testing.T) {
func TestServerReadiness(t *testing.T) {
	// 详细注释: mux := NewServer(Config{DBDSN: "postgres://demo", JWTSecret: "secret", Version: "v1"})
	mux := NewServer(Config{DBDSN: "postgres://demo", JWTSecret: "secret", Version: "v1"})

	// 详细注释: healthReq := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	healthReq := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	// 详细注释: healthRec := httptest.NewRecorder()
	healthRec := httptest.NewRecorder()
	// 详细注释: mux.ServeHTTP(healthRec, healthReq)
	mux.ServeHTTP(healthRec, healthReq)
	// 详细注释: if healthRec.Code != http.StatusOK {
	if healthRec.Code != http.StatusOK {
		// 详细注释: t.Fatalf("healthz want 200 got %d", healthRec.Code)
		t.Fatalf("healthz want 200 got %d", healthRec.Code)
		// 详细注释: }
	}

	// 详细注释: readyReq := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	readyReq := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	// 详细注释: readyRec := httptest.NewRecorder()
	readyRec := httptest.NewRecorder()
	// 详细注释: mux.ServeHTTP(readyRec, readyReq)
	mux.ServeHTTP(readyRec, readyReq)
	// 详细注释: if readyRec.Code != http.StatusOK {
	if readyRec.Code != http.StatusOK {
		// 详细注释: t.Fatalf("readyz want 200 got %d", readyRec.Code)
		t.Fatalf("readyz want 200 got %d", readyRec.Code)
		// 详细注释: }
	}
	// 详细注释: }
}
