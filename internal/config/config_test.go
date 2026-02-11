package config

import (
	"testing"
	"time"
)

func TestLoad_Defaults(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("REDIS_ENABLED", "")
	t.Setenv("REDIS_ADDR", "")
	t.Setenv("REDIS_PASSWORD", "")
	t.Setenv("REDIS_DB", "")
	t.Setenv("REDIS_CACHE_TTL_SECONDS", "")
	t.Setenv("REDIS_DIAL_TIMEOUT_MS", "")
	t.Setenv("REDIS_IO_TIMEOUT_MS", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.Server.Port != "8080" {
		t.Fatalf("expected default port 8080, got %q", cfg.Server.Port)
	}
	if cfg.Redis.Enabled {
		t.Fatal("expected redis disabled by default")
	}
	if cfg.Redis.Addr != "localhost:6379" {
		t.Fatalf("expected default redis addr, got %q", cfg.Redis.Addr)
	}
	if cfg.Redis.DB != 0 {
		t.Fatalf("expected default redis db 0, got %d", cfg.Redis.DB)
	}
	if cfg.Redis.CacheTTL != 30*time.Second {
		t.Fatalf("expected default cache ttl 30s, got %s", cfg.Redis.CacheTTL)
	}
}

func TestLoad_InvalidValues(t *testing.T) {
	t.Setenv("REDIS_ENABLED", "not-bool")
	if _, err := Load(); err == nil {
		t.Fatal("expected error for invalid REDIS_ENABLED")
	}

	t.Setenv("REDIS_ENABLED", "true")
	t.Setenv("REDIS_DB", "bad-int")
	if _, err := Load(); err == nil {
		t.Fatal("expected error for invalid REDIS_DB")
	}

	t.Setenv("REDIS_DB", "0")
	t.Setenv("REDIS_CACHE_TTL_SECONDS", "0")
	if _, err := Load(); err == nil {
		t.Fatal("expected error for REDIS_CACHE_TTL_SECONDS <= 0")
	}
}

func TestLoad_CustomValues(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("REDIS_ENABLED", "true")
	t.Setenv("REDIS_ADDR", "redis:6379")
	t.Setenv("REDIS_PASSWORD", "secret")
	t.Setenv("REDIS_DB", "3")
	t.Setenv("REDIS_CACHE_TTL_SECONDS", "90")
	t.Setenv("REDIS_DIAL_TIMEOUT_MS", "2000")
	t.Setenv("REDIS_IO_TIMEOUT_MS", "1500")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}

	if cfg.Server.Port != "9090" {
		t.Fatalf("expected port 9090, got %q", cfg.Server.Port)
	}
	if !cfg.Redis.Enabled {
		t.Fatal("expected redis enabled")
	}
	if cfg.Redis.Addr != "redis:6379" {
		t.Fatalf("unexpected redis addr: %q", cfg.Redis.Addr)
	}
	if cfg.Redis.Password != "secret" {
		t.Fatalf("unexpected redis password: %q", cfg.Redis.Password)
	}
	if cfg.Redis.DB != 3 {
		t.Fatalf("unexpected redis db: %d", cfg.Redis.DB)
	}
	if cfg.Redis.CacheTTL != 90*time.Second {
		t.Fatalf("unexpected redis cache ttl: %s", cfg.Redis.CacheTTL)
	}
	if cfg.Redis.DialTimeout != 2*time.Second {
		t.Fatalf("unexpected redis dial timeout: %s", cfg.Redis.DialTimeout)
	}
	if cfg.Redis.IOTimeout != 1500*time.Millisecond {
		t.Fatalf("unexpected redis io timeout: %s", cfg.Redis.IOTimeout)
	}
}

func TestLoad_RedisEnabledRequiresAddr(t *testing.T) {
	t.Setenv("REDIS_ENABLED", "true")
	t.Setenv("REDIS_ADDR", "   ")

	if _, err := Load(); err == nil {
		t.Fatal("expected error when REDIS_ENABLED=true and REDIS_ADDR empty")
	}
}
