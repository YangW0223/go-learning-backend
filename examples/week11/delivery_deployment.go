package week11

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	// ErrMissingConfig 表示关键配置缺失。
	ErrMissingConfig = errors.New("missing required config")
)

// Config 表示部署时由环境变量注入的配置。
type Config struct {
	Port      string
	DBDSN     string
	JWTSecret string
	Env       string
	Version   string
}

// LoadConfigFromEnv 从 getenv 读取配置。
func LoadConfigFromEnv(getenv func(string) string) (Config, error) {
	cfg := Config{
		Port:      fallback(getenv("APP_PORT"), "8080"),
		DBDSN:     strings.TrimSpace(getenv("DB_DSN")),
		JWTSecret: strings.TrimSpace(getenv("JWT_SECRET")),
		Env:       fallback(getenv("APP_ENV"), "dev"),
		Version:   fallback(getenv("APP_VERSION"), "unknown"),
	}
	if cfg.DBDSN == "" || cfg.JWTSecret == "" {
		return Config{}, ErrMissingConfig
	}
	return cfg, nil
}

// NewServer 构建最小可交付服务。
// /healthz 反映进程存活，/readyz 反映依赖配置是否就绪。
func NewServer(cfg Config) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		if cfg.DBDSN == "" || cfg.JWTSecret == "" {
			http.Error(w, "not ready", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))
	})
	mux.HandleFunc("/version", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(cfg.Version))
	})
	return mux
}

// BuildRollbackPlan 生成最小回滚说明。
func BuildRollbackPlan(currentImage, previousImage string) string {
	return fmt.Sprintf("rollback from %s to %s", currentImage, previousImage)
}

func fallback(value, def string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return def
	}
	return value
}
