package week11

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestLoadConfigMissing 验证缺少关键配置时报错。
func TestLoadConfigMissing(t *testing.T) {
	_, err := LoadConfigFromEnv(func(_ string) string { return "" })
	if !errors.Is(err, ErrMissingConfig) {
		t.Fatalf("expected ErrMissingConfig, got %v", err)
	}
}

// TestLoadConfigSuccess 验证完整配置可正常加载。
func TestLoadConfigSuccess(t *testing.T) {
	values := map[string]string{
		"APP_PORT":    "18080",
		"DB_DSN":      "postgres://demo",
		"JWT_SECRET":  "secret",
		"APP_ENV":     "test",
		"APP_VERSION": "v1.0.0",
	}
	cfg, err := LoadConfigFromEnv(func(k string) string { return values[k] })
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if cfg.Port != "18080" || cfg.Version != "v1.0.0" {
		t.Fatalf("unexpected cfg: %+v", cfg)
	}
}

// TestServerReadiness 验证 healthz 和 readyz 行为。
func TestServerReadiness(t *testing.T) {
	mux := NewServer(Config{DBDSN: "postgres://demo", JWTSecret: "secret", Version: "v1"})

	healthReq := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	healthRec := httptest.NewRecorder()
	mux.ServeHTTP(healthRec, healthReq)
	if healthRec.Code != http.StatusOK {
		t.Fatalf("healthz want 200 got %d", healthRec.Code)
	}

	readyReq := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	readyRec := httptest.NewRecorder()
	mux.ServeHTTP(readyRec, readyReq)
	if readyRec.Code != http.StatusOK {
		t.Fatalf("readyz want 200 got %d", readyRec.Code)
	}
}
