package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	t.Setenv("APP_NAME", "")
	t.Setenv("APP_ENV", "")
	t.Setenv("HTTP_PORT", "")
	t.Setenv("PG_DSN", "")
	t.Setenv("JWT_SECRET", "change-me-test")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if cfg.App.Name != "gin-backend" {
		t.Fatalf("unexpected app name: %s", cfg.App.Name)
	}
	if cfg.HTTP.Port != "8081" {
		t.Fatalf("unexpected http port: %s", cfg.HTTP.Port)
	}
}

func TestLoadInvalidBool(t *testing.T) {
	t.Setenv("REDIS_ENABLED", "bad")
	t.Setenv("JWT_SECRET", "change-me-test")
	if _, err := Load(); err == nil {
		t.Fatal("expected error for invalid REDIS_ENABLED")
	}
}
